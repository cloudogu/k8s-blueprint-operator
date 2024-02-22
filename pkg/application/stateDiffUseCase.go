package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type StateDiffUseCase struct {
	blueprintSpecRepo         blueprintSpecRepository
	doguInstallationRepo      doguInstallationRepository
	componentInstallationRepo componentInstallationRepository
	globalConfigRepo          globalConfigEntryRepository
	doguConfigRepo            doguConfigEntryRepository
	sensitiveDoguConfigRepo   sensitiveDoguConfigEntryRepository
	encryptionAdapter         configEncryptionAdapter
}

func NewStateDiffUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguInstallationRepo domainservice.DoguInstallationRepository,
	componentInstallationRepo domainservice.ComponentInstallationRepository,
	globalConfigRepo domainservice.GlobalConfigEntryRepository,
	doguConfigRepo domainservice.DoguConfigEntryRepository,
	sensitiveDoguConfigRepo domainservice.SensitiveDoguConfigEntryRepository,
	encryptionAdapter configEncryptionAdapter,
) *StateDiffUseCase {
	return &StateDiffUseCase{
		blueprintSpecRepo:         blueprintSpecRepo,
		doguInstallationRepo:      doguInstallationRepo,
		componentInstallationRepo: componentInstallationRepo,
		globalConfigRepo:          globalConfigRepo,
		doguConfigRepo:            doguConfigRepo,
		sensitiveDoguConfigRepo:   sensitiveDoguConfigRepo,
		encryptionAdapter:         encryptionAdapter,
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

	logger.Info("collect ecosystem state for state diff")
	ecosystemState, err := useCase.collectEcosystemState(ctx, blueprintSpec.EffectiveBlueprint)
	if err != nil {
		return fmt.Errorf("could not determine state diff: %w", err)
	}

	//determine state diff
	logger.Info("determine state diff to the cloudogu ecosystem", "blueprintStatus", blueprintSpec.Status)
	stateDiffError := blueprintSpec.DetermineStateDiff(ecosystemState)
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

	//TODO: collect ecosystem state in parallel (like for ecosystem health) if we have time
	//load current dogus and components
	logger.Info("collect installed dogus")
	installedDogus, doguErr := useCase.doguInstallationRepo.GetAll(ctx)
	logger.Info("collect installed components")
	installedComponents, componentErr := useCase.componentInstallationRepo.GetAll(ctx)
	// load current config
	logger.Info("collect needed global config")
	globalConfig, globalConfigErr := useCase.globalConfigRepo.GetAllByKey(ctx, effectiveBlueprint.Config.Global.GetGlobalConfigKeys())
	logger.Info("collect needed dogu config")
	doguConfig, doguConfigErr := useCase.doguConfigRepo.GetAllByKey(ctx, effectiveBlueprint.Config.GetDoguConfigKeys())
	logger.Info("collect needed sensitive dogu config")
	sensitiveDoguConfig, sensitiveConfigErr := useCase.sensitiveDoguConfigRepo.GetAllByKey(ctx, effectiveBlueprint.Config.GetSensitiveDoguConfigKeys())

	joinedError := errors.Join(doguErr, componentErr, globalConfigErr, doguConfigErr, sensitiveConfigErr)

	var internalErrorType *domainservice.InternalError
	if errors.As(joinedError, &internalErrorType) {
		return ecosystem.EcosystemState{}, fmt.Errorf("could not collect ecosystem state: %w", joinedError)
	}

	encryptedConfig := map[common.SensitiveDoguConfigKey]common.EncryptedDoguConfigValue{}
	for key, entry := range sensitiveDoguConfig {
		encryptedConfig[key] = entry.Value
	}

	logger.Info("decrypt sensitive dogu config")
	decryptedConfig, err := useCase.encryptionAdapter.DecryptAll(ctx, encryptedConfig)
	if err != nil {
		// we cannot ignore any error type here:
		// - InternalError -> there could be a network error -> retry by reconciliation
		// - NotFoundError -> we only have encrypted values to decrypt, therefore the encryption key pair should be present
		//                    if the key pair is not present, we could have a serious problem or there is config for a not installed dogu
		return ecosystem.EcosystemState{}, fmt.Errorf("could not decrypt sensitive dogu config: %w", err)
	}

	return ecosystem.EcosystemState{
		InstalledDogus:               installedDogus,
		InstalledComponents:          installedComponents,
		GlobalConfig:                 globalConfig,
		DoguConfig:                   doguConfig,
		EncryptedDoguConfig:          sensitiveDoguConfig,
		DecryptedSensitiveDoguConfig: decryptedConfig,
	}, nil
}
