package pkg

import (
	"fmt"
	configAdapter "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/remote"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/doguregistry"
	kubernetes2 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/blueprintcr"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/componentcr"
	config2 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/config"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/dogucr"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/maintenance"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintV2"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/config"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	componentEcoClient "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	doguEcoClient "github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
)

// ApplicationContext contains vital application parts for this operator.
type ApplicationContext struct {
	RemoteDoguRegistry             domainservice.RemoteDoguRegistry
	DoguInstallationRepository     domainservice.DoguInstallationRepository
	BlueprintSpecRepository        domainservice.BlueprintSpecRepository
	BlueprintSpecDomainUseCase     *domainservice.ValidateDependenciesDomainUseCase
	DoguInstallationUseCase        *application.DoguInstallationUseCase
	EcosystemHealthUseCase         *application.EcosystemHealthUseCase
	ApplyBlueprintSpecUseCase      *application.ApplyBlueprintSpecUseCase
	BlueprintSpecChangeUseCase     *application.BlueprintSpecChangeUseCase
	BlueprintSpecValidationUseCase *application.BlueprintSpecValidationUseCase
	EffectiveBlueprintUseCase      *application.EffectiveBlueprintUseCase
	StateDiffUseCase               *application.StateDiffUseCase
	BlueprintSerializer            serializer.BlueprintSerializer
	BlueprintMaskSerializer        serializer.BlueprintMaskSerializer
	Reconciler                     *reconciler.BlueprintReconciler
	configEncryptionAdapter        domainservice.ConfigEncryptionAdapter
	doguConfigAdapter              domainservice.DoguConfigEntryRepository
	sensitiveDoguConfigAdapter     domainservice.SensitiveDoguConfigEntryRepository
	globalConfigAdapter            domainservice.GlobalConfigEntryRepository
}

// Bootstrap creates the ApplicationContext.
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

	configRegistry, err := createConfigRegistry(namespace)
	if err != nil {
		return nil, err
	}

	maintenanceMode := maintenance.New(configRegistry.GlobalConfig())

	remoteDoguRegistry, err := createRemoteDoguRegistry()
	if err != nil {
		return nil, err
	}

	configEncryptionAdapter := configAdapter.NewPublicKeyConfigEncryptionAdapter(ecosystemClientSet.CoreV1().Secrets(namespace), configRegistry, namespace)
	doguConfigAdapter := configAdapter.NewEtcdDoguConfigRepository(configRegistry)
	sensitiveDoguConfigAdapter := configAdapter.NewEtcdSensitiveDoguConfigRepository(configRegistry)
	globalConfigAdapter := configAdapter.NewEtcdGlobalConfigRepository(configRegistry.GlobalConfig())

	doguInstallationRepo := dogucr.NewDoguInstallationRepo(dogusInterface.Dogus(namespace), ecosystemClientSet.CoreV1().PersistentVolumeClaims(namespace))
	componentInstallationRepo := componentcr.NewComponentInstallationRepo(componentsInterface.Components(namespace))
	healthConfigRepo := config2.NewHealthConfigProvider(ecosystemClientSet.CoreV1().ConfigMaps(namespace))

	blueprintSpecDomainUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	blueprintValidationUseCase := application.NewBlueprintSpecValidationUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase)
	effectiveBlueprintUseCase := application.NewEffectiveBlueprintUseCase(blueprintSpecRepository)
	stateDiffUseCase := application.NewStateDiffUseCase(blueprintSpecRepository, doguInstallationRepo, componentInstallationRepo, globalConfigAdapter, doguConfigAdapter, sensitiveDoguConfigAdapter, configEncryptionAdapter)
	doguInstallationUseCase := application.NewDoguInstallationUseCase(blueprintSpecRepository, doguInstallationRepo, healthConfigRepo)
	componentInstallationUseCase := application.NewComponentInstallationUseCase(blueprintSpecRepository, componentInstallationRepo, healthConfigRepo)
	ecosystemHealthUseCase := application.NewEcosystemHealthUseCase(doguInstallationUseCase, componentInstallationUseCase, healthConfigRepo)
	applyBlueprintSpecUseCase := application.NewApplyBlueprintSpecUseCase(blueprintSpecRepository, doguInstallationUseCase, ecosystemHealthUseCase, componentInstallationUseCase, maintenanceMode)
	registryConfigUseCase := application.NewEcosystemRegistryUseCase(blueprintSpecRepository, doguConfigAdapter, sensitiveDoguConfigAdapter, globalConfigAdapter)

	blueprintChangeUseCase := application.NewBlueprintSpecChangeUseCase(
		blueprintSpecRepository, blueprintValidationUseCase,
		effectiveBlueprintUseCase, stateDiffUseCase,
		applyBlueprintSpecUseCase, registryConfigUseCase,
	)
	blueprintReconciler := reconciler.NewBlueprintReconciler(blueprintChangeUseCase)

	return &ApplicationContext{
		RemoteDoguRegistry:             remoteDoguRegistry,
		DoguInstallationRepository:     doguInstallationRepo,
		BlueprintSpecRepository:        blueprintSpecRepository,
		BlueprintSpecDomainUseCase:     blueprintSpecDomainUseCase,
		EcosystemHealthUseCase:         ecosystemHealthUseCase,
		ApplyBlueprintSpecUseCase:      applyBlueprintSpecUseCase,
		BlueprintSpecChangeUseCase:     blueprintChangeUseCase,
		BlueprintSpecValidationUseCase: blueprintValidationUseCase,
		EffectiveBlueprintUseCase:      effectiveBlueprintUseCase,
		StateDiffUseCase:               stateDiffUseCase,
		DoguInstallationUseCase:        doguInstallationUseCase,
		BlueprintSerializer:            blueprintSerializer,
		BlueprintMaskSerializer:        blueprintMaskSerializer,
		Reconciler:                     blueprintReconciler,
		configEncryptionAdapter:        configEncryptionAdapter,
		doguConfigAdapter:              doguConfigAdapter,
		sensitiveDoguConfigAdapter:     sensitiveDoguConfigAdapter,
		globalConfigAdapter:            globalConfigAdapter,
	}, nil
}

func createConfigRegistry(namespace string) (registry.Registry, error) {
	configRegistry, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", namespace)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CES configuration registry: %w", err)
	}

	return configRegistry, nil
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

	doguRemoteRegistry, err := remote.New(remoteConfig, remoteCreds)
	if err != nil {
		return nil, fmt.Errorf("failed to create new remote dogu registry: %w", err)
	}

	return doguregistry.NewRemote(doguRemoteRegistry), nil
}

func createEcosystemClientSet(restConfig *rest.Config) (*kubernetes2.ClientSet, error) {
	k8sClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s clientset: %w", err)
	}

	ecosystemClientSet, err := kubernetes2.NewClientSet(restConfig, k8sClientSet)
	if err != nil {
		return nil, fmt.Errorf("unable to create ecosystem clientset: %w", err)
	}
	return ecosystemClientSet, nil
}
