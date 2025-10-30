package application

import (
	"context"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type RestoreInProgressUseCase struct {
	restoreRepo restoreRepository
}

func NewRestoreInProgressUseCase(restoreRepo restoreRepository) *RestoreInProgressUseCase {
	return &RestoreInProgressUseCase{
		restoreRepo: restoreRepo,
	}
}

// CheckRestoreInProgress checks if a restore is currently in progress.
// returns a domain.RestoreInProgressError if a restore is in progress.
// returns a domainservice.InternalError if there was any other problem.
func (useCase *RestoreInProgressUseCase) CheckRestoreInProgress(ctx context.Context) error {
	restoreInProgress, err := useCase.restoreRepo.IsRestoreInProgress(ctx)
	if err != nil {
		return domainservice.NewInternalError(err, "error while checking if a restore is in progress")
	}

	if restoreInProgress {
		return &domain.RestoreInProgressError{Message: "cannot apply blueprint because a restore is in progress"}
	}

	return nil
}
