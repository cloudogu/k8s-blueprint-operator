package application

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type BlueprintSpecChangeUseCase struct {
	repo                   blueprintSpecRepository
	validation             blueprintSpecValidationUseCase
	effectiveBlueprint     effectiveBlueprintUseCase
	stateDiff              stateDiffUseCase
	applyUseCase           applyBlueprintSpecUseCase
	ecosystemConfigUseCase ecosystemConfigUseCase
	doguRestartUseCase     doguRestartUseCase
	selfUpgradeUseCase     selfUpgradeUseCase
	applyComponentUseCase  applyComponentUseCase
	healthUseCase          ecosystemHealthUseCase
}

func NewBlueprintSpecChangeUseCase(
	repo domainservice.BlueprintSpecRepository,
	validation blueprintSpecValidationUseCase,
	effectiveBlueprint effectiveBlueprintUseCase,
	stateDiff stateDiffUseCase,
	applyUseCase applyBlueprintSpecUseCase,
	ecosystemConfigUseCase ecosystemConfigUseCase,
	doguRestartUseCase doguRestartUseCase,
	selfUpgradeUseCase selfUpgradeUseCase,
	applyComponentUseCase applyComponentUseCase,
	ecosystemHealthUseCase ecosystemHealthUseCase,
) *BlueprintSpecChangeUseCase {
	return &BlueprintSpecChangeUseCase{
		repo:                   repo,
		validation:             validation,
		effectiveBlueprint:     effectiveBlueprint,
		stateDiff:              stateDiff,
		applyUseCase:           applyUseCase,
		ecosystemConfigUseCase: ecosystemConfigUseCase,
		doguRestartUseCase:     doguRestartUseCase,
		selfUpgradeUseCase:     selfUpgradeUseCase,
		applyComponentUseCase:  applyComponentUseCase,
		healthUseCase:          ecosystemHealthUseCase,
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
		WithName("BlueprintSpecChangeUseCase.HandleUntilApplied").
		WithValues("blueprintId", blueprintId)
	// set the logger in the context to make use of structured logging
	// we will give this ctx in every use case, therefore all of them will include the values given here
	ctx := log.IntoContext(givenCtx, logger)

	logger.Info("getting changed blueprint") // log with id
	blueprint, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		errMsg := "cannot load blueprint spec"
		logger.Error(err, errMsg)
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	logger.Info("handle blueprint")

	err = useCase.validation.ValidateBlueprintSpecStatically(ctx, blueprint)
	if err != nil {
		return err
	}
	err = useCase.effectiveBlueprint.CalculateEffectiveBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}
	err = useCase.validation.ValidateBlueprintSpecDynamically(ctx, blueprint)
	if err != nil {
		return err
	}
	err = useCase.stateDiff.DetermineStateDiff(ctx, blueprint)
	if err != nil {
		// error could be either a technical error from a repository or an InvalidBlueprintError from the domain
		// both cases can be handled the same way as the calling method (reconciler) can handle the error type itself.
		return err
	}
	// always check health here, even if we already know here, that we don't need to apply anything
	// because we need to update the health condition
	_, err = useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
	if err != nil {
		return err
	}

	if !blueprint.ShouldBeApplied() {
		// just stop the loop here on dry run or early exit
		return nil
	}

	// === Apply from here on ===
	err = useCase.selfUpgradeUseCase.HandleSelfUpgrade(ctx, blueprint)
	if err != nil {
		// could be a domain.AwaitSelfUpgradeError to trigger another reconcile
		return err
	}
	err = useCase.ecosystemConfigUseCase.ApplyConfig(ctx, blueprint)
	if err != nil {
		return err
	}
	err = useCase.applyComponentUseCase.ApplyComponents(ctx, blueprint)
	if err != nil {
		return err
	}
	// check after applying components
	_, err = useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
	if err != nil {
		return err
	}
	err = useCase.applyUseCase.ApplyBlueprintSpec(ctx, blueprint)
	if err != nil {
		return err
	}
	// check after installing or updating dogus
	_, err = useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
	if err != nil {
		return err
	}

	//TODO: remove this loop, when all use cases are reworked without the use of status
	for blueprint.Status != domain.StatusPhaseCompleted {
		err := useCase.handleChange(ctx, blueprint)
		if err != nil {
			return err
		}
	}

	logger.Info("blueprint successfully applied")
	return nil
}

func (useCase *BlueprintSpecChangeUseCase) handleChange(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	switch blueprint.Status {
	case domain.StatusPhaseBlueprintApplied:
		return useCase.doguRestartUseCase.TriggerDoguRestarts(ctx, blueprint)
	case domain.StatusPhaseRestartsTriggered:
		_, err := useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
		if err != nil {
			return err
		}
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprint)
	case domain.StatusPhaseBlueprintApplicationFailed:
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprint)
	case domain.StatusPhaseCompleted:
		return nil
	case domain.StatusPhaseFailed:
		return nil
	default:
		return fmt.Errorf("could not handle unknown status of blueprint")
	}
}
