package pkg

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/doguregistry"
	kubernetes2 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/dogucr"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/maintenance"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintV2"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/config"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

// ApplicationContext contains vital application parts for this operator.
type ApplicationContext struct {
	RemoteDoguRegistry             domainservice.RemoteDoguRegistry
	DoguInstallationRepository     domainservice.DoguInstallationRepository
	BlueprintSpecRepository        domainservice.BlueprintSpecRepository
	BlueprintSpecDomainUseCase     *domainservice.ValidateDependenciesDomainUseCase
	DoguInstallationUseCase        *application.DoguInstallationUseCase
	BlueprintSpecChangeUseCase     *application.BlueprintSpecChangeUseCase
	BlueprintSpecValidationUseCase *application.BlueprintSpecValidationUseCase
	EffectiveBlueprintUseCase      *application.EffectiveBlueprintUseCase
	StateDiffUseCase               *application.StateDiffUseCase
	MaintenanceModeUseCase         *domainservice.MaintenanceModeUseCase
	BlueprintSerializer            serializer.BlueprintSerializer
	BlueprintMaskSerializer        serializer.BlueprintMaskSerializer
	Reconciler                     *reconciler.BlueprintReconciler
}

// Bootstrap creates the ApplicationContext.
func Bootstrap(restConfig *rest.Config, eventRecorder record.EventRecorder, namespace string) (*ApplicationContext, error) {
	blueprintSerializer := blueprintV2.Serializer{}
	blueprintMaskSerializer := blueprintMaskV1.Serializer{}
	ecosystemClientSet, err := createEcosystemClientSet(restConfig)
	if err != nil {
		return nil, err
	}

	dogusInterface, err := ecoSystem.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dogus interface: %w", err)
	}

	blueprintSpecRepository := kubernetes2.NewBlueprintSpecRepository(
		ecosystemClientSet.EcosystemV1Alpha1().Blueprints(namespace),
		blueprintSerializer,
		blueprintMaskSerializer,
		eventRecorder,
	)

	configRegistry, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", namespace)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CES configuration registry: %w", err)
	}

	maintenanceMode := maintenance.NewSwitch(configRegistry.GlobalConfig())
	maintenanceUseCase := domainservice.NewMaintenanceModeUseCase(maintenanceMode)

	remoteDoguRegistry, err := createRemoteDoguRegistry()
	if err != nil {
		return nil, err
	}

	doguInstallationRepo := dogucr.NewDoguInstallationRepo(dogusInterface.Dogus(namespace))

	blueprintSpecDomainUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	blueprintValidationUseCase := application.NewBlueprintSpecValidationUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase)
	effectiveBlueprintUseCase := application.NewEffectiveBlueprintUseCase(blueprintSpecRepository)
	stateDiffUseCase := application.NewStateDiffUseCase(blueprintSpecRepository, doguInstallationRepo)
	doguInstallationUseCase := application.NewDoguInstallationUseCase(blueprintSpecRepository, doguInstallationRepo)
	blueprintChangeUseCase := application.NewBlueprintSpecChangeUseCase(
		blueprintSpecRepository, blueprintValidationUseCase,
		effectiveBlueprintUseCase, stateDiffUseCase,
		doguInstallationUseCase,
	)
	blueprintReconciler := reconciler.NewBlueprintReconciler(blueprintChangeUseCase)

	return &ApplicationContext{
		RemoteDoguRegistry:             remoteDoguRegistry,
		DoguInstallationRepository:     doguInstallationRepo,
		BlueprintSpecRepository:        blueprintSpecRepository,
		BlueprintSpecDomainUseCase:     blueprintSpecDomainUseCase,
		BlueprintSpecChangeUseCase:     blueprintChangeUseCase,
		BlueprintSpecValidationUseCase: blueprintValidationUseCase,
		EffectiveBlueprintUseCase:      effectiveBlueprintUseCase,
		StateDiffUseCase:               stateDiffUseCase,
		MaintenanceModeUseCase:         maintenanceUseCase,
		DoguInstallationUseCase:        doguInstallationUseCase,
		BlueprintSerializer:            blueprintSerializer,
		BlueprintMaskSerializer:        blueprintMaskSerializer,
		Reconciler:                     blueprintReconciler,
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
