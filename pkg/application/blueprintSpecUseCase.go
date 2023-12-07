package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type BlueprintSpecUseCase struct {
	repo                        domainservice.BlueprintSpecRepository
	validateDependenciesUseCase *domainservice.ValidateDependenciesDomainUseCase
	doguUseCase                 *DoguInstallationUseCase
}

func NewBlueprintSpecUseCase(
	repo domainservice.BlueprintSpecRepository,
	validateDependenciesUseCase *domainservice.ValidateDependenciesDomainUseCase,
	doguUseCase *DoguInstallationUseCase,
) *BlueprintSpecUseCase {
	return &BlueprintSpecUseCase{repo: repo, validateDependenciesUseCase: validateDependenciesUseCase, doguUseCase: doguUseCase}
}

func (useCase *BlueprintSpecUseCase) HandleBlueprintSpecChange(ctx context.Context, blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec: %w", err)
	}
	switch blueprintSpec.Status {
	case domain.StatusPhaseNew:
		err := useCase.ValidateBlueprintSpecStatically(ctx, blueprintId)
		if err != nil {
			return err
		}
	case domain.StatusPhaseInvalid:
		return nil
	case domain.StatusPhaseValidated:
		err := useCase.calculateEffectiveBlueprint(ctx, blueprintId)
		if err != nil {
			return err
		}
	case domain.StatusPhaseInProgress:
		return nil
	case domain.StatusPhaseCompleted:
		return nil
	default:
		return fmt.Errorf("could not handle change in blueprint spec with ID '%s': unknown status '%s'", blueprintId, blueprintSpec.Status)

	}

	return nil
}

func (useCase *BlueprintSpecUseCase) ValidateBlueprintSpecStatically(ctx context.Context, blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	validationError := blueprintSpec.Validate()
	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec after validation: %w", err)
	}

	return validationError
}

func (useCase *BlueprintSpecUseCase) ValidateBlueprintSpecDynamically(ctx context.Context, blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	errorList := []error{
		useCase.validateDependenciesUseCase.ValidateDependenciesForAllDogus(blueprintSpec.EffectiveBlueprint),
	}
	validationError := errors.Join(errorList...)
	if validationError != nil {
		validationError = fmt.Errorf("blueprint spec is invalid: %w", validationError)
	}
	//TODO: Maybe we should persist the status change, but then we have two validated status (statically and dynamically)
	return validationError
}

func (useCase *BlueprintSpecUseCase) calculateEffectiveBlueprint(ctx context.Context, blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to calculate effective blueprint: %w", err)
	}

	calcError := blueprintSpec.CalculateEffectiveBlueprint()
	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return err
	}

	return calcError
}
