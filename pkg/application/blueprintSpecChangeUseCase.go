package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type BlueprintPreparationUseCases struct {
	initialStatus      initialBlueprintStatusUseCase
	validation         blueprintSpecValidationUseCase
	effectiveBlueprint effectiveBlueprintUseCase
	stateDiff          stateDiffUseCase
	healthUseCase      ecosystemHealthUseCase
}

func NewBlueprintPreparationUseCases(
	initialStatus initialBlueprintStatusUseCase,
	validation blueprintSpecValidationUseCase,
	effectiveBlueprint effectiveBlueprintUseCase,
	stateDiff stateDiffUseCase,
	ecosystemHealthUseCase ecosystemHealthUseCase,
) BlueprintPreparationUseCases {
	return BlueprintPreparationUseCases{
		initialStatus:      initialStatus,
		validation:         validation,
		effectiveBlueprint: effectiveBlueprint,
		stateDiff:          stateDiff,
		healthUseCase:      ecosystemHealthUseCase,
	}
}

type BlueprintApplyUseCases struct {
	completeUseCase        completeBlueprintUseCase
	ecosystemConfigUseCase ecosystemConfigUseCase
	selfUpgradeUseCase     selfUpgradeUseCase
	applyComponentUseCase  applyComponentsUseCase
	applyDogusUseCase      applyDogusUseCase
	healthUseCase          ecosystemHealthUseCase
}

func NewBlueprintApplyUseCases(
	completeUseCase completeBlueprintUseCase,
	ecosystemConfigUseCase ecosystemConfigUseCase,
	selfUpgradeUseCase selfUpgradeUseCase,
	applyComponentUseCase applyComponentsUseCase,
	applyDogusUseCase applyDogusUseCase,
	healthUseCase ecosystemHealthUseCase,
) BlueprintApplyUseCases {
	return BlueprintApplyUseCases{
		completeUseCase:        completeUseCase,
		ecosystemConfigUseCase: ecosystemConfigUseCase,
		selfUpgradeUseCase:     selfUpgradeUseCase,
		applyComponentUseCase:  applyComponentUseCase,
		applyDogusUseCase:      applyDogusUseCase,
		healthUseCase:          healthUseCase,
	}
}

type BlueprintSpecChangeUseCase struct {
	repo                blueprintSpecRepository
	preparationUseCases BlueprintPreparationUseCases
	applyUseCases       BlueprintApplyUseCases
}

func NewBlueprintSpecChangeUseCase(
	repo blueprintSpecRepository,
	preparationUseCases BlueprintPreparationUseCases,
	applyUseCases BlueprintApplyUseCases,
) *BlueprintSpecChangeUseCase {
	return &BlueprintSpecChangeUseCase{
		repo:                repo,
		preparationUseCases: preparationUseCases,
		applyUseCases:       applyUseCases,
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

	err = useCase.preparationUseCases.prepareBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}

	if !blueprint.ShouldBeApplied() {
		// just stop the loop here on dry run or early exit
		return nil
	}

	// === Apply from here on ===
	err = useCase.applyUseCases.applyBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}

	logger.Info("blueprint successfully applied")
	return nil
}

func (useCase *BlueprintPreparationUseCases) prepareBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
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

func (useCase *BlueprintApplyUseCases) applyBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.selfUpgradeUseCase.HandleSelfUpgrade(ctx, blueprint)
	if err != nil {
		// could be a domain.AwaitSelfUpgradeError to trigger another reconcile
		return err
	}
	err = useCase.ecosystemConfigUseCase.ApplyConfig(ctx, blueprint)
	// TODO: prevent Dogu restarts here to avoid starting Dogus with config not fitting to the expected dogu version
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
	// TODO: check if config in dogus is already up to date and if installed Version is up to date

	err = useCase.completeUseCase.CompleteBlueprint(ctx, blueprint)
	if err != nil {
		return err
	}
	return nil
}
