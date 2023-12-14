package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecUseCase.HandleBlueprintSpecChange").
		WithValues("blueprintId", blueprintId, "blueprintStatus", blueprintSpec.Status)

	logger.Info("handle blueprint") //log with id and status values.
	if err != nil {
		logger.Error(err, "cannot load blueprint spec")
		return fmt.Errorf("cannot load blueprint spec: %w", err)
	}
	// without any error, the blueprint spec is always ready to be further evaluated, therefore call this function again to do that.
	switch blueprintSpec.Status {
	case domain.StatusPhaseNew:
		err := useCase.ValidateBlueprintSpecStatically(ctx, blueprintId)
		if err != nil {
			return err
		}
		return useCase.HandleBlueprintSpecChange(ctx, blueprintId)
	case domain.StatusPhaseInvalid:
		return nil
	case domain.StatusPhaseStaticallyValidated:
		err := useCase.ValidateBlueprintSpecDynamically(ctx, blueprintId)
		if err != nil {
			return err
		}
		return useCase.HandleBlueprintSpecChange(ctx, blueprintId)
	case domain.StatusPhaseValidated:
		err := useCase.calculateEffectiveBlueprint(ctx, blueprintId)
		if err != nil {
			return err
		}
		return nil
	case domain.StatusPhaseInProgress:
		return nil
	case domain.StatusPhaseCompleted:
		return nil
	case domain.StatusPhaseFailed:
		return nil
	default:
		return fmt.Errorf("could not handle unknown status of blueprint")
	}
}

// ValidateBlueprintSpecStatically checks the blueprintSpec for semantic errors and persists it.
// returns a domain.InvalidBlueprintError if blueprint is invalid or
// a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *BlueprintSpecUseCase) ValidateBlueprintSpecStatically(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("BlueprintSpecUseCase.ValidateBlueprintSpecStatically")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	logger.Info("statically validate blueprint spec", "blueprintId", blueprintId, "blueprintStatus", blueprintSpec.Status)
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
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec.
func (useCase *BlueprintSpecUseCase) ValidateBlueprintSpecDynamically(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("BlueprintSpecUseCase.ValidateBlueprintSpecStatically")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	logger.Info("dynamically validate blueprint spec", "blueprintId", blueprintId, "blueprintStatus", blueprintSpec.Status)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to validate it: %w", err)
	}

	errorList := []error{
		useCase.validateDependenciesUseCase.ValidateDependenciesForAllDogus(blueprintSpec.EffectiveBlueprint),
	}
	validationError := errors.Join(errorList...)
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

// calculateEffectiveBlueprint load the blueprintSpec, lets it calculate the effective blueprint and persists it again.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec.
func (useCase *BlueprintSpecUseCase) calculateEffectiveBlueprint(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("BlueprintSpecUseCase.ValidateBlueprintSpecStatically")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	logger.Info("calculate effective blueprint", "blueprintId", blueprintId, "blueprintStatus", blueprintSpec.Status)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to calculate effective blueprint: %w", err)
	}

	calcError := blueprintSpec.CalculateEffectiveBlueprint()
	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec after calculating the effective blueprint: %w", err)
	}

	return calcError
}
