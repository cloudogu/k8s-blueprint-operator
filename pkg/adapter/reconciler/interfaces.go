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
	HandleChange(ctx context.Context, blueprintId string) error
}
