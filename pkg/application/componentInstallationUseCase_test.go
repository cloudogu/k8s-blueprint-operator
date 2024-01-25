package application

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNewComponentInstallationUseCase(t *testing.T) {
	componentRepo := newMockComponentInstallationRepository(t)
	healthConfigRepo := newMockHealthConfigRepository(t)
	interval := 1 * time.Second
	componentUseCase := NewComponentInstallationUseCase(componentRepo, healthConfigRepo, interval)
	assert.Same(t, componentRepo, componentUseCase.componentRepo)
	assert.Same(t, healthConfigRepo, componentUseCase.healthConfigRepo)
	assert.Equal(t, interval, componentUseCase.healthCheckInterval)
}

func TestComponentInstallationUseCase_CheckComponentHealth(t *testing.T) {
	type fields struct {
		componentRepoFn    func(t *testing.T) componentInstallationRepository
		healthConfigRepoFn func(t *testing.T) healthConfigRepository
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
				healthConfigRepoFn: func(t *testing.T) healthConfigRepository {
					return newMockHealthConfigRepository(t)
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "cannot retrieve installed components", i)
			},
		},
		{
			name: "should fail to get health config",
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					componentMock.EXPECT().GetAll(testCtx).
						Return(map[string]*ecosystem.ComponentInstallation{}, nil)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigRepository {
					healthConfigMock := newMockHealthConfigRepository(t)
					healthConfigMock.EXPECT().Get(testCtx).
						Return(domain.HealthConfig{}, assert.AnError)
					return healthConfigMock
				},
			},
			want: ecosystem.ComponentHealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "cannot retrieve health config", i)
			},
		},
		{
			name: "should succeed",
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					componentMock := newMockComponentInstallationRepository(t)
					componentMock.EXPECT().GetAll(testCtx).
						Return(map[string]*ecosystem.ComponentInstallation{
							"k8s-component-operator": {Name: "k8s-component-operator",
								Health: ecosystem.UnavailableHealthStatus},
						}, nil)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigRepository {
					healthConfigMock := newMockHealthConfigRepository(t)
					healthConfigMock.EXPECT().Get(testCtx).
						Return(domain.HealthConfig{ComponentHealthConfig: domain.ComponentHealthConfig{
							RequiredComponents: []ecosystem.RequiredComponent{
								{Name: "k8s-dogu-operator"},
							},
						}}, nil)
					return healthConfigMock
				},
			},
			want: ecosystem.ComponentHealthResult{
				ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
					ecosystem.NotInstalledHealthStatus: {"k8s-dogu-operator"},
					ecosystem.UnavailableHealthStatus:  {"k8s-component-operator"},
				}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &ComponentInstallationUseCase{
				componentRepo:    tt.fields.componentRepoFn(t),
				healthConfigRepo: tt.fields.healthConfigRepoFn(t),
			}
			got, err := useCase.CheckComponentHealth(testCtx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComponentInstallationUseCase_WaitForHealthyComponents(t *testing.T) {
	type fields struct {
		componentRepoFn     func(t *testing.T) componentInstallationRepository
		healthConfigRepoFn  func(t *testing.T) healthConfigRepository
		healthCheckInterval time.Duration
	}
	tests := []struct {
		name    string
		ctx     context.Context
		fields  fields
		want    ecosystem.ComponentHealthResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to check component health",
			ctx:  testCtx,
			fields: fields{
				componentRepoFn: func(t *testing.T) componentInstallationRepository {
					repoMock := newMockComponentInstallationRepository(t)
					repoMock.EXPECT().GetAll(mock.Anything).Return(nil, assert.AnError)
					return repoMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigRepository {
					return newMockHealthConfigRepository(t)
				},
				healthCheckInterval: 1 * time.Second,
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
						Return(map[string]*ecosystem.ComponentInstallation{}, nil)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigRepository {
					healthConfigMock := newMockHealthConfigRepository(t)
					healthConfigMock.EXPECT().Get(mock.Anything).
						Return(domain.HealthConfig{ComponentHealthConfig: domain.ComponentHealthConfig{
							RequiredComponents: []ecosystem.RequiredComponent{
								{Name: "k8s-dogu-operator"},
							},
						}}, nil)
					return healthConfigMock
				},
				healthCheckInterval: 1,
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
						Return(map[string]*ecosystem.ComponentInstallation{}, nil).Once()
					componentMock.EXPECT().GetAll(mock.Anything).
						Return(map[string]*ecosystem.ComponentInstallation{
							"k8s-dogu-operator": {
								Name: "k8s-dogu-operator", Health: ecosystem.AvailableHealthStatus,
							},
						}, nil).
						Once().NotBefore(unsuccessfulCall)
					return componentMock
				},
				healthConfigRepoFn: func(t *testing.T) healthConfigRepository {
					healthConfigMock := newMockHealthConfigRepository(t)
					healthConfigMock.EXPECT().Get(mock.Anything).
						Return(domain.HealthConfig{ComponentHealthConfig: domain.ComponentHealthConfig{
							RequiredComponents: []ecosystem.RequiredComponent{
								{Name: "k8s-dogu-operator"},
							},
						}}, nil)
					return healthConfigMock
				},
				healthCheckInterval: 1,
			},
			want: ecosystem.ComponentHealthResult{ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
				ecosystem.AvailableHealthStatus: {"k8s-dogu-operator"},
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &ComponentInstallationUseCase{
				componentRepo:       tt.fields.componentRepoFn(t),
				healthConfigRepo:    tt.fields.healthConfigRepoFn(t),
				healthCheckInterval: tt.fields.healthCheckInterval,
			}
			got, err := useCase.WaitForHealthyComponents(tt.ctx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
