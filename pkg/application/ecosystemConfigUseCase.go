package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-registry-lib/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EcosystemConfigUseCase struct {
	blueprintRepository           blueprintSpecRepository
	doguConfigRepository          doguConfigEntryRepository
	doguSensitiveConfigRepository sensitiveDoguConfigEntryRepository
	globalConfigEntryRepository   globalConfigEntryRepository
	globalConfigRepository        globalConfigRepository
}

var errSensitiveDoguConfigEntry = fmt.Errorf("sensitive dogu config error")

func NewEcosystemConfigUseCase(
	blueprintRepository blueprintSpecRepository,
	doguConfigRepository doguConfigEntryRepository,
	doguSensitiveConfigRepository sensitiveDoguConfigEntryRepository,
	globalConfigEntryRepository globalConfigEntryRepository,
	globalConfigRepository globalConfigRepository,
) *EcosystemConfigUseCase {
	return &EcosystemConfigUseCase{
		blueprintRepository:           blueprintRepository,
		doguConfigRepository:          doguConfigRepository,
		doguSensitiveConfigRepository: doguSensitiveConfigRepository,
		globalConfigEntryRepository:   globalConfigEntryRepository,
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

	var errs []error
	errs = append(errs, useCase.applyDoguConfigDiffs(ctx, blueprintSpec.StateDiff.GetDoguConfigDiffsByAction()))
	errs = append(errs, useCase.applySensitiveDoguConfigDiffs(ctx, blueprintSpec.StateDiff.GetSensitiveDoguConfigDiffsByAction()))
	errs = append(errs, useCase.applyGlobalConfigDiffs(ctx, globalConfigDiffs.GetGlobalConfigDiffsByAction()))

	joinedErr := errors.Join(errs...)
	if joinedErr != nil {
		return useCase.handleFailedApplyRegistryConfig(ctx, blueprintSpec, joinedErr)
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

func (useCase *EcosystemConfigUseCase) applyDoguConfigDiffs(ctx context.Context, doguConfigDiffsByAction map[domain.ConfigAction]domain.DoguConfigDiffs) error {
	var errs []error
	var entriesToSet []*ecosystem.DoguConfigEntry
	var keysToDelete []common.DoguConfigKey

	for _, diff := range doguConfigDiffsByAction[domain.ConfigActionSet] {
		entry := &ecosystem.DoguConfigEntry{
			Key:   common.DoguConfigKey{DoguName: diff.Key.DoguName, Key: diff.Key.Key},
			Value: common.DoguConfigValue(diff.Expected.Value),
		}
		entriesToSet = append(entriesToSet, entry)
	}

	for _, diff := range doguConfigDiffsByAction[domain.ConfigActionRemove] {
		keysToDelete = append(keysToDelete, common.DoguConfigKey{DoguName: diff.Key.DoguName, Key: diff.Key.Key})
	}

	errs = append(errs, callIfNotEmpty(ctx, entriesToSet, useCase.doguConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.doguConfigRepository.DeleteAllByKeys))

	return errors.Join(errs...)
}

func (useCase *EcosystemConfigUseCase) applySensitiveDoguConfigDiffs(ctx context.Context, sensitiveDoguConfigDiffsByAction map[domain.ConfigAction]domain.SensitiveDoguConfigDiffs) error {
	var errs []error

	var keysToSet []*ecosystem.SensitiveDoguConfigEntry
	var keysToDelete []common.SensitiveDoguConfigKey

	for _, diff := range sensitiveDoguConfigDiffsByAction[domain.ConfigActionSet] {
		entry := getSensitiveDoguConfigEntry(diff.Key.DoguName, diff)
		keysToSet = append(keysToSet, entry)
	}

	for _, diff := range sensitiveDoguConfigDiffsByAction[domain.ConfigActionRemove] {
		keysToDelete = append(keysToDelete, common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: diff.Key.DoguName, Key: diff.Key.Key}})
	}

	errs = append(errs, callIfNotEmpty(ctx, keysToSet, useCase.doguSensitiveConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.doguSensitiveConfigRepository.DeleteAllByKeys))

	return errors.Join(errs...)
}

func callIfNotEmpty[T ecosystem.RegistryConfigEntry | common.RegistryConfigKey](ctx context.Context, collection []T, fn func(context.Context, []T) error) error {
	if len(collection) > 0 {
		return fn(ctx, collection)
	}

	return nil
}

func getSensitiveDoguConfigEntry(doguName common.SimpleDoguName, diff domain.SensitiveDoguConfigEntryDiff) *ecosystem.SensitiveDoguConfigEntry {
	return &ecosystem.SensitiveDoguConfigEntry{
		Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}},
		Value: common.SensitiveDoguConfigValue(diff.Expected.Value),
	}
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
