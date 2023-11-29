package domainservice

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type DoguInstallationDomainUseCase struct{}

func (useCase DoguInstallationDomainUseCase) ValidateDoguHealth(installedDogus []ecosystem.DoguInstallation) error {
	var healthErrors []error
	for _, dogu := range installedDogus {
		if dogu.Health != ecosystem.Healhty {
			healthErrors = append(healthErrors, fmt.Errorf("dogu %s is unhealthy", dogu.Name))
		}
	}

	err := errors.Join(healthErrors...)
	if err != nil {
		err = fmt.Errorf("dogus are unhealthy: %w", err)
	}
	return err
}
