package application

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDoguRestartUseCase_TriggerDoguRestarts(t *testing.T) {
	t.Run("no dogu restarts triggered on blueprint with empty state diff", func(t *testing.T) {
		// given
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:         domain.DoguDiffs{},
			ComponentDiffs:    domain.ComponentDiffs{},
			DoguConfigDiffs:   map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{},
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(nil)

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("dogu restarts triggered on blueprint with non-empty state diff", func(t *testing.T) {
		// given
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:       domain.DoguDiffs{},
			ComponentDiffs:  domain.ComponentDiffs{},
			DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{{
				Key:          "testkey",
				Actual:       domain.GlobalConfigValueState{Value: "changed", Exists: true},
				Expected:     domain.GlobalConfigValueState{Value: "initial", Exists: true},
				NeededAction: domain.ConfigActionSet,
			}},
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		installedDogu := ecosystem.DoguInstallation{
			Name:               common.QualifiedDoguName{Namespace: "testing", SimpleName: testDoguSimpleName},
			Version:            core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1},
			Status:             "installed",
			Health:             ecosystem.AvailableHealthStatus,
			UpgradeConfig:      ecosystem.UpgradeConfig{AllowNamespaceSwitch: false},
			PersistenceContext: nil,
		}
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			testDoguSimpleName: &installedDogu,
		}
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		installationRepository.EXPECT().GetAll(testContext).Return(installedDogus, nil)
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(nil)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(nil)

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("fail on get all dogus from repository error", func(t *testing.T) {
		// given
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:       domain.DoguDiffs{},
			ComponentDiffs:  domain.ComponentDiffs{},
			DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{{
				Key:          "testkey",
				Actual:       domain.GlobalConfigValueState{Value: "changed", Exists: true},
				Expected:     domain.GlobalConfigValueState{Value: "initial", Exists: true},
				NeededAction: domain.ConfigActionSet,
			}},
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		installationRepository.EXPECT().GetAll(testContext).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, assert.AnError)

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not restart all installed Dogus: could not get all installed Dogus:")
	})

	t.Run("fail on repository restart all error", func(t *testing.T) {
		// given
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:       domain.DoguDiffs{},
			ComponentDiffs:  domain.ComponentDiffs{},
			DoguConfigDiffs: map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{{
				Key:          "testkey",
				Actual:       domain.GlobalConfigValueState{Value: "changed", Exists: true},
				Expected:     domain.GlobalConfigValueState{Value: "initial", Exists: true},
				NeededAction: domain.ConfigActionSet,
			}},
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		installedDogu := ecosystem.DoguInstallation{
			Name:               common.QualifiedDoguName{Namespace: "testing", SimpleName: testDoguSimpleName},
			Version:            core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1},
			Status:             "installed",
			Health:             ecosystem.AvailableHealthStatus,
			UpgradeConfig:      ecosystem.UpgradeConfig{AllowNamespaceSwitch: false},
			PersistenceContext: nil,
		}
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			testDoguSimpleName: &installedDogu,
		}
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		installationRepository.EXPECT().GetAll(testContext).Return(installedDogus, nil)
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(assert.AnError)
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not restart all installed Dogus")
	})

	t.Run("restart some dogus", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{Value: "initial", Exists: true},
			NeededAction: domain.ConfigActionSet}},
		}
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:         domain.DoguDiffs{},
			ComponentDiffs:    domain.ComponentDiffs{},
			DoguConfigDiffs:   doguConfigDiff,
			GlobalConfigDiffs: domain.GlobalConfigDiffs{},
		}
		testDogu := domain.Dogu{
			Name:        common.QualifiedDoguName{SimpleName: testDoguSimpleName, Namespace: "testing"},
			Version:     core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1},
			TargetState: 0,
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.Dogu{testDogu}},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(nil)
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(nil)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("fail on dogu restart for some dogus", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{Value: "initial", Exists: true},
			NeededAction: domain.ConfigActionSet}},
		}
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:         domain.DoguDiffs{},
			ComponentDiffs:    domain.ComponentDiffs{},
			DoguConfigDiffs:   doguConfigDiff,
			GlobalConfigDiffs: domain.GlobalConfigDiffs{},
		}
		testDogu := domain.Dogu{
			Name:        common.QualifiedDoguName{SimpleName: testDoguSimpleName, Namespace: "testing"},
			Version:     core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1},
			TargetState: 0,
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.Dogu{testDogu}},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(assert.AnError)
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not restart Dogus")
	})

	t.Run("fail on error in blueprint spec update", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{Value: "initial", Exists: true},
			NeededAction: domain.ConfigActionSet}},
		}
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:         domain.DoguDiffs{},
			ComponentDiffs:    domain.ComponentDiffs{},
			DoguConfigDiffs:   doguConfigDiff,
			GlobalConfigDiffs: domain.GlobalConfigDiffs{},
		}
		testDogu := domain.Dogu{
			Name:        common.QualifiedDoguName{SimpleName: testDoguSimpleName, Namespace: "testing"},
			Version:     core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1},
			TargetState: 0,
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.Dogu{testDogu}},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(nil)
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(assert.AnError)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not update blueprint spec")
	})

	t.Run("fail on error when getting blueprint", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{Value: "initial", Exists: true},
			NeededAction: domain.ConfigActionSet}},
		}
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:         domain.DoguDiffs{},
			ComponentDiffs:    domain.ComponentDiffs{},
			DoguConfigDiffs:   doguConfigDiff,
			GlobalConfigDiffs: domain.GlobalConfigDiffs{},
		}
		testDogu := domain.Dogu{
			Name:        common.QualifiedDoguName{SimpleName: testDoguSimpleName, Namespace: "testing"},
			Version:     core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1},
			TargetState: 0,
		}
		testBlueprint := domain.BlueprintSpec{
			Id:                 testBlueprintId,
			Blueprint:          domain.Blueprint{},
			BlueprintMask:      domain.BlueprintMask{},
			EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.Dogu{testDogu}},
			StateDiff:          testStateDiff,
			Config:             domain.BlueprintConfiguration{},
			Status:             "",
			PersistenceContext: nil,
			Events:             nil,
		}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartRepository := newMockDoguRestartRepository(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, assert.AnError)
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not get blueprint spec by id")
	})
}
