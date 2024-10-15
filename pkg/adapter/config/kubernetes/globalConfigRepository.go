package kubernetes

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type GlobalConfigRepository struct {
	repo k8sGlobalConfigRepo
}

func NewGlobalConfigRepository(repo k8sGlobalConfigRepo) *GlobalConfigRepository {
	return &GlobalConfigRepository{repo: repo}
}

func (e GlobalConfigRepository) Get(ctx context.Context) (config.GlobalConfig, error) {
	loadedConfig, err := e.repo.Get(ctx)
	if err != nil {
		return loadedConfig, fmt.Errorf("could not load global config: %w", mapToBlueprintError(err))
	}
	return loadedConfig, nil
}

func (e GlobalConfigRepository) Update(ctx context.Context, config config.GlobalConfig) (config.GlobalConfig, error) {
	updatedConfig, err := e.repo.Update(ctx, config)
	if err != nil {
		return updatedConfig, fmt.Errorf("could not update global config: %w", mapToBlueprintError(err))
	}
	return updatedConfig, nil
}
