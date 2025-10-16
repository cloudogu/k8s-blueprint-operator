package reconciler

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// used for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type controllerManager interface {
	manager.Manager
}

type BlueprintChangeHandler interface {
	HandleUntilApplied(ctx context.Context, blueprintId string) error
	CheckForMultipleBlueprintResources(ctx context.Context) error
}
