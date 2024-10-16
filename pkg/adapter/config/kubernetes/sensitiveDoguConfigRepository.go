package kubernetes

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type SensitiveDoguConfigRepository struct {
	repo k8sDoguConfigRepo
}

func NewSensitiveDoguConfigRepository(repo k8sDoguConfigRepo) *SensitiveDoguConfigRepository {
	return &SensitiveDoguConfigRepository{repo: repo}
}

func (repo *SensitiveDoguConfigRepository) GetAll(ctx context.Context, doguNames []common.SimpleDoguName) (map[common.SimpleDoguName]config.DoguConfig, error) {
	var configByDogus = map[common.SimpleDoguName]config.DoguConfig{}
	for _, doguName := range doguNames {
		loaded, err := repo.Get(ctx, doguName)
		if err != nil {
			return nil, fmt.Errorf("could not load sensitive config for all given dogus: %w", err)
		}
		configByDogus[doguName] = loaded
	}
	return configByDogus, nil
}

func (repo *SensitiveDoguConfigRepository) Get(ctx context.Context, doguName common.SimpleDoguName) (config.DoguConfig, error) {
	loadedConfig, err := repo.repo.Get(ctx, doguName)
	if err != nil {
		return loadedConfig, fmt.Errorf("could not load sensitive dogu config for %s: %w", doguName, mapToBlueprintError(err))
	}
	return loadedConfig, nil
}

func (repo *SensitiveDoguConfigRepository) Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	updatedConfig, err := repo.repo.Update(ctx, config)
	if err != nil {
		return updatedConfig, fmt.Errorf("could not update sensitive dogu config for %s: %w", config.DoguName, mapToBlueprintError(err))
	}
	return updatedConfig, nil
}

func (repo *SensitiveDoguConfigRepository) Create(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	createdConfig, err := repo.repo.Create(ctx, config)
	if err != nil {
		return createdConfig, fmt.Errorf("could not create sensitive dogu config for %s: %w", config.DoguName, mapToBlueprintError(err))
	}
	return createdConfig, nil
}
