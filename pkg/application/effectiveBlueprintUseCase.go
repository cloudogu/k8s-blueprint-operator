package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type EffectiveBlueprintUseCase struct {
	blueprintSpecRepo blueprintSpecRepository
}

func NewEffectiveBlueprintUseCase(blueprintSpecRepo domainservice.BlueprintSpecRepository) *EffectiveBlueprintUseCase {
	return &EffectiveBlueprintUseCase{blueprintSpecRepo: blueprintSpecRepo}
}

// CalculateEffectiveBlueprint loads the blueprintSpec, lets it calculate the effective blueprint and persists it again.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *EffectiveBlueprintUseCase) CalculateEffectiveBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	calcError := blueprint.CalculateEffectiveBlueprint()
	err := useCase.blueprintSpecRepo.Update(ctx, blueprint)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec after calculating the effective blueprint: %w", err)
	}

	return calcError
}
