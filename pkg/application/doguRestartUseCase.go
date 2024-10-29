package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguRestartUseCase struct {
	doguInstallationRepository doguInstallationRepository
	blueprintSpecRepo          blueprintSpecRepository
	restartRepository          doguRestartRepository
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
	var dogusThatNeedARestart []common.SimpleDoguName

	if blueprintSpec.StateDiff.GlobalConfigDiffs.HasChanges() {
		logger.Info("restarting all installed Dogus...")
		err = useCase.restartAllInstalledDogus(ctx)
		if err != nil {
			return domainservice.NewInternalError(err, "could not restart all installed Dogus")
		}
	} else {
		dogusThatNeedARestart = blueprintSpec.GetDogusThatNeedARestart()
		if len(dogusThatNeedARestart) > 0 {
			logger.Info("restarting Dogus...")
			restartError := useCase.restartRepository.RestartAll(ctx, dogusThatNeedARestart)
			if restartError != nil {
				return domainservice.NewInternalError(err, "could not restart Dogus")
			}
		} else {
			logger.Info("no Dogu restarts necessary")
		}
	}

	blueprintSpec.Status = domain.StatusPhaseRestartsTriggered
	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return domainservice.NewInternalError(err, "could not update blueprint spec")
	}
	return nil
}

func (useCase *DoguRestartUseCase) restartAllInstalledDogus(ctx context.Context) error {
	installedDogus, getInstalledDogusError := useCase.doguInstallationRepository.GetAll(ctx)
	if getInstalledDogusError != nil {
		return domainservice.NewInternalError(getInstalledDogusError, "could not get all installed Dogus")
	}
	var installedDogusSimpleNames []common.SimpleDoguName
	for _, installation := range installedDogus {
		installedDogusSimpleNames = append(installedDogusSimpleNames, installation.Name.SimpleName)
	}
	return useCase.restartRepository.RestartAll(ctx, installedDogusSimpleNames)
}
