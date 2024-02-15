package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

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

func (useCase *EcosystemRegistryUseCase) ApplyConfig(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("EcosystemRegistryUseCase.ApplyConfig").
		WithValues("blueprintId", blueprintId)

	blueprintSpec, err := useCase.blueprintRepository.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint to apply dogu config: %w", err)
	}

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
		return useCase.MarkConfigApplied(ctx, blueprintSpec)
	}

	err = useCase.StartApplyConfig(ctx, blueprintSpec)
	if err != nil {
		return err
	}

	var errs []error
	for doguName, doguDiffs := range doguConfigDiffs {
		errs = append(errs, useCase.applyDoguConfigDiffs(ctx, doguName, doguDiffs.DoguConfigDiff))
		errs = append(errs, useCase.applySensitiveDoguConfigDiffs(ctx, doguName, doguDiffs.SensitiveDoguConfigDiff))
	}

	errs = append(errs, useCase.applyGlobalConfigDiffs(ctx, globalConfigDiffs))

	if len(errs) > 0 {
		return useCase.handleFailedApplyRegistryConfig(ctx, blueprintSpec, errors.Join(errs...))
	}

	return useCase.MarkConfigApplied(ctx, blueprintSpec)
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
			entry := &ecosystem.DoguConfigEntry{
				Key:   common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key},
				Value: common.DoguConfigValue(diff.Expected.Value),
			}
			errs = append(errs, useCase.doguConfigRepository.Save(ctx, entry))
		case domain.ConfigActionRemove:
			errs = append(errs, useCase.doguConfigRepository.Delete(ctx, common.DoguConfigKey{DoguName: doguName, Key: diff.Key.Key}))
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, fmt.Errorf("cannot perform unknown action %q for dogu %q with key %q", diff.Action, doguName, diff.Key))
		}
	}

	return errors.Join(errs...)
}

// Values of sensitiveConfig from existing Dogus are already encrypted???
// TODO Check if Dogus are not installed. Maybe sensitive config has to be created as secrets.
func (useCase *EcosystemRegistryUseCase) applySensitiveDoguConfigDiffs(ctx context.Context, doguName common.SimpleDoguName, diffs domain.SensitiveDoguConfigDiff) error {
	var errs []error

	for _, diff := range diffs {
		switch diff.Action {
		case domain.ConfigActionSet:
			entry := &ecosystem.SensitiveDoguConfigEntry{
				Key:   common.SensitiveDoguConfigKey{DoguName: doguName, Key: diff.Key.Key},
				Value: common.EncryptedDoguConfigValue(diff.Expected.Value),
			}
			errs = append(errs, useCase.doguSensitiveConfigRepository.Save(ctx, entry))
		case domain.ConfigActionRemove:
			errs = append(errs, useCase.doguSensitiveConfigRepository.Delete(ctx, common.SensitiveDoguConfigKey{DoguName: doguName, Key: diff.Key.Key}))
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, fmt.Errorf("cannot perform unknown action %q for dogu %q with key %q", diff.Action, doguName, diff.Key))
		}
	}

	return errors.Join(errs...)
}

func (useCase *EcosystemRegistryUseCase) StartApplyConfig(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.StartApplyRegistryConfig()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as in progress: %w", err)
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
		logger.Error(repoErr, "cannot mark blueprint as failed")
		return fmt.Errorf("cannot mark blueprint as failed while handling %q status: %w", blueprintSpec.Status, repoErr)
	}
	return nil
}

func (useCase *EcosystemRegistryUseCase) MarkConfigApplied(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkRegistryConfigApplied()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("failed to mark registry config applied: %w", err)
	}
	return nil
}
