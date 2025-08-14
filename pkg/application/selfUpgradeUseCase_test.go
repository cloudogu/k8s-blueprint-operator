package application

import (
	"context"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/meta"
	"testing"
)

func TestSelfUpgradeUseCase_HandleSelfUpgrade(t *testing.T) {
	var blueprintOperatorName = common.SimpleComponentName("k8s-blueprint-operator")
	blueprintId := "myBlueprint"
	testCtx := context.TODO()
	version1, _ := semver.NewVersion("1.0")
	version2, _ := semver.NewVersion("2.0")
	internalTestError := domainservice.NewInternalError(assert.AnError, "internal error")

	NoActionV2ComponentDiff := domain.ComponentDiff{
		Name: blueprintOperatorName,
		Actual: domain.ComponentDiffState{
			Version: version2,
		},
		Expected: domain.ComponentDiffState{
			Version: version2,
		},
		NeededActions: []domain.Action{},
	}
	NoActionV2StateDiff := domain.StateDiff{
		ComponentDiffs: []domain.ComponentDiff{
			NoActionV2ComponentDiff,
		},
	}
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
	InstallV2ComponentDiff := domain.ComponentDiff{
		Name:   blueprintOperatorName,
		Actual: domain.ComponentDiffState{},
		Expected: domain.ComponentDiffState{
			Version: version2,
		},
		NeededActions: []domain.Action{domain.ActionInstall},
	}
	InstallV2StateDiff := domain.StateDiff{
		ComponentDiffs: []domain.ComponentDiff{
			InstallV2ComponentDiff,
		},
	}

	t.Run("apply upgrade and trigger reconcile", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff:  upgradeToV2StateDiff,
			Conditions: &[]domain.Condition{},
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version1,
			ActualVersion:   version1,
		}

		// only once as the operator will terminate and will set status completed later.
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil).Once()
		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(component, nil).Once()
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, component).Return(nil)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		var awaitError *domain.AwaitSelfUpgradeError
		assert.ErrorAs(t, err, &awaitError)
		assert.ErrorContains(t, err, awaitSelfUpgradeErrorMsg)
		assert.True(t, meta.IsStatusConditionFalse(*blueprint.Conditions, domain.ConditionSelfUpgradeCompleted))
	})

	t.Run("apply upgrade with missing component cr", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff:  InstallV2StateDiff,
			Conditions: &[]domain.Condition{},
		}
		var nilComponent *ecosystem.ComponentInstallation

		// only once as the operator will terminate and will set status completed later.
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil).Once()
		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(nil, domainservice.NewNotFoundError(assert.AnError, "test-error")).Once()
		componentUseCase.EXPECT().applyComponentState(testCtx, InstallV2ComponentDiff, nilComponent).Return(nil)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		var awaitError *domain.AwaitSelfUpgradeError
		assert.ErrorAs(t, err, &awaitError)
		assert.ErrorContains(t, err, awaitSelfUpgradeErrorMsg)
		assert.True(t, meta.IsStatusConditionFalse(*blueprint.Conditions, domain.ConditionSelfUpgradeCompleted))
	})

	t.Run("check if self-upgrade is done -> not yet", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff:  NoActionV2StateDiff,
			Conditions: &[]domain.Condition{},
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version1,
		}

		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(component, nil).Once()
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		var awaitError *domain.AwaitSelfUpgradeError
		assert.ErrorAs(t, err, &awaitError)
		assert.ErrorContains(t, err, awaitSelfUpgradeErrorMsg)
		assert.True(t, meta.IsStatusConditionFalse(*blueprint.Conditions, domain.ConditionSelfUpgradeCompleted))
	})

	t.Run("check if self-upgrade is done -> done", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff:  NoActionV2StateDiff,
			Conditions: &[]domain.Condition{},
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version2,
		}

		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(component, nil).Once()
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		assert.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, domain.ConditionSelfUpgradeCompleted))
	})

	t.Run("cannot load component", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
		}

		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(nil, internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "cannot load component installation for \""+string(blueprintOperatorName)+"\" from ecosystem")
	})

	t.Run("cannot save blueprint in doSelfUpgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			Id:        blueprintId,
			StateDiff: upgradeToV2StateDiff,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version1,
		}

		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "cannot persist blueprint spec \""+blueprintId+"\" to mark it waiting for self upgrade")
	})

	t.Run("cannot apply self upgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			StateDiff: upgradeToV2StateDiff,
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version1,
		}

		componentRepo.EXPECT().GetByName(mock.Anything, blueprintOperatorName).Return(component, nil)
		blueprintRepo.EXPECT().Update(mock.Anything, blueprint).Return(nil)
		componentUseCase.EXPECT().applyComponentState(mock.Anything, UpgradeToV2ComponentDiff, component).Return(internalTestError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		assert.ErrorIs(t, err, internalTestError)
		assert.ErrorContains(t, err, "an error occurred while applying the self-upgrade to the ecosystem")
	})

	t.Run("cannot save blueprint to complete self upgrade", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		componentRepo := newMockComponentInstallationRepository(t)
		componentUseCase := newMockComponentInstallationUseCase(t)
		useCase := NewSelfUpgradeUseCase(blueprintRepo, componentRepo, componentUseCase, blueprintOperatorName)

		blueprint := &domain.BlueprintSpec{
			Id:         blueprintId,
			StateDiff:  NoActionV2StateDiff,
			Conditions: &[]domain.Condition{},
		}

		component := &ecosystem.ComponentInstallation{
			ExpectedVersion: version2,
			ActualVersion:   version2,
		}

		componentRepo.EXPECT().GetByName(testCtx, blueprintOperatorName).Return(component, nil).Once() // no reload for actual version check
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		err := useCase.HandleSelfUpgrade(testCtx, blueprint)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"myBlueprint\" after self upgrading the operator")
	})
}
