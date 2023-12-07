package kubernetes

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

type DoguCRRepository struct{}

func (repo DoguCRRepository) getByName(doguName string) (ecosystem.DoguInstallation, error) {
	//TODO
	return ecosystem.DoguInstallation{}, nil
}

func (repo DoguCRRepository) getAll() ([]ecosystem.DoguInstallation, error) {
	//TODO
	return []ecosystem.DoguInstallation{}, nil
}

func (repo DoguCRRepository) save(ecosystem.DoguInstallation) error {
	//TODO
	return nil
}
