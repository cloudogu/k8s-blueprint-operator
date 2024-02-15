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

type saveDeleteDoguConfigRepository interface {
	Save(context.Context, *ecosystem.DoguConfigEntry) error
	Delete(ctx context.Context, key ecosystem.DoguConfigKey) error
}

type DoguConfigUseCase struct {
	blueprintRepository          blueprintSpecRepository
	doguConfigRepository         doguConfigRepository
	doguSensibleConfigRepository doguSensibleConfigRepository
	globalConfigRepository       domainservice.GlobalConfigKeyRepository
}

func NewDoguConfigUseCase(doguConfigRepository doguConfigRepository, doguSensibleConfigRepository doguSensibleConfigRepository) *DoguConfigUseCase {
	return &DoguConfigUseCase{
		doguConfigRepository:         doguConfigRepository,
		doguSensibleConfigRepository: doguSensibleConfigRepository,
	}
}

func (useCase *DoguConfigUseCase) ApplyDoguConfig(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguConfigUseCase.ApplyDoguConfig").
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
		return useCase.MarkDoguConfigApplied(ctx, blueprintSpec)
	}

	err = useCase.StartApplyDoguConfig(ctx, blueprintSpec)
	if err != nil {
		return err
	}

	var errs []error
	for doguName, doguDiffs := range doguConfigDiffs {
		errs = append(errs, useCase.applyNormalDoguConfig(ctx, doguName, doguDiffs))
		// Values of sensitiveConfig from existing Dogus are already encrypted???
		// TODO Check if Dogus are not installed. Maybe sensitive config hast to be created as secrets.
		errs = append(errs, useCase.applySensibleDoguConfig(ctx, doguName, doguDiffs))
	}

	errs = append(errs, useCase.applyGlobalConfig(ctx, globalConfigDiffs))

	if len(errs) > 0 {
		return useCase.handleFailedApplyDoguConfig(ctx, blueprintSpec, errors.Join(errs...))
	}

	return useCase.MarkDoguConfigApplied(ctx, blueprintSpec)
}

func (useCase *DoguConfigUseCase) applyGlobalConfig(ctx context.Context, diffs domain.GlobalConfigDiff) error {
	var errs []error

	// for _, diff := range diffs {
	// 	switch diff.Action {
	// 	case domain.ConfigActionSet:
	// 		entry := &ecosystem.GlobalConfigEntry{
	// 			Key:                diff.Key,
	// 			Value:              diff.Expected.Value,
	// 			PersistenceContext: nil,
	// 		}
	// 		errs = append(errs, useCase.globalConfigRepository.Save(ctx, entry))
	// 	case domain.ConfigActionRemove:
	// 		errs = append(errs, useCase.globalConfigRepository.Delete(ctx, ecosystem.GlobalConfigKey(diff.Key)))
	// 	case domain.ConfigActionNone:
	// 		continue
	// 	default:
	// 		errs = append(errs, fmt.Errorf("cannot perform unknown action %q for dogu %q with key %q", config.Action, doguName, config.Key))
	// 	}
	// }

	return errors.Join(errs...)
}

func (useCase *DoguConfigUseCase) StartApplyDoguConfig(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.StartApplyDoguConfig()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot mark blueprint as in progress: %w", err)
	}
	return nil
}

func (useCase *DoguConfigUseCase) handleFailedApplyDoguConfig(ctx context.Context, blueprintSpec *domain.BlueprintSpec, err error) error {
	logger := log.FromContext(ctx).
		WithName("DoguConfigUseCase.handleFailedApplyDoguConfig").
		WithValues("blueprintId", blueprintSpec.Id)

	blueprintSpec.MarkApplyDoguConfigFailed(err)
	repoErr := useCase.blueprintRepository.Update(ctx, blueprintSpec)

	if repoErr != nil {
		repoErr = errors.Join(repoErr, err)
		logger.Error(repoErr, "cannot mark blueprint as failed")
		return fmt.Errorf("cannot mark blueprint as failed while handling %q status: %w", blueprintSpec.Status, repoErr)
	}
	return nil
}

func (useCase *DoguConfigUseCase) MarkDoguConfigApplied(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	blueprintSpec.MarkDoguConfigApplied()
	err := useCase.blueprintRepository.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("failed to mark dogu config applied: %w", err)
	}
	return nil
}

func (useCase *DoguConfigUseCase) applyNormalDoguConfig(ctx context.Context, doguName common.SimpleDoguName, doguDiff domain.DoguConfigDiff) error {
	return applyDoguConfigForRepository(ctx, doguName, doguDiff.NormalDoguConfigDiff, useCase.doguConfigRepository)
}

func (useCase *DoguConfigUseCase) applySensibleDoguConfig(ctx context.Context, doguName common.SimpleDoguName, doguDiff domain.DoguConfigDiff) error {
	return applyDoguConfigForRepository(ctx, doguName, doguDiff.SensibleDoguConfigDiff, useCase.doguSensibleConfigRepository)
}

func applyDoguConfigForRepository(ctx context.Context, doguName common.SimpleDoguName, diffs []domain.ConfigKeyDiff, saveDeleteRepository saveDeleteDoguConfigRepository) error {
	var errs []error

	for _, config := range diffs {
		switch config.Action {
		case domain.ConfigActionSet:
			entry := &ecosystem.DoguConfigEntry{
				Key:                ecosystem.DoguConfigKey{DoguName: doguName, Key: config.Key},
				Value:              ecosystem.DoguConfigValue(config.Expected.Value),
				PersistenceContext: nil,
			}
			errs = append(errs, saveDeleteRepository.Save(ctx, entry))
		case domain.ConfigActionRemove:
			errs = append(errs, saveDeleteRepository.Delete(ctx, ecosystem.DoguConfigKey{DoguName: doguName, Key: config.Key}))
		case domain.ConfigActionNone:
			continue
		default:
			errs = append(errs, fmt.Errorf("cannot perform unknown action %q for dogu %q with key %q", config.Action, doguName, config.Key))
		}
	}

	return errors.Join(errs...)
}
