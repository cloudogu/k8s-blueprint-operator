package application

import (
	"fmt"
)

type BlueprintSpecUseCase struct {
	repo        BlueprintSpecRepository
	doguUseCase DoguInstallationUseCase
}

func (useCase BlueprintSpecUseCase) ValidateBlueprintSpecStatically(blueprintId string) error {
	blueprintSpec, err := useCase.repo.getById(blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	validationError := blueprintSpec.Validate()
	err = useCase.repo.update(blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec after validation: %w", err)
	}

	return validationError
}

func (useCase BlueprintSpecUseCase) calculateEffectiveBlueprint(blueprintId string) error {
	blueprintSpec, err := useCase.repo.getById(blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to calculate effective blueprint it: %w", err)
	}

	calcError := blueprintSpec.CalculateEffectiveBlueprint()
	err = useCase.repo.update(blueprintSpec)
	if err != nil {
		return err
	}

	return calcError
}
