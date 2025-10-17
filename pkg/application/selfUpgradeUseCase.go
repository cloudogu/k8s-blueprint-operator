package application

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const awaitSelfUpgradeErrorMsg = "await self upgrade"

type SelfUpgradeUseCase struct {
	blueprintRepo         domainservice.BlueprintSpecRepository
	componentRepo         componentInstallationRepository
	componentUseCase      componentInstallationUseCase
	blueprintOperatorName common.SimpleComponentName
}

func NewSelfUpgradeUseCase(
	blueprintRepo domainservice.BlueprintSpecRepository,
	componentRepo componentInstallationRepository,
	componentUseCase componentInstallationUseCase,
	blueprintOperatorName common.SimpleComponentName,
) *SelfUpgradeUseCase {
	return &SelfUpgradeUseCase{
		blueprintRepo:         blueprintRepo,
		componentRepo:         componentRepo,
		componentUseCase:      componentUseCase,
		blueprintOperatorName: blueprintOperatorName,
	}
}

// HandleSelfUpgrade checks if a self upgrade is necessary, executes all needed steps and
// can check if the self upgrade was successful after a restart.
// It always sets the fitting status in the blueprint spec.
func (useCase *SelfUpgradeUseCase) HandleSelfUpgrade(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	ownDiff := blueprint.StateDiff.ComponentDiffs.GetComponentDiffByName(useCase.blueprintOperatorName)
	ownComponent, err := useCase.componentRepo.GetByName(ctx, useCase.blueprintOperatorName)

	if err != nil && !domainservice.IsNotFoundError(err) {
		// ignore not found errors as this is ok if the component was not installed via a component CR before
		// only return if other errors happen, e.g. InternalError
		return fmt.Errorf("cannot load component installation for %q from ecosystem: %w", useCase.blueprintOperatorName, err)
	}

	needsToApply, isCompleted := checkStateOfSelfUpgrade(ownDiff, ownComponent)

	if isCompleted {
		blueprint.MarkSelfUpgradeCompleted()
		err = useCase.blueprintRepo.Update(ctx, blueprint)
		if err != nil {
			return fmt.Errorf("cannot save blueprint spec %q after self upgrading the operator: %w", blueprint.Id, err)
		}
		return nil
	}
	// if not done, set conditions accordingly
	blueprint.MarkWaitingForSelfUpgrade()
	err = useCase.blueprintRepo.Update(ctx, blueprint)
	if err != nil {
		return fmt.Errorf("cannot persist blueprint spec %q to mark it waiting for self upgrade: %w", blueprint.Id, err)
	}

	if needsToApply {
		return useCase.doSelfUpgrade(ctx, ownDiff, ownComponent)
	}
	return &domain.AwaitSelfUpgradeError{Message: awaitSelfUpgradeErrorMsg}
}

func checkStateOfSelfUpgrade(ownDiff domain.ComponentDiff, ownComponent *ecosystem.ComponentInstallation) (needsToApply, isCompleted bool) {
	// use extra vars to avoid nil pointer dereferences of the component
	var versionSetForInstallation, installedVersion *semver.Version
	if ownComponent != nil {
		versionSetForInstallation = ownComponent.ExpectedVersion
		installedVersion = ownComponent.ActualVersion
	}

	if ownDiff.IsExpectedVersion(installedVersion) {
		// if component CR status.installedVersion already says: our wanted version is installed
		return false, true
	}
	if ownDiff.IsExpectedVersion(versionSetForInstallation) {
		// update is already triggered but not done
		return false, false
	} else {
		// no update triggered yet
		// trigger update and trigger later reconciliation via error
		return true, false
	}
}

func (useCase *SelfUpgradeUseCase) doSelfUpgrade(ctx context.Context, ownDiff domain.ComponentDiff, ownComponent *ecosystem.ComponentInstallation) error {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.doSelfUpgrade")
	logger.Info("self upgrade needed, apply self upgrade")
	err := useCase.componentUseCase.applyComponentState(ctx, ownDiff, ownComponent)
	if err != nil {
		return fmt.Errorf("an error occurred while applying the self-upgrade to the ecosystem: %w", err)
	}
	return &domain.AwaitSelfUpgradeError{Message: awaitSelfUpgradeErrorMsg}
}
