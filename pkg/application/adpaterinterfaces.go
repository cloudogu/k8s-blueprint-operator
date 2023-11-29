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

type BlueprintRepository interface {
	getById(doguName string) (domain.BlueprintV2, error)
	getAll() ([]domain.BlueprintV2, error)
	create(domain.BlueprintV2) error
	update(domain.BlueprintV2) error
	delete(domain.BlueprintV2) error
}
