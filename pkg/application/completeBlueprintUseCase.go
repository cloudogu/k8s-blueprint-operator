package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

// CompleteBlueprintUseCase contains all use cases which are needed for or around applying
// the new ecosystem state after the determining the state diff.
type CompleteBlueprintUseCase struct {
	repo blueprintSpecRepository
}

func NewCompleteBlueprintUseCase(
	repo blueprintSpecRepository,
) *CompleteBlueprintUseCase {
	return &CompleteBlueprintUseCase{
		repo: repo,
	}
}

// CompleteBlueprint handles the completion of the blueprint after all other steps were successful.
// returns a domainservice.InternalError on any error.
func (useCase *CompleteBlueprintUseCase) CompleteBlueprint(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	// Only complete if there are no changes left
	if blueprint.StateDiff.HasChanges() {
		return &domain.StateDiffNotEmptyError{Message: "cannot complete blueprint because the StateDiff has still changes"}
	}
	changed := blueprint.Complete()
	if changed {
		err := useCase.repo.Update(ctx, blueprint)
		if err != nil {
			return fmt.Errorf("cannot update blueprint to complete it: %w", err)
		}
	}
	return nil
}
