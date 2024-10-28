package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
	"maps"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"slices"
)

type EcosystemConfigUseCase struct {
	blueprintRepository           blueprintSpecRepository
	doguConfigRepository          doguConfigRepository
	sensitiveDoguConfigRepository sensitiveDoguConfigRepository
	globalConfigRepository        globalConfigRepository
}

func NewEcosystemConfigUseCase(
	blueprintRepository blueprintSpecRepository,
	doguConfigRepository doguConfigRepository,
	sensitiveDoguConfigRepository sensitiveDoguConfigRepository,
	globalConfigRepository globalConfigRepository,
) *EcosystemConfigUseCase {
	return &EcosystemConfigUseCase{
		blueprintRepository:           blueprintRepository,
		doguConfigRepository:          doguConfigRepository,
		sensitiveDoguConfigRepository: sensitiveDoguConfigRepository,
		globalConfigRepository:        globalConfigRepository,
	}
}

// ApplyConfig fetches the dogu and global config stateDiff of the blueprint and applies these keys to the repositories.
func (useCase *EcosystemConfigUseCase) ApplyConfig(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("EcosystemConfigUseCase.ApplyConfig").
		WithValues("blueprintId", blueprintId)

	blueprintSpec, err := useCase.blueprintRepository.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint to apply config: %w", err)
	}

	doguConfigDiffs := blueprintSpec.StateDiff.DoguConfigDiffs
	isEmptyDoguDiff := len(doguConfigDiffs) == 0
	if isEmptyDoguDiff {
		logger.Info("dogu config diffs are empty...")
	}

	globalConfigDiffs := blueprintSpec.StateDiff.GlobalConfigDiffs
	isEmptyGlobalDiff := len(globalConfigDiffs) == 0
	if isEmptyGlobalDiff {
		logger.Info("global config diffs are empty...")
	}

	if isEmptyDoguDiff && isEmptyGlobalDiff {
		return useCase.markConfigApplied(ctx, blueprintSpec)
	}

	err = useCase.markApplyConfigStart(ctx, blueprintSpec)
	if err != nil {
		return useCase.handleFailedApplyRegistryConfig(ctx, blueprintSpec, err)
	}

	// do not apply further configs if error happens, we don't want to corrupt the system more than needed
	// apply normal and sensitive config with this
	err = useCase.applyDoguConfigDiffs(ctx, blueprintSpec.StateDiff.DoguConfigDiffs)
	if err != nil {
		return useCase.handleFailedApplyRegistryConfig(ctx, blueprintSpec, err)
	}
	err = useCase.applyGlobalConfigDiffs(ctx, globalConfigDiffs.GetGlobalConfigDiffsByAction())
	if err != nil {
		return useCase.handleFailedApplyRegistryConfig(ctx, blueprintSpec, fmt.Errorf("could not apply global config: %w", err))
	}

	return useCase.markConfigApplied(ctx, blueprintSpec)
}

