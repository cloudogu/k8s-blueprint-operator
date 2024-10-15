package kubernetes

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type k8sGlobalConfigRepo interface {
	Get(ctx context.Context) (config.GlobalConfig, error)
	Update(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error)
}

type k8sDoguConfigRepo interface {
	Get(ctx context.Context, doguName common.SimpleDoguName) (config.DoguConfig, error)
	Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error)
}
