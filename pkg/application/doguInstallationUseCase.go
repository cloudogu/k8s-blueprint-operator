package application

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguInstallationUseCase struct {
	blueprintSpecRepo domainservice.BlueprintSpecRepository
	doguRepo          domainservice.DoguInstallationRepository
}

func NewDoguInstallationUseCase(
	blueprintSpecRepo domainservice.BlueprintSpecRepository,
	doguRepo domainservice.DoguInstallationRepository,
) *DoguInstallationUseCase {
	return &DoguInstallationUseCase{
		blueprintSpecRepo: blueprintSpecRepo,
		doguRepo:          doguRepo,
	}
}

func (useCase *DoguInstallationUseCase) CheckDoguHealth(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.CheckDoguHealth").
		WithValues("blueprintId", blueprintId)

	logger.Info("getting blueprint spec for checking dogu health")
	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to check dogu health: %w", blueprintId, err)
	}

	installedDogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot evaluate dogu health states for blueprint spec %q: %w", blueprintId, err)
	}

	blueprintSpec.CheckDoguHealth(installedDogus)

	err = useCase.blueprintSpecRepo.Update(ctx, blueprintSpec)
	if err != nil {
		return fmt.Errorf("cannot save blueprint spec %q after checking the dogu health: %w", blueprintId, err)
	}

	return nil
}

// ApplyDoguStates applies the expected dogu state from the Blueprint to the ecosystem.
// Fail-fast here, so that the possible damage is as small as possible.
func (useCase *DoguInstallationUseCase) ApplyDoguStates(ctx context.Context, blueprintId string) error {
	logger := log.FromContext(ctx).WithName("DoguInstallationUseCase.ApplyDoguChanges").
		WithValues("blueprintId", blueprintId)

	blueprintSpec, err := useCase.blueprintSpecRepo.GetById(ctx, blueprintId)
	if err != nil {
		return fmt.Errorf("cannot load blueprint spec %q to install dogus: %w", blueprintId, err)
	}

	//DoguDiff contains all installed dogus anyway (but some with action none) so we can load them all at once
	dogus, err := useCase.doguRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cannot load dogu installations to apply dogu state: %w", err)
	}

	for _, doguDiff := range blueprintSpec.StateDiff.DoguDiffs {
		logger.Info("apply dogu state", "dogu", doguDiff.DoguName, "diff", doguDiff)
		err = useCase.applyDoguState(ctx, doguDiff, dogus[doguDiff.DoguName], blueprintSpec.Config)
		if err != nil {
			logger.Error(err, "an error occurred while applying dogu state to the ecosystem")
			return err
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
	switch doguDiff.NeededAction {
	case domain.ActionNone:
		return nil
	case domain.ActionInstall:
		newDogu := ecosystem.InstallDogu(doguDiff.Expected.Namespace, doguDiff.DoguName, doguDiff.Expected.Version)
		return useCase.doguRepo.Create(ctx, newDogu)
	case domain.ActionUninstall:
		return useCase.doguRepo.Delete(ctx, doguInstallation.Name)
	case domain.ActionUpgrade:
		doguInstallation.Upgrade(doguDiff.Expected.Version)
		return useCase.doguRepo.Update(ctx, doguInstallation)
	case domain.ActionDowngrade:
		return fmt.Errorf(noDowngradesExplanationText)
	case domain.ActionSwitchNamespace:
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

var noDowngradesExplanationText = "downgrades are not allowed as the data model of the dogu could have changed and " +
	"doing rollbacks to older models is not supported. " +
	"You can downgrade dogus by restoring a backup. " +
	"If you want an 'allow-downgrades' flag, issue a feature request"
