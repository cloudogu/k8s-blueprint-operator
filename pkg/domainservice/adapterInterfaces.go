package domainservice

import (
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
	GetById(doguName string) (domain.BlueprintSpec, error)
	GetAll() ([]domain.BlueprintSpec, error)
	Create(domain.BlueprintSpec) error
	Update(domain.BlueprintSpec) error
	Delete(domain.BlueprintSpec) error
}

type RemoteDoguRegistry interface {
	GetDogu(qualifiedDoguName string, version string) (core.Dogu, error)
}
