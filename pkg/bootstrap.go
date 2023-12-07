package pkg

import (
	"fmt"
	kubernetes2 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ApplicationContext contains vital application parts for this operator.
type ApplicationContext struct {
	remoteDoguRegistry         domainservice.RemoteDoguRegistry
	blueprintSpecRepository    domainservice.BlueprintSpecRepository
	blueprintSpecDomainUseCase *domainservice.ValidateDependenciesDomainUseCase
	DoguInstallationUseCase    *application.DoguInstallationUseCase
	blueprintSpecUseCase       *application.BlueprintSpecUseCase
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

	var remoteDoguRegistry domainservice.RemoteDoguRegistry
	blueprintSpecRepository := kubernetes2.NewBlueprintSpecRepository(namespace, ecosystemClient)
	blueprintSpecDomainUseCase := domainservice.NewValidateDependenciesDomainUseCase(remoteDoguRegistry)
	doguInstallationUseCase := &application.DoguInstallationUseCase{}
	blueprintUseCase := application.NewBlueprintSpecUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase, doguInstallationUseCase)

	return &ApplicationContext{
		remoteDoguRegistry:         remoteDoguRegistry,
		blueprintSpecRepository:    blueprintSpecRepository,
		blueprintSpecDomainUseCase: blueprintSpecDomainUseCase,
		DoguInstallationUseCase:    doguInstallationUseCase,
		blueprintSpecUseCase:       blueprintUseCase,
	}, nil
}
