package pkg

import (
	"fmt"
	adapterconfigk8s "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/config/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/restartcr"
	"github.com/cloudogu/k8s-registry-lib/repository"
	remotedogudescriptor "github.com/cloudogu/remote-dogu-descriptor-lib/repository"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	adapterk8s "github.com/cloudogu/k8s-blueprint-lib/client"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/doguregistry"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/componentcr"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/dogucr"
	adapterhealthconfig "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/healthConfig"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/maintenance"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer/blueprintV2"
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

var maintenanceModeOwner = "blueprint-operator"

// Bootstrap creates the ApplicationContext and does all dependency injection of the whole application.
func Bootstrap(restConfig *rest.Config, eventRecorder record.EventRecorder, namespace string) (*ApplicationContext, error) {
	blueprintSerializer := blueprintV2.Serializer{}
	blueprintMaskSerializer := blueprintMaskV1.Serializer{}
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
	blueprintSpecRepository := blueprintcr.NewBlueprintSpecRepository(
		ecosystemClientSet.EcosystemV1Alpha1().Blueprints(namespace),
		blueprintSerializer,
		blueprintMaskSerializer,
		eventRecorder,
	)

	libMaintenanceAdapter := repository.NewMaintenanceModeAdapter(maintenanceModeOwner, ecosystemClientSet.CoreV1().ConfigMaps(namespace))
	maintenanceMode := maintenance.NewMaintenanceModeAdapter(libMaintenanceAdapter)

	remoteDoguRegistry, err := createRemoteDoguRegistry()
	if err != nil {
		return nil, err
	}

	k8sDoguConfigRepo := repository.NewDoguConfigRepository(ecosystemClientSet.CoreV1().ConfigMaps(namespace))
	doguConfigRepo := adapterconfigk8s.NewDoguConfigRepository(*k8sDoguConfigRepo)
	k8sSensitiveDoguConfigRepo := repository.NewSensitiveDoguConfigRepository(ecosystemClientSet.CoreV1().Secrets(namespace))
	sensitiveDoguConfigRepo := adapterconfigk8s.NewSensitiveDoguConfigRepository(*k8sSensitiveDoguConfigRepo)
	k8sGlobalConfigRepo := repository.NewGlobalConfigRepository(ecosystemClientSet.CoreV1().ConfigMaps(namespace))
	globalConfigRepoAdapter := adapterconfigk8s.NewGlobalConfigRepository(*k8sGlobalConfigRepo)

	doguInstallationRepo := dogucr.NewDoguInstallationRepo(dogusInterface.Dogus(namespace))
	componentInstallationRepo := componentcr.NewComponentInstallationRepo(componentsInterface.Components(namespace))
	healthConfigRepo := adapterhealthconfig.NewHealthConfigProvider(ecosystemClientSet.CoreV1().ConfigMaps(namespace))
	doguRestartAdapter := dogusInterface.DoguRestarts(namespace)
	restartRepository := restartcr.NewDoguRestartRepository(doguRestartAdapter)

	validateDependenciesUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	validateMountsUseCase := domainservice.NewValidateAdditionalMountsDomainUseCase(remoteDoguRegistry)
	blueprintValidationUseCase := application.NewBlueprintSpecValidationUseCase(blueprintSpecRepository, validateDependenciesUseCase, validateMountsUseCase)
	effectiveBlueprintUseCase := application.NewEffectiveBlueprintUseCase(blueprintSpecRepository)
	stateDiffUseCase := application.NewStateDiffUseCase(blueprintSpecRepository, doguInstallationRepo, componentInstallationRepo, globalConfigRepoAdapter, doguConfigRepo, sensitiveDoguConfigRepo)
	doguInstallationUseCase := application.NewDoguInstallationUseCase(blueprintSpecRepository, doguInstallationRepo, healthConfigRepo)
	componentInstallationUseCase := application.NewComponentInstallationUseCase(blueprintSpecRepository, componentInstallationRepo, healthConfigRepo)
	ecosystemHealthUseCase := application.NewEcosystemHealthUseCase(doguInstallationUseCase, componentInstallationUseCase, healthConfigRepo)
	applyBlueprintSpecUseCase := application.NewApplyBlueprintSpecUseCase(blueprintSpecRepository, doguInstallationUseCase, ecosystemHealthUseCase, componentInstallationUseCase, maintenanceMode)
	ConfigUseCase := application.NewEcosystemConfigUseCase(blueprintSpecRepository, doguConfigRepo, sensitiveDoguConfigRepo, globalConfigRepoAdapter)
	doguRestartUseCase := application.NewDoguRestartUseCase(doguInstallationRepo, blueprintSpecRepository, restartRepository)

	selfUpgradeUseCase := application.NewSelfUpgradeUseCase(blueprintSpecRepository, componentInstallationRepo, componentInstallationUseCase, blueprintOperatorName.SimpleName, healthConfigRepo)

	blueprintChangeUseCase := application.NewBlueprintSpecChangeUseCase(
		blueprintSpecRepository, blueprintValidationUseCase,
		effectiveBlueprintUseCase, stateDiffUseCase,
		applyBlueprintSpecUseCase, ConfigUseCase,
		doguRestartUseCase,
		selfUpgradeUseCase,
	)
	blueprintReconciler := reconciler.NewBlueprintReconciler(blueprintChangeUseCase)

	return &ApplicationContext{
		Reconciler: blueprintReconciler,
	}, nil
}

func createRemoteDoguRegistry() (*doguregistry.Remote, error) {
	remoteConfig, err := config.GetRemoteConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote dogu registry config: %w", err)
	}

	remoteCreds, err := config.GetRemoteCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote dogu registry credentials: %w", err)
	}

	doguRemoteRepository, err := remotedogudescriptor.NewRemoteDoguDescriptorRepository(remoteConfig, remoteCreds)
	if err != nil {
		return nil, fmt.Errorf("failed to create new remote dogu repository: %w", err)
	}

	return doguregistry.NewRemote(doguRemoteRepository), nil
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
