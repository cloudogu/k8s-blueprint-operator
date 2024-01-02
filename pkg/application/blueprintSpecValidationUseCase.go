package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type BlueprintSpecValidationUseCase struct {
	repo                        domainservice.BlueprintSpecRepository
	validateDependenciesUseCase *domainservice.ValidateDependenciesDomainUseCase
}

func NewBlueprintSpecValidationUseCase(
	repo domainservice.BlueprintSpecRepository,
	validateDependenciesUseCase *domainservice.ValidateDependenciesDomainUseCase,
) *BlueprintSpecValidationUseCase {
	return &BlueprintSpecValidationUseCase{repo: repo, validateDependenciesUseCase: validateDependenciesUseCase}
}

// ValidateBlueprintSpecStatically checks the blueprintSpec for semantic errors and persists it.
// returns a domain.InvalidBlueprintError if blueprint is invalid or
// a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *BlueprintSpecValidationUseCase) ValidateBlueprintSpecStatically(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecValidationUseCase.ValidateBlueprintSpecStatically").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for static validation")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		var invalidError *domain.InvalidBlueprintError
		if errors.As(err, &invalidError) {
			blueprintSpec.MarkInvalid(err)
			updateErr := useCase.repo.Update(ctx, blueprintSpec)
			if updateErr != nil {
				return updateErr
			}
			return fmt.Errorf("blueprint spec syntax is invalid: %w", err)
		} else {
			return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
		}
	}

	logger.Info("statically validate blueprint spec", "blueprintStatus", blueprintSpec.Status)

	invalidBlueprintError := blueprintSpec.ValidateStatically()
	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		// InternalError or ConflictError, both should be handled by the caller
		return fmt.Errorf("cannot update blueprint spec after static validation: %w", err)
	}

	return invalidBlueprintError
}

// ValidateBlueprintSpecDynamically checks the blueprintSpec for semantic errors in combination with external data like dogu specs.
// returns a domain.InvalidBlueprintError if blueprint is invalid or
// a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *BlueprintSpecValidationUseCase) ValidateBlueprintSpecDynamically(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecValidationUseCase.ValidateBlueprintSpecDynamically").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for dynamic validation")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	logger.Info("dynamically validate blueprint spec", "blueprintStatus", blueprintSpec.Status)

	validationError := useCase.validateDependenciesUseCase.ValidateDependenciesForAllDogus(ctx, blueprintSpec.EffectiveBlueprint)
	if validationError != nil {
		validationError = &domain.InvalidBlueprintError{
			WrappedError: validationError,
			Message:      "blueprint spec is invalid",
		}
	}

	blueprintSpec.ValidateDynamically(validationError)
	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec after dynamic validation: %w", err)
	}
	return validationError
}
