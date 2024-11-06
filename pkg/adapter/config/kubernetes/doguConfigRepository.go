package kubernetes

import (
	"context"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
	liberrors "github.com/cloudogu/k8s-registry-lib/errors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ConfigRepoType string

var (
	normalConfig    ConfigRepoType = "normal dogu config"
	sensitiveConfig ConfigRepoType = "sensitive dogu config"
)

type DoguConfigRepository struct {
	repo     k8sDoguConfigRepo
	repoType ConfigRepoType
}

func NewDoguConfigRepository(repo k8sDoguConfigRepo) *DoguConfigRepository {
	return &DoguConfigRepository{repo: repo, repoType: normalConfig}
}

func NewSensitiveDoguConfigRepository(repo k8sDoguConfigRepo) *DoguConfigRepository {
	return &DoguConfigRepository{repo: repo, repoType: sensitiveConfig}
}

func (repo *DoguConfigRepository) GetAll(ctx context.Context, doguNames []cescommons.SimpleDoguName) (map[cescommons.SimpleDoguName]config.DoguConfig, error) {
	var configByDogus = map[cescommons.SimpleDoguName]config.DoguConfig{}
	for _, doguName := range doguNames {
		loaded, err := repo.Get(ctx, doguName)
		if err != nil {
			return nil, fmt.Errorf("could not load %s for all given dogus: %w", repo.repoType, err)
		}
		configByDogus[doguName] = loaded
	}
	return configByDogus, nil
}

func (repo *DoguConfigRepository) GetAllExisting(ctx context.Context, doguNames []cescommons.SimpleDoguName) (map[cescommons.SimpleDoguName]config.DoguConfig, error) {
	var configByDogus = map[cescommons.SimpleDoguName]config.DoguConfig{}
	for _, doguName := range doguNames {
		loaded, err := repo.Get(ctx, doguName)
		if liberrors.IsNotFoundError(err) {
			// if notFoundError happens, the dogu is not yet installed. Therefore, the config is empty
			loaded = config.CreateDoguConfig(doguName, map[config.Key]config.Value{})
		} else if err != nil {
			return nil, fmt.Errorf("could not load %s for all given dogus: %w", repo.repoType, err)
		}
		configByDogus[doguName] = loaded
	}
	return configByDogus, nil
}

func (repo *DoguConfigRepository) Get(ctx context.Context, doguName cescommons.SimpleDoguName) (config.DoguConfig, error) {
	loadedConfig, err := repo.repo.Get(ctx, doguName)
	if err != nil {
		return loadedConfig, fmt.Errorf("could not load %s for %s: %w", repo.repoType, doguName, mapToBlueprintError(err))
	}
	return loadedConfig, nil
}

func (repo *DoguConfigRepository) Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	updatedConfig, err := repo.repo.Update(ctx, config)
	if err != nil {
		return updatedConfig, fmt.Errorf("could not update %s for %s: %w", repo.repoType, config.DoguName, mapToBlueprintError(err))
	}
	return updatedConfig, nil
}

func (repo *DoguConfigRepository) Create(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	createdConfig, err := repo.repo.Create(ctx, config)
	if err != nil {
		return createdConfig, fmt.Errorf("could not create %s for %s: %w", repo.repoType, config.DoguName, mapToBlueprintError(err))
	}
	return createdConfig, nil
}

func (repo *DoguConfigRepository) UpdateOrCreate(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	logger := log.FromContext(ctx).
		WithName("DoguConfigRepository.UpdateOrCreate").
		WithValues("dogu", config.DoguName)

	updatedConfig, err := repo.repo.Update(ctx, config)
	if err != nil {
		if liberrors.IsNotFoundError(err) {
			logger.Info("dogu config is not present, try to create it", "error", err)
			return repo.Create(ctx, config)
		}
		return updatedConfig, fmt.Errorf("could not update %s for %s: %w", repo.repoType, config.DoguName, mapToBlueprintError(err))
	}
	return updatedConfig, nil
}
