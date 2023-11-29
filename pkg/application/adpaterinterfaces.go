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
	getById(doguName string) (domain.Blueprint, error)
	getAll() ([]domain.Blueprint, error)
	create(domain.Blueprint) error
	update(domain.Blueprint) error
	delete(domain.Blueprint) error
}
