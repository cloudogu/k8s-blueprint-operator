package pkg

import (
	"fmt"

	adapterconfigk8s "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/config/kubernetes"
	v2 "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/sensitiveconfigref"
	"github.com/cloudogu/k8s-registry-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/repository"
	remotedogudescriptor "github.com/cloudogu/remote-dogu-descriptor-lib/repository"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	adapterk8s "github.com/cloudogu/k8s-blueprint-lib/v2/client"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/doguregistry"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/componentcr"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/dogucr"
	adapterhealthconfig "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/healthConfig"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/config"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	componentEcoClient "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	doguEcoClient "github.com/cloudogu/k8s-dogu-lib/v2/client"
)

// ApplicationContext contains vital application parts for this operator.
type ApplicationContext struct {
	Reconciler *reconciler.BlueprintReconciler
}

var blueprintOperatorName = common.QualifiedComponentName{
	Namespace:  "k8s",
	SimpleName: "k8s-blueprint-operator",
}

// Bootstrap creates the ApplicationContext and does all dependency injection of the whole application.
func Bootstrap(restConfig *rest.Config, eventRecorder record.EventRecorder, namespace string) (*ApplicationContext, error) {
	ecosystemClientSet, err := createEcosystemClientSet(restConfig)
	if err != nil {
		return nil, err
	}

	dogusInterface, err := doguEcoClient.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dogus interface: %w", err)
	}
	componentsInterface, err := componentEcoClient.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create components interface: %w", err)
	}
	blueprintRepo := v2.NewBlueprintSpecRepository(
		ecosystemClientSet.EcosystemV1Alpha1().Blueprints(namespace),
		eventRecorder,
	)

	remoteDoguRegistry, err := createRemoteDoguRegistry(ecosystemClientSet, namespace)
	if err != nil {
		return nil, err
	}

	k8sDoguConfigRepo := repository.NewDoguConfigRepository(ecosystemClientSet.CoreV1().ConfigMaps(namespace))
	doguConfigRepo := adapterconfigk8s.NewDoguConfigRepository(*k8sDoguConfigRepo)
	k8sSensitiveDoguConfigRepo := repository.NewSensitiveDoguConfigRepository(ecosystemClientSet.CoreV1().Secrets(namespace))
	sensitiveDoguConfigRepo := adapterconfigk8s.NewSensitiveDoguConfigRepository(*k8sSensitiveDoguConfigRepo)
	sensitiveConfigRefReader := sensitiveconfigref.NewSecretRefReader(ecosystemClientSet.CoreV1().Secrets(namespace))
	k8sGlobalConfigRepo := repository.NewGlobalConfigRepository(ecosystemClientSet.CoreV1().ConfigMaps(namespace))
	globalConfigRepoAdapter := adapterconfigk8s.NewGlobalConfigRepository(*k8sGlobalConfigRepo)

	doguRepo := dogucr.NewDoguInstallationRepo(dogusInterface.Dogus(namespace))
	componentRepo := componentcr.NewComponentInstallationRepo(componentsInterface.Components(namespace))
	healthConfigRepo := adapterhealthconfig.NewHealthConfigProvider(ecosystemClientSet.CoreV1().ConfigMaps(namespace))

	initialBlueprintStateUseCase := application.NewInitiateBlueprintStatusUseCase(blueprintRepo)
	validateDependenciesUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	validateMountsUseCase := domainservice.NewValidateAdditionalMountsDomainUseCase(remoteDoguRegistry)
	blueprintValidationUseCase := application.NewBlueprintSpecValidationUseCase(blueprintRepo, validateDependenciesUseCase, validateMountsUseCase)
	effectiveBlueprintUseCase := application.NewEffectiveBlueprintUseCase(blueprintRepo)
	stateDiffUseCase := application.NewStateDiffUseCase(blueprintRepo, doguRepo, componentRepo, globalConfigRepoAdapter, doguConfigRepo, sensitiveDoguConfigRepo, sensitiveConfigRefReader)
	doguInstallationUseCase := application.NewDoguInstallationUseCase(blueprintRepo, doguRepo, healthConfigRepo)
	componentInstallationUseCase := application.NewComponentInstallationUseCase(blueprintRepo, componentRepo, healthConfigRepo)
	ecosystemHealthUseCase := application.NewEcosystemHealthUseCase(doguInstallationUseCase, componentInstallationUseCase, blueprintRepo)
	applyBlueprintSpecUseCase := application.NewCompleteBlueprintUseCase(blueprintRepo)
	applyComponentUseCase := application.NewApplyComponentsUseCase(blueprintRepo, componentInstallationUseCase)
	applyDogusUseCase := application.NewApplyDogusUseCase(blueprintRepo, doguInstallationUseCase)
	ConfigUseCase := application.NewEcosystemConfigUseCase(blueprintRepo, doguConfigRepo, sensitiveDoguConfigRepo, globalConfigRepoAdapter)
	selfUpgradeUseCase := application.NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentInstallationUseCase, blueprintOperatorName.SimpleName)

	blueprintChangeUseCase := application.NewBlueprintSpecChangeUseCase(
		blueprintRepo, initialBlueprintStateUseCase, blueprintValidationUseCase,
		effectiveBlueprintUseCase, stateDiffUseCase,
		applyBlueprintSpecUseCase, ConfigUseCase,
		selfUpgradeUseCase,
		applyComponentUseCase,
		applyDogusUseCase,
		ecosystemHealthUseCase,
	)
	blueprintReconciler := reconciler.NewBlueprintReconciler(blueprintChangeUseCase)

	return &ApplicationContext{
		Reconciler: blueprintReconciler,
	}, nil
}

func createRemoteDoguRegistry(clientSet *adapterk8s.ClientSet, namespace string) (*doguregistry.DoguDescriptorRepository, error) {
	remoteConfig, err := config.GetRemoteConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote dogu registry config: %w", err)
	}

	remoteCreds, err := config.GetRemoteCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote dogu registry credentials: %w", err)
	}

	doguRemoteRepository, err := remotedogudescriptor.NewRemoteDoguDescriptorRepository(remoteConfig, remoteCreds)
	doguLocalRepository := dogu.NewLocalDoguDescriptorRepository(clientSet.CoreV1().ConfigMaps(namespace))
	if err != nil {
		return nil, fmt.Errorf("failed to create new remote dogu repository: %w", err)
	}
	return doguregistry.NewDoguDescriptorRepository(doguRemoteRepository, doguLocalRepository), nil
}

func createEcosystemClientSet(restConfig *rest.Config) (*adapterk8s.ClientSet, error) {
	k8sClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s clientset: %w", err)
	}

	ecosystemClientSet, err := adapterk8s.NewClientSet(restConfig, k8sClientSet)
	if err != nil {
		return nil, fmt.Errorf("unable to create ecosystem clientset: %w", err)
	}
	return ecosystemClientSet, nil
}
