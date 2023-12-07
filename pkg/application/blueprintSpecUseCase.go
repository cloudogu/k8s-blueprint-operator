package application

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type BlueprintSpecUseCase struct {
	repo          domainservice.BlueprintSpecRepository
	domainUseCase *domainservice.BlueprintSpecDomainUseCase
	doguUseCase   *DoguInstallationUseCase
}

func NewBlueprintSpecUseCase(
	repo domainservice.BlueprintSpecRepository,
	domainUseCase *domainservice.BlueprintSpecDomainUseCase,
	doguUseCase *DoguInstallationUseCase,
) *BlueprintSpecUseCase {
	return &BlueprintSpecUseCase{repo: repo, domainUseCase: domainUseCase, doguUseCase: doguUseCase}
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
		useCase.domainUseCase.ValidateDependenciesForAllDogus(blueprintSpec.EffectiveBlueprint),
	}
	validationError := errors.Join(errorList...)
	if validationError != nil {
		validationError = fmt.Errorf("blueprint spec is invalid: %w", validationError)
	}
	return validationError
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
