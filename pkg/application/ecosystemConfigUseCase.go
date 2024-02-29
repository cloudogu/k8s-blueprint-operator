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
	for doguName, doguDiffs := range doguConfigDiffs {
		errs = append(errs, useCase.applyDoguConfigDiffs(ctx, doguName, doguDiffs.DoguConfigDiff))
		errs = append(errs, useCase.applySensitiveDoguConfigDiffs(ctx, doguName, doguDiffs.SensitiveDoguConfigDiff))
	}

	errs = append(errs, useCase.applyGlobalConfigDiffs(ctx, globalConfigDiffs))

	joinedErr := errors.Join(errs...)
	if joinedErr != nil {
		return useCase.handleFailedApplyRegistryConfig(ctx, blueprintSpec, joinedErr)
	}

	return useCase.markConfigApplied(ctx, blueprintSpec)
}

func (useCase *EcosystemConfigUseCase) applyGlobalConfigDiffs(ctx context.Context, diffs domain.GlobalConfigDiffs) error {
	var errs []error
	var entriesToSet []*ecosystem.GlobalConfigEntry
	var keysToDelete []common.GlobalConfigKey

	for _, diff := range diffs {
		switch diff.NeededAction {
		case domain.ConfigActionSet:
			entry := &ecosystem.GlobalConfigEntry{
				Key:   diff.Key,
				Value: common.GlobalConfigValue(diff.Expected.Value),
			}
			entriesToSet = append(entriesToSet, entry)
		case domain.ConfigActionRemove:
			keysToDelete = append(keysToDelete, diff.Key)
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, fmt.Errorf("cannot perform unknown action %q for global config with key %q", diff.NeededAction, diff.Key))
		}
	}

	errs = append(errs, callIfNotEmpty(ctx, entriesToSet, useCase.globalConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.globalConfigRepository.DeleteAllByKeys))

	return errors.Join(errs...)
}

func (useCase *EcosystemConfigUseCase) applyDoguConfigDiffs(ctx context.Context, doguName common.SimpleDoguName, diffs domain.DoguConfigDiffs) error {
	var errs []error
	var entriesToSet []*ecosystem.DoguConfigEntry
	var keysToDelete []common.DoguConfigKey

	for _, diff := range diffs {
		switch diff.NeededAction {
		case domain.ConfigActionSet:
			entry := &ecosystem.DoguConfigEntry{
				Key:   common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key},
				Value: common.DoguConfigValue(diff.Expected.Value),
			}
			entriesToSet = append(entriesToSet, entry)
		case domain.ConfigActionRemove:
			keysToDelete = append(keysToDelete, common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key})
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, doguUnknownConfigActionError(diff.NeededAction, diff.Key.Key, doguName))
		}
	}

	errs = append(errs, callIfNotEmpty(ctx, entriesToSet, useCase.doguConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.doguConfigRepository.DeleteAllByKeys))

	return errors.Join(errs...)
}

func (useCase *EcosystemConfigUseCase) applySensitiveDoguConfigDiffs(ctx context.Context, doguName common.SimpleDoguName, diffs domain.SensitiveDoguConfigDiffs) error {
	var errs []error

	var encryptedEntriesToSet []*ecosystem.SensitiveDoguConfigEntry
	var entriesToEncrypt []*ecosystem.SensitiveDoguConfigEntry
	var keysToDelete []common.SensitiveDoguConfigKey

	encryptedEntryValues, err := useCase.encryptSensitiveDoguDiffs(ctx, diffs)
	if err != nil {
		return err
	}

	for _, diff := range diffs {
		switch diff.NeededAction {
		case domain.ConfigActionSetEncrypted:
			entry, createEncryptedEntryErr := getSensitiveDoguConfigEntryWithEncryption(doguName, diff, encryptedEntryValues)
			if createEncryptedEntryErr != nil {
				errs = append(errs, createEncryptedEntryErr)
				continue
			}
			entriesToEncrypt = append(entriesToEncrypt, entry)
		case domain.ConfigActionSetToEncrypt:
			entry := getSensitiveDoguConfigEntry(doguName, diff)
			encryptedEntriesToSet = append(encryptedEntriesToSet, entry)
		case domain.ConfigActionRemove:
			keysToDelete = append(keysToDelete, common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}})
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, doguUnknownConfigActionError(diff.NeededAction, diff.Key.Key, doguName))
		}
	}

	errs = append(errs, callIfNotEmpty(ctx, entriesToEncrypt, useCase.doguSensitiveConfigRepository.SaveAll))
	errs = append(errs, callIfNotEmpty(ctx, encryptedEntriesToSet, useCase.doguSensitiveConfigRepository.SaveAllForNotInstalledDogus))
	errs = append(errs, callIfNotEmpty(ctx, keysToDelete, useCase.doguSensitiveConfigRepository.DeleteAllByKeys))

	return errors.Join(errs...)
}

// Only encrypt diffs with action domain.ConfigActionSetEncrypted. Diffs with action domain.ConfigActionSetToEncrypt will
// be encrypted by other components in further procedure.
func (useCase *EcosystemConfigUseCase) encryptSensitiveDoguDiffs(ctx context.Context, diffs domain.SensitiveDoguConfigDiffs) (map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue, error) {
	toEncryptEntries := map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{}
	for _, diff := range diffs {
		if diff.NeededAction == domain.ConfigActionSetEncrypted {
			toEncryptEntries[diff.Key] = common.SensitiveDoguConfigValue(diff.Expected.Value)
		}
	}

	if len(toEncryptEntries) > 0 {
		return useCase.encryptionAdapter.EncryptAll(ctx, toEncryptEntries)
	}

	return nil, nil
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

func doguUnknownConfigActionError(action domain.ConfigAction, key string, doguName common.SimpleDoguName) error {
	return fmt.Errorf("cannot perform unknown action %q for dogu %q with key %q", action, doguName, key)
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
