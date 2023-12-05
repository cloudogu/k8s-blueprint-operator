package application

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type BlueprintSpecUseCase struct {
	repo          domainservice.BlueprintSpecRepository
	domainUseCase domainservice.BlueprintSpecDomainUseCase
	doguUseCase   DoguInstallationUseCase
}

func (useCase *BlueprintSpecUseCase) ValidateBlueprintSpecStatically(blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	validationError := blueprintSpec.Validate()
	err = useCase.repo.Update(blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec after validation: %w", err)
	}

	return validationError
}

func (useCase *BlueprintSpecUseCase) ValidateBlueprintSpecDynamically(blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	errorList := []error{
		useCase.domainUseCase.ValidateDoguDependencies(blueprintSpec),
	}
	validationError := errors.Join(errorList...)
	if validationError != nil {
		validationError = fmt.Errorf("")
	}
	//TODO
	return nil
}

func (useCase *BlueprintSpecUseCase) calculateEffectiveBlueprint(blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to calculate effective blueprint: %w", err)
	}

	calcError := blueprintSpec.CalculateEffectiveBlueprint()
	err = useCase.repo.Update(blueprintSpec)
	if err != nil {
		return err
	}

	return calcError
}
