package application

import (
	"context"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

type BlueprintPreparationUseCase struct {
	initialStatus      initialBlueprintStatusUseCase
	validation         blueprintSpecValidationUseCase
	effectiveBlueprint effectiveBlueprintUseCase
	stateDiff          stateDiffUseCase
	healthUseCase      ecosystemHealthUseCase
}

func NewBlueprintPreparationUseCase(
	initialStatus initialBlueprintStatusUseCase,
	validation blueprintSpecValidationUseCase,
	effectiveBlueprint effectiveBlueprintUseCase,
	stateDiff stateDiffUseCase,
	ecosystemHealthUseCase ecosystemHealthUseCase,
) BlueprintPreparationUseCase {
	return BlueprintPreparationUseCase{
		initialStatus:      initialStatus,
		validation:         validation,
		effectiveBlueprint: effectiveBlueprint,
		stateDiff:          stateDiff,
		healthUseCase:      ecosystemHealthUseCase,
	}
}

func (useCase *BlueprintPreparationUseCase) prepareBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.initialStatus.InitateConditions(ctx, blueprint)
	if err != nil {
		return err
	}
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
	// always check health here, even if we already know here, that we don't need to apply anything
	// because we need to update the health condition
	_, err = useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
	if err != nil {
		return err
	}
	err = useCase.stateDiff.DetermineStateDiff(ctx, blueprint)
	if err != nil {
		// error could be either a technical error from a repository or an InvalidBlueprintError from the domain
		// both cases can be handled the same way as the calling method (reconciler) can handle the error type itself.
		return err
	}
	return nil
}
