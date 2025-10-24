package restorecr

import (
	"context"

	restorev1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// interface replication for generating mocks

//nolint:unused
type RestoreInterface interface {
	// List takes label and field selectors, and returns the list of Restores that match those selectors.
	List(ctx context.Context, opts metav1.ListOptions) (*restorev1.RestoreList, error)
}
