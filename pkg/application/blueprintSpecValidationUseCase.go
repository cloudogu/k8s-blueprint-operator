package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type BlueprintSpecValidationUseCase struct {
	repo                        blueprintSpecRepository
	validateDependenciesUseCase validateDependenciesDomainUseCase
	ValidateMountsUseCase       validateAdditionalMountsDomainUseCase
}

func NewBlueprintSpecValidationUseCase(
	repo domainservice.BlueprintSpecRepository,
	validateDependenciesUseCase validateDependenciesDomainUseCase,
	ValidateMountsUseCase validateAdditionalMountsDomainUseCase,
) *BlueprintSpecValidationUseCase {
	return &BlueprintSpecValidationUseCase{
		repo:                        repo,
		validateDependenciesUseCase: validateDependenciesUseCase,
		ValidateMountsUseCase:       ValidateMountsUseCase,
	}
}

// ValidateBlueprintSpecStatically checks the blueprintSpec for semantic errors and persists it.
// returns a domain.InvalidBlueprintError if blueprint is invalid or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *BlueprintSpecValidationUseCase) ValidateBlueprintSpecStatically(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecValidationUseCase.ValidateBlueprintSpecStatically")

	logger.Info("statically validate blueprint spec")

	invalidBlueprintError := blueprint.ValidateStatically()
	err := useCase.repo.Update(ctx, blueprint)
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
func (useCase *BlueprintSpecValidationUseCase) ValidateBlueprintSpecDynamically(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecValidationUseCase.ValidateBlueprintSpecDynamically")
	logger.Info("dynamically validate blueprint spec")

	validationError := errors.Join(
		useCase.validateDependenciesUseCase.ValidateDependenciesForAllDogus(ctx, blueprint.EffectiveBlueprint),
		useCase.ValidateMountsUseCase.ValidateAdditionalMounts(ctx, blueprint.EffectiveBlueprint),
	)

	if validationError != nil {
		validationError = &domain.InvalidBlueprintError{
			WrappedError: validationError,
			Message:      "blueprint spec is invalid",
		}
	}

	blueprint.ValidateDynamically(validationError)
	err := useCase.repo.Update(ctx, blueprint)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec after dynamic validation: %w", err)
	}
	return validationError
}
