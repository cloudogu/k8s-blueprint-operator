package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguRestartUseCase struct {
	doguInstallationRepository doguInstallationRepository
	blueprintSpecRepo          blueprintSpecRepository
	doguRestartAdapter         doguRestartAdapter
}

func NewDoguRestartUseCase(doguInstallationRepository doguInstallationRepository, blueprintSpecRepo blueprintSpecRepository, doguRestartAdapter doguRestartAdapter) *DoguRestartUseCase {
	return &DoguRestartUseCase{doguInstallationRepository: doguInstallationRepository, blueprintSpecRepo: blueprintSpecRepo, doguRestartAdapter: doguRestartAdapter}
}

func (useCase *DoguRestartUseCase) TriggerDoguRestarts(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguRestartUseCase.TriggerDoguRestarts")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("could not get blueprint spec by id: %q", err)
	}

	logger.Info("searching for Dogus that need a restart...")
	dogusThatNeedARestart := []common.QualifiedDoguName{}
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
		installedDogusQualifiedNames := []common.QualifiedDoguName{}
		for _, installation := range installedDogus {
			installedDogusQualifiedNames = append(installedDogusQualifiedNames, installation.Name)
		}
		restartAllError := useCase.doguRestartAdapter.RestartAll(ctx, installedDogusQualifiedNames)
		if restartAllError != nil {
			logger.Error(restartAllError, "could not restart all Dogus")
		}
	} else {
		dogusThatNeedARestart = getDogusThatNeedARestart(blueprintSpec)
		if len(dogusThatNeedARestart) > 0 {
			logger.Info("restarting Dogus...")
			restartError := useCase.doguRestartAdapter.RestartAll(ctx, dogusThatNeedARestart)
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

func getDogusThatNeedARestart(blueprintSpec *domain.BlueprintSpec) []common.QualifiedDoguName {
	dogusThatNeedRestart := []common.QualifiedDoguName{}
	dogusInEffectiveBlueprint := blueprintSpec.EffectiveBlueprint.Dogus
	for _, dogu := range dogusInEffectiveBlueprint {
		if blueprintSpec.StateDiff.DoguConfigDiffs[dogu.Name.SimpleName].HasChanges() {
			// only restart Dogu if config changed and action is none
			// if action is not none, the Dogu is already restarted
			for _, diff := range blueprintSpec.StateDiff.DoguDiffs {
				if diff.DoguName == dogu.Name.SimpleName && diff.NeededAction == domain.ActionNone {
					dogusThatNeedRestart = append(dogusThatNeedRestart, dogu.Name)
					break
				}
			}
		}
	}
	return dogusThatNeedRestart
}
