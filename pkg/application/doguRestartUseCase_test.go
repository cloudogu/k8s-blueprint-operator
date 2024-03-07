package application

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
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
		restartAdapter := newMockDoguRestartAdapter(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(nil)

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartAdapter)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.NoError(t, err)
	})

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
		restartAdapter := newMockDoguRestartAdapter(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(nil)

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartAdapter)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("dogu restarts triggered on blueprint with non-empty state diff", func(t *testing.T) {
		// given
		testContext := context.Background()
		testStateDiff := domain.StateDiff{
			DoguDiffs:         domain.DoguDiffs{},
			ComponentDiffs:    domain.ComponentDiffs{},
			DoguConfigDiffs:   map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs{},
			GlobalConfigDiffs: domain.GlobalConfigDiffs{{"testkey", domain.GlobalConfigValueState{"changed", true}, domain.GlobalConfigValueState{"initial", true}, domain.ConfigActionSet}},
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
		installedDogu := ecosystem.DoguInstallation{Name: common.QualifiedDoguName{Namespace: "testing", SimpleName: testDoguSimpleName}, Version: core.Version{Raw: "1.0.0-1", Major: 1, Extra: 1}, Status: "installed", Health: ecosystem.AvailableHealthStatus, UpgradeConfig: ecosystem.UpgradeConfig{AllowNamespaceSwitch: false}, PersistenceContext: nil}
		installedDogus := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			testDoguSimpleName: &installedDogu,
		}
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		installationRepository := newMockDoguInstallationRepository(t)
		blueprintSpecRepo := newMockBlueprintSpecRepository(t)
		restartAdapter := newMockDoguRestartAdapter(t)
		blueprintSpecRepo.EXPECT().GetById(testContext, testBlueprintId).Return(&testBlueprint, nil)
		installationRepository.EXPECT().GetAll(testContext).Return(installedDogus, nil)
		restartAdapter.EXPECT().RestartAll(testContext, dogusThatNeedARestart).Return(nil)
		blueprintSpecRepo.EXPECT().Update(testContext, &testBlueprint).Return(nil)

		restartUseCase := NewDoguRestartUseCase(installationRepository, blueprintSpecRepo, restartAdapter)

		// when
		err := restartUseCase.TriggerDoguRestarts(testContext, testBlueprintId)

		// then
		require.NoError(t, err)
	})
}
