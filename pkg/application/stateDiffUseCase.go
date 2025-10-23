package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const REFERENCED_CONFIG_NOT_FOUND = "could not load referenced sensitive config"

type StateDiffUseCase struct {
	blueprintSpecRepo        blueprintSpecRepository
	doguInstallationRepo     doguInstallationRepository
	globalConfigRepo         globalConfigRepository
	doguConfigRepo           doguConfigRepository
	sensitiveDoguConfigRepo  sensitiveDoguConfigRepository
	sensitiveConfigRefReader sensitiveConfigRefReader
	debugModeRepo            debugModeRepository
}

func NewStateDiffUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguInstallationRepo domainservice.DoguInstallationRepository,
	globalConfigRepo domainservice.GlobalConfigRepository,
	doguConfigRepo domainservice.DoguConfigRepository,
	sensitiveDoguConfigRepo domainservice.SensitiveDoguConfigRepository,
	sensitiveConfigRefReader domainservice.SensitiveConfigRefReader,
	debugModeRepo domainservice.DebugModeRepository,
) *StateDiffUseCase {
	return &StateDiffUseCase{
		blueprintSpecRepo:        blueprintSpecRepo,
		doguInstallationRepo:     doguInstallationRepo,
		globalConfigRepo:         globalConfigRepo,
		doguConfigRepo:           doguConfigRepo,
		sensitiveDoguConfigRepo:  sensitiveDoguConfigRepo,
		sensitiveConfigRefReader: sensitiveConfigRefReader,
		debugModeRepo:            debugModeRepo,
	}
}

// DetermineStateDiff loads the state of the ecosystem and compares it to the blueprint. It creates a declarative diff.
// returns:
//   - a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or while collecting the ecosystem state or
//   - a domainservice.ConflictError if there was a concurrent write to the blueprint or
//   - a domain.InvalidBlueprintError if there are any forbidden actions in the stateDiff.
//   - any error if there is any other error.
func (useCase *StateDiffUseCase) DetermineStateDiff(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("StateDiffUseCase.DetermineStateDiff")

	logger.V(2).Info("load referenced sensitive config")
	// load referenced config before collecting ecosystem state
	// if an error happens here, we save a lot of heavy work
	referencedSensitiveConfig, err := useCase.sensitiveConfigRefReader.GetValues(
		ctx, blueprint.EffectiveBlueprint.Config.GetSensitiveConfigReferences(),
	)
	if err != nil {
		err = fmt.Errorf("%s: %w", REFERENCED_CONFIG_NOT_FOUND, err)
		blueprint.MissingConfigReferences(err)
		updateError := useCase.blueprintSpecRepo.Update(ctx, blueprint)
		if updateError != nil {
			return errors.Join(updateError, err)
		}
		return err
	}

	logger.V(2).Info("collect ecosystem state for state diff")
	ecosystemState, err := useCase.collectEcosystemState(ctx, blueprint.EffectiveBlueprint)
	if err != nil {
		return fmt.Errorf("could not determine state diff: %w", err)
	}

	logger.V(2).Info("determine state diff to the cloudogu ecosystem")
	isDebugModeActive, err := useCase.determineDebugModeState(ctx, logger)
	if err != nil {
		return err
	}
	stateDiffError := blueprint.DetermineStateDiff(ecosystemState, referencedSensitiveConfig, isDebugModeActive)
	var invalidError *domain.InvalidBlueprintError
	if errors.As(stateDiffError, &invalidError) {
		// do not return here as with this error the blueprint status and events should be persisted as normal.
	} else if stateDiffError != nil {
		return fmt.Errorf("failed to determine state diff: %w", stateDiffError)
	}

	err = useCase.blueprintSpecRepo.Update(ctx, blueprint)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after determining the state diff to the ecosystem: %w", blueprint.Id, err)
	}

	// return this error back here to persist the blueprint status and events first.
	// return it to signal that a repeated call to this function will not result in any progress.
	return stateDiffError
}

func (useCase *StateDiffUseCase) determineDebugModeState(ctx context.Context, logger logr.Logger) (bool, error) {
	debugMode, err := useCase.debugModeRepo.GetSingleton(ctx)
	if err != nil {
		// ignore not found error, no debug mode cr means we are not in debug mode
		if !domainservice.IsNotFoundError(err) {
			return false, fmt.Errorf("cannot calculate effective blueprint due to an error when loading the debug mode cr: %w", err)
		}
	}
	isDebugModeActive := debugMode != nil && debugMode.IsActive()
	if isDebugModeActive {
		logger.Info("debug mode is active, will ignore loglevel changes until deactivated")
	}
	return isDebugModeActive, nil
}

func (useCase *StateDiffUseCase) collectEcosystemState(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) (ecosystem.EcosystemState, error) {
	logger := log.FromContext(ctx).WithName("StateDiffUseCase.collectEcosystemState")

	// TODO: collect ecosystem state in parallel (like for ecosystem health) if we have time
	// load current dogus
	logger.V(2).Info("collect installed dogus")
	installedDogus, doguErr := useCase.doguInstallationRepo.GetAll(ctx)
	// load current config
	logger.V(2).Info("collect needed global config")
	globalConfig, globalConfigErr := useCase.globalConfigRepo.Get(ctx)

	logger.V(2).Info("collect needed dogu config")
	configByDogu, doguConfigErr := useCase.doguConfigRepo.GetAllExisting(ctx, effectiveBlueprint.Config.GetDogusWithChangedConfig())

	logger.V(2).Info("collect needed sensitive dogu config")
	sensitiveConfigByDogu, sensitiveConfigErr := useCase.sensitiveDoguConfigRepo.GetAllExisting(ctx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig())

	joinedError := errors.Join(doguErr, globalConfigErr, doguConfigErr, sensitiveConfigErr)
	if joinedError != nil {
		return ecosystem.EcosystemState{}, fmt.Errorf("could not collect ecosystem state: %w", joinedError)
	}

	return ecosystem.EcosystemState{
		InstalledDogus:        installedDogus,
		GlobalConfig:          globalConfig,
		ConfigByDogu:          configByDogu,
		SensitiveConfigByDogu: sensitiveConfigByDogu,
	}, nil
}
