package application

import (
	"context"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

type BlueprintApplyUseCase struct {
	completeUseCase        completeBlueprintUseCase
	ecosystemConfigUseCase ecosystemConfigUseCase
	applyDogusUseCase      applyDogusUseCase
	healthUseCase          ecosystemHealthUseCase
	dogusUpToDateUseCase   dogusUpToDateUseCase
}

func NewBlueprintApplyUseCase(
	completeUseCase completeBlueprintUseCase,
	ecosystemConfigUseCase ecosystemConfigUseCase,
	applyDogusUseCase applyDogusUseCase,
	healthUseCase ecosystemHealthUseCase,
	dogusUpToDateUseCase dogusUpToDateUseCase,
) BlueprintApplyUseCase {
	return BlueprintApplyUseCase{
		completeUseCase:        completeUseCase,
		ecosystemConfigUseCase: ecosystemConfigUseCase,
		applyDogusUseCase:      applyDogusUseCase,
		healthUseCase:          healthUseCase,
		dogusUpToDateUseCase:   dogusUpToDateUseCase,
	}
}

func (useCase *BlueprintApplyUseCase) applyBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.ecosystemConfigUseCase.ApplyConfig(ctx, blueprint)
	if err != nil {
		return err
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

	err = useCase.dogusUpToDateUseCase.CheckDogus(ctx, blueprint)
	if err != nil {
		// could be a domain.AwaitSelfUpgradeError to trigger another reconcile
		return err
	}

	// Only complete if there are no changes left
	if blueprint.StateDiff.HasChanges() {
		return &domain.StateDiffNotEmptyError{Message: "cannot complete blueprint because the StateDiff has still changes"}
	} else {
		err = useCase.completeUseCase.CompleteBlueprint(ctx, blueprint)
		if err != nil {
			return err
		}
	}
	return nil
}
