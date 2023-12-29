package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EffectiveBlueprintUseCase struct {
	blueprintSpecRepo domainservice.BlueprintSpecRepository
}

func NewEffectiveBlueprintUseCase(blueprintSpecRepo domainservice.BlueprintSpecRepository) *EffectiveBlueprintUseCase {
	return &EffectiveBlueprintUseCase{blueprintSpecRepo: blueprintSpecRepo}
}

// CalculateEffectiveBlueprint loads the blueprintSpec, lets it calculate the effective blueprint and persists it again.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *EffectiveBlueprintUseCase) CalculateEffectiveBlueprint(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("EffectiveBlueprintUseCase.CalculateEffectiveBlueprint").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for effective blueprint calculation")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec to calculate effective blueprint: %w", err)
	}

	logger.Info("calculate effective blueprint", "blueprintStatus", blueprintSpec.Status)
	calcError := blueprintSpec.CalculateEffectiveBlueprint()
	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec after calculating the effective blueprint: %w", err)
	}

	return calcError
}
