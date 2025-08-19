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
func (useCase *ApplyDogusUseCase) ApplyDogus(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	err := useCase.doguInstallUseCase.ApplyDoguStates(ctx, blueprint)
	changed := blueprint.SetDogusAppliedCondition(err)

	if changed {
		updateErr := useCase.repo.Update(ctx, blueprint)
		if updateErr != nil {
			return fmt.Errorf("cannot update condition while applying dogus: %w", errors.Join(updateErr, err))
		}
	}
	return err
}
