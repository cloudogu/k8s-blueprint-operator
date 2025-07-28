package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type StateDiffUseCase struct {
	blueprintSpecRepo         blueprintSpecRepository
	doguInstallationRepo      doguInstallationRepository
	componentInstallationRepo componentInstallationRepository
	globalConfigRepo          globalConfigRepository
	doguConfigRepo            doguConfigRepository
	sensitiveDoguConfigRepo   sensitiveDoguConfigRepository
	sensitiveConfigRefReader  sensitiveConfigRefReader
}

func NewStateDiffUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguInstallationRepo domainservice.DoguInstallationRepository,
	componentInstallationRepo domainservice.ComponentInstallationRepository,
	globalConfigRepo domainservice.GlobalConfigRepository,
	doguConfigRepo domainservice.DoguConfigRepository,
	sensitiveDoguConfigRepo domainservice.SensitiveDoguConfigRepository,
	sensitiveConfigRefReader domainservice.SensitiveConfigRefReader,
) *StateDiffUseCase {
	return &StateDiffUseCase{
		blueprintSpecRepo:         blueprintSpecRepo,
		doguInstallationRepo:      doguInstallationRepo,
		componentInstallationRepo: componentInstallationRepo,
		globalConfigRepo:          globalConfigRepo,
		doguConfigRepo:            doguConfigRepo,
		sensitiveDoguConfigRepo:   sensitiveDoguConfigRepo,
		sensitiveConfigRefReader:  sensitiveConfigRefReader,
	}
}

// DetermineStateDiff loads the state of the ecosystem and compares it to the blueprint. It creates a declarative diff.
// returns:
//   - a domainservice.NotFoundError if the blueprint was not found or could not found dogu decryption keys or
//   - a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or while collecting the ecosystem state or
//   - a domainservice.ConflictError if there was a concurrent write to the blueprint or
//   - a domain.InvalidBlueprintError if there are any forbidden actions in the stateDiff.
//   - any error if there is any other error.
func (useCase *StateDiffUseCase) DetermineStateDiff(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("StateDiffUseCase.DetermineStateDiff").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for determining state diff")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to determine state diff: %w", blueprintId, err)
	}

	logger.Info("load referenced sensitive config")
	// load referenced config before collecting ecosystem state
	// if an error happens here, we save a lot of heavy work
	referencedSensitiveConfig, err := useCase.sensitiveConfigRefReader.GetValues(
		ctx, blueprintSpec.EffectiveBlueprint.Config.GetSensitiveConfigReferences(),
	)
	if err != nil {
		return fmt.Errorf("could not load referenced sensitive config: %w", err)
	}

	logger.Info("collect ecosystem state for state diff")
	ecosystemState, err := useCase.collectEcosystemState(ctx, blueprintSpec.EffectiveBlueprint)
	if err != nil {
		return fmt.Errorf("could not determine state diff: %w", err)
	}

	logger.Info("determine state diff to the cloudogu ecosystem", "blueprintStatus", blueprintSpec.Status)
	stateDiffError := blueprintSpec.DetermineStateDiff(ecosystemState, referencedSensitiveConfig)
	var invalidError *domain.InvalidBlueprintError
	if errors.As(stateDiffError, &invalidError) {
		// do not return here as with this error the blueprint status and events should be persisted as normal.
	} else if stateDiffError != nil {
		return fmt.Errorf("failed to determine state diff for blueprint %q: %w", blueprintId, stateDiffError)
	}
	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after determining the state diff to the ecosystem: %w", blueprintId, err)
	}

	// return this error back here to persist the blueprint status and events first.
	// return it to signal that a repeated call to this function will not result in any progress.
	return stateDiffError
}

func (useCase *StateDiffUseCase) collectEcosystemState(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) (ecosystem.EcosystemState, error) {
	logger := log.FromContext(ctx).WithName("StateDiffUseCase.collectEcosystemState")

	// TODO: collect ecosystem state in parallel (like for ecosystem health) if we have time
	// load current dogus and components
	logger.Info("collect installed dogus")
	installedDogus, doguErr := useCase.doguInstallationRepo.GetAll(ctx)
	logger.Info("collect installed components")
	installedComponents, componentErr := useCase.componentInstallationRepo.GetAll(ctx)
	// load current config
	logger.Info("collect needed global config")
	globalConfig, globalConfigErr := useCase.globalConfigRepo.Get(ctx)

	logger.Info("collect needed dogu config")
	configByDogu, doguConfigErr := useCase.doguConfigRepo.GetAllExisting(ctx, effectiveBlueprint.Config.GetDogusWithChangedConfig())

	logger.Info("collect needed sensitive dogu config")
	sensitiveConfigByDogu, sensitiveConfigErr := useCase.sensitiveDoguConfigRepo.GetAllExisting(ctx, effectiveBlueprint.Config.GetDogusWithChangedSensitiveConfig())

	joinedError := errors.Join(doguErr, componentErr, globalConfigErr, doguConfigErr, sensitiveConfigErr)
	if joinedError != nil {
		return ecosystem.EcosystemState{}, fmt.Errorf("could not collect ecosystem state: %w", joinedError)
	}

	return ecosystem.EcosystemState{
		InstalledDogus:        installedDogus,
		InstalledComponents:   installedComponents,
		GlobalConfig:          globalConfig,
		ConfigByDogu:          configByDogu,
		SensitiveConfigByDogu: sensitiveConfigByDogu,
	}, nil
}