func (useCase *EcosystemConfigUseCase) applyGlobalConfigDiffs(ctx context.Context, globalConfigDiffsByAction map[domain.ConfigAction][]domain.GlobalConfigEntryDiff) error {
	var errs []error

	globalConfig, err := useCase.globalConfigRepository.Get(ctx)
	if err != nil {
		return err
	}

	updatedEntries := globalConfig.Config
	entryDiffsToSet := globalConfigDiffsByAction[domain.ConfigActionSet]
	for _, diff := range entryDiffsToSet {
		var err error
		updatedEntries, err = updatedEntries.Set(diff.Key, common.GlobalConfigValue(diff.Expected.Value))
		errs = append(errs, err)
	}

	entryDiffsToRemove := globalConfigDiffsByAction[domain.ConfigActionRemove]
	for _, diff := range entryDiffsToRemove {
		updatedEntries = updatedEntries.Delete(diff.Key)
	}

	if len(entryDiffsToSet) != 0 || len(entryDiffsToRemove) != 0 {
		_, err = useCase.globalConfigRepository.Update(ctx, config.GlobalConfig{Config: updatedEntries})
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (useCase *EcosystemConfigUseCase) applyDoguConfigDiffs(
	ctx context.Context,
	diffsByDogu map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs,
) error {
	var doguConfigDiffs = map[common.SimpleDoguName]domain.DoguConfigDiffs{}
	var sensitiveDoguConfigDiffs = map[common.SimpleDoguName]domain.SensitiveDoguConfigDiffs{}

	for dogu, combinedDiff := range diffsByDogu {
		// only collect doguConfigs with changes, so we don't need to load all.
		if combinedDiff.DoguConfigDiff.HasChangesForDogu(dogu) {
			doguConfigDiffs[dogu] = combinedDiff.DoguConfigDiff
		}
		if combinedDiff.SensitiveDoguConfigDiff.HasChangesForDogu(dogu) {
			sensitiveDoguConfigDiffs[dogu] = combinedDiff.SensitiveDoguConfigDiff
		}
	}

	err := saveDoguConfigs(ctx, useCase.doguConfigRepository, doguConfigDiffs)
	if err != nil {
		return fmt.Errorf("could not apply normal dogu config: %w", err)
	}

	err = saveDoguConfigs(ctx, useCase.sensitiveDoguConfigRepository, sensitiveDoguConfigDiffs)
	if err != nil {
		return fmt.Errorf("could not apply sensitive dogu config: %w", err)
	}

	return nil
}

func saveDoguConfigs(
	ctx context.Context,
	repo doguConfigRepository,
	diffsByDogu map[common.SimpleDoguName]domain.DoguConfigDiffs,
) error {
	// sort to simplify tests
	// this has no real performance impact as we only have a very limited amount of dogus
	dogus := slices.Sorted(maps.Keys(diffsByDogu))
	// has an entry even for not yet existing dogu configs
	configByDogu, err := repo.GetAllExisting(ctx, dogus)
	if err != nil {
		return err
	}

	for dogu, doguDiff := range diffsByDogu {
		newConfig := configByDogu[dogu]
		updatedConfig, err := applyDiff(newConfig, doguDiff)
		if err != nil {
			return err
		}
		_, err = repo.UpdateOrCreate(ctx, config.DoguConfig{
			DoguName: dogu,
			Config:   updatedConfig,
		})
		if err != nil {
			return fmt.Errorf("could not persist config for dogu %s: %w", dogu, err)
		}
	}
	return nil
}

func (useCase *EcosystemConfigUseCase) markApplyConfigStart(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.StartApplyRegistryConfig()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as applying config: %w", err)
	}
	return nil
}

func (useCase *EcosystemConfigUseCase) handleFailedApplyRegistryConfig(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) error {
	logger := log.FromContext(ctx).
		WithName("EcosystemConfigUseCase.handleFailedApplyRegistryConfig").
		WithValues("blueprintId", blueprintSpec.Id)

	blueprintSpec.MarkApplyRegistryConfigFailed(err)
	repoErr := useCase.blueprintRepository.Update(ctx, blueprintSpec)

	if repoErr != nil {
		repoErr = errors.Join(repoErr, err)
		logger.Error(repoErr, "cannot mark blueprint config apply as failed")
		return fmt.Errorf("cannot mark blueprint config apply as failed while handling %q status: %w", blueprintSpec.Status, repoErr)
	}
	return nil
}

func (useCase *EcosystemConfigUseCase) markConfigApplied(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkRegistryConfigApplied()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("failed to mark registry config applied: %w", err)
	}
	return nil
}

// applyDiff merges the given changes from the doguConfigDiff in the DoguConfig.
// Works with normal dogu config and with sensitive config as well.
func applyDiff(doguConfig config.DoguConfig, diffs []domain.DoguConfigEntryDiff) (config.Config, error) {
	updatedEntries := doguConfig.Config

	for _, diff := range diffs {
		var err error
		if diff.NeededAction == domain.ConfigActionSet {
			updatedEntries, err = updatedEntries.Set(diff.Key.Key, config.Value(diff.Expected.Value))
		} else if diff.NeededAction == domain.ConfigActionRemove {
			updatedEntries = updatedEntries.Delete(diff.Key.Key)
		}

		if err != nil {
			return config.Config{}, err
		}
	}
	return updatedEntries, nil
}
