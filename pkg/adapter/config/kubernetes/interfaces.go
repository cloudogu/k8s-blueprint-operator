package kubernetes

import (
	"context"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type k8sGlobalConfigRepo interface {
	Get(ctx context.Context) (config.GlobalConfig, error)
	Update(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error)
}
