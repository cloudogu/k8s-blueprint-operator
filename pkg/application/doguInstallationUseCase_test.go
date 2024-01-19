package application

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const blueprintId = "blueprint1"

var version3_2_1_1, _ = core.ParseVersion("3.2.1-1")
var version3_2_1_2, _ = core.ParseVersion("3.2.1-2")

const healthCheckInterval = 10 * time.Second

func TestDoguInstallationUseCase_applyDoguState(t *testing.T) {
	t.Run("action none", func(t *testing.T) {
		// given
		sut := NewDoguInstallationUseCase(nil, nil, healthCheckInterval)

		// when
		err := sut.applyDoguState(testCtx, domain.DoguDiff{
			DoguName: "postgresql",
			Actual: domain.DoguDiffState{
				Namespace:         "official",
				Version:           version3_2_1_1,
				InstallationState: domain.TargetStatePresent,
			},
			Expected: domain.DoguDiffState{
				Namespace:         "official",
				Version:           version3_2_1_1,
				InstallationState: domain.TargetStatePresent,
			},
			NeededAction: domain.ActionNone,
		}, &ecosystem.DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version3_2_1_1,
		}, domain.BlueprintConfiguration{})

		// then
		require.NoError(t, err)
	})

	t.Run("action install", func(t *testing.T) {
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Create(testCtx, ecosystem.InstallDogu("official", "postgresql", version3_2_1_1)).
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Actual: domain.DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: domain.TargetStateAbsent,
				},
				Expected: domain.DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: domain.TargetStatePresent,
				},
				NeededAction: domain.ActionInstall,
			},
			nil,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action uninstall", func(t *testing.T) {
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Delete(testCtx, "postgresql").
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName:     "postgresql",
				NeededAction: domain.ActionUninstall,
			},
			&ecosystem.DoguInstallation{
				Namespace: "official",
				Name:      "postgresql",
				Version:   version3_2_1_1,
			},
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
	})

	t.Run("action upgrade", func(t *testing.T) {
		dogu := &ecosystem.DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version3_2_1_1,
		}
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().
			Update(testCtx, dogu).
			Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Version: version3_2_1_2,
				},
				NeededAction: domain.ActionUpgrade,
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, version3_2_1_2, dogu.Version)
	})

	t.Run("action downgrade", func(t *testing.T) {

		dogu := &ecosystem.DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version3_2_1_2,
		}

		sut := NewDoguInstallationUseCase(nil, nil, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Version: version3_2_1_1,
				},
				NeededAction: domain.ActionDowngrade,
			},
			dogu,
			domain.BlueprintConfiguration{},
		)

		// then
		require.ErrorContains(t, err, noDowngradesExplanationText)
		assert.Equal(t, version3_2_1_2, dogu.Version)
	})

	t.Run("action SwitchNamespace not allowed", func(t *testing.T) {
		dogu := &ecosystem.DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version3_2_1_2,
		}

		sut := NewDoguInstallationUseCase(nil, nil, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededAction: domain.ActionSwitchDoguNamespace,
			},
			dogu,
			domain.BlueprintConfiguration{
				AllowDoguNamespaceSwitch: false,
			},
		)

		// then
		require.ErrorContains(t, err, "not allowed to switch dogu namespace")
	})

	t.Run("action SwitchNamespace allowed", func(t *testing.T) {
		dogu := &ecosystem.DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version3_2_1_2,
		}
		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().Update(testCtx, dogu).Return(nil)

		sut := NewDoguInstallationUseCase(nil, doguRepoMock, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededAction: domain.ActionSwitchDoguNamespace,
			},
			dogu,
			domain.BlueprintConfiguration{
				AllowDoguNamespaceSwitch: true,
			},
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, "premium", dogu.Namespace)
	})

	t.Run("unknown action", func(t *testing.T) {
		//given
		sut := NewDoguInstallationUseCase(nil, nil, healthCheckInterval)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededAction: "unknown",
			},
			nil,
			domain.BlueprintConfiguration{},
		)

		// then
		require.ErrorContains(t, err, "cannot perform unknown action \"unknown\"")
	})

}

func TestDoguInstallationUseCase_ApplyDoguStates(t *testing.T) {
	t.Run("cannot load blueprintSpec", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		doguRepoMock := newMockDoguInstallationRepository(t)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, healthCheckInterval)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("cannot load doguInstallations", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, healthCheckInterval)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load dogu installations")
	})

	t.Run("success", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:     "postgresql",
						NeededAction: domain.ActionNone,
					},
				},
			},
			Config: domain.BlueprintConfiguration{},
		}, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, healthCheckInterval)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				DoguDiffs: []domain.DoguDiff{
					{
						DoguName:     "postgresql",
						NeededAction: domain.ActionDowngrade,
					},
				},
			},
			Config: domain.BlueprintConfiguration{},
		}, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{
			"postgresql": {
				Namespace:     "official",
				Name:          "postgresql",
				Version:       version3_2_1_1,
				UpgradeConfig: ecosystem.UpgradeConfig{},
			},
		}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock, healthCheckInterval)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorContains(t, err, noDowngradesExplanationText)
		require.ErrorContains(t, err, "an error occurred while applying dogu state to the ecosystem")
	})
}

func TestDoguInstallationUseCase_WaitForHealthyDogus(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		timedCtx, cancel := context.WithTimeout(testCtx, 10*time.Millisecond)
		defer cancel()
		doguRepoMock.EXPECT().GetAll(timedCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:   nil,
			doguRepo:            doguRepoMock,
			healthCheckInterval: 1 * time.Millisecond,
		}

		// when
		result, err := sut.WaitForHealthyDogus(timedCtx)

		// then
		require.NoError(t, err)
		assert.True(t, result.AllHealthy())
	})

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		timedCtx, cancel := context.WithTimeout(testCtx, 0*time.Millisecond)
		defer cancel()
		// return unhealthy result
		doguRepoMock.EXPECT().GetAll(timedCtx).Return(map[string]*ecosystem.DoguInstallation{
			"postgresql": {Health: ecosystem.DoguStatusInstalling},
		}, nil).Maybe()

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:   nil,
			doguRepo:            doguRepoMock,
			healthCheckInterval: 5 * time.Millisecond,
		}

		// when
		result, err := sut.WaitForHealthyDogus(timedCtx)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, ecosystem.DoguHealthResult{}, result)
	})

	t.Run("cannot load dogus", func(t *testing.T) {
		t.Parallel()
		// given
		doguRepoMock := newMockDoguInstallationRepository(t)
		timedCtx, cancel := context.WithTimeout(testCtx, 10*time.Millisecond)
		defer cancel()
		doguRepoMock.EXPECT().GetAll(timedCtx).Return(nil, assert.AnError).Maybe()

		sut := DoguInstallationUseCase{
			blueprintSpecRepo:   nil,
			doguRepo:            doguRepoMock,
			healthCheckInterval: 1 * time.Millisecond,
		}

		// when
		result, err := sut.WaitForHealthyDogus(timedCtx)

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, ecosystem.DoguHealthResult{}, result)
	})

}
