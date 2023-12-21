package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type DoguInstallationUseCase struct {
	doguRepo          domainservice.DoguInstallationRepository
	doguDomainUseCase domainservice.DoguInstallationDomainUseCase
}

func (useCase *DoguInstallationUseCase) validateDoguHealth(ctx context.Context) error {
	//TODO: this is only a stub to get an idea of the upcoming implementation
	installedDogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot evaluate dogu health states: %w", err)
	}

	return useCase.doguDomainUseCase.ValidateDoguHealth(installedDogus)
}

func (useCase *DoguInstallationUseCase) installDogu(ctx context.Context, doguName string, version string) error {
	//TODO
	return nil
}

func (useCase *DoguInstallationUseCase) uninstallDogu(ctx context.Context, doguName string) error {
	//TODO
	return nil
}
