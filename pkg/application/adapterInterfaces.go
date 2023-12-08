package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

// interface duplication for mocks

//nolint:all
type doguInstallationRepository interface {
	domainservice.DoguInstallationRepository
}

//nolint:all
type blueprintSpecRepository interface {
	domainservice.BlueprintSpecRepository
}

//nolint:all
type remoteDoguRegistry interface {
	domainservice.RemoteDoguRegistry
}
