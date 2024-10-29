package kubernetes

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

//only to generate mocks
//see k8s-registry-lib for possible go docs

type k8sGlobalConfigRepo interface {
	Get(ctx context.Context) (config.GlobalConfig, error)
	Update(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error)
}

type k8sDoguConfigRepo interface {
	Get(context.Context, common.SimpleDoguName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
	Create(context.Context, config.DoguConfig) (config.DoguConfig, error)
}
