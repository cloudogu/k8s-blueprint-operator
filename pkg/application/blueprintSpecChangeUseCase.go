package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type BlueprintSpecChangeUseCase struct {
	repo                   blueprintSpecRepository
	validation             blueprintSpecValidationUseCase
	effectiveBlueprint     effectiveBlueprintUseCase
	stateDiff              stateDiffUseCase
	applyUseCase           completeBlueprintUseCase
	ecosystemConfigUseCase ecosystemConfigUseCase
	selfUpgradeUseCase     selfUpgradeUseCase
	applyComponentUseCase  applyComponentsUseCase
	applyDogusUseCase      applyDogusUseCase
	healthUseCase          ecosystemHealthUseCase
}

func NewBlueprintSpecChangeUseCase(
	repo blueprintSpecRepository,
	validation blueprintSpecValidationUseCase,
	effectiveBlueprint effectiveBlueprintUseCase,
	stateDiff stateDiffUseCase,
	applyUseCase completeBlueprintUseCase,
	ecosystemConfigUseCase ecosystemConfigUseCase,
	selfUpgradeUseCase selfUpgradeUseCase,
	applyComponentUseCase applyComponentsUseCase,
	applyDogusUseCase applyDogusUseCase,
	ecosystemHealthUseCase ecosystemHealthUseCase,
) *BlueprintSpecChangeUseCase {
	return &BlueprintSpecChangeUseCase{
		repo:                   repo,
		validation:             validation,
		effectiveBlueprint:     effectiveBlueprint,
		stateDiff:              stateDiff,
		applyUseCase:           applyUseCase,
		ecosystemConfigUseCase: ecosystemConfigUseCase,
		selfUpgradeUseCase:     selfUpgradeUseCase,
		applyComponentUseCase:  applyComponentUseCase,
		applyDogusUseCase:      applyDogusUseCase,
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

	err = useCase.prepareBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}

	if !blueprint.ShouldBeApplied() {
		// just stop the loop here on dry run or early exit
		return nil
	}

	// === Apply from here on ===
	err = useCase.applyBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}

	logger.Info("blueprint successfully applied")
	return nil
}

func (useCase *BlueprintSpecChangeUseCase) prepareBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.validation.ValidateBlueprintSpecStatically(ctx, blueprint)
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
	return nil
}

func (useCase *BlueprintSpecChangeUseCase) applyBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.selfUpgradeUseCase.HandleSelfUpgrade(ctx, blueprint)
	if err != nil {
		// could be a domain.AwaitSelfUpgradeError to trigger another reconcile
		return err
	}
	err = useCase.ecosystemConfigUseCase.ApplyConfig(ctx, blueprint)
	if err != nil {
		return err
	}
	changedComponents, err := useCase.applyComponentUseCase.ApplyComponents(ctx, blueprint)
	if err != nil {
		return err
	}
	// check after applying components
	if changedComponents {
		_, err = useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
		if err != nil {
			return err
		}
	}
	changedDogus, err := useCase.applyDogusUseCase.ApplyDogus(ctx, blueprint)
	if err != nil {
		return err
	}
	// check after installing or updating dogus
	if changedDogus {
		_, err = useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprint)
		if err != nil {
			return err
		}
	}

	err = useCase.applyUseCase.CompleteBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}
	return nil
}
