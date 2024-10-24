package maintenance

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

// mock for repository.MaintenanceModeAdapter
type libMaintenanceModeAdapter interface {
	Activate(ctx context.Context, content repository.MaintenanceModeDescription) error
	Deactivate(ctx context.Context) error
}
