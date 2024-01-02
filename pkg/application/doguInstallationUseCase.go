package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguInstallationUseCase struct {
	blueprintSpecRepo domainservice.BlueprintSpecRepository
	doguRepo          domainservice.DoguInstallationRepository
}

func NewDoguInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguRepo domainservice.DoguInstallationRepository,
) *DoguInstallationUseCase {
	return &DoguInstallationUseCase{
		blueprintSpecRepo: blueprintSpecRepo,
		doguRepo:          doguRepo,
	}
}

func (useCase *DoguInstallationUseCase) CheckDoguHealth(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.CheckDoguHealth").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for checking dogu health")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to check dogu health: %w", blueprintId, err)
	}

	installedDogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot evaluate dogu health states: %w", err)
	}

	blueprintSpec.CheckDoguHealth(installedDogus)

	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after Determining the state diff to the ecosystem: %w", blueprintId, err)
	}

	return nil
}
