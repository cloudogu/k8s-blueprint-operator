package debugmodecr

import (
	"context"

	v1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// interface replication for generating mocks

//nolint:unused
type DebugModeInterface interface {
	// Get takes name of the debugMode, and returns the corresponding debugMode object, and an error if there is any.
	Get(ctx context.Context, name string, opts metav1.GetOptions) (result *v1.DebugMode, err error)
}
