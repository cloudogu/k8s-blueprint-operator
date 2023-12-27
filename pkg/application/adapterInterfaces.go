package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

// interface duplication for mocks

//nolint:unused
type doguInstallationRepository interface {
	domainservice.DoguInstallationRepository
}

//nolint:unused
type blueprintSpecRepository interface {
	domainservice.BlueprintSpecRepository
}

//nolint:unused
type remoteDoguRegistry interface {
	domainservice.RemoteDoguRegistry
}
