package pkg

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type ApplicationContext struct {
	remoteDoguRegistry         domainservice.RemoteDoguRegistry
	blueprintSpecRepository    domainservice.BlueprintSpecRepository
	blueprintSpecDomainUseCase *domainservice.BlueprintSpecDomainUseCase
	DoguInstallationUseCase    *application.DoguInstallationUseCase
	blueprintSpecUseCase       *application.BlueprintSpecUseCase
}

var ApplicationContextContainer *ApplicationContext

func Bootstrap() {
	var remoteDoguRegistry domainservice.RemoteDoguRegistry
	var blueprintSpecRepository domainservice.BlueprintSpecRepository
	blueprintSpecDomainUseCase := domainservice.NewBlueprintSpecDomainUseCase(remoteDoguRegistry)
	doguInstallationUseCase := &application.DoguInstallationUseCase{}

	ApplicationContextContainer = &ApplicationContext{
		remoteDoguRegistry:         remoteDoguRegistry,
		blueprintSpecRepository:    blueprintSpecRepository,
		blueprintSpecDomainUseCase: blueprintSpecDomainUseCase,
		DoguInstallationUseCase:    doguInstallationUseCase,
		blueprintSpecUseCase:       application.NewBlueprintSpecUseCase(blueprintSpecRepository, blueprintSpecDomainUseCase, doguInstallationUseCase),
	}
}
