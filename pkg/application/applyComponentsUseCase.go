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
func (useCase *ApplyComponentsUseCase) ApplyComponents(ctx context.Context, blueprint *domain.BlueprintSpec) (bool, error) {
	err := useCase.componentUseCase.ApplyComponentStates(ctx, blueprint)
	isComponentsApplied := blueprint.StateDiff.ComponentDiffs.HasChanges() && err == nil
	if isComponentsApplied {
		blueprint.Events = append(blueprint.Events, domain.ComponentsAppliedEvent{Diffs: blueprint.StateDiff.ComponentDiffs})
	}
	conditionChanged := blueprint.SetLastApplySucceededCondition(domain.ReasonLastApplyErrorAtComponents, err)

	if isComponentsApplied || conditionChanged {
		updateErr := useCase.repo.Update(ctx, blueprint)
		if updateErr != nil {
			return isComponentsApplied, fmt.Errorf("cannot update status while applying components: %w", errors.Join(updateErr, err))
		}
	}
	return isComponentsApplied, err
}
