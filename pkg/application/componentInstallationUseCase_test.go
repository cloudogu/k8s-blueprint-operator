package application

import (
	"context"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		healthConfigRepo := newMockHealthConfigProvider(t)

		// when
		useCase := NewComponentInstallationUseCase(blueprintSpecRepoMock, componentRepoMock, healthConfigRepo)

		// then
		assert.Equal(t, blueprintSpecRepoMock, useCase.blueprintSpecRepo)
		assert.Equal(t, componentRepoMock, useCase.componentRepo)
		assert.Same(t, healthConfigRepo, useCase.healthConfigProvider)
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
						Name:          componentName1,
						NeededActions: []domain.Action{},
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

		allComponents := map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
			componentName1: nil,
		}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)
		componentRepoMock.EXPECT().GetAll(testCtx).Return(allComponents, nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
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
		}

		// when
		err := sut.ApplyComponentStates(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"blueprint1\" to apply components")
	})

	t.Run("should return nil and do nothing with no component diffs", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)

		expectedBlueprintSpec := &domain.BlueprintSpec{}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
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
						Name:          componentName1,
						NeededActions: []domain.Action{"unknown"},
					},
				},
			},
		}

		allComponents := map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
			componentName1: nil,
		}

		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(expectedBlueprintSpec, nil)
		componentRepoMock.EXPECT().GetAll(testCtx).Return(allComponents, nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
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
			Name:          componentName1,
			NeededActions: []domain.Action{domain.ActionInstall},
			Expected: domain.ComponentDiffState{
				Namespace: testNamespace,
				Version:   semVer3212,
			},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name:            common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
			ExpectedVersion: semVer3212,
		}

		componentRepoMock.EXPECT().Create(testCtx, componentInstallation).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
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
			Name:          componentName1,
			NeededActions: []domain.Action{domain.ActionUninstall},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
		}

		componentRepoMock.EXPECT().Delete(testCtx, componentInstallation.Name.SimpleName).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
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
				Namespace: testNamespace,
				Version:   semVer3212,
			},
			NeededActions: []domain.Action{domain.ActionUpgrade},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name:            common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
			ExpectedVersion: semVer3212,
		}

		componentRepoMock.EXPECT().Update(testCtx, componentInstallation).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.NoError(t, err)
	})

	t.Run("should update component with multiple actions", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name: componentName1,
			Expected: domain.ComponentDiffState{
				Namespace: testNamespace,
				Version:   semVer3212,
				DeployConfig: map[string]interface{}{
					"deployNamespace": "longhorn-system",
					"overwriteConfig": map[string]string{
						"key": "value",
					},
				},
			},
			NeededActions: []domain.Action{domain.ActionUpgrade, domain.ActionUpdateComponentDeployConfig},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name:            common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
			ExpectedVersion: semVer3212,
			DeployConfig: map[string]interface{}{
				"deployNamespace": "longhorn-system",
				"overwriteConfig": map[string]string{
					"key": "value",
				},
			},
		}

		componentRepoMock.EXPECT().Update(testCtx, componentInstallation).Return(nil)

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
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
			Name:          componentName1,
			NeededActions: []domain.Action{domain.ActionDowngrade},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
		}

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
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
			Name:          componentName1,
			NeededActions: []domain.Action{domain.ActionSwitchComponentNamespace},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
		}

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, noDistributionNamespaceSwitchExplanationText)
	})

	t.Run("should return no error on empty actions in diff", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		componentRepoMock := newMockComponentInstallationRepository(t)

		componentDiff := domain.ComponentDiff{
			Name:          componentName1,
			NeededActions: []domain.Action{},
		}

		componentInstallation := &ecosystem.ComponentInstallation{
			Name: common.QualifiedComponentName{SimpleName: componentName1, Namespace: testNamespace},
		}

		sut := &ComponentInstallationUseCase{
			blueprintSpecRepo: blueprintSpecRepoMock,
			componentRepo:     componentRepoMock,
		}

		// when
		err := sut.applyComponentState(testCtx, componentDiff, componentInstallation)

		// then
		require.NoError(t, err)
	})
}

