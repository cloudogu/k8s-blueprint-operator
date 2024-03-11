package application

import (
	"context"
	"errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
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
				Expected:     domain.GlobalConfigValueState{"initial", true},
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
				Expected:     domain.GlobalConfigValueState{"initial", true},
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
		installationRepository.EXPECT().GetAll(testContext).Return(map[common.SimpleDoguName]*ecosystem.DoguInstallation{}, errors.New("testerror"))

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.Equal(t, "could not get all installed Dogus: \"testerror\"", err.Error())
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
				Expected:     domain.GlobalConfigValueState{"initial", true},
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
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(errors.New("testerror"))
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.Equal(t, "testerror", err.Error())
	})

	t.Run("restart some dogus", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{"initial", true},
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
			Expected:     domain.DoguConfigValueState{"initial", true},
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
		restartRepository.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(errors.New("testerror"))
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.Equal(t, "testerror", err.Error())
	})

	t.Run("fail on error in blueprint spec update", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{"initial", true},
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
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(errors.New("testerror"))

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.Equal(t, "testerror", err.Error())
	})

	t.Run("fail on error when getting blueprint", func(t *testing.T) {
		// given
		doguConfigDiff := map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{}
		testDoguSimpleName := common.SimpleDoguName("testdogu1")
		doguConfigDiff[testDoguSimpleName] = domain.CombinedDoguConfigDiffs{DoguConfigDiff: domain.DoguConfigDiffs{{
			Key:          common.DoguConfigKey{DoguName: testDoguSimpleName, Key: "testkey"},
			Actual:       domain.DoguConfigValueState{Value: "changed", Exists: true},
			Expected:     domain.DoguConfigValueState{"initial", true},
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
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, errors.New("testerror"))
		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartRepository)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.Error(t, err)
		assert.Equal(t, "could not get blueprint spec by id: \"testerror\"", err.Error())
	})
}

func Test_getDogusThatNeedARestart(t *testing.T) {
	type args struct {
		blueprintSpec *domain.BlueprintSpec
	}

	testdogu1 := domain.Dogu{Name: common.QualifiedDoguName{Namespace: "testnamespace", SimpleName: "testdogu1"}}
	testBlueprint1 := domain.Blueprint{Dogus: []domain.Dogu{testdogu1}}
	testDoguConfigDiffsChanged := domain.CombinedDoguConfigDiffs{
		DoguConfigDiff: []domain.DoguConfigEntryDiff{{
			Actual:       domain.DoguConfigValueState{"", false},
			Expected:     domain.DoguConfigValueState{"testvalue", true},
			NeededAction: domain.ConfigActionSet,
		}},
		SensitiveDoguConfigDiff: nil,
	}
	testDoguConfigDiffsActionNone := domain.CombinedDoguConfigDiffs{
		DoguConfigDiff: []domain.DoguConfigEntryDiff{{
			NeededAction: domain.ActionNone,
		}},
		SensitiveDoguConfigDiff: nil,
	}

	testDoguConfigChangeDiffChanged := domain.StateDiff{
		DoguDiffs:         nil,
		ComponentDiffs:    nil,
		DoguConfigDiffs:   map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{testdogu1.Name.SimpleName: testDoguConfigDiffsChanged},
		GlobalConfigDiffs: nil,
	}
	testDoguConfigChangeDiffActionNone := domain.StateDiff{
		DoguDiffs:         nil,
		ComponentDiffs:    nil,
		DoguConfigDiffs:   map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{testdogu1.Name.SimpleName: testDoguConfigDiffsActionNone},
		GlobalConfigDiffs: nil,
	}

	tests := []struct {
		name string
		args args
		want []common.SimpleDoguName
	}{
		{
			name: "return nothing on empty blueprint",
			args: args{blueprintSpec: &domain.BlueprintSpec{}},
			want: []common.SimpleDoguName{},
		},
		{
			name: "return nothing on no config change",
			args: args{blueprintSpec: &domain.BlueprintSpec{Blueprint: testBlueprint1}},
			want: []common.SimpleDoguName{},
		},
		{
			name: "return dogu on dogu config change",
			args: args{
				blueprintSpec: &domain.BlueprintSpec{
					Blueprint:          testBlueprint1,
					EffectiveBlueprint: domain.EffectiveBlueprint(testBlueprint1),
					StateDiff:          testDoguConfigChangeDiffChanged,
				},
			},
			want: []common.SimpleDoguName{testdogu1.Name.SimpleName},
		},
		{
			name: "return nothing on dogu config unchanged",
			args: args{
				blueprintSpec: &domain.BlueprintSpec{
					Blueprint:          testBlueprint1,
					EffectiveBlueprint: domain.EffectiveBlueprint(testBlueprint1),
					StateDiff:          testDoguConfigChangeDiffActionNone,
				},
			},
			want: []common.SimpleDoguName{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getDogusThatNeedARestart(tt.args.blueprintSpec), "getDogusThatNeedARestart(%v)", tt.args.blueprintSpec)
		})
	}
}

func Test_checkForAllDoguRestart(t *testing.T) {
	type args struct {
		blueprintSpec *domain.BlueprintSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false on empty input",
			args: args{blueprintSpec: &domain.BlueprintSpec{}},
			want: false,
		},
		{
			name: "true on GlobalConfigDiff with action",
			args: args{
				blueprintSpec: &domain.BlueprintSpec{
					StateDiff: domain.StateDiff{
						GlobalConfigDiffs: []domain.GlobalConfigEntryDiff{
							{
								Key:          "testkey",
								Actual:       domain.GlobalConfigValueState{Value: "changed", Exists: true},
								Expected:     domain.GlobalConfigValueState{"initial", true},
								NeededAction: domain.ConfigActionSet,
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "false on GlobalConfigDiff without action",
			args: args{
				blueprintSpec: &domain.BlueprintSpec{
					StateDiff: domain.StateDiff{
						GlobalConfigDiffs: []domain.GlobalConfigEntryDiff{
							{
								Key:          "testkey",
								Actual:       domain.GlobalConfigValueState{Value: "changed", Exists: true},
								Expected:     domain.GlobalConfigValueState{"initial", true},
								NeededAction: domain.ActionNone,
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, checkForAllDoguRestart(tt.args.blueprintSpec), "checkForAllDoguRestart(%v)", tt.args.blueprintSpec)
		})
	}
}
