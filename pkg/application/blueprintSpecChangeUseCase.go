package application

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type BlueprintSpecChangeUseCase struct {
	repo               domainservice.BlueprintSpecRepository
	validation         blueprintSpecValidationUseCase
	effectiveBlueprint effectiveBlueprintUseCase
	stateDiff          stateDiffUseCase
	doguInstall        doguInstallationUseCase
}

func NewBlueprintSpecChangeUseCase(
	repo domainservice.BlueprintSpecRepository,
	validation blueprintSpecValidationUseCase,
	effectiveBlueprint effectiveBlueprintUseCase,
	stateDiff stateDiffUseCase,
	doguInstall doguInstallationUseCase,
) *BlueprintSpecChangeUseCase {
	return &BlueprintSpecChangeUseCase{
		repo:               repo,
		validation:         validation,
		effectiveBlueprint: effectiveBlueprint,
		stateDiff:          stateDiff,
		doguInstall:        doguInstall,
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
		return useCase.checkDoguHealth(ctx, blueprintId)
	case domain.StatusPhaseIgnoreDoguHealth:
		fallthrough
	case domain.StatusPhaseDogusHealthy:
		return nil
	case domain.StatusPhaseDogusUnhealthy:
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

func (useCase *BlueprintSpecChangeUseCase) checkDoguHealth(ctx context.Context, blueprintId string) error {
	err := useCase.doguInstall.CheckDoguHealth(ctx, blueprintId)
	if err != nil {
		return err
	}

	return useCase.HandleChange(ctx, blueprintId)
}
