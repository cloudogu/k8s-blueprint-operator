package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

// ApplyComponentsUseCase can handle component installations, updates and deletions.
type ApplyComponentsUseCase struct {
	repo             blueprintSpecRepository
	componentUseCase componentInstallationUseCase
}

func NewApplyComponentsUseCase(
	repo blueprintSpecRepository,
	componentUseCase componentInstallationUseCase,
) *ApplyComponentsUseCase {
	return &ApplyComponentsUseCase{
		repo:             repo,
		componentUseCase: componentUseCase,
	}
}

// ApplyComponents applies components if necessary.
// The conditions in the blueprint will be set accordingly.
// returns domainservice.ConflictError if there was a concurrent update to the blueprint or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
func (useCase *ApplyComponentsUseCase) ApplyComponents(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.componentUseCase.ApplyComponentStates(ctx, blueprint)
	changed := blueprint.SetComponentsAppliedCondition(err)

	if changed {
		updateErr := useCase.repo.Update(ctx, blueprint)
		if updateErr != nil {
			return fmt.Errorf("cannot update condition while applying components: %w", errors.Join(updateErr, err))
		}
	}
	return err
}
