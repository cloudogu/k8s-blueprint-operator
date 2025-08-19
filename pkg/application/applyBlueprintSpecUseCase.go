package application

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
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
