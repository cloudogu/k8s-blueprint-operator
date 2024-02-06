package application

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	componentName1 = "operator1"
	testNamespace  = "k8s"
)

var (
	semVer3212, _ = semver.NewVersion("3.2.1-2")
)

func TestNewComponentInstallationUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		// when
		useCase := NewComponentInstallationUseCase(blueprintSpecRepoMock, componentRepoMock, time.Second)

		// then
		assert.Equal(t, blueprintSpecRepoMock, useCase.blueprintSpecRepo)
		assert.Equal(t, componentRepoMock, useCase.componentRepo)
		assert.Equal(t, time.Second, useCase.healthCheckInterval)
	})
}

func TestComponentInstallationUseCase_ApplyComponentStates(t *testing.T) {
	t.Run("success with no needed action", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		expectedBlueprintSpec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				ComponentDiffs: []domain.ComponentDiff{
					{
						Name:         componentName1,
						NeededAction: domain.ActionNone,
						Actual: domain.ComponentDiffState{
							Version: semVer3212,
						},
						Expected: domain.ComponentDiffState{
							Version: semVer3212,
						},
					},
				},
			},
		}

		allComponents := map[string]*ecosystem.ComponentInstallation{
			componentName1: nil,
		}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)
		componentRepoMock.EXPECT().GetAll(testCtx).Return(allComponents, nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.ApplyComponentStates(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error getting blueprint spec", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.ApplyComponentStates(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"blueprint1\" to install components")
	})

	t.Run("should return nil and do nothing with no component diffs", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)

		expectedBlueprintSpec := &domain.BlueprintSpec{}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.ApplyComponentStates(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error getting all components", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		expectedBlueprintSpec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				ComponentDiffs: []domain.ComponentDiff{
					{},
				},
			},
		}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)
		componentRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.ApplyComponentStates(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load component installations to apply component state")
	})

	t.Run("should return error with unknown action", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		expectedBlueprintSpec := &domain.BlueprintSpec{
			StateDiff: domain.StateDiff{
				ComponentDiffs: []domain.ComponentDiff{
					{
						Name:         componentName1,
						NeededAction: "unknown",
					},
				},
			},
		}

		allComponents := map[string]*ecosystem.ComponentInstallation{
			componentName1: nil,
		}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)
		componentRepoMock.EXPECT().GetAll(testCtx).Return(allComponents, nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.ApplyComponentStates(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot perform unknown action \"unknown\"")
	})
}

func TestComponentInstallationUseCase_applyComponentState(t *testing.T) {
	t.Run("should create component on action install", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name:         componentName1,
			NeededAction: domain.ActionInstall,
			Expected: domain.ComponentDiffState{
				DistributionNamespace: testNamespace,
				Version:               semVer3212,
			},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name:                  componentName1,
			DistributionNamespace: testNamespace,
			Version:               semVer3212,
		}

		componentRepoMock.EXPECT().Create(testCtx, componentInstallation).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should delete component on action uninstall", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name:         componentName1,
			NeededAction: domain.ActionUninstall,
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: componentName1,
		}

		componentRepoMock.EXPECT().Delete(testCtx, componentInstallation.Name).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should update component on action upgrade", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name: componentName1,
			Expected: domain.ComponentDiffState{
				DistributionNamespace: testNamespace,
				Version:               semVer3212,
			},
			NeededAction: domain.ActionUpgrade,
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name:                  componentName1,
			Version:               semVer3212,
			DistributionNamespace: testNamespace,
		}

		componentRepoMock.EXPECT().Update(testCtx, componentInstallation).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on action downgrade", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name:         componentName1,
			NeededAction: domain.ActionDowngrade,
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: componentName1,
		}

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, getNoDowngradesExplanationTextForComponents())
	})

	t.Run("should return error on action distribution namespace switch", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name:         componentName1,
			NeededAction: domain.ActionSwitchComponentDistributionNamespace,
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: componentName1,
		}

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, noDistributionNamespaceSwitchExplanationText)
	})

	t.Run("should return error on action deploy namespace switch", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name:         componentName1,
			NeededAction: domain.ActionSwitchComponentDeployNamespace,
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: componentName1,
		}

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
			// TODO
			// healthCheckInterval:
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, noDeployNamespaceSwitchExplanationText)
	})
}
