package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EffectiveBlueprintUseCase struct {
	blueprintSpecRepo blueprintSpecRepository
	debugModeRepo     debugModeRepository
}

func NewEffectiveBlueprintUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	debugModeRepo domainservice.DebugModeRepository,
) *EffectiveBlueprintUseCase {
	return &EffectiveBlueprintUseCase{
		blueprintSpecRepo: blueprintSpecRepo,
		debugModeRepo:     debugModeRepo,
	}
}

// CalculateEffectiveBlueprint loads the blueprintSpec, lets it calculate the effective blueprint and persists it again.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
func (useCase *EffectiveBlueprintUseCase) CalculateEffectiveBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("EffectiveBlueprintUseCase.CalculateEffectiveBlueprint")
	debugMode, err := useCase.debugModeRepo.GetSingleton(ctx)
	if err != nil {
		// ignore not found error, no debug mode cr means we are not in debug mode
		if !domainservice.IsNotFoundError(err) {
			return fmt.Errorf("cannot calculate effective blueprint due to an error when loading the debug mode cr: %w", err)
		}
	}
	isDebugModeActive := debugMode != nil && debugMode.IsActive()
	if isDebugModeActive {
		logger.Info("debug mode is active, will ignore loglevel changes until deactivated")
	}
	calcError := blueprint.CalculateEffectiveBlueprint(isDebugModeActive)
	err = useCase.blueprintSpecRepo.Update(ctx, blueprint)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec after calculating the effective blueprint: %w", err)
	}

	return calcError
}
