package doguregistry

import (
	"context"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
)

type remoteDoguDescriptorRepository interface {
	cescommons.RemoteDoguDescriptorRepository
}

type localDoguDescriptorRepository interface {
	Get(ctx context.Context, doguVersion cescommons.SimpleNameVersion) (*core.Dogu, error)
	Add(ctx context.Context, name cescommons.SimpleName, dogu *core.Dogu) error
}
