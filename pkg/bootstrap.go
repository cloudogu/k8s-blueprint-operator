package pkg

import (
	"fmt"

	adapterconfigk8s "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/config/kubernetes"
	v2 "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/debugmodecr"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/sensitiveconfigref"
	"github.com/cloudogu/k8s-registry-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/repository"
	remotedogudescriptor "github.com/cloudogu/remote-dogu-descriptor-lib/repository"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	adapterk8s "github.com/cloudogu/k8s-blueprint-lib/v2/client"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/doguregistry"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/dogucr"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/config"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	debugModeClient "github.com/cloudogu/k8s-debug-mode-cr-lib/pkg/client/v1"
	doguEcoClient "github.com/cloudogu/k8s-dogu-lib/v2/client"
)

// ApplicationContext contains vital application parts for this operator.
type ApplicationContext struct {
	Reconciler *reconciler.BlueprintReconciler
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
	debugModeInterface, err := debugModeClient.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create debug mode interface: %w", err)
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
	debugModeRepo := debugmodecr.NewDebugModeRepo(debugModeInterface.DebugMode(namespace))

	initialBlueprintStateUseCase := application.NewInitiateBlueprintStatusUseCase(blueprintRepo)
	validateDependenciesUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	validateMountsUseCase := domainservice.NewValidateAdditionalMountsDomainUseCase(remoteDoguRegistry)
	blueprintValidationUseCase := application.NewBlueprintSpecValidationUseCase(blueprintRepo, validateDependenciesUseCase, validateMountsUseCase)
	effectiveBlueprintUseCase := application.NewEffectiveBlueprintUseCase(blueprintRepo, debugModeRepo)
	stateDiffUseCase := application.NewStateDiffUseCase(blueprintRepo, doguRepo, globalConfigRepoAdapter, doguConfigRepo, sensitiveDoguConfigRepo, sensitiveConfigRefReader)
	doguInstallationUseCase := application.NewDoguInstallationUseCase(blueprintRepo, doguRepo, doguConfigRepo, globalConfigRepoAdapter)
	ecosystemHealthUseCase := application.NewEcosystemHealthUseCase(doguInstallationUseCase, blueprintRepo)
	completeBlueprintSpecUseCase := application.NewCompleteBlueprintUseCase(blueprintRepo)
	applyDogusUseCase := application.NewApplyDogusUseCase(blueprintRepo, doguInstallationUseCase)
	ConfigUseCase := application.NewEcosystemConfigUseCase(blueprintRepo, doguConfigRepo, sensitiveDoguConfigRepo, globalConfigRepoAdapter, doguRepo)
	dogusUpToDateUseCase := application.NewDogusUpToDateUseCase(blueprintRepo, doguInstallationUseCase)

	preparationUseCases := application.NewBlueprintPreparationUseCases(
		initialBlueprintStateUseCase,
		blueprintValidationUseCase,
		effectiveBlueprintUseCase,
		stateDiffUseCase,
		ecosystemHealthUseCase,
	)
	applyUseCases := application.NewBlueprintApplyUseCases(
		completeBlueprintSpecUseCase,
		ConfigUseCase,
		applyDogusUseCase,
		ecosystemHealthUseCase,
		dogusUpToDateUseCase,
	)
	blueprintChangeUseCase := application.NewBlueprintSpecChangeUseCase(blueprintRepo, preparationUseCases, applyUseCases)
	debounceWindow, err := config.GetDebounceWindow()
	if err != nil {
		return nil, err
	}
	blueprintReconciler := reconciler.NewBlueprintReconciler(blueprintChangeUseCase, blueprintRepo, namespace, debounceWindow)

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
