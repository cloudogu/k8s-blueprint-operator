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

// HandleChange further executes a blueprint spec given by the blueprintId until it is fully applied or an error occurred.
// Returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write or
// a domain.InvalidBlueprintError if the blueprint is invalid.
func (useCase *BlueprintSpecChangeUseCase) HandleChange(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecChangeUseCase.HandleChange").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting changed blueprint") // log with id
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		errMsg := "cannot load blueprint spec"
		logger.Error(err, errMsg)
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	logger = logger.WithValues("blueprintStatus", blueprintSpec.Status)
	logger.Info("handle blueprint") // log with id and status values.

	// without any error, the blueprint spec is always ready to be further evaluated, therefore call this function again to do that.
	switch blueprintSpec.Status {
	case domain.StatusPhaseNew:
		return useCase.validateStatically(ctx, blueprintId)
	case domain.StatusPhaseInvalid:
		return nil
	case domain.StatusPhaseStaticallyValidated:
		return useCase.calculateEffectiveBlueprint(ctx, blueprintId)
	case domain.StatusPhaseEffectiveBlueprintGenerated:
		return useCase.validateDynamically(ctx, blueprintId)
	case domain.StatusPhaseValidated:
		return useCase.determineStateDiff(ctx, blueprintId)
	case domain.StatusPhaseStateDiffDetermined:
		return useCase.checkEcosystemHealthUpfront(ctx, blueprintId)
	case domain.StatusPhaseEcosystemHealthyUpfront:
		return useCase.preProcessBlueprintApplication(ctx, blueprintSpec)
	case domain.StatusPhaseEcosystemUnhealthyUpfront:
		return nil
	case domain.StatusPhaseBlueprintApplicationPreProcessed:
		return useCase.handleSelfUpgrade(ctx, blueprintId)
	case domain.StatusPhaseAwaitSelfUpgrade:
		return useCase.handleSelfUpgrade(ctx, blueprintId)
	case domain.StatusPhaseSelfUpgradeCompleted:
		return useCase.applyEcosystemConfig(ctx, blueprintId)
	case domain.StatusPhaseEcosystemConfigApplied:
		return useCase.applyBlueprintSpec(ctx, blueprintId)
	case domain.StatusPhaseApplyEcosystemConfigFailed:
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprintId)
	case domain.StatusPhaseInProgress:
		// should only happen if the system was interrupted, normally this state will be updated to blueprintApplied or BlueprintApplicationFailed
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprintId)
	case domain.StatusPhaseBlueprintApplied:
		return useCase.triggerDoguRestarts(ctx, blueprintId)
	case domain.StatusPhaseRestartsTriggered:
		return useCase.checkEcosystemHealthAfterwards(ctx, blueprintId)
	case domain.StatusPhaseBlueprintApplicationFailed:
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprintId)
	case domain.StatusPhaseEcosystemHealthyAfterwards:
		// deactivate maintenance mode and set status to completed
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprintId)
	case domain.StatusPhaseEcosystemUnhealthyAfterwards:
		// deactivate maintenance mode and set status to failed
		return useCase.applyUseCase.PostProcessBlueprintApplication(ctx, blueprintId)
	case domain.StatusPhaseCompleted:
		return nil
	case domain.StatusPhaseFailed:
		return nil
	default:
		return fmt.Errorf("could not handle unknown status of blueprint")
	}
}

func (useCase *BlueprintSpecChangeUseCase) validateStatically(ctx context.Context, blueprintId string) error {
	err := useCase.validation.ValidateBlueprintSpecStatically(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) calculateEffectiveBlueprint(ctx context.Context, blueprintId string) error {
	err := useCase.effectiveBlueprint.CalculateEffectiveBlueprint(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) validateDynamically(ctx context.Context, blueprintId string) error {
	err := useCase.validation.ValidateBlueprintSpecDynamically(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) determineStateDiff(ctx context.Context, blueprintId string) error {
	err := useCase.stateDiff.DetermineStateDiff(ctx, blueprintId)

	// error could be either a technical error from a repository or an InvalidBlueprintError from the domain
	// both cases can be handled the same way as the calling method (reconciler) can handle the error type itself.
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) checkEcosystemHealthUpfront(ctx context.Context, blueprintId string) error {
	err := useCase.applyUseCase.CheckEcosystemHealthUpfront(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) preProcessBlueprintApplication(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	err := useCase.applyUseCase.PreProcessBlueprintApplication(ctx, blueprintSpec.Id)
	if err != nil {
		return err
	}
	if !blueprintSpec.ShouldBeApplied() {
		// event recording and so on happen in PreProcessBlueprintApplication
		// just stop the loop here on dry run or early exit
		return nil
	}
	return useCase.HandleChange(ctx, blueprintSpec.Id)
}

func (useCase *BlueprintSpecChangeUseCase) handleSelfUpgrade(ctx context.Context, blueprintId string) error {
	err := useCase.selfUpgradeUseCase.HandleSelfUpgrade(ctx, blueprintId)
	if err != nil {
		return err
	}
	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) applyBlueprintSpec(ctx context.Context, blueprintId string) error {
	err := useCase.applyUseCase.ApplyBlueprintSpec(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) checkEcosystemHealthAfterwards(ctx context.Context, blueprintId string) error {
	err := useCase.applyUseCase.CheckEcosystemHealthAfterwards(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) applyEcosystemConfig(ctx context.Context, blueprintId string) error {
	err := useCase.ecosystemConfigUseCase.ApplyConfig(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) triggerDoguRestarts(ctx context.Context, blueprintId string) error {
	err := useCase.doguRestartUseCase.TriggerDoguRestarts(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}
