package kubernetes

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type DoguConfigRepository struct {
	repo k8sDoguConfigRepo
}

func NewDoguConfigRepository(repo k8sDoguConfigRepo) *DoguConfigRepository {
	return &DoguConfigRepository{repo: repo}
}

func (repo *DoguConfigRepository) GetAll(ctx context.Context, doguNames []common.SimpleDoguName) (map[common.SimpleDoguName]config.DoguConfig, error) {
	var configByDogus = map[common.SimpleDoguName]config.DoguConfig{}
	for _, doguName := range doguNames {
		loaded, err := repo.Get(ctx, doguName)
		if err != nil {
			return nil, fmt.Errorf("could not load config for all given dogus: %w", err)
		}
		configByDogus[doguName] = loaded
	}
	return configByDogus, nil
}

func (repo *DoguConfigRepository) Get(ctx context.Context, doguName common.SimpleDoguName) (config.DoguConfig, error) {
	loadedConfig, err := repo.repo.Get(ctx, doguName)
	if err != nil {
		return loadedConfig, fmt.Errorf("could not load dogu config for %s: %w", doguName, mapToBlueprintError(err))
	}
	return loadedConfig, nil
}

func (repo *DoguConfigRepository) Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	updatedConfig, err := repo.repo.Update(ctx, config)
	if err != nil {
		return updatedConfig, fmt.Errorf("could not update dogu config for %s: %w", config.DoguName, mapToBlueprintError(err))
	}
	return updatedConfig, nil
}
