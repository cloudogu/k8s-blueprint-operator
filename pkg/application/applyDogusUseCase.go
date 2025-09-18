package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

// ApplyDogusUseCase can handle dogu installations, updates and deletions.
type ApplyDogusUseCase struct {
	repo               blueprintSpecRepository
	doguInstallUseCase doguInstallationUseCase
}

func NewApplyDogusUseCase(
	repo blueprintSpecRepository,
	doguInstallUseCase doguInstallationUseCase,
) *ApplyDogusUseCase {
	return &ApplyDogusUseCase{
		repo:               repo,
		doguInstallUseCase: doguInstallUseCase,
	}
}

// ApplyDogus applies dogus if necessary.
// The conditions in the blueprint will be set accordingly.
// returns domainservice.ConflictError if there was a concurrent update to the blueprint or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
func (useCase *ApplyDogusUseCase) ApplyDogus(ctx context.Context, blueprint *domain.BlueprintSpec) (bool, error) {
	err := useCase.doguInstallUseCase.ApplyDoguStates(ctx, blueprint)
	isDogusApplied := blueprint.StateDiff.DoguDiffs.HasChanges() && err == nil
	if isDogusApplied {
		blueprint.Events = append(blueprint.Events, domain.DogusAppliedEvent{Diffs: blueprint.StateDiff.DoguDiffs})
	}
	conditionChanged := blueprint.SetLastApplySucceededCondition(domain.ReasonLastApplyErrorAtDogus, err)

	if isDogusApplied || conditionChanged {
		updateErr := useCase.repo.Update(ctx, blueprint)
		if updateErr != nil {
			return isDogusApplied, fmt.Errorf("cannot update status while applying dogus: %w", errors.Join(updateErr, err))
		}
	}
	return isDogusApplied, err
}
