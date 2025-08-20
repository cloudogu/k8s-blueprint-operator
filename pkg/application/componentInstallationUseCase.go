package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	noDowngradesExplanationTextFmt = "downgrades are not allowed as the data model of the %s could have changed and " +
		"doing rollbacks to older models is not supported. " +
		"You can downgrade %s by restoring a backup. " +
		"If you want an 'allow-downgrades' flag, issue a feature request"
	noDistributionNamespaceSwitchExplanationText = "switching distribution namespace of components is not allowed. If you want an " +
		"`allow-switch-distribution-namespace` flag, issue a feature request"
)

type ComponentInstallationUseCase struct {
	blueprintSpecRepo    domainservice.BlueprintSpecRepository
	componentRepo        domainservice.ComponentInstallationRepository
	healthConfigProvider healthConfigProvider
}

func NewComponentInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	componentRepo domainservice.ComponentInstallationRepository,
	healthConfigProvider healthConfigProvider,
) *ComponentInstallationUseCase {
	return &ComponentInstallationUseCase{
		blueprintSpecRepo:    blueprintSpecRepo,
		componentRepo:        componentRepo,
		healthConfigProvider: healthConfigProvider,
	}
}

func (useCase *ComponentInstallationUseCase) CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.CheckComponentHealth")
	logger.V(2).Info("check component health...")
	installedComponents, err := useCase.componentRepo.GetAll(ctx)
	if err != nil {
		return ecosystem.ComponentHealthResult{}, fmt.Errorf("cannot retrieve installed components: %w", err)
	}

	requiredComponents, err := useCase.healthConfigProvider.GetRequiredComponents(ctx)
	if err != nil {
		return ecosystem.ComponentHealthResult{}, fmt.Errorf("cannot retrieve required components: %w", err)
	}

	return ecosystem.CalculateComponentHealthResult(installedComponents, requiredComponents), nil
}

// ApplyComponentStates applies the expected component state from the Blueprint to the ecosystem.
// Fail-fast here, so that the possible damage is as small as possible.
func (useCase *ComponentInstallationUseCase) ApplyComponentStates(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.ApplyComponentStates")

	if len(blueprint.StateDiff.ComponentDiffs) == 0 {
		logger.Info("apply no components because blueprint has no component state differences")
		return nil
	}

	// ComponentDiff contains all installed components anyway (but some with action none) so we can load them all at once
	components, err := useCase.componentRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot load component installations to apply component state: %w", err)
	}

	for _, componentDiff := range blueprint.StateDiff.ComponentDiffs {
		err = useCase.applyComponentState(ctx, componentDiff, components[componentDiff.Name])
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
) error {
	logger := log.FromContext(ctx).
		WithName("ComponentInstallationUseCase.applyComponentState").
		WithValues("component", componentDiff.Name, "diff", componentDiff.String())

	for _, action := range componentDiff.NeededActions {
		switch action {
		case domain.ActionInstall:
			logger.Info("install component")
			newComponent := ecosystem.InstallComponent(common.QualifiedComponentName{
				Namespace:  componentDiff.Expected.Namespace,
				SimpleName: componentDiff.Name,
			}, componentDiff.Expected.Version, componentDiff.Expected.DeployConfig)
			return useCase.componentRepo.Create(ctx, newComponent)
		case domain.ActionUninstall:
			logger.Info("uninstall component")
			return useCase.componentRepo.Delete(ctx, componentInstallation.Name.SimpleName)
		case domain.ActionUpgrade:
			componentInstallation.Upgrade(componentDiff.Expected.Version)
		case domain.ActionUpdateComponentDeployConfig:
			componentInstallation.UpdateDeployConfig(componentDiff.Expected.DeployConfig)
		case domain.ActionSwitchComponentNamespace:
			logger.Info("switch distribution namespace")
			return errors.New(noDistributionNamespaceSwitchExplanationText)
		case domain.ActionDowngrade:
			logger.Info("downgrade component")
			return fmt.Errorf(noDowngradesExplanationTextFmt, "components", "components")
		default:
			return fmt.Errorf("cannot perform unknown action %q", action)
		}
	}

	// If this routine did not terminate until this point, it is always an update.
	if len(componentDiff.NeededActions) > 0 {
		logger.Info("upgrade component")
		return useCase.componentRepo.Update(ctx, componentInstallation)
	}

	return nil
}
