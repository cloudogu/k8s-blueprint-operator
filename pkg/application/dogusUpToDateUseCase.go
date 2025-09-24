package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DogusUpToDateUseCase checks if all dogus are up to date, meaning they are on the desired version and configuration.
type DogusUpToDateUseCase struct {
	repo               blueprintSpecRepository
	doguInstallUseCase doguInstallationUseCase
}

func NewDogusUpToDateUseCase(
	repo blueprintSpecRepository,
	doguInstallUseCase doguInstallationUseCase,
) *DogusUpToDateUseCase {
	return &DogusUpToDateUseCase{
		repo:               repo,
		doguInstallUseCase: doguInstallUseCase,
	}
}

// CheckDogus checks that all dogs are up to date.
// returns domainservice.ConflictError if there was a concurrent update to the blueprint or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
func (useCase *DogusUpToDateUseCase) CheckDogus(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("DogusUpToDateUseCase.CheckDogus")

	dogusNotUpToDate, err := useCase.doguInstallUseCase.CheckDogusUpToDate(ctx)
	if err != nil {
		return err
	}
	if len(dogusNotUpToDate) > 0 {
		// event and error
		blueprint.Events = append(blueprint.Events, domain.DogusNotUpToDateEvent{DogusNotUpToDate: dogusNotUpToDate})
		updateErr := useCase.repo.Update(ctx, blueprint)
		if updateErr != nil {
			return fmt.Errorf("cannot update status while checking dogus: %w", errors.Join(updateErr, err))
		}
		return &domain.DogusNotUpToDateError{Message: fmt.Sprintf("following dogus are not up to date yet: %v", dogusNotUpToDate)}
	}

	logger.V(2).Info("all dogus are up to date")
	return nil
}
