package application

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type DoguInstallationUseCase struct {
	doguRepo          domainservice.DoguInstallationRepository
	doguDomainUseCase domainservice.DoguInstallationDomainUseCase
}

func (useCase *DoguInstallationUseCase) validateDoguHealth() error {
	installedDogus, err := useCase.doguRepo.GetAll()
	if err != nil {
		return fmt.Errorf("cannot evaluate dogu health states: %w", err)
	}

	return useCase.doguDomainUseCase.ValidateDoguHealth(installedDogus)
}

func (useCase *DoguInstallationUseCase) installDogu(doguName string) error {
	//TODO
	return nil
}

func (useCase *DoguInstallationUseCase) uninstallDogu(doguName string) error {
	//TODO
	return nil
}
