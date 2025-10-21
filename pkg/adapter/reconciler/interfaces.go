package reconciler

import (
	"context"

	bpv2client "github.com/cloudogu/k8s-blueprint-lib/v2/client"
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

type blueprintMaskInterface interface {
	bpv2client.BlueprintMaskInterface
}

type blueprintInterface interface {
	bpv2client.BlueprintInterface
}
