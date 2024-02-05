package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type StateDiffUseCase struct {
	blueprintSpecRepo    domainservice.BlueprintSpecRepository
	doguInstallationRepo domainservice.DoguInstallationRepository
}

func NewStateDiffUseCase(blueprintSpecRepo domainservice.BlueprintSpecRepository, doguInstallationRepo domainservice.DoguInstallationRepository) *StateDiffUseCase {
	return &StateDiffUseCase{blueprintSpecRepo: blueprintSpecRepo, doguInstallationRepo: doguInstallationRepo}
}

// DetermineStateDiff loads the state of the ecosystem and compares it to the blueprint. It creates a declarative diff.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
// a domain.InvalidBlueprintError if there are any forbidden actions in the stateDiff.
// any error if there is any other error.
func (useCase *StateDiffUseCase) DetermineStateDiff(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("StateDiffUseCase.DetermineStateDiff").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for determining state diff")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to determine state diff: %w", blueprintId, err)
	}

	logger.Info("determine state diff to the cloudogu ecosystem", "blueprintStatus", blueprintSpec.Status)
	installedDogus, err := useCase.doguInstallationRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot get installed dogus to determine state diff: %w", err)
	}

	// for now, state diff only takes dogus, but there will be components and registry keys as well
	stateDiffError := blueprintSpec.DetermineStateDiff(installedDogus)

	var invalidError *domain.InvalidBlueprintError
	if errors.As(stateDiffError, &invalidError) {
		// do not return here as with this error the blueprint status and events should be persisted as normal.
	} else if stateDiffError != nil {
		return fmt.Errorf("failed to determine state diff for blueprint %q: %w", blueprintId, stateDiffError)
	}

	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after determining the state diff to the ecosystem: %w", blueprintId, err)
	}

	// return this error back here to persist the blueprint status and events first.
	// return it to signal that a repeated call to this function will not result in any progress.
	return stateDiffError
}
