package domainservice

import (
	"context"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type DoguInstallationRepository interface {
	GetByName(doguName string) (ecosystem.DoguInstallation, error)
	GetAll() ([]ecosystem.DoguInstallation, error)
	Create(ecosystem.DoguInstallation) error
	Update(ecosystem.DoguInstallation) error
	Delete(ecosystem.DoguInstallation) error
}

type BlueprintSpecRepository interface {
	// GetById returns a Blueprint identified by its ID.
	GetById(ctx context.Context, blueprintId string) (domain.BlueprintSpec, error)
	// Update updates a given BlueprintSpec.
	Update(ctx context.Context, blueprintSpec domain.BlueprintSpec) error
}

type RemoteDoguRegistry interface {
	GetDogu(qualifiedDoguName string, version string) (core.Dogu, error)
}
