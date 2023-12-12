package pkg

import (
	"fmt"
	kubernetes2 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/reconciler"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintV2"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

// ApplicationContext contains vital application parts for this operator.
type ApplicationContext struct {
	RemoteDoguRegistry         domainservice.RemoteDoguRegistry
	BlueprintSpecRepository    domainservice.BlueprintSpecRepository
	BlueprintSpecDomainUseCase *domainservice.ValidateDependenciesDomainUseCase
	DoguInstallationUseCase    *application.DoguInstallationUseCase
	BlueprintSpecUseCase       *application.BlueprintSpecUseCase
	BlueprintSerializer        serializer.BlueprintSerializer
	BlueprintMaskSerializer    serializer.BlueprintMaskSerializer
	Reconciler                 *reconciler.BlueprintReconciler
}

// Bootstrap creates the ApplicationContext.
func Bootstrap(restConfig *rest.Config, eventRecorder record.EventRecorder, namespace string) (*ApplicationContext, error) {
	blueprintSerializer := blueprintV2.Serializer{}
	blueprintMaskSerializer := blueprintMaskV1.Serializer{}
	ecosystemClientSet, err := createEcosystemClientSet(restConfig)
	if err != nil {
		return nil, err
	}

	blueprintSpecRepository := kubernetes2.NewBlueprintSpecRepository(
		ecosystemClientSet.EcosystemV1Alpha1().Blueprints(namespace),
		blueprintSerializer,
		blueprintMaskSerializer,
		eventRecorder,
	)
	var remoteDoguRegistry domainservice.RemoteDoguRegistry
	blueprintSpecDomainUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	doguInstallationUseCase := &application.DoguInstallationUseCase{}
	blueprintUseCase := application.NewBlueprintSpecUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase, doguInstallationUseCase)
	blueprintReconciler := reconciler.NewBlueprintReconciler(blueprintUseCase)

	return &ApplicationContext{
		RemoteDoguRegistry:         remoteDoguRegistry,
		BlueprintSpecRepository:    blueprintSpecRepository,
		BlueprintSpecDomainUseCase: blueprintSpecDomainUseCase,
		DoguInstallationUseCase:    doguInstallationUseCase,
		BlueprintSpecUseCase:       blueprintUseCase,
		BlueprintSerializer:        blueprintSerializer,
		BlueprintMaskSerializer:    blueprintMaskSerializer,
		Reconciler:                 blueprintReconciler,
	}, nil
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
