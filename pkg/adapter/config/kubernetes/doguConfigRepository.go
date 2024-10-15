package kubernetes

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

type DoguConfigRepository struct {
	repo repository.DoguConfigRepository
}

func (e DoguConfigRepository) GetAll(ctx context.Context, doguNames []common.SimpleDoguName) (map[common.SimpleDoguName]config.DoguConfig, error) {
	var configByDogus map[common.SimpleDoguName]config.DoguConfig
	for _, doguName := range doguNames {
		loaded, err := e.Get(ctx, doguName)
		if err != nil {
			return nil, fmt.Errorf("could not load config for all given dogus: %w", err)
		}
		configByDogus[doguName] = loaded
	}
	return configByDogus, nil
}

func NewDoguConfigRepository(repo repository.DoguConfigRepository) *DoguConfigRepository {
	return &DoguConfigRepository{repo: repo}
}

func (e DoguConfigRepository) Get(ctx context.Context, doguName common.SimpleDoguName) (config.DoguConfig, error) {
	loadedConfig, err := e.repo.Get(ctx, doguName)
	if err != nil {
		return loadedConfig, fmt.Errorf("could not load dogu config: %w", mapToBlueprintError(err))
	}
	return loadedConfig, nil
}

func (e DoguConfigRepository) Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	updatedConfig, err := e.repo.Update(ctx, config)
	if err != nil {
		return updatedConfig, fmt.Errorf("could not update dogu config: %w", mapToBlueprintError(err))
	}
	return updatedConfig, nil
}
