package application

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type StateDiffUseCase struct {
	blueprintSpecRepo         domainservice.BlueprintSpecRepository
	doguInstallationRepo      domainservice.DoguInstallationRepository
	componentInstallationRepo domainservice.ComponentInstallationRepository
}

func NewStateDiffUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguInstallationRepo domainservice.DoguInstallationRepository,
	componentInstallationRepo domainservice.ComponentInstallationRepository,
) *StateDiffUseCase {
	return &StateDiffUseCase{
		blueprintSpecRepo:         blueprintSpecRepo,
		doguInstallationRepo:      doguInstallationRepo,
		componentInstallationRepo: componentInstallationRepo,
	}
}

// DetermineStateDiff loads the state of the ecosystem and compares it to the blueprint. It creates a declarative diff.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
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

	installedComponents, err := useCase.componentInstallationRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot get installed components to determine state diff: %w", err)
	}

	// for now, state diff only takes dogus and components, but there will be registry keys as well
	stateDiffError := blueprintSpec.DetermineStateDiff(installedDogus, installedComponents)
	if stateDiffError != nil {
		return fmt.Errorf("failed to determine state diff for blueprint %q: %w", blueprintId, stateDiffError)
	}

	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after Determining the state diff to the ecosystem: %w", blueprintId, err)
	}

	return nil
}
