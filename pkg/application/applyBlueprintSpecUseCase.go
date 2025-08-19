package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
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
	repo blueprintSpecRepository,
	doguInstallUseCase doguInstallationUseCase,
	healthUseCase ecosystemHealthUseCase,
) *ApplyBlueprintSpecUseCase {
	return &ApplyBlueprintSpecUseCase{
		repo:               repo,
		doguInstallUseCase: doguInstallUseCase,
		healthUseCase:      healthUseCase,
	}
}

// PostProcessBlueprintApplication makes changes to the environment after applying the blueprint.
// returns a domainservice.InternalError on any error.
func (useCase *ApplyBlueprintSpecUseCase) PostProcessBlueprintApplication(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	blueprint.CompletePostProcessing()

	err := useCase.repo.Update(ctx, blueprint)
	if err != nil {
		return fmt.Errorf("cannot update blueprint spec %q while post-processing blueprint application: %w", blueprint.Id, err)
	}

	return nil
}

// ApplyBlueprintSpec applies the expected state to the ecosystem. It will stop if any unexpected error happens and sets blueprint status.
// Returns domainservice.ConflictError if there was a concurrent update to the blueprint spec or other resources or
// returns a domainservice.InternalError if there was an unspecified error while collecting or modifying the ecosystem state.
// There is no error, if the ecosystem is unhealthy as this gets reflected in the blueprint spec status.
func (useCase *ApplyBlueprintSpecUseCase) ApplyBlueprintSpec(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("ApplyBlueprintSpecUseCase.ApplyBlueprintSpec")

	applyError := useCase.doguInstallUseCase.ApplyDoguStates(ctx, blueprint)
	if applyError != nil {
		return useCase.handleApplyFailedError(ctx, blueprint, applyError)
	}

	// FIXME: this health check is blocking. I think we need to split the apply logic into multiple steps to
	// we have to wait for all dogus to be healthy
	// otherwise service account creation might fail because dogus are restarted right after this step
	_, err := useCase.doguInstallUseCase.WaitForHealthyDogus(ctx)
	if err != nil {
		return useCase.handleApplyFailedError(ctx, blueprint, err)
	}

	logger.Info("blueprint successfully applied to the cluster")
	return useCase.markBlueprintApplied(ctx, blueprint)
}

func (useCase *ApplyBlueprintSpecUseCase) handleApplyFailedError(ctx context.Context, blueprintSpec *domain.BlueprintSpec, applyError error) error {
	err := useCase.markBlueprintApplicationFailed(ctx, blueprintSpec, applyError)
	if err != nil {
		return err
	}
	return applyError
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
