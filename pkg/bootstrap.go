package pkg

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/blueprint"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ApplicationContext struct {
	remoteDoguRegistry         domainservice.RemoteDoguRegistry
	blueprintSpecRepository    domainservice.BlueprintSpecRepository
	blueprintSpecDomainUseCase *domainservice.BlueprintSpecDomainUseCase
	DoguInstallationUseCase    *application.DoguInstallationUseCase
	blueprintSpecUseCase       *application.BlueprintSpecUseCase
}

var ApplicationContextContainer *ApplicationContext

func Bootstrap(restConfig *rest.Config, namespace string) error {
	k8sClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("unable to create k8s clientset: %w", err)
	}

	//TODO where will be the eventRecorder interface be located? here? Events must probably be written overall the operator.

	ecosystemClient, err := blueprint.NewClientSet(restConfig, k8sClientSet)
	if err != nil {
		return fmt.Errorf("unable to create ecosystem client: %w", err)
	}

	var remoteDoguRegistry domainservice.RemoteDoguRegistry
	blueprintSpecRepository := blueprint.NewRepository(namespace, ecosystemClient)
	blueprintSpecDomainUseCase := domainservice.NewBlueprintSpecDomainUseCase(remoteDoguRegistry)
	doguInstallationUseCase := &application.DoguInstallationUseCase{}

	ApplicationContextContainer = &ApplicationContext{
		remoteDoguRegistry:         remoteDoguRegistry,
		blueprintSpecRepository:    blueprintSpecRepository,
		blueprintSpecDomainUseCase: blueprintSpecDomainUseCase,
		DoguInstallationUseCase:    doguInstallationUseCase,
		blueprintSpecUseCase:       application.NewBlueprintSpecUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase, doguInstallationUseCase),
	}
	return nil
}
