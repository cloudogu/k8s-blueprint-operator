package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type BlueprintSpecChangeUseCase struct {
	repo               blueprintSpecRepository
	preparationUseCase BlueprintPreparationUseCase
	applyUseCase       BlueprintApplyUseCase
}

func NewBlueprintSpecChangeUseCase(
	repo blueprintSpecRepository,
	preparationUseCase BlueprintPreparationUseCase,
	applyUseCase BlueprintApplyUseCase,
) *BlueprintSpecChangeUseCase {
	return &BlueprintSpecChangeUseCase{
		repo:               repo,
		preparationUseCase: preparationUseCase,
		applyUseCase:       applyUseCase,
	}
}

// HandleUntilApplied further executes a blueprint given by the blueprintId until it is as far applied as possible or an error occurred.
// If the process needs to wait for something, this function will return.
// Another call of this function is necessary to proceed.
// Returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write or
// a domain.AwaitSelfUpgradeError if we need to wait for a self-upgrade or
// a domain.InvalidBlueprintError if the blueprint is invalid.
func (useCase *BlueprintSpecChangeUseCase) HandleUntilApplied(givenCtx context.Context, blueprintId string) error {
	logger := log.FromContext(givenCtx).
		WithValues("blueprintId", blueprintId)
	// set the logger in the context to make use of structured logging
	// we will give this ctx in every use case, therefore all of them will include the values given here
	ctx := log.IntoContext(givenCtx, logger)
	logger = logger.WithName("BlueprintSpecChangeUseCase.HandleUntilApplied")

	logger.V(2).Info("getting changed blueprint") // log with id
	blueprint, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		errMsg := "cannot load blueprint spec"
		logger.Error(err, errMsg)
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	logger.V(1).Info("handle blueprint")

	err = useCase.preparationUseCase.prepareBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}

	if !blueprint.ShouldBeApplied() {
		// stop the loop here on stopped-flag or early exit
		return useCase.handleShouldNotBeApplied(ctx, logger, blueprint)
	}

	// === Apply from here on ===
	err = useCase.applyUseCase.applyBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}

	logger.Info("blueprint successfully applied")
	return nil
}

func (useCase *BlueprintSpecChangeUseCase) handleShouldNotBeApplied(ctx context.Context, logger logr.Logger, blueprint *domain.BlueprintSpec) error {
	// post event and log only if blueprint is stopped, all other cases are just NoOps
	if blueprint.Config.Stopped {
		logger.Info("blueprint is currently set as stopped and will not be applied")
		blueprint.Events = append(blueprint.Events, domain.BlueprintStoppedEvent{})
		err := useCase.repo.Update(ctx, blueprint)
		if err != nil {
			return fmt.Errorf("cannot update status to set stopped event: %w", err)
		}
	} else {
		logger.V(1).Info("no diff detected, no changes required")
	}
	return nil
}

func (useCase *BlueprintSpecChangeUseCase) CheckForMultipleBlueprintResources(ctx context.Context) error {
	logger := log.FromContext(ctx).WithName("BlueprintSpecChangeUseCase.CheckForMultipleBlueprintResources")

	logger.V(2).Info("check for multiple blueprints")
	err := useCase.repo.CheckSingleton(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", "check for multiple blueprints not successful", err)
	}

	return nil
}
