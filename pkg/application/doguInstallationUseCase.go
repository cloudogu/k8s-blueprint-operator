package application

import (
	"context"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"golang.org/x/exp/maps"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguInstallationUseCase struct {
	blueprintSpecRepo  blueprintSpecRepository
	doguRepo           doguInstallationRepository
	waitConfigProvider healthWaitConfigProvider
	doguConfigRepo     doguConfigRepository
	globalConfigRepo   globalConfigRepository
}

func NewDoguInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguRepo domainservice.DoguInstallationRepository,
	waitConfigProvider domainservice.HealthWaitConfigProvider,
	doguConfigRepo doguConfigRepository,
	globalConfigRepo globalConfigRepository,
) *DoguInstallationUseCase {
	return &DoguInstallationUseCase{
		blueprintSpecRepo:  blueprintSpecRepo,
		doguRepo:           doguRepo,
		waitConfigProvider: waitConfigProvider,
		doguConfigRepo:     doguConfigRepo,
		globalConfigRepo:   globalConfigRepo,
	}
}

func (useCase *DoguInstallationUseCase) CheckDoguHealth(ctx context.Context) (ecosystem.DoguHealthResult, error) {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.CheckDoguHealth")
	logger.V(2).Info("check dogu health...")
	installedDogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return ecosystem.DoguHealthResult{}, fmt.Errorf("cannot evaluate dogu health states: %w", err)
	}
	// accept experimental maps.Values as we can implement it ourselves in a minute
	return ecosystem.CalculateDoguHealthResult(maps.Values(installedDogus)), nil
}

func (useCase *DoguInstallationUseCase) CheckDogusUpToDate(ctx context.Context) ([]cescommons.SimpleName, error) {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.CheckDoguHealth")
	logger.V(2).Info("check if dogus are up to date...")
	installedDogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	globalConfig, err := useCase.globalConfigRepo.Get(ctx)
	if err != nil {
		return nil, err
	}
	globalConfigUpdateTime := globalConfig.LastUpdated

	var dogusNotUpToDate []cescommons.SimpleName

	for doguName, dogu := range installedDogus {
		versionUpToDate := dogu.IsVersionUpToDate()
		if !versionUpToDate {
			dogusNotUpToDate = append(dogusNotUpToDate, doguName)
			continue
		}

		doguConfig, err := useCase.doguConfigRepo.Get(ctx, doguName)
		if err != nil {
			return nil, err
		}
		doguConfigUpdateTime := doguConfig.LastUpdated
		configUpToDate := dogu.IsConfigUpToDate(globalConfigUpdateTime, doguConfigUpdateTime)
		if !configUpToDate {
			dogusNotUpToDate = append(dogusNotUpToDate, doguName)
			continue
		}
	}

	return dogusNotUpToDate, nil
}

// ApplyDoguStates applies the expected dogu state from the Blueprint to the ecosystem.
// Fail-fast here, so that the possible damage is as small as possible.
func (useCase *DoguInstallationUseCase) ApplyDoguStates(ctx context.Context, blueprint *domain.BlueprintSpec) error {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.ApplyDoguChanges")
	logger.V(2).Info("apply dogu states")
	// DoguDiff contains all installed dogus anyway (but some with action none) so we can load them all at once
	dogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot load dogu installations to apply dogu state: %w", err)
	}

	for _, doguDiff := range blueprint.StateDiff.DoguDiffs {
		err = useCase.applyDoguState(ctx, doguDiff, dogus[doguDiff.DoguName], blueprint.Config)
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
			newDogu := ecosystem.InstallDogu(
				cescommons.QualifiedName{Namespace: doguDiff.Expected.Namespace, SimpleName: doguDiff.DoguName},
				doguDiff.Expected.Version,
				doguDiff.Expected.MinVolumeSize,
				doguDiff.Expected.ReverseProxyConfig,
				doguDiff.Expected.AdditionalMounts,
			)
			return useCase.doguRepo.Create(ctx, newDogu)
		case domain.ActionUninstall:
			if doguInstallation == nil {
				return &domainservice.NotFoundError{Message: fmt.Sprintf("dogu %q not found", doguDiff.DoguName)}
			}
			logger.Info("uninstall dogu")
			return useCase.doguRepo.Delete(ctx, doguInstallation.Name.SimpleName)
		case domain.ActionUpgrade:
			doguInstallation.Upgrade(doguDiff.Expected.Version)
			continue
		case domain.ActionDowngrade:
			logger.Info("downgrade dogu")
			return fmt.Errorf(noDowngradesExplanationTextFmt, "dogu", "dogus")
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
		case domain.ActionUpdateAdditionalMounts:
			logger.Info("update additional mounts")
			doguInstallation.UpdateAdditionalMounts(doguDiff.Expected.AdditionalMounts)
			continue
		default:
			return fmt.Errorf("cannot perform unknown action %q for dogu %q", action, doguDiff.DoguName)
		}
	}

	// If this routine did not terminate until this point, it is always an update.
	if len(doguDiff.NeededActions) > 0 {
		logger.Info("upgrade dogu")
		// remove potential pause reconciliation flags here so that the dogu gets updates again
		doguInstallation.PauseReconciliation = false
		return useCase.doguRepo.Update(ctx, doguInstallation)
	}

	return nil
}
