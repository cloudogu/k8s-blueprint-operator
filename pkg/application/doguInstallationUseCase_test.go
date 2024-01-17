package application

import (
	"errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const blueprintId = "blueprint1"

var version3_2_1_1, _ = core.ParseVersion("3.2.1-1")
var version3_2_1_2, _ = core.ParseVersion("3.2.1-2")

func TestDoguInstallationUseCase_CheckDoguHealth(t *testing.T) {
	t.Run("should fail to get blueprint spec", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		doguRepoMock := newMockDoguInstallationRepository(t)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"blueprint1\" to check dogu health")
	})
	t.Run("should fail to get dogus", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{}, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot evaluate dogu health states for blueprint spec \"blueprint1\"")
	})
	t.Run("should fail to update blueprint spec", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		blueprintSpecRepoMock.EXPECT().Update(testCtx, blueprintSpec).Return(assert.AnError)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the dogu health")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		blueprintSpecRepoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})
}

func TestDoguInstallationUseCase_applyDoguState(t *testing.T) {
	t.Run("action none", func(t *testing.T) {
		// given
		sut := NewDoguInstallationUseCase(nil, nil)

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

		sut := NewDoguInstallationUseCase(nil, doguRepoMock)

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

		sut := NewDoguInstallationUseCase(nil, doguRepoMock)

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

		sut := NewDoguInstallationUseCase(nil, doguRepoMock)

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

		sut := NewDoguInstallationUseCase(nil, nil)

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

		sut := NewDoguInstallationUseCase(nil, nil)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededAction: domain.ActionSwitchNamespace,
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

		sut := NewDoguInstallationUseCase(nil, doguRepoMock)

		// when
		err := sut.applyDoguState(
			testCtx,
			domain.DoguDiff{
				DoguName: "postgresql",
				Expected: domain.DoguDiffState{
					Namespace: "premium",
				},
				NeededAction: domain.ActionSwitchNamespace,
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
		sut := NewDoguInstallationUseCase(nil, nil)

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
		expectedError := errors.New("test-error")
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, expectedError)

		doguRepoMock := newMockDoguInstallationRepository(t)
		//doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("cannot load doguInstallations", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		expectedError := errors.New("test-error")
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, nil)
		//blueprintSpecRepoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(nil, expectedError)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorIs(t, err, expectedError)
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

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

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

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.ApplyDoguStates(testCtx, blueprintId)

		// then
		require.ErrorContains(t, err, noDowngradesExplanationText)
		require.ErrorContains(t, err, "an error occurred while applying dogu state to the ecosystem")
	})
}
