package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type StateDiffUseCase struct {
	blueprintSpecRepo         blueprintSpecRepository
	doguInstallationRepo      doguInstallationRepository
	componentInstallationRepo componentInstallationRepository
	globalConfigRepo          GlobalConfigEntryRepository
	doguConfigRepo            DoguConfigEntryRepository
	sensitiveDoguConfigRepo   SensitiveDoguConfigEntryRepository
}

func NewStateDiffUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguInstallationRepo domainservice.DoguInstallationRepository,
	componentInstallationRepo domainservice.ComponentInstallationRepository,
	globalConfigRepo domainservice.GlobalConfigEntryRepository,
	doguConfigRepo domainservice.DoguConfigEntryRepository,
	sensitiveDoguConfigRepo domainservice.SensitiveDoguConfigEntryRepository,
) *StateDiffUseCase {
	return &StateDiffUseCase{
		blueprintSpecRepo:         blueprintSpecRepo,
		doguInstallationRepo:      doguInstallationRepo,
		componentInstallationRepo: componentInstallationRepo,
		globalConfigRepo:          globalConfigRepo,
		doguConfigRepo:            doguConfigRepo,
		sensitiveDoguConfigRepo:   sensitiveDoguConfigRepo,
	}
}

// DetermineStateDiff loads the state of the ecosystem and compares it to the blueprint. It creates a declarative diff.
// returns a domainservice.NotFoundError if the blueprintId does not correspond to a blueprintSpec or
// a domainservice.InternalError if there is any error while loading or persisting the blueprintSpec or
// a domainservice.ConflictError if there was a concurrent write.
// a domain.InvalidBlueprintError if there are any forbidden actions in the stateDiff.
// any error if there is any other error.
func (useCase *StateDiffUseCase) DetermineStateDiff(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("StateDiffUseCase.DetermineStateDiff").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for determining state diff")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to determine state diff: %w", blueprintId, err)
	}

	//load current dogus and components
	logger.Info("determine state diff to the cloudogu ecosystem", "blueprintStatus", blueprintSpec.Status)
	installedDogus, err := useCase.doguInstallationRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot get installed dogus to determine state diff: %w", err)
	}

	installedComponents, err := useCase.componentInstallationRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot get installed components to determine state diff: %w", err)
	}

	// load current config
	actualGlobalConfig, err := useCase.globalConfigRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot get global config to determine state diff: %w", err)
	}
	globalConfigByKey := util.ToMap(actualGlobalConfig, func(entry *ecosystem.GlobalConfigEntry) common.GlobalConfigKey { return entry.Key })

	actualDoguConfig, err := useCase.doguConfigRepo.GetAllByKey2(ctx, blueprintSpec.EffectiveBlueprint.Config.GetDoguConfigKeys())
	if err != nil {
		return fmt.Errorf("cannot get dogu config to determine state diff: %w", err)
	}
	actualSensitiveDoguConfig, err := useCase.sensitiveDoguConfigRepo.GetAllByKey2(ctx, blueprintSpec.EffectiveBlueprint.Config.GetSensitiveDoguConfigKeys())
	if err != nil {
		return fmt.Errorf("cannot get sensitive dogu config to determine state diff: %w", err)
	}

	//determine state diff
	stateDiffError := blueprintSpec.DetermineStateDiff(installedDogus, installedComponents, globalConfigByKey, actualDoguConfig, actualSensitiveDoguConfig)
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
