package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguRestartUseCase struct {
	doguInstallationRepository doguInstallationRepository
	blueprintSpecRepo          blueprintSpecRepository
	restartRepository          domainservice.DoguRestartRepository
}

func NewDoguRestartUseCase(doguInstallationRepository doguInstallationRepository, blueprintSpecRepo blueprintSpecRepository, restartRepository doguRestartRepository) *DoguRestartUseCase {
	return &DoguRestartUseCase{doguInstallationRepository: doguInstallationRepository, blueprintSpecRepo: blueprintSpecRepo, restartRepository: restartRepository}
}

func (useCase *DoguRestartUseCase) TriggerDoguRestarts(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguRestartUseCase.TriggerDoguRestarts")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("could not get blueprint spec by id: %q", err)
	}

	logger.Info("searching for Dogus that need a restart...")
	dogusThatNeedARestart := []common.SimpleDoguName{}
	allDogusNeedARestart := false

	for _, globalConfigDiff := range blueprintSpec.StateDiff.GlobalConfigDiffs {
		if globalConfigDiff.NeededAction != domain.ActionNone {
			logger.Info("global config has changed, we need to restart all Dogus")
			allDogusNeedARestart = true
			break
		}
	}

	if allDogusNeedARestart {
		logger.Info("restarting all Dogus...")
		installedDogus, getInstalledDogusError := useCase.doguInstallationRepository.GetAll(ctx)
		if getInstalledDogusError != nil {
			return fmt.Errorf("could not get all installed Dogus: %q", getInstalledDogusError)
		}
		installedDogusQualifiedNames := []common.SimpleDoguName{}
		for _, installation := range installedDogus {
			installedDogusQualifiedNames = append(installedDogusQualifiedNames, installation.Name.SimpleName)
		}
		restartAllError := useCase.restartRepository.RestartAll(ctx, installedDogusQualifiedNames)
		if restartAllError != nil {
			logger.Error(restartAllError, "could not restart all Dogus")
		}
	} else {
		dogusThatNeedARestart = getDogusThatNeedARestart(blueprintSpec)
		if len(dogusThatNeedARestart) > 0 {
			logger.Info("restarting Dogus...")
			restartError := useCase.restartRepository.RestartAll(ctx, dogusThatNeedARestart)
			if restartError != nil {
				logger.Error(restartError, "could not restart Dogus")
			}
		} else {
			logger.Info("no Dogu restarts necessary")
		}
	}

	blueprintSpec.Status = domain.StatusPhaseRestartsTriggered
	useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	return nil
}

func getDogusThatNeedARestart(blueprintSpec *domain.BlueprintSpec) []common.SimpleDoguName {
	dogusThatNeedRestart := []common.SimpleDoguName{}
	dogusInEffectiveBlueprint := blueprintSpec.EffectiveBlueprint.Dogus
	for _, dogu := range dogusInEffectiveBlueprint {
		if blueprintSpec.StateDiff.DoguConfigDiffs[dogu.Name.SimpleName].HasChanges() {
			dogusThatNeedRestart = append(dogusThatNeedRestart, dogu.Name.SimpleName)
		}
	}
	return dogusThatNeedRestart
}