func TestComponentInstallationUseCase_CheckComponentHealth(t *testing.T) {
	type fields struct {
		componentRepoFn    func(t *testing.T) componentInstallationRepository
		healthConfigRepoFn func(t *testing.T) healthConfigProvider
	}
	tests := []struct {
		name    string
		fields  fields
		want    ecosystem.ComponentHealthResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get installed components",
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					componentMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					return newMockHealthConfigProvider(t)
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "cannot retrieve installed components", i)
			},
		},
		{
			name: "should fail to get required components",
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					componentMock.EXPECT().GetAll(testCtx).
						Return(map[common.SimpleComponentName]*ecosystem.ComponentInstallation{}, nil)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					healthConfigMock := newMockHealthConfigProvider(t)
					healthConfigMock.EXPECT().GetRequiredComponents(testCtx).
						Return(nil, assert.AnError)
					return healthConfigMock
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "cannot retrieve required components", i)
			},
		},
		{
			name: "should succeed",
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					componentMock.EXPECT().GetAll(testCtx).
						Return(map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
							"k8s-component-operator": {
								Name:   common.QualifiedComponentName{SimpleName: "k8s-component-operator", Namespace: testNamespace},
								Health: ecosystem.UnavailableHealthStatus},
						}, nil)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					healthConfigMock := newMockHealthConfigProvider(t)
					healthConfigMock.EXPECT().GetRequiredComponents(testCtx).
						Return([]ecosystem.RequiredComponent{{Name: "k8s-dogu-operator"}}, nil)
					return healthConfigMock
				},
			},
			want: ecosystem.ComponentHealthResult{
				ComponentsByStatus: map[ecosystem.HealthStatus][]common.SimpleComponentName{
					ecosystem.NotInstalledHealthStatus: {"k8s-dogu-operator"},
					ecosystem.UnavailableHealthStatus:  {"k8s-component-operator"},
				}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &ComponentInstallationUseCase{
				componentRepo:        tt.fields.componentRepoFn(t),
				healthConfigProvider: tt.fields.healthConfigRepoFn(t),
			}
			got, err := useCase.CheckComponentHealth(testCtx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComponentInstallationUseCase_WaitForHealthyComponents(t *testing.T) {
	type fields struct {
		componentRepoFn    func(t *testing.T) componentInstallationRepository
		healthConfigRepoFn func(t *testing.T) healthConfigProvider
	}
	tests := []struct {
		name    string
		ctx     context.Context
		fields  fields
		want    ecosystem.ComponentHealthResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get health check interval",
			ctx:  testCtx,
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					repoMock := newMockComponentInstallationRepository(t)
					return repoMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					providerMock := newMockHealthConfigProvider(t)
					providerMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{}, assert.AnError)
					return providerMock
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get health check interval", i)
			},
		},
		{
			name: "should fail to check component health",
			ctx:  testCtx,
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					repoMock := newMockComponentInstallationRepository(t)
					repoMock.EXPECT().GetAll(mock.Anything).Return(nil, assert.AnError)
					return repoMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					providerMock := newMockHealthConfigProvider(t)
					providerMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Interval: time.Second}, nil)
					return providerMock
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "stop waiting for component health", i)
			},
		},
		{
			name: "should fail after context cancellation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(testCtx)
				cancel()
				return ctx
			}(),
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					componentMock.EXPECT().GetAll(mock.Anything).
						Return(map[common.SimpleComponentName]*ecosystem.ComponentInstallation{}, nil)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					healthConfigMock := newMockHealthConfigProvider(t)
					healthConfigMock.EXPECT().GetWaitConfig(mock.Anything).Return(ecosystem.WaitConfig{Interval: 1}, nil)
					healthConfigMock.EXPECT().GetRequiredComponents(mock.Anything).
						Return([]ecosystem.RequiredComponent{{Name: "k8s-dogu-operator"}}, nil)
					return healthConfigMock
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "stop waiting for component health: context canceled", i)
			},
		},
		{
			name: "should be successful after retry",
			ctx:  testCtx,
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					unsuccessfulCall := componentMock.EXPECT().GetAll(mock.Anything).
						Return(map[common.SimpleComponentName]*ecosystem.ComponentInstallation{}, nil).Once()
					componentMock.EXPECT().GetAll(mock.Anything).
						Return(map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
							"k8s-dogu-operator": {

								Name: common.QualifiedComponentName{SimpleName: "k8s-dogu-operator", Namespace: testNamespace}, Health: ecosystem.AvailableHealthStatus,
							},
						}, nil).
						Once().NotBefore(unsuccessfulCall)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigProvider {
					healthConfigMock := newMockHealthConfigProvider(t)
					healthConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Interval: time.Second}, nil)
					healthConfigMock.EXPECT().GetRequiredComponents(mock.Anything).
						Return([]ecosystem.RequiredComponent{{Name: "k8s-dogu-operator"}}, nil)
					return healthConfigMock
				},
			},
			want: ecosystem.ComponentHealthResult{ComponentsByStatus: map[ecosystem.HealthStatus][]common.SimpleComponentName{
				ecosystem.AvailableHealthStatus: {"k8s-dogu-operator"},
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &ComponentInstallationUseCase{
				componentRepo:        tt.fields.componentRepoFn(t),
				healthConfigProvider: tt.fields.healthConfigRepoFn(t),
			}
			got, err := useCase.WaitForHealthyComponents(tt.ctx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
