package application

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EcosystemConfigUseCase struct {
	blueprintRepository           blueprintSpecRepository
	doguConfigRepository          doguConfigRepository
	sensitiveDoguConfigRepository sensitiveDoguConfigRepository
	globalConfigRepository        globalConfigRepository
	doguInstallationRepository    doguInstallationRepository
}

func NewEcosystemConfigUseCase(blueprintRepository blueprintSpecRepository, doguConfigRepository doguConfigRepository, sensitiveDoguConfigRepository sensitiveDoguConfigRepository, globalConfigRepository globalConfigRepository, doguInstallationRepository domainservice.DoguInstallationRepository) *EcosystemConfigUseCase {
	return &EcosystemConfigUseCase{
		blueprintRepository:           blueprintRepository,
		doguConfigRepository:          doguConfigRepository,
		sensitiveDoguConfigRepository: sensitiveDoguConfigRepository,
		globalConfigRepository:        globalConfigRepository,
		doguInstallationRepository:    doguInstallationRepository,
	}
}

// ApplyConfig fetches the dogu and global config stateDiff of the blueprint and applies these keys to the repositories.
func (useCase *EcosystemConfigUseCase) ApplyConfig(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("EcosystemConfigUseCase.ApplyConfig")

	err := pauseReconciliationForDogus(ctx, useCase.doguInstallationRepository, blueprint.StateDiff)
	if err != nil {
		return useCase.handleFailedApplyEcosystemConfig(ctx, blueprint, fmt.Errorf("could not pause reconciliation for some dogus: %w", err))
	}
	err = applyDoguConfigDiffs(ctx, useCase.doguConfigRepository, blueprint.StateDiff.DoguConfigDiffs)
	if err != nil {
		return useCase.handleFailedApplyEcosystemConfig(ctx, blueprint, fmt.Errorf("could not apply normal dogu config: %w", err))
	}
	err = applyDoguConfigDiffs(ctx, useCase.sensitiveDoguConfigRepository, blueprint.StateDiff.SensitiveDoguConfigDiffs)
	if err != nil {
		return useCase.handleFailedApplyEcosystemConfig(ctx, blueprint, fmt.Errorf("could not apply sensitive dogu config: %w", err))
	}
	err = useCase.applyGlobalConfigDiffs(ctx, blueprint.StateDiff.GlobalConfigDiffs.GetGlobalConfigDiffsByAction())
	if err != nil {
		return useCase.handleFailedApplyEcosystemConfig(ctx, blueprint, fmt.Errorf("could not apply global config: %w", err))
	}

	blueprint.Events = append(blueprint.Events, domain.EcosystemConfigAppliedEvent{})
	repoErr := useCase.blueprintRepository.Update(ctx, blueprint)

	if repoErr != nil {
		repoErr = errors.Join(repoErr, err)
		logger.Error(repoErr, "cannot update blueprint events")
		return fmt.Errorf("cannot update blueprint events: %w", repoErr)
	}
	return nil
}

func pauseReconciliationForDogus(ctx context.Context, repository doguInstallationRepository, diff domain.StateDiff) error {
	allDogus, err := repository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("error while attempting to load dogus: %w", err)
	}
	globalConfigChanges := diff.GlobalConfigDiffs.HasChanges()
	for _, dogu := range allDogus {
		for _, doguDiff := range diff.DoguDiffs {
			if doguDiff.DoguName != dogu.Name.SimpleName {
				continue
			}
			if slices.Contains(doguDiff.NeededActions, domain.ActionUpgrade) &&
				(globalConfigChanges ||
					diff.DoguConfigDiffs[dogu.Name.SimpleName].HasChanges() ||
					diff.SensitiveDoguConfigDiffs[dogu.Name.SimpleName].HasChanges()) {
				dogu.PauseReconciliation = true
				err = repository.Update(ctx, dogu)
				if err != nil {
					return fmt.Errorf("could not pause reconciliation for dogu: %w", err)
				}
			}
		}
	}
	return nil
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
		val := ""
		if diff.Expected.Value != nil {
			val = *diff.Expected.Value
		}
		updatedEntries, err = updatedEntries.Set(diff.Key, common.GlobalConfigValue(val))
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

func applyDoguConfigDiffs(
	ctx context.Context,
	repo doguConfigRepository,
	diffsByDogu map[cescommons.SimpleName]domain.DoguConfigDiffs,
) error {
	var doguConfigDiffs = map[cescommons.SimpleName]domain.DoguConfigDiffs{}

	for dogu, entryDiffs := range diffsByDogu {
		// only collect doguConfigs with changes, so we don't need to load all.
		if entryDiffs.HasChanges() {
			doguConfigDiffs[dogu] = entryDiffs
		}
	}

	return saveDoguConfigs(ctx, repo, doguConfigDiffs)
}

func saveDoguConfigs(
	ctx context.Context,
	repo doguConfigRepository,
	diffsByDogu map[cescommons.SimpleName]domain.DoguConfigDiffs,
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

func (useCase *EcosystemConfigUseCase) handleFailedApplyEcosystemConfig(ctx context.Context, blueprint *domain.BlueprintSpec, err error) error {
	logger := log.FromContext(ctx).
		WithName("EcosystemConfigUseCase.handleFailedApplyEcosystemConfig").
		WithValues("blueprintId", blueprint.Id)

	// sets condition
	changed := blueprint.SetLastApplySucceededConditionOnError(domain.ReasonLastApplyErrorAtConfig, err)
	if changed {
		repoErr := useCase.blueprintRepository.Update(ctx, blueprint)

		if repoErr != nil {
			repoErr = errors.Join(repoErr, err)
			logger.Error(repoErr, "cannot mark blueprint config apply as failed")
			return fmt.Errorf("cannot mark blueprint config apply as failed: %w", repoErr)
		}
	}
	return err
}

// applyDiff merges the given changes from the doguConfigDiff in the DoguConfig.
// Works with normal dogu config and with sensitive config as well.
func applyDiff(doguConfig config.DoguConfig, diffs []domain.DoguConfigEntryDiff) (config.Config, error) {
	updatedEntries := doguConfig.Config

	for _, diff := range diffs {
		var err error
		switch diff.NeededAction {
		case domain.ConfigActionSet:
			val := ""
			if diff.Expected.Value != nil {
				val = *diff.Expected.Value
			}
			updatedEntries, err = updatedEntries.Set(diff.Key.Key, config.Value(val))
		case domain.ConfigActionRemove:
			updatedEntries = updatedEntries.Delete(diff.Key.Key)
		}

		if err != nil {
			return config.Config{}, err
		}
	}
	return updatedEntries, nil
}
