package restorecr

import (
	"context"
	"fmt"

	restorev1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type restoreRepo struct {
	restoreClient RestoreInterface
}

// NewRestoreRepo returns a new restoreRepo to interact with the restore CR.
func NewRestoreRepo(restoreClient RestoreInterface) domainservice.RestoreRepository {
	return &restoreRepo{restoreClient: restoreClient}
}

func (repo *restoreRepo) IsRestoreInProgress(ctx context.Context) (bool, error) {
	list, err := repo.restoreClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// no restores found, so no restore can be in progress
			return false, nil
		}

		return false, fmt.Errorf("error while listing restore CRs: %w", err)
	}

	for _, restore := range list.Items {
		if restore.Status.Status == restorev1.RestoreStatusInProgress {
			return true, nil
		}
	}

	return false, nil
}
