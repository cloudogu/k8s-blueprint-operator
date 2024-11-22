package application

import (
	"context"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"golang.org/x/exp/maps"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguInstallationUseCase struct {
	blueprintSpecRepo  blueprintSpecRepository
	doguRepo           doguInstallationRepository
	waitConfigProvider healthWaitConfigProvider
}

func NewDoguInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguRepo domainservice.DoguInstallationRepository,
	waitConfigProvider domainservice.HealthWaitConfigProvider,
) *DoguInstallationUseCase {
	return &DoguInstallationUseCase{
		blueprintSpecRepo:  blueprintSpecRepo,
		doguRepo:           doguRepo,
		waitConfigProvider: waitConfigProvider,
	}
}

func (useCase *DoguInstallationUseCase) CheckDoguHealth(ctx context.Context) (ecosystem.DoguHealthResult, error) {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.CheckDoguHealth")
	logger.Info("check dogu health...")
	installedDogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return ecosystem.DoguHealthResult{}, fmt.Errorf("cannot evaluate dogu health states: %w", err)
	}
	// accept experimental maps.Values as we can implement it ourselves in a minute
	return ecosystem.CalculateDoguHealthResult(maps.Values(installedDogus)), nil
}

func (useCase *DoguInstallationUseCase) WaitForHealthyDogus(ctx context.Context) (ecosystem.DoguHealthResult, error) {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.WaitForHealthyDogus")

	waitConfig, err := useCase.waitConfigProvider.GetWaitConfig(ctx)
	if err != nil {
		return ecosystem.DoguHealthResult{}, fmt.Errorf("failed to get health check interval: %w", err)
	}

	logger.Info("start waiting for dogu health")
	healthResult, err := util.RetryUntilSuccessOrCancellation(
		ctx,
		waitConfig.Interval,
		useCase.checkDoguHealthStatesRetryable,
	)
	var result ecosystem.DoguHealthResult
	if healthResult == nil {
		result = ecosystem.DoguHealthResult{}
	} else {
		result = *healthResult
	}

	if err != nil {
		err = fmt.Errorf("stop waiting for dogu health: %w", err)
		logger.Error(err, "stop waiting for dogu health because of an error or time out")
	}

	return result, err
}

func (useCase *DoguInstallationUseCase) checkDoguHealthStatesRetryable(ctx context.Context) (result *ecosystem.DoguHealthResult, err error, shouldRetry bool) {
	// use named return values to make their meaning clear
	health, err := useCase.CheckDoguHealth(ctx)
	if err != nil {
		// no retry if error while loading dogus
		return &ecosystem.DoguHealthResult{}, err, false
	}
	result = &health
	shouldRetry = !health.AllHealthy()
	return
}

// ApplyDoguStates applies the expected dogu state from the Blueprint to the ecosystem.
// Fail-fast here, so that the possible damage is as small as possible.
func (useCase *DoguInstallationUseCase) ApplyDoguStates(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.ApplyDoguChanges").
		WithValues("blueprintId", blueprintId)
	log.IntoContext(ctx, logger)

	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to install dogus: %w", blueprintId, err)
	}

	// DoguDiff contains all installed dogus anyway (but some with action none) so we can load them all at once
	dogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot load dogu installations to apply dogu state: %w", err)
	}

	for _, doguDiff := range blueprintSpec.StateDiff.DoguDiffs {
		err = useCase.applyDoguState(ctx, doguDiff, dogus[doguDiff.DoguName], blueprintSpec.Config)
		if err != nil {
			return fmt.Errorf("an error occurred while applying dogu state to the ecosystem: %w", err)
		}
	}
	return nil
}

func (useCase *DoguInstallationUseCase) applyDoguState(
	ctx context.Context,
	doguDiff domain.DoguDiff,
	doguInstallation *ecosystem.DoguInstallation,
	blueprintConfig domain.BlueprintConfiguration,
) error {
	logger := log.FromContext(ctx).
		WithName("DoguInstallationUseCase.applyDoguState").
		WithValues("dogu", doguDiff.DoguName, "diff", doguDiff.String())
	for _, action := range doguDiff.NeededActions {
		switch action {
		case domain.ActionInstall:
			logger.Info("install dogu")
			newDogu := ecosystem.InstallDogu(cescommons.QualifiedName{
				Namespace:  doguDiff.Expected.Namespace,
				SimpleName: doguDiff.DoguName,
			}, doguDiff.Expected.Version, doguDiff.Expected.MinVolumeSize, doguDiff.Expected.ReverseProxyConfig)
			return useCase.doguRepo.Create(ctx, newDogu)
		case domain.ActionUninstall:
			logger.Info("uninstall dogu")
			return useCase.doguRepo.Delete(ctx, doguInstallation.Name.SimpleName)
		case domain.ActionUpgrade:
			doguInstallation.Upgrade(doguDiff.Expected.Version)
			continue
		case domain.ActionDowngrade:
			logger.Info("downgrade dogu")
			return fmt.Errorf(getNoDowngradesExplanationTextForDogus())
		case domain.ActionSwitchDoguNamespace:
			logger.Info("do namespace switch for dogu")
			err := doguInstallation.SwitchNamespace(
				doguDiff.Expected.Namespace,
				blueprintConfig.AllowDoguNamespaceSwitch,
			)
			if err != nil {
				return err
			}
			continue
		case domain.ActionUpdateDoguResourceMinVolumeSize:
			logger.Info("update minimum volume size for dogu")
			doguInstallation.UpdateMinVolumeSize(doguDiff.Expected.MinVolumeSize)
			continue
		case domain.ActionUpdateDoguProxyBodySize:
			logger.Info("update proxy body size for dogu")
			doguInstallation.UpdateProxyBodySize(doguDiff.Expected.ReverseProxyConfig.MaxBodySize)
			continue
		case domain.ActionUpdateDoguProxyRewriteTarget:
			logger.Info("update proxy rewrite target for dogu")
			doguInstallation.UpdateProxyRewriteTarget(doguDiff.Expected.ReverseProxyConfig.RewriteTarget)
			continue
		case domain.ActionUpdateDoguProxyAdditionalConfig:
			logger.Info("update proxy additional config for dogu")
			doguInstallation.UpdateProxyAdditionalConfig(doguDiff.Expected.ReverseProxyConfig.AdditionalConfig)
			continue
		default:
			return fmt.Errorf("cannot perform unknown action %q for dogu %q", action, doguDiff.DoguName)
		}
	}

	// If this routine did not terminate until this point, it is always an update.
	if len(doguDiff.NeededActions) > 0 {
		logger.Info("upgrade dogu")
		return useCase.doguRepo.Update(ctx, doguInstallation)
	}

	return nil
}

func getNoDowngradesExplanationTextForDogus() string {
	return fmt.Sprintf(noDowngradesExplanationTextFmt, "dogu", "dogus")
}
