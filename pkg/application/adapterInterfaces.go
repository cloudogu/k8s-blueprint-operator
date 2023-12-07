package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

// interface duplication for mocks
type doguInstallationRepository interface {
	domainservice.DoguInstallationRepository
}

type blueprintSpecRepository interface {
	domainservice.BlueprintSpecRepository
}

type remoteDoguRegistry interface {
	domainservice.RemoteDoguRegistry
}
