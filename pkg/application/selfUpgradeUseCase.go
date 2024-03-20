package application

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.ApplySelfUpgrade").
		WithValues("blueprintId", blueprintId)
	ctx = log.IntoContext(ctx, logger)

	blueprintSpec, err := useCase.blueprintRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to possibly self upgrade the operator: %w", blueprintId, err)
	}
	ownComponent, err := useCase.componentRepo.GetByName(ctx, useCase.blueprintOperatorName)

	if err != nil && !domainservice.IsNotFoundError(err) {
		// ignore not found errors as this is ok if the component was not installed via a component CR
		// only return if other errors happen, e.g. InternalError
		return fmt.Errorf("cannot load component installation for %q from ecosystem: %w", useCase.blueprintOperatorName, err)
	}
	// FIXME: we need to use the version which the component operator really has installed, not what is just in the spec.
	//  the actual version is not yet implemented in the component operator.
	var actualInstalledVersion *semver.Version
	if ownComponent != nil {
		actualInstalledVersion = ownComponent.Version
	}
	ownDiff := blueprintSpec.HandleSelfUpgrade(useCase.blueprintOperatorName, actualInstalledVersion)
	err = useCase.blueprintRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q while possibly self upgrading the operator: %w", blueprintId, err)
	}

	if blueprintSpec.Status == domain.StatusPhaseAwaitSelfUpgrade {
		logger.Info("self upgrade needed, apply self upgrade")
		err = useCase.applySelfUpgrade(ctx, ownDiff, ownComponent)
		if err != nil {
			return err
		}
		logger.Info("await termination for self upgrade. You can check the component-CR for the installation status")
		useCase.waitForTermination(ctx)
		// nothing can come after this as the operator gets terminated while waiting.
	} else {
		// we don't need to check the health status as this code would not run if the operator is not healthy.
		logger.Info("self upgrade successful or not needed")
	}
	return nil
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
