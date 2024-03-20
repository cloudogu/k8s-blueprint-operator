package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

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
func (useCase *SelfUpgradeUseCase) HandleSelfUpgrade(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.ApplySelfUpgrade")

	blueprintSpec, err := useCase.blueprintRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to possibly self upgrade the operator: %w", blueprintId, err)
	}

	ownDiff := blueprintSpec.StateDiff.ComponentDiffs.GetComponentDiffByName(useCase.blueprintOperatorName)

	if !ownDiff.HasChanges() {
		logger.Info("self upgrade not needed")
		blueprintSpec.MarkSelfUpgradeCompleted()
		err = useCase.blueprintRepo.Update(ctx, blueprintSpec)
		if err != nil {
			return fmt.Errorf("cannot save blueprint spec %q to skip self upgrade: %w", blueprintSpec.Id, err)
		}
		return nil
	}

	ownComponent, err := useCase.componentRepo.GetByName(ctx, useCase.blueprintOperatorName)

	if err != nil && !domainservice.IsNotFoundError(err) {
		// ignore not found errors as this is ok if the component was not installed via a component CR
		// only return if other errors happen, e.g. InternalError
		return fmt.Errorf("cannot load component installation for %q from ecosystem: %w", useCase.blueprintOperatorName, err)
	}
	// use extra vars to avoid nil pointer dereferences of the component
	var expectedVersion, actualVersion *semver.Version
	if ownComponent != nil {
		expectedVersion = ownComponent.ExpectedVersion
		actualVersion = ownComponent.ActualVersion
	}

	if !ownDiff.IsExpectedVersion(expectedVersion) {
		return useCase.doSelfUpgrade(ctx, blueprintSpec, ownDiff, ownComponent)
		// the operator waits for termination, unless there was an error, so we can return here
	}

	if !ownDiff.IsExpectedVersion(actualVersion) {
		err = useCase.awaitInstallationConfirmation(ctx, blueprintSpec)
		if err != nil {
			return err
		}
	}

	blueprintSpec.MarkSelfUpgradeCompleted()
	err = useCase.blueprintRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after self upgrading the operator: %w", blueprintSpec.Id, err)
	}
	logger.Info("self upgrade successful")

	return nil
}

func (useCase *SelfUpgradeUseCase) doSelfUpgrade(ctx context.Context, blueprintSpec *domain.BlueprintSpec, ownDiff domain.ComponentDiff, ownComponent *ecosystem.ComponentInstallation) error {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.doSelfUpgrade")
	logger.Info("self upgrade needed, apply self upgrade")
	blueprintSpec.MarkWaitingForSelfUpgrade()
	err := useCase.blueprintRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot persist blueprint spec %q to mark it waiting for self upgrade: %w", blueprintSpec.Id, err)
	}
	err = useCase.applySelfUpgrade(ctx, ownDiff, ownComponent)
	if err != nil {
		return err
	}
	logger.Info("await termination for self upgrade. Check the component-CR for the installation status")
	useCase.waitForTermination(ctx)
	return nil // this code is never reached as we wait for termination before
}

func (useCase *SelfUpgradeUseCase) applySelfUpgrade(ctx context.Context, ownDiff domain.ComponentDiff, ownComponent *ecosystem.ComponentInstallation) error {
	err := useCase.componentUseCase.applyComponentState(ctx, ownDiff, ownComponent)
	if err != nil {
		return fmt.Errorf("an error occurred while applying the self-upgrade to the ecosystem: %w", err)
	}
	return nil
}

func (useCase *SelfUpgradeUseCase) waitForTermination(ctx context.Context) {
	<-ctx.Done()
}

func (useCase *SelfUpgradeUseCase) awaitInstallationConfirmation(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	//TODO: extract retryInterval
	_, err := util.RetryUntilSuccessOrCancellation(ctx, 5*time.Second, func(ctx context.Context) (*interface{}, error, bool) {
		ownComponent, err := useCase.componentRepo.GetByName(ctx, useCase.blueprintOperatorName)
		if err != nil {
			return nil, err, false
		}
		ownDiff := blueprintSpec.StateDiff.ComponentDiffs.GetComponentDiffByName(useCase.blueprintOperatorName)
		return nil, nil, !ownDiff.IsExpectedVersion(ownComponent.ActualVersion) // retry if true
	})
	if err != nil && !errors.Is(err, ctx.Err()) {
		// ignore cancellation error as this can happen, if the operator is getting restarted more than once (e.g. maybe because of a cluster failure)
		return fmt.Errorf("error while waiting for version confirmation: %w", err)
	}
	return nil
}
