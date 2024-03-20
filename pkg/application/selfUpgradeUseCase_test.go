package application

import (
	"context"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSelfUpgradeUseCase_HandleSelfUpgrade(t *testing.T) {
	var blueprintOperatorName = common.SimpleComponentName("k8s-blueprint-operator")
	blueprintId := "myBlueprint"
	testCtx := context.TODO()
	version1, _ := semver.NewVersion("1.0")
	version2, _ := semver.NewVersion("2.0")
	internalTestError := domainservice.NewInternalError(assert.AnError, "internal error")
	waitConfig := ecosystem.WaitConfig{Interval: 5 * time.Second}

	UpgradeToV2ComponentDiff := domain.ComponentDiff{
		Name: blueprintOperatorName,
		Actual: domain.ComponentDiffState{
			Version: version1,
		},
		Expected: domain.ComponentDiffState{
			Version: version2,
		},
		NeededActions: []domain.Action{domain.ActionUpgrade},
	}
	upgradeToV2StateDiff := domain.StateDiff{
		ComponentDiffs: []domain.ComponentDiff{
			UpgradeToV2ComponentDiff,
		},
	}

	t.Run("nothing to do", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{},
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		})

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("apply upgrade until termination", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version1,
			ActualVersion:   version1,
		}

		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // but usually cancel
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).
			Return(nil).
			Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
				assert.Equal(t, domain.StatusPhaseAwaitSelfUpgrade, blueprint.Status)
			}).Once() // only once as the operator will terminate and will set status completed later.

		componentRepo.EXPECT().GetByName(timeoutCtx, blueprintOperatorName).Return(component, nil).Once()
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, component).Return(nil).
			Run(func(_ context.Context, _ domain.ComponentDiff, _ *ecosystem.ComponentInstallation) {
				// check that the status is set beforehand, as we cannot guarantee that we can set it afterward before termination
				assert.Equal(t, domain.StatusPhaseAwaitSelfUpgrade, blueprint.Status)
				cancelCtx()
			})

		err := useCase.HandleSelfUpgrade(timeoutCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("verify installation after termination", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseAwaitSelfUpgrade,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version2,
		}

		blueprintRepo.EXPECT().GetById(testCtx, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(component, nil).Once() // no reload for actual version check
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		})

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("await installation confirmation after termination", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseAwaitSelfUpgrade,
		}

		component1 := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version1,
		}
		component2 := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version2,
		}
		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // no timeout should happen as
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(timeoutCtx, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(timeoutCtx, blueprintOperatorName).Return(component1, nil).Once()
		componentRepo.EXPECT().GetByName(timeoutCtx, blueprintOperatorName).Return(component2, nil).Once()
		configProvider.EXPECT().GetWaitConfig(timeoutCtx).Return(waitConfig, nil)
		blueprintRepo.EXPECT().Update(timeoutCtx, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		})

		err := useCase.HandleSelfUpgrade(timeoutCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("apply upgrade with missing component cr", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // but usually cancel
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(nil, domainservice.NewNotFoundError(assert.AnError, "test-error"))
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil)
		var nilComponent *ecosystem.ComponentInstallation
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, nilComponent).Return(nil).Run(
			func(_ context.Context, _ domain.ComponentDiff, _ *ecosystem.ComponentInstallation) {
				cancelCtx()
			},
		)

		err := useCase.HandleSelfUpgrade(timeoutCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("could not load blueprint", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, assert.AnError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \""+blueprintId+"\" to possibly self upgrade the operator")
	})

	t.Run("cannot load component", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(nil, internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "cannot load component installation for \""+string(blueprintOperatorName)+"\" from ecosystem")
	})

	t.Run("cannot save blueprint in doSelfUpgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			Id:        blueprintId,
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version1,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "cannot persist blueprint spec \""+blueprintId+"\" to mark it waiting for self upgrade")
	})

	t.Run("cannot apply self upgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version1,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil)
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, component).Return(internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "an error occurred while applying the self-upgrade to the ecosystem")
	})

	t.Run("cannot save blueprint to skip self upgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			Id:        blueprintId,
			StateDiff: domain.StateDiff{},
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(assert.AnError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \""+blueprintId+"\" to skip self upgrade")
	})

	t.Run("error awaiting version confirmation, cannot load component", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseAwaitSelfUpgrade,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version1,
		}
		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // no timeout should happen as
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(timeoutCtx, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(timeoutCtx, blueprintOperatorName).Return(component, nil).Once()
		componentRepo.EXPECT().GetByName(timeoutCtx, blueprintOperatorName).Return(nil, assert.AnError).Once()
		configProvider.EXPECT().GetWaitConfig(timeoutCtx).Return(waitConfig, nil)

		err := useCase.HandleSelfUpgrade(timeoutCtx, blueprintId)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while waiting for version confirmation")
		assert.ErrorContains(t, err, "could not reload component for version confirmation")
	})

	t.Run("error awaiting version confirmation, cannot load wait config", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseAwaitSelfUpgrade,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version1,
		}
		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // no timeout should happen as
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(timeoutCtx, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(timeoutCtx, blueprintOperatorName).Return(component, nil).Once()
		configProvider.EXPECT().GetWaitConfig(timeoutCtx).Return(waitConfig, assert.AnError)

		err := useCase.HandleSelfUpgrade(timeoutCtx, blueprintId)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not retrieve wait interval config for self upgrade")
	})

	t.Run("error saving blueprint after awaiting version confirmation", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		configProvider := newMockHealthConfigProvider(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName, configProvider)

		blueprint := &domain.BlueprintSpec{
			Id:        blueprintId,
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseAwaitSelfUpgrade,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version2,
		}

		blueprintRepo.EXPECT().GetById(testCtx, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(component, nil).Once() // no reload for actual version check
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"myBlueprint\" after self upgrading the operator")
	})
}
