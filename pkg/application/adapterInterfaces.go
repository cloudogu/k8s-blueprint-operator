package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type DoguInstallationRepository interface {
	getByName(doguName string) (ecosystem.DoguInstallation, error)
	getAll() ([]ecosystem.DoguInstallation, error)
	create(ecosystem.DoguInstallation) error
	update(ecosystem.DoguInstallation) error
	delete(ecosystem.DoguInstallation) error
}

type BlueprintSpecRepository interface {
	getById(doguName string) (domain.BlueprintSpec, error)
	getAll() ([]domain.BlueprintSpec, error)
	create(domain.BlueprintSpec) error
	update(domain.BlueprintSpec) error
	delete(domain.BlueprintSpec) error
}
