package application

import (
	"context"
	"fmt"
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
	if err != nil {
		return err
	}
	// FIXME: we need to use the version which the component operator really has installed, not what is just in the spec.
	//  the actual version is not yet implemented in the component operator.
	ownDiff := blueprintSpec.HandleSelfUpgrade(useCase.blueprintOperatorName, ownComponent.Version)
	err = useCase.blueprintRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return err
	}

	if blueprintSpec.Status == domain.StatusPhaseAwaitSelfUpgrade {
		logger.Info("self upgrade needed")
		err = useCase.applySelfUpgrade(ctx, ownDiff, ownComponent)
		logger.Info("apply self upgrade")
		if err != nil {
			return err
		}
		logger.Info("await self upgrade")
		useCase.waitForTermination(ctx)
		// nothing can come after this as the operator gets terminated while waiting.
	} else {
		// we don't need to check the health status as this code would not run if the operator is not healthy.
		logger.Info("self upgrade successful or not needed")
	}
	return nil
}

func (useCase *SelfUpgradeUseCase) applySelfUpgrade(ctx context.Context, ownDiff domain.ComponentDiff, ownComponent *ecosystem.ComponentInstallation) error {
	err := useCase.componentUseCase.ApplyComponentState(ctx, ownDiff, ownComponent)
	if err != nil {
		return fmt.Errorf("an error occurred while applying the self-upgrade to the ecosystem: %w", err)
	}
	return nil
}

func (useCase *SelfUpgradeUseCase) waitForTermination(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	}
}
