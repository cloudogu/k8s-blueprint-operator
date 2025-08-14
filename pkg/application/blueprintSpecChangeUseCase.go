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
	}
}

// HandleUntilApplied further executes a blueprint given by the blueprintId until it is as far applied as possible or an error occurred.
// If the process needs to wait for something, this function will return.
// Another call of this function is necessary to proceed.
// Returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write or
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
	err = useCase.applyUseCase.CheckEcosystemHealthUpfront(ctx, blueprint)
	if err != nil {
		return err
	}

	if !blueprint.ShouldBeApplied() {
		// just stop the loop here on dry run or early exit
		return nil
	}

	err = useCase.selfUpgradeUseCase.HandleSelfUpgrade(ctx, blueprint)
	if err != nil {
		return err
	}

	// without any error, the blueprint spec is always ready to be further evaluated, therefore call this function again to do that.
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
	case domain.StatusPhaseAwaitSelfUpgrade:
		return useCase.selfUpgradeUseCase.HandleSelfUpgrade(ctx, blueprint)
	case domain.StatusPhaseSelfUpgradeCompleted:
		return useCase.ecosystemConfigUseCase.ApplyConfig(ctx, blueprint)
	case domain.StatusPhaseEcosystemConfigApplied:
		return useCase.applyUseCase.ApplyBlueprintSpec(ctx, blueprint)
	case domain.StatusPhaseApplyEcosystemConfigFailed:
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprint)
	case domain.StatusPhaseInProgress:
		// should only happen if the system was interrupted, normally this state will be updated to blueprintApplied or BlueprintApplicationFailed
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprint)
	case domain.StatusPhaseBlueprintApplied:
		return useCase.doguRestartUseCase.TriggerDoguRestarts(ctx, blueprint)
	case domain.StatusPhaseRestartsTriggered:
		err := useCase.applyUseCase.CheckEcosystemHealthAfterwards(ctx, blueprint)
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
