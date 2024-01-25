package application

import (
	"context"
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type BlueprintSpecChangeUseCase struct {
	repo               blueprintSpecRepository
	validation         blueprintSpecValidationUseCase
	effectiveBlueprint effectiveBlueprintUseCase
	stateDiff          stateDiffUseCase
	doguInstallUseCase doguInstallationUseCase
	applyUseCase       applyBlueprintSpecUseCase
}

func NewBlueprintSpecChangeUseCase(
	repo domainservice.BlueprintSpecRepository,
	validation blueprintSpecValidationUseCase,
	effectiveBlueprint effectiveBlueprintUseCase,
	stateDiff stateDiffUseCase,
	doguInstallUseCase doguInstallationUseCase,
	applyUseCase applyBlueprintSpecUseCase,
) *BlueprintSpecChangeUseCase {
	return &BlueprintSpecChangeUseCase{
		repo:               repo,
		validation:         validation,
		effectiveBlueprint: effectiveBlueprint,
		stateDiff:          stateDiff,
		doguInstallUseCase: doguInstallUseCase,
		applyUseCase:       applyUseCase,
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

	logger.Info("getting changed blueprint") //log with id
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		errMsg := "cannot load blueprint spec"
		logger.Error(err, errMsg)
		return fmt.Errorf("%s: %w", errMsg, err)
	}

	logger = logger.WithValues("blueprintStatus", blueprintSpec.Status)
	logger.Info("handle blueprint") //log with id and status values.

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
		// activate maintenance mode
		// applyBlueprintSpec should happen in a new statusPhase then
		return useCase.applyBlueprintSpec(ctx, blueprintId)
	case domain.StatusPhaseEcosystemUnhealthyUpfront:
		return nil
	case domain.StatusPhaseInProgress:
		//should only happen if the system was interrupted, normally this state will be updated to completed or failed
		return useCase.handleInProgress(ctx, blueprintSpec)
	case domain.StatusPhaseBlueprintApplied:
		return useCase.applyUseCase.CheckEcosystemHealthAfterwards(ctx, blueprintId)
	case domain.StatusPhaseEcosystemHealthyAfterwards:
		//deactivate maintenance mode
		return nil
	case domain.StatusPhaseEcosystemUnhealthyAfterwards:
		//deactivate maintenance mode and set status to failed
		return nil
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

func (useCase *BlueprintSpecChangeUseCase) applyBlueprintSpec(ctx context.Context, blueprintId string) error {
	err := useCase.applyUseCase.ApplyBlueprintSpec(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}

func (useCase *BlueprintSpecChangeUseCase) handleInProgress(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).
		WithName("BlueprintSpecChangeUseCase.HandleChange").
		WithValues("blueprintId", blueprintSpec.Id)

	err := errors.New(handleInProgressMsg)
	logger.Error(err, "mark the blueprint as failed as the inProgress status should never be handled here")
	// do not return the inProgressError as this would lead to a reconcile, but this is not necessary if the status is failed afterward
	return useCase.applyUseCase.MarkFailed(ctx, blueprintSpec, err)
}

const handleInProgressMsg = "cannot handle blueprint in state " + string(domain.StatusPhaseInProgress) +
	" as this state shows that the appliance of the blueprint was interrupted before it could update the state " +
	"to either " + string(domain.StatusPhaseFailed) + " or " + string(domain.StatusPhaseCompleted)
