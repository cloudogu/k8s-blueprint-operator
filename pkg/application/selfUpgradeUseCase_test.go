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

	UpgradeToV2ComponentDiff := domain.ComponentDiff{
		Name: blueprintOperatorName,
		Actual: domain.ComponentDiffState{
			Version: version1,
		},
		Expected: domain.ComponentDiffState{
			Version: version2,
		},
		NeededAction: domain.ActionUpgrade,
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
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{},
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			Version: version1,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		})

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("upgrade already done, but not completed", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseAwaitSelfUpgrade,
		}

		component := &ecosystem.ComponentInstallation{
			Version: version2,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		})

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("upgrade already completed and verified", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseSelfUpgradeCompleted,
		}

		component := &ecosystem.ComponentInstallation{
			Version: version2,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		})

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("apply upgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			Version: version1,
		}

		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // but usually cancel
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseAwaitSelfUpgrade, blueprint.Status)
		})
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, component).Return(nil).Run(
			func(_ context.Context, _ domain.ComponentDiff, _ *ecosystem.ComponentInstallation) {
				cancelCtx()
			},
		)

		err := useCase.HandleSelfUpgrade(timeoutCtx, blueprintId)

		assert.NoError(t, err)
	})

	t.Run("apply upgrade with missing component cr", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		timeoutCtx, cancelCtx := context.WithTimeout(testCtx, time.Second) // but usually cancel
		defer cancelCtx()

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(nil, domainservice.NewNotFoundError(assert.AnError, "test-error"))
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Run(func(ctx context.Context, blueprintSpec *domain.BlueprintSpec) {
			require.Equal(t, domain.StatusPhaseAwaitSelfUpgrade, blueprint.Status)
		})
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
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

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
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

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

	t.Run("cannot save blueprint", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			Version: version1,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \""+blueprintId+"\" while possibly self upgrading the operator")
	})

	t.Run("cannot apply self upgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
			Status:    domain.StatusPhaseBlueprintApplicationPreProcessed,
		}

		component := &ecosystem.ComponentInstallation{
			Version: version1,
		}

		blueprintRepo.EXPECT().GetById(mock.Anything, blueprintId).Return(blueprint, nil)
		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil)
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, component).Return(internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprintId)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "an error occurred while applying the self-upgrade to the ecosystem")
	})
}
