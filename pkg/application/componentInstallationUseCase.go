package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type ComponentInstallationUseCase struct {
	blueprintSpecRepo   domainservice.BlueprintSpecRepository
	componentRepo       domainservice.ComponentInstallationRepository
	healthCheckInterval time.Duration
}

func NewComponentInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	componentRepo domainservice.ComponentInstallationRepository,
	healthCheckInterval time.Duration,
) *ComponentInstallationUseCase {
	return &ComponentInstallationUseCase{
		blueprintSpecRepo:   blueprintSpecRepo,
		componentRepo:       componentRepo,
		healthCheckInterval: healthCheckInterval,
	}
}

// ApplyComponentStates applies the expected component state from the Blueprint to the ecosystem.
// Fail-fast here, so that the possible damage is as small as possible.
func (useCase *ComponentInstallationUseCase) ApplyComponentStates(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.ApplyComponentStates").
		WithValues("blueprintId", blueprintId)
	log.IntoContext(ctx, logger)

	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to install components: %w", blueprintId, err)
	}

	if len(blueprintSpec.StateDiff.ComponentDiffs) == 0 {
		logger.Info("apply no components because blueprint has no component state differences")
		return nil
	}

	// ComponentDiff contains all installed components anyway (but some with action none) so we can load them all at once
	components, err := useCase.componentRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot load component installations to apply component state: %w", err)
	}

	for _, componentDiff := range blueprintSpec.StateDiff.ComponentDiffs {
		err = useCase.applyComponentState(ctx, componentDiff, components[componentDiff.ComponentName], blueprintSpec.Config)
		if err != nil {
			return fmt.Errorf("an error occurred while applying component state to the ecosystem: %w", err)
		}
	}
	return nil
}

func (useCase *ComponentInstallationUseCase) applyComponentState(
	ctx context.Context,
	componentDiff domain.ComponentDiff,
	componentInstallation *ecosystem.ComponentInstallation,
	blueprintConfig domain.BlueprintConfiguration,
) error {
	logger := log.FromContext(ctx).
		WithName("ComponentInstallationUseCase.applyComponentState").
		WithValues("component", componentDiff.ComponentName, "diff", componentDiff.String())

	switch componentDiff.NeededAction {
	case domain.ActionNone:
		logger.Info("apply nothing for component")
		return nil
	case domain.ActionInstall:
		logger.Info("install component")
		// TODO wait for namespace, deployNamespace and valuesYamlOverwrite in diff
		newComponent := ecosystem.InstallComponent("k8s", componentDiff.ComponentName, componentDiff.Expected.Version)
		return useCase.componentRepo.Create(ctx, newComponent)
	case domain.ActionUninstall:
		logger.Info("uninstall component")
		return useCase.componentRepo.Delete(ctx, componentInstallation.Name)
	case domain.ActionUpgrade:
		logger.Info("upgrade component")
		componentInstallation.Upgrade(componentDiff.Expected.Version)
		return useCase.componentRepo.Update(ctx, componentInstallation)
	case domain.ActionDowngrade:
		logger.Info("downgrade component")
		return fmt.Errorf(noDowngradesExplanationText)
	case domain.ActionSwitchDoguNamespace:
		logger.Info("do namespace switch for component")
		// TODO
		// err := componentInstallation.SwitchNamespace(
		// 	componentDiff.Expected.Namespace,
		// 	componentDiff.Expected.Version,
		// 	blueprintConfig.AllowDoguNamespaceSwitch,
		// )
		// if err != nil {
		// 	return err
		// }
		return useCase.componentRepo.Update(ctx, componentInstallation)
	default:
		return fmt.Errorf("cannot perform unknown action %q", componentDiff.NeededAction)
	}
}
