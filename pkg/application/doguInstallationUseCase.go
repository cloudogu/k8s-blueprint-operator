package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
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

	switch doguDiff.NeededAction {
	case domain.ActionNone:
		logger.Info("apply nothing for dogu")
		return nil
	case domain.ActionInstall:
		logger.Info("install dogu")
		newDogu := ecosystem.InstallDogu(common.QualifiedDoguName{
			Namespace: doguDiff.Expected.Namespace,
			Name:      doguDiff.DoguName,
		}, doguDiff.Expected.Version)
		return useCase.doguRepo.Create(ctx, newDogu)
	case domain.ActionUninstall:
		logger.Info("uninstall dogu")
		return useCase.doguRepo.Delete(ctx, doguInstallation.Name.Name)
	case domain.ActionUpgrade:
		logger.Info("upgrade dogu")
		doguInstallation.Upgrade(doguDiff.Expected.Version)
		return useCase.doguRepo.Update(ctx, doguInstallation)
	case domain.ActionDowngrade:
		logger.Info("downgrade dogu")
		return fmt.Errorf(getNoDowngradesExplanationTextForDogus())
	case domain.ActionSwitchDoguNamespace:
		logger.Info("do namespace switch for dogu")
		err := doguInstallation.SwitchNamespace(
			doguDiff.Expected.Namespace,
			doguDiff.Expected.Version,
			blueprintConfig.AllowDoguNamespaceSwitch,
		)
		if err != nil {
			return err
		}
		return useCase.doguRepo.Update(ctx, doguInstallation)
	default:
		return fmt.Errorf("cannot perform unknown action %q", doguDiff.NeededAction)
	}
}

func getNoDowngradesExplanationTextForDogus() string {
	return fmt.Sprintf(noDowngradesExplanationTextFmt, "dogu", "dogus")
}
