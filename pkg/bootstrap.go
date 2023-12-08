package pkg

import (
	"fmt"
	kubernetes2 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintV2"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
}

// Bootstrap creates the ApplicationContext.
func Bootstrap(restConfig *rest.Config, namespace string) (*ApplicationContext, error) {
	k8sClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s clientset: %w", err)
	}

	//TODO where will be the eventRecorder interface be located? here? Events must probably be written overall the operator.

	ecosystemClient, err := kubernetes2.NewClientSet(restConfig, k8sClientSet)
	if err != nil {
		return nil, fmt.Errorf("unable to create ecosystem client: %w", err)
	}

	blueprintSerializer := blueprintV2.Serializer{}
	blueprintMaskSerializer := blueprintMaskV1.Serializer{}

	var remoteDoguRegistry domainservice.RemoteDoguRegistry
	blueprintSpecRepository := kubernetes2.NewBlueprintSpecRepository(namespace, ecosystemClient)
	blueprintSpecDomainUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	doguInstallationUseCase := &application.DoguInstallationUseCase{}
	blueprintUseCase := application.NewBlueprintSpecUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase, doguInstallationUseCase)

	return &ApplicationContext{
		RemoteDoguRegistry:         remoteDoguRegistry,
		BlueprintSpecRepository:    blueprintSpecRepository,
		BlueprintSpecDomainUseCase: blueprintSpecDomainUseCase,
		DoguInstallationUseCase:    doguInstallationUseCase,
		BlueprintSpecUseCase:       blueprintUseCase,
		BlueprintSerializer:        blueprintSerializer,
		BlueprintMaskSerializer:    blueprintMaskSerializer,
	}, nil
}
