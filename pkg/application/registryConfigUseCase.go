package application

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

//go:embed testdata/doguConfigDiffMockData.yaml
var doguConfigDiffMockDataBytes []byte

//go:embed testdata/globalConfigDiffMockData.yaml
var globalConfigDiffMockDataBytes []byte

type EcosystemRegistryUseCase struct {
	blueprintRepository           blueprintSpecRepository
	doguConfigRepository          doguConfigRepository
	doguSensitiveConfigRepository doguSensitiveConfigRepository
	globalConfigRepository        globalConfigRepository
}

func NewEcosystemRegistryUseCase(blueprintRepository blueprintSpecRepository, doguConfigRepository doguConfigRepository, doguSensitiveConfigRepository doguSensitiveConfigRepository, globalConfigRepository globalConfigRepository) *EcosystemRegistryUseCase {
	return &EcosystemRegistryUseCase{
		blueprintRepository:           blueprintRepository,
		doguConfigRepository:          doguConfigRepository,
		doguSensitiveConfigRepository: doguSensitiveConfigRepository,
		globalConfigRepository:        globalConfigRepository,
	}
}

// ApplyConfig fetches the dogu and global config statediff of the blueprint and applies these keys to the repositories.
func (useCase *EcosystemRegistryUseCase) ApplyConfig(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("EcosystemRegistryUseCase.ApplyConfig").
		WithValues("blueprintId", blueprintId)

	blueprintSpec, err := useCase.blueprintRepository.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint to apply config: %w", err)
	}

	// TODO Remove this before merge.
	// stage := os.Getenv("STAGE")
	// if stage == "development" {
	// 	logger.Info("set config diffs from mock data...")
	// 	data := parseDoguConfigDiffMockData()
	// 	logger.Info("dogu config diffs:")
	// 	logger.Info(fmt.Sprintf("%+v", data))
	// 	blueprintSpec.StateDiff.DoguConfigDiff = data
	// 	logger.Info("dogu config diffs in statediff:")
	// 	logger.Info(fmt.Sprintf("%+v", blueprintSpec.StateDiff.DoguConfigDiff))
	// 	blueprintSpec.StateDiff.GlobalConfigDiff = parseGlobalConfigDiffMockData()
	// }

	doguConfigDiffs := blueprintSpec.StateDiff.DoguConfigDiff
	isEmptyDoguDiff := len(doguConfigDiffs) == 0
	if isEmptyDoguDiff {
		logger.Info("dogu config diffs are empty...")
	}

	globalConfigDiffs := blueprintSpec.StateDiff.GlobalConfigDiff
	isEmptyGlobalDiff := len(globalConfigDiffs) == 0
	if isEmptyGlobalDiff {
		logger.Info("global config diffs are empty...")
	}

	if isEmptyDoguDiff && isEmptyGlobalDiff {
		// TODO Correct Status or create new for no Action needed?
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

func parseGlobalConfigDiffMockData() domain.GlobalConfigDiff {
	object := &domain.GlobalConfigDiff{}
	err := yaml.Unmarshal(globalConfigDiffMockDataBytes, object)
	if err != nil {
		panic(fmt.Errorf("error during mock data deserialization"))
	}

	return *object
}

func parseDoguConfigDiffMockData() map[common.SimpleDoguName]domain.CombinedDoguConfigDiff {
	object := &map[common.SimpleDoguName]domain.CombinedDoguConfigDiff{}
	err := yaml.Unmarshal(doguConfigDiffMockDataBytes, object)
	if err != nil {
		panic(fmt.Errorf("error during mock data deserialization"))
	}

	return *object
}

func (useCase *EcosystemRegistryUseCase) applyGlobalConfigDiffs(ctx context.Context, diffs domain.GlobalConfigDiff) error {
	var errs []error

	for _, diff := range diffs {
		switch diff.Action {
		case domain.ConfigActionSet:
			entry := &ecosystem.GlobalConfigEntry{
				Key:   diff.Key,
				Value: common.GlobalConfigValue(diff.Expected.Value),
			}
			errs = append(errs, useCase.globalConfigRepository.Save(ctx, entry))
		case domain.ConfigActionRemove:
			errs = append(errs, useCase.globalConfigRepository.Delete(ctx, diff.Key))
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, fmt.Errorf("cannot perform unknown action %q for global config with key %q", diff.Action, diff.Key))
		}
	}

	return errors.Join(errs...)
}

func (useCase *EcosystemRegistryUseCase) applyDoguConfigDiffs(ctx context.Context, doguName common.SimpleDoguName, diffs domain.DoguConfigDiff) error {
	var errs []error

	for _, diff := range diffs {
		switch diff.Action {
		case domain.ConfigActionSet:
			entry := getDoguConfigEntry(doguName, diff.Key.Key, diff.Expected.Value)
			errs = append(errs, useCase.doguConfigRepository.Save(ctx, entry))
		case domain.ConfigActionRemove:
			errs = append(errs, useCase.doguConfigRepository.Delete(ctx, common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}))
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, doguUnknownConfigActionError(diff.Action, diff.Key.Key, doguName))
		}
	}

	return errors.Join(errs...)
}

// Values of sensitiveConfig from existing Dogus are already encrypted???
func (useCase *EcosystemRegistryUseCase) applySensitiveDoguConfigDiffs(ctx context.Context, doguName common.SimpleDoguName, diffs domain.SensitiveDoguConfigDiff) error {
	var errs []error

	for _, diff := range diffs {
		switch diff.Action {
		case domain.ConfigActionSet:
			entry := &ecosystem.SensitiveDoguConfigEntry{
				Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}},
				Value: common.EncryptedDoguConfigValue(diff.Expected.Value),
			}
			errs = append(errs, useCase.doguSensitiveConfigRepository.Save(ctx, entry))
		case domain.ConfigActionSetToEncrypt:
			entry := getDoguConfigEntry(doguName, diff.Key.Key, diff.Expected.Value)
			errs = append(errs, useCase.doguSensitiveConfigRepository.SaveForNotInstalledDogu(ctx, entry))
		case domain.ConfigActionRemove:
			errs = append(errs, useCase.doguSensitiveConfigRepository.Delete(ctx, common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}}))
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, doguUnknownConfigActionError(diff.Action, diff.Key.Key, doguName))
		}
	}

	return errors.Join(errs...)
}

func getDoguConfigEntry(doguName common.SimpleDoguName, key, value string) *ecosystem.DoguConfigEntry {
	return &ecosystem.DoguConfigEntry{
		Key:   common.DoguConfigKey{DoguName: doguName, Key: key},
		Value: common.DoguConfigValue(value),
	}
}

func doguUnknownConfigActionError(action domain.ConfigAction, key string, doguName common.SimpleDoguName) error {
	return fmt.Errorf("cannot perform unknown action %q for dogu %q with key %q", action, doguName, key)
}

func (useCase *EcosystemRegistryUseCase) markApplyConfigStart(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.StartApplyRegistryConfig()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as applying config: %w", err)
	}
	return nil
}

func (useCase *EcosystemRegistryUseCase) handleFailedApplyRegistryConfig(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) error {
	logger := log.FromContext(ctx).
		WithName("EcosystemRegistryUseCase.handleFailedApplyRegistryConfig").
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

func (useCase *EcosystemRegistryUseCase) markConfigApplied(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkRegistryConfigApplied()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("failed to mark registry config applied: %w", err)
	}
	return nil
}
