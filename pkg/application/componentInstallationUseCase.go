package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	noDowngradesExplanationTextFmt = "downgrades are not allowed as the data model of the %s could have changed and " +
		"doing rollbacks to older models is not supported. " +
		"You can downgrade %s by restoring a backup. " +
		"If you want an 'allow-downgrades' flag, issue a feature request"
	noDistributionNamespaceSwitchExplanationText = "switching distribution namespace is not allowed. If you want an " +
		"`allow-switch-distribution-namespace` flag, issue a feature request"
)

type ComponentInstallationUseCase struct {
	blueprintSpecRepo   domainservice.BlueprintSpecRepository
	componentRepo       domainservice.ComponentInstallationRepository
	healthConfigProvider healthConfigProvider
}

func NewComponentInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	componentRepo domainservice.ComponentInstallationRepository,
	healthConfigProvider healthConfigProvider,
) *ComponentInstallationUseCase {
	return &ComponentInstallationUseCase{
		blueprintSpecRepo:   blueprintSpecRepo,
		componentRepo:       componentRepo,
		healthConfigProvider: healthConfigProvider,
	}
}


func (useCase *ComponentInstallationUseCase) CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.CheckComponentHealth")
	logger.Info("check component health...")
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

func (useCase *ComponentInstallationUseCase) WaitForHealthyComponents(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	logger := log.FromContext(ctx).WithName("ComponentInstallationUseCase.WaitForHealthyComponents")

	waitConfig, err := useCase.healthConfigProvider.GetWaitConfig(ctx)
	if err != nil {
		return ecosystem.ComponentHealthResult{}, fmt.Errorf("failed to get health check interval: %w", err)
	}

	logger.Info("start waiting for component health")
	healthResult, err := util.RetryUntilSuccessOrCancellation(
		ctx,
		waitConfig.Interval,
		useCase.checkComponentHealthStatesRetryable,
	)
	var result ecosystem.ComponentHealthResult
	if healthResult == nil {
		result = ecosystem.ComponentHealthResult{}
	} else {
		result = *healthResult
	}

	if err != nil {
		err = fmt.Errorf("stop waiting for component health: %w", err)
		logger.Error(err, "stop waiting for component health because of an error or time out")
	}

	return result, err
}

func (useCase *ComponentInstallationUseCase) checkComponentHealthStatesRetryable(ctx context.Context) (result *ecosystem.ComponentHealthResult, err error, shouldRetry bool) {
	// use named return values to make their meaning clear
	health, err := useCase.CheckComponentHealth(ctx)
	if err != nil {
		// no retry on error while loading components
		return &ecosystem.ComponentHealthResult{}, err, false
	}
	result = &health
	shouldRetry = !health.AllHealthy()
	return
}

func (useCase *ComponentInstallationUseCase) applyComponentState(
	ctx context.Context,
	componentDiff domain.ComponentDiff,
	componentInstallation *ecosystem.ComponentInstallation,
) error {
	logger := log.FromContext(ctx).
		WithName("ComponentInstallationUseCase.applyComponentState").
		WithValues("component", componentDiff.Name, "diff", componentDiff.String())

	switch componentDiff.NeededAction {
	case domain.ActionNone:
		logger.Info("apply nothing for component")
		return nil
	case domain.ActionInstall:
		logger.Info("install component")
		// TODO apply valuesYamlOverwrite
		newComponent := ecosystem.InstallComponent(componentDiff.Expected.DistributionNamespace, componentDiff.Name, componentDiff.Expected.Version)
		return useCase.componentRepo.Create(ctx, newComponent)
	case domain.ActionUninstall:
		logger.Info("uninstall component")
		return useCase.componentRepo.Delete(ctx, componentInstallation.Name)
	case domain.ActionUpgrade:
		logger.Info("upgrade component")
		// TODO apply valuesYamlOverwrite
		componentInstallation.Upgrade(componentDiff.Expected.Version)
		return useCase.componentRepo.Update(ctx, componentInstallation)
	case domain.ActionSwitchComponentDistributionNamespace:
		logger.Info("switch distribution namespace")
		return fmt.Errorf(noDistributionNamespaceSwitchExplanationText)
	case domain.ActionDowngrade:
		logger.Info("downgrade component")
		return fmt.Errorf(getNoDowngradesExplanationTextForComponents())
	default:
		return fmt.Errorf("cannot perform unknown action %q", componentDiff.NeededAction)
	}
}

func getNoDowngradesExplanationTextForComponents() string {
	return fmt.Sprintf(noDowngradesExplanationTextFmt, "components", "components")
}
