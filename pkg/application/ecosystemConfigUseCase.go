package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EcosystemConfigUseCase struct {
	blueprintRepository           blueprintSpecRepository
	doguConfigRepository          doguConfigEntryRepository
	doguSensitiveConfigRepository sensitiveDoguConfigEntryRepository
	globalConfigRepository        globalConfigEntryRepository
	encryptionAdapter             configEncryptionAdapter
}

var errSensitiveDoguConfigEntry = fmt.Errorf("sensitive dogu config error")

func NewEcosystemConfigUseCase(blueprintRepository blueprintSpecRepository, doguConfigRepository doguConfigEntryRepository, doguSensitiveConfigRepository sensitiveDoguConfigEntryRepository, globalConfigRepository globalConfigEntryRepository, encryptionAdapter configEncryptionAdapter) *EcosystemConfigUseCase {
	return &EcosystemConfigUseCase{
		blueprintRepository:           blueprintRepository,
		doguConfigRepository:          doguConfigRepository,
		doguSensitiveConfigRepository: doguSensitiveConfigRepository,
		globalConfigRepository:        globalConfigRepository,
		encryptionAdapter:             encryptionAdapter,
	}
}

// ApplyConfig fetches the dogu and global config statediff of the blueprint and applies these keys to the repositories.
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

	entryDiffsToSet := globalConfigDiffsByAction[domain.ConfigActionSet]
	var entriesToSet = make([]*ecosystem.GlobalConfigEntry, 0, len(entryDiffsToSet))
	for _, diff := range entryDiffsToSet {
		entry := &ecosystem.GlobalConfigEntry{
			Key:   diff.Key,
			Value: common.GlobalConfigValue(diff.Expected.Value),
		}
		entriesToSet = append(entriesToSet, entry)
	}

	entryDiffsToRemove := globalConfigDiffsByAction[domain.ConfigActionRemove]
	var keysToDelete = make([]common.GlobalConfigKey, 0, len(entryDiffsToRemove))
	for _, diff := range entryDiffsToRemove {
		keysToDelete = append(keysToDelete, diff.Key)
	}

	errs = append(errs, callIfNotEmpty(ctx, entriesToSet, useCase.globalConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.globalConfigRepository.DeleteAllByKeys))

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

	var encryptedEntriesToSet []*ecosystem.SensitiveDoguConfigEntry
	var entriesToEncrypt []*ecosystem.SensitiveDoguConfigEntry
	var keysToDelete []common.SensitiveDoguConfigKey

	encryptedEntryValues, err := useCase.encryptSensitiveDoguDiffs(ctx, sensitiveDoguConfigDiffsByAction)
	if err != nil {
		errs = append(errs, err)
	}

	for _, diff := range sensitiveDoguConfigDiffsByAction[domain.ConfigActionSetEncrypted] {
		entry, createEncryptedEntryErr := getSensitiveDoguConfigEntryWithEncryption(diff.Key.DoguName, diff, encryptedEntryValues)
		if createEncryptedEntryErr != nil {
			errs = append(errs, createEncryptedEntryErr)
			continue
		}
		entriesToEncrypt = append(entriesToEncrypt, entry)
	}

	for _, diff := range sensitiveDoguConfigDiffsByAction[domain.ConfigActionSetToEncrypt] {
		entry := getSensitiveDoguConfigEntry(diff.Key.DoguName, diff)
		encryptedEntriesToSet = append(encryptedEntriesToSet, entry)
	}

	for _, diff := range sensitiveDoguConfigDiffsByAction[domain.ConfigActionRemove] {
		keysToDelete = append(keysToDelete, common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: diff.Key.DoguName, Key: diff.Key.Key}})
	}

	errs = append(errs, callIfNotEmpty(ctx, entriesToEncrypt, useCase.doguSensitiveConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, encryptedEntriesToSet, useCase.doguSensitiveConfigRepository.SaveAllForNotInstalledDogus))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.doguSensitiveConfigRepository.DeleteAllByKeys))

	return errors.Join(errs...)
}

// Only encrypt diffs with action domain.ConfigActionSetEncrypted. Diffs with action domain.ConfigActionSetToEncrypt will
// be encrypted by other components in further procedure.
func (useCase *EcosystemConfigUseCase) encryptSensitiveDoguDiffs(ctx context.Context, sensitiveDoguConfigDiffsByAction map[domain.ConfigAction]domain.SensitiveDoguConfigDiffs) (map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, error) {
	valuesToEncrypt := map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}

	for _, diff := range sensitiveDoguConfigDiffsByAction[domain.ConfigActionSetEncrypted] {
		valuesToEncrypt[diff.Key] = common.SensitiveDoguConfigValue(diff.Expected.Value)
	}

	if len(valuesToEncrypt) > 0 {
		return useCase.encryptionAdapter.EncryptAll(ctx, valuesToEncrypt)
	}

	return map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}, nil
}

func callIfNotEmpty[T ecosystem.RegistryConfigEntry | common.RegistryConfigKey](ctx context.Context, collection []T, fn func(context.Context, []T) error) error {
	if len(collection) > 0 {
		return fn(ctx, collection)
	}

	return nil
}

func getSensitiveDoguConfigEntryWithEncryption(doguName common.SimpleDoguName, diff domain.SensitiveDoguConfigEntryDiff, encryptedEntryValues map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue) (*ecosystem.SensitiveDoguConfigEntry, error) {
	entry := getSensitiveDoguConfigEntry(doguName, diff)
	if encryptedEntryValues == nil {
		return nil, domainservice.NewInternalError(errSensitiveDoguConfigEntry, "encrypted entry value map is nil")
	}
	value, ok := encryptedEntryValues[entry.Key]
	if !ok {
		return nil, domainservice.NewNotFoundError(errSensitiveDoguConfigEntry, "did not find encrypted value for key %s", entry.Key.Key)
	}
	entry.Value = value

	return entry, nil
}

func getSensitiveDoguConfigEntry(doguName common.SimpleDoguName, diff domain.SensitiveDoguConfigEntryDiff) *ecosystem.SensitiveDoguConfigEntry {
	return &ecosystem.SensitiveDoguConfigEntry{
		Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}},
		Value: common.EncryptedDoguConfigValue(diff.Expected.Value),
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
