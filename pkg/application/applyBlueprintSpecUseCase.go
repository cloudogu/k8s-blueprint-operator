package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ApplyBlueprintSpecUseCase contains all use cases which are needed for or around applying
// the new ecosystem state after the determining the state diff.
type ApplyBlueprintSpecUseCase struct {
	repo               blueprintSpecRepository
	doguInstallUseCase doguInstallationUseCase
	healthUseCase      ecosystemHealthUseCase
}

func NewApplyBlueprintSpecUseCase(
	repo domainservice.BlueprintSpecRepository,
	doguInstallUseCase doguInstallationUseCase,
	healthUseCase ecosystemHealthUseCase,
) *ApplyBlueprintSpecUseCase {
	return &ApplyBlueprintSpecUseCase{
		repo:               repo,
		doguInstallUseCase: doguInstallUseCase,
		healthUseCase:      healthUseCase,
	}
}

// CheckEcosystemHealthUpfront checks the ecosystem health before applying the blueprint and sets the related status in the blueprint.
// Returns domainservice.ConflictError if there was a concurrent update to the blueprint spec or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state or
// There is no error, if the ecosystem is unhealthy as this gets reflected in the blueprint spec status.
func (useCase *ApplyBlueprintSpecUseCase) CheckEcosystemHealthUpfront(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ApplyBlueprintSpecUseCase.CheckEcosystemHealthUpfront").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for checking ecosystem health")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to check ecosystem health: %w", blueprintId, err)
	}

	healthResult, err := useCase.healthUseCase.CheckEcosystemHealth(ctx, blueprintSpec.Config.IgnoreDoguHealth, blueprintSpec.Config.IgnoreComponentHealth)
	if err != nil {
		return fmt.Errorf("cannot check ecosystem health upfront of applying the blueprint %q: %w", blueprintId, err)
	}
	blueprintSpec.CheckEcosystemHealthUpfront(healthResult)

	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after checking the ecosystem health: %w", blueprintId, err)
	}

	return nil
}

// CheckEcosystemHealthAfterwards waits for a healthy ecosystem health after applying the blueprint and sets the related status in the blueprint.
// Returns domainservice.ConflictError if there was a concurrent update to the blueprint spec or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
// There is no error, if the ecosystem is unhealthy as this gets reflected in the blueprint spec status.
func (useCase *ApplyBlueprintSpecUseCase) CheckEcosystemHealthAfterwards(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ApplyBlueprintSpecUseCase.CheckEcosystemHealthAfterwards").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for checking ecosystem health afterwards")
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to check ecosystem health: %w", blueprintId, err)
	}

	healthResult, err := useCase.healthUseCase.WaitForHealthyEcosystem(ctx)
	if err != nil {
		return fmt.Errorf("cannot check ecosystem health after applying the blueprint %q: %w", blueprintId, err)
	}
	blueprintSpec.CheckEcosystemHealthAfterwards(healthResult)

	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after checking the ecosystem health: %w", blueprintId, err)
	}

	return nil
}

// TODO: activate maintenance mode

// ApplyBlueprintSpec applies the expected state to the ecosystem. It will stop if any unexpected error happens and sets blueprint status.
// Returns domainservice.ConflictError if there was a concurrent update to the blueprint spec or other resources or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
// There is no error, if the ecosystem is unhealthy as this gets reflected in the blueprint spec status.
func (useCase *ApplyBlueprintSpecUseCase) ApplyBlueprintSpec(ctx context.Context, blueprintId string) error {
	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint to apply blueprint spec: %w", err)
	}

	err = useCase.markInProgress(ctx, blueprintSpec)
	if err != nil {
		return err
	}

	applyError := useCase.doguInstallUseCase.ApplyDoguStates(ctx, blueprintId)
	if applyError != nil {
		err := useCase.MarkFailed(ctx, blueprintSpec, err)
		if err != nil {
			return err
		}
		return applyError
	}

	return useCase.markBlueprintApplied(ctx, blueprintSpec)
}

//TODO: deactivate maintenance mode

func (useCase *ApplyBlueprintSpecUseCase) markInProgress(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkInProgress()
	err := useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as in progress: %w", err)
	}
	return nil
}

// MarkFailed marks the blueprint as failed. The error which leads to the failed blueprint needs to be provided.
func (useCase *ApplyBlueprintSpecUseCase) MarkFailed(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) error {
	logger := log.FromContext(ctx).
		WithName("ApplyBlueprintSpecUseCase.MarkFailed").
		WithValues("blueprintId", blueprintSpec.Id)

	blueprintSpec.MarkFailed(err)
	repoErr := useCase.repo.Update(ctx, blueprintSpec)

	if repoErr != nil {
		repoErr = errors.Join(repoErr, err)
		logger.Error(repoErr, "cannot mark blueprint as failed")
		return fmt.Errorf("cannot mark blueprint as failed while handling %q status: %w", blueprintSpec.Status, repoErr)
	}
	return nil
}

func (useCase *ApplyBlueprintSpecUseCase) markBlueprintApplied(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkBlueprintApplied()
	err := useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as waiting for a healthy ecosystem: %w", err)
	}
	return nil
}

func (useCase *ApplyBlueprintSpecUseCase) markCompleted(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkCompleted()
	err := useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as completed: %w", err)
	}
	return nil
}
