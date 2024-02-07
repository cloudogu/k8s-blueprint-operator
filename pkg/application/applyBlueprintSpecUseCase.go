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
	repo                    blueprintSpecRepository
	doguInstallUseCase      doguInstallationUseCase
	healthUseCase           ecosystemHealthUseCase
	componentInstallUseCase componentInstallationUseCase
	maintenanceModeAdapter  maintenanceMode
}

func NewApplyBlueprintSpecUseCase(
	repo blueprintSpecRepository,
	doguInstallUseCase doguInstallationUseCase,
	healthUseCase ecosystemHealthUseCase,
	componentInstallUseCase componentInstallationUseCase,
	maintenanceModeAdapter maintenanceMode,
) *ApplyBlueprintSpecUseCase {
	return &ApplyBlueprintSpecUseCase{
		repo:                    repo,
		doguInstallUseCase:      doguInstallUseCase,
		healthUseCase:           healthUseCase,
		componentInstallUseCase: componentInstallUseCase,
		maintenanceModeAdapter:  maintenanceModeAdapter,
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

// PreProcessBlueprintApplication prepares the environment for applying the blueprint, e.g. activating the maintenance mode.
// returns a domainservice.ConflictError if another party activated the maintenance mode or
// returns a domainservice.InternalError on any other error.
func (useCase *ApplyBlueprintSpecUseCase) PreProcessBlueprintApplication(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ApplyBlueprintSpecUseCase.PreProcessBlueprintApplication").
		WithValues("blueprintId", blueprintId)

	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to activate maintenance mode: %w", blueprintId, err)
	}

	if !blueprintSpec.ShouldBeApplied() {
		logger.Info("stop before activating maintenance mode as blueprint should not be applied")
	} else {
		logger.Info("activate maintenance mode")
		err = useCase.maintenanceModeAdapter.Activate(domainservice.MaintenancePageModel{
			Title: "Blueprint getting applied",
			Text:  "A new Blueprint with updates for the Cloudogu Ecosystem is getting applied.",
		})
		if err != nil {
			return fmt.Errorf("could not activate maintenance mode before applying the blueprint: %w", err)
		}
	}

	blueprintSpec.CompletePreProcessing()

	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after activating the maintenance mode: %w", blueprintId, err)
	}

	return nil
}

// PostProcessBlueprintApplication makes changes to the environment after applying the blueprint, e.g. deactivating the maintenance mode.
// returns a domainservice.ConflictError if another party holds the lock to the maintenance mode or
// returns a domainservice.InternalError on any other error.
func (useCase *ApplyBlueprintSpecUseCase) PostProcessBlueprintApplication(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ApplyBlueprintSpecUseCase.PostProcessBlueprintApplication").
		WithValues("blueprintId", blueprintId)

	logger.Info("deactivate maintenance mode")
	err := useCase.maintenanceModeAdapter.Deactivate()
	if err != nil {
		return fmt.Errorf("could not deactivate maintenance mode after applying the blueprint: %w", err)
	}

	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q while post-processing blueprint application: %w", blueprintId, err)
	}

	blueprintSpec.CompletePostProcessing()

	err = useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec %q while post-processing blueprint application: %w", blueprintId, err)
	}

	return nil
}

// ApplyBlueprintSpec applies the expected state to the ecosystem. It will stop if any unexpected error happens and sets blueprint status.
// Returns domainservice.ConflictError if there was a concurrent update to the blueprint spec or other resources or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
// There is no error, if the ecosystem is unhealthy as this gets reflected in the blueprint spec status.
func (useCase *ApplyBlueprintSpecUseCase) ApplyBlueprintSpec(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ApplyBlueprintSpecUseCase.ApplyBlueprintSpec").
		WithValues("blueprintId", blueprintId)

	blueprintSpec, err := useCase.repo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint to apply blueprint spec: %w", err)
	}

	logger.Info("start applying blueprint to the cluster")
	err = useCase.startApplying(ctx, blueprintSpec)
	if err != nil {
		return err
	}

	applyError := useCase.componentInstallUseCase.ApplyComponentStates(ctx, blueprintId)
	if applyError != nil {
		return useCase.handleApplyFailedError(ctx, blueprintSpec, applyError)
	}

	applyError = useCase.doguInstallUseCase.ApplyDoguStates(ctx, blueprintId)
	if applyError != nil {
		return useCase.handleApplyFailedError(ctx, blueprintSpec, applyError)
	}

	logger.Info("blueprint successfully applied to the cluster")
	return useCase.markBlueprintApplied(ctx, blueprintSpec)
}

func (useCase *ApplyBlueprintSpecUseCase) handleApplyFailedError(ctx context.Context, blueprintSpec *domain.BlueprintSpec, applyError error) error {
	err := useCase.markBlueprintApplicationFailed(ctx, blueprintSpec, applyError)
	if err != nil {
		return err
	}
	return applyError
}

func (useCase *ApplyBlueprintSpecUseCase) startApplying(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.StartApplying()
	err := useCase.repo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as in progress: %w", err)
	}
	return nil
}

// markBlueprintApplicationFailed marks the blueprint application as failed.
// Returns the error which leads to the failed blueprint needs to be provided.
func (useCase *ApplyBlueprintSpecUseCase) markBlueprintApplicationFailed(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) error {
	logger := log.FromContext(ctx).
		WithName("ApplyBlueprintSpecUseCase.markBlueprintApplicationFailed").
		WithValues("blueprintId", blueprintSpec.Id)

	blueprintSpec.MarkBlueprintApplicationFailed(err)
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
