package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

func TestNewEcosystemHealthUseCase(t *testing.T) {
	doguUseCase := newMockDoguInstallationUseCase(t)
	componentUseCase := newMockComponentInstallationUseCase(t)
	waitConfigMock := newMockHealthWaitConfigProvider(t)
	useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, waitConfigMock)

	assert.Same(t, doguUseCase, useCase.doguUseCase)
	assert.Same(t, componentUseCase, useCase.componentUseCase)
	assert.Same(t, waitConfigMock, useCase.waitConfigProvider)
}

func TestEcosystemHealthUseCase_CheckEcosystemHealth(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
				ecosystem.AvailableHealthStatus:   {"postgresql"},
				ecosystem.UnavailableHealthStatus: {"postfix"},
				ecosystem.PendingHealthStatus:     {"scm"},
			},
		}
		componentHealth := ecosystem.ComponentHealthResult{
			ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
				ecosystem.NotInstalledHealthStatus: {"k8s-dogu-operator"},
				ecosystem.UnavailableHealthStatus:  {"k8s-etcd"},
				ecosystem.PendingHealthStatus:      {"k8s-service-discovery"},
				ecosystem.AvailableHealthStatus:    {"k8s-component-operator"},
			},
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		componentUseCase := newMockComponentInstallationUseCase(t)
		componentUseCase.EXPECT().CheckComponentHealth(testCtx).Return(componentHealth, nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, nil)

		health, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth, ComponentHealth: componentHealth}, health)
	})

	t.Run("ok, ignore dogu health", func(t *testing.T) {
		componentHealth := ecosystem.ComponentHealthResult{
			ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
				ecosystem.NotInstalledHealthStatus: {"k8s-dogu-operator"},
				ecosystem.UnavailableHealthStatus:  {"k8s-etcd"},
				ecosystem.PendingHealthStatus:      {"k8s-service-discovery"},
				ecosystem.AvailableHealthStatus:    {"k8s-component-operator"},
			},
		}
		componentUseCase := newMockComponentInstallationUseCase(t)
		componentUseCase.EXPECT().CheckComponentHealth(testCtx).Return(componentHealth, nil)
		useCase := NewEcosystemHealthUseCase(nil, componentUseCase, nil)

		health, err := useCase.CheckEcosystemHealth(testCtx, true, false)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{ComponentHealth: componentHealth}, health)
	})

	t.Run("ok, ignore component health", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
				ecosystem.AvailableHealthStatus:   {"postgresql"},
				ecosystem.UnavailableHealthStatus: {"postfix"},
				ecosystem.PendingHealthStatus:     {"scm"},
			},
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, nil, nil)

		health, err := useCase.CheckEcosystemHealth(testCtx, false, true)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth}, health)
	})

	t.Run("error checking dogu health", func(t *testing.T) {
		componentHealth := ecosystem.ComponentHealthResult{
			ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
				ecosystem.NotInstalledHealthStatus: {"k8s-dogu-operator"},
				ecosystem.UnavailableHealthStatus:  {"k8s-etcd"},
				ecosystem.PendingHealthStatus:      {"k8s-service-discovery"},
				ecosystem.AvailableHealthStatus:    {"k8s-component-operator"},
			},
		}
		componentUseCase := newMockComponentInstallationUseCase(t)
		componentUseCase.EXPECT().CheckComponentHealth(testCtx).Return(componentHealth, nil)
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(ecosystem.DoguHealthResult{}, assert.AnError)
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, nil)

		_, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("error checking component health", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
				ecosystem.AvailableHealthStatus:   {"postgresql"},
				ecosystem.UnavailableHealthStatus: {"postfix"},
				ecosystem.PendingHealthStatus:     {"scm"},
			},
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		componentUseCase := newMockComponentInstallationUseCase(t)
		componentUseCase.EXPECT().CheckComponentHealth(testCtx).Return(ecosystem.ComponentHealthResult{}, assert.AnError)
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, nil)

		_, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestEcosystemHealthUseCase_WaitForHealthyEcosystem(t *testing.T) {
	type fields struct {
		doguUseCaseFn        func(t *testing.T) doguInstallationUseCase
		componentUseCaseFn   func(t *testing.T) componentInstallationUseCase
		waitConfigProviderFn func(t *testing.T) healthWaitConfigProvider
	}
	tests := []struct {
		name    string
		fields  fields
		want    ecosystem.HealthResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get health check timeout",
			fields: fields{
				doguUseCaseFn: func(t *testing.T) doguInstallationUseCase {
					return newMockDoguInstallationUseCase(t)
				},
				componentUseCaseFn: func(t *testing.T) componentInstallationUseCase {
					return newMockComponentInstallationUseCase(t)
				},
				waitConfigProviderFn: func(t *testing.T) healthWaitConfigProvider {
					waitConfigMock := newMockHealthWaitConfigProvider(t)
					waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{}, assert.AnError)
					return waitConfigMock
				},
			},
			want: ecosystem.HealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get health check timeout", i)
			},
		},
		{
			name: "dogu and component health check should fail with error",
			fields: fields{
				doguUseCaseFn: func(t *testing.T) doguInstallationUseCase {
					doguMock := newMockDoguInstallationUseCase(t)
					doguMock.EXPECT().WaitForHealthyDogus(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.DoguHealthResult, error) {
							return ecosystem.DoguHealthResult{}, assert.AnError
						})
					return doguMock
				},
				componentUseCaseFn: func(t *testing.T) componentInstallationUseCase {
					componentMock := newMockComponentInstallationUseCase(t)
					componentMock.EXPECT().WaitForHealthyComponents(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
							return ecosystem.ComponentHealthResult{}, assert.AnError
						})
					return componentMock
				},
				waitConfigProviderFn: func(t *testing.T) healthWaitConfigProvider {
					waitConfigMock := newMockHealthWaitConfigProvider(t)
					waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Timeout: time.Second}, nil)
					return waitConfigMock
				},
			},
			want: ecosystem.HealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to wait for healthy components", i) &&
					assert.ErrorContains(t, err, "failed to wait for healthy dogus", i)
			},
		},
		{
			name: "context should time out for dogu and component health check",
			fields: fields{
				doguUseCaseFn: func(t *testing.T) doguInstallationUseCase {
					doguMock := newMockDoguInstallationUseCase(t)
					doguMock.EXPECT().WaitForHealthyDogus(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.DoguHealthResult, error) {
							select {
							case <-ctx.Done():
								return ecosystem.DoguHealthResult{}, assert.AnError
							case <-time.After(1 * time.Second):
								return ecosystem.DoguHealthResult{}, fmt.Errorf("test failed with timeout in dogu use case")
							}
						})
					return doguMock
				},
				componentUseCaseFn: func(t *testing.T) componentInstallationUseCase {
					componentMock := newMockComponentInstallationUseCase(t)
					componentMock.EXPECT().WaitForHealthyComponents(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
							select {
							case <-ctx.Done():
								return ecosystem.ComponentHealthResult{}, assert.AnError
							case <-time.After(1 * time.Second):
								return ecosystem.ComponentHealthResult{}, fmt.Errorf("test failed with timeout in component use case")
							}
						})
					return componentMock
				},
				waitConfigProviderFn: func(t *testing.T) healthWaitConfigProvider {
					waitConfigMock := newMockHealthWaitConfigProvider(t)
					waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Timeout: 0}, nil)
					return waitConfigMock
				},
			},
			want: ecosystem.HealthResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i)
			},
		},
		{
			name: "waiting for healthy dogus should fail",
			fields: fields{
				doguUseCaseFn: func(t *testing.T) doguInstallationUseCase {
					doguMock := newMockDoguInstallationUseCase(t)
					doguMock.EXPECT().WaitForHealthyDogus(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.DoguHealthResult, error) {
							return ecosystem.DoguHealthResult{}, assert.AnError
						})
					return doguMock
				},
				componentUseCaseFn: func(t *testing.T) componentInstallationUseCase {
					componentMock := newMockComponentInstallationUseCase(t)
					componentMock.EXPECT().WaitForHealthyComponents(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
							return ecosystem.ComponentHealthResult{ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
								ecosystem.AvailableHealthStatus: {"k8s-dogu-operator"},
							}}, nil
						})
					return componentMock
				},
				waitConfigProviderFn: func(t *testing.T) healthWaitConfigProvider {
					waitConfigMock := newMockHealthWaitConfigProvider(t)
					waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Timeout: time.Second}, nil)
					return waitConfigMock
				},
			},
			want: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{},
				ComponentHealth: ecosystem.ComponentHealthResult{ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
					ecosystem.AvailableHealthStatus: {"k8s-dogu-operator"},
				}},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i)
			},
		},
		{
			name: "waiting for healthy components should fail",
			fields: fields{
				doguUseCaseFn: func(t *testing.T) doguInstallationUseCase {
					doguMock := newMockDoguInstallationUseCase(t)
					doguMock.EXPECT().WaitForHealthyDogus(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.DoguHealthResult, error) {
							return ecosystem.DoguHealthResult{DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
								ecosystem.UnavailableHealthStatus: {"nginx-ingress"},
							}}, nil
						})
					return doguMock
				},
				componentUseCaseFn: func(t *testing.T) componentInstallationUseCase {
					componentMock := newMockComponentInstallationUseCase(t)
					componentMock.EXPECT().WaitForHealthyComponents(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
							return ecosystem.ComponentHealthResult{}, assert.AnError
						})
					return componentMock
				},
				waitConfigProviderFn: func(t *testing.T) healthWaitConfigProvider {
					waitConfigMock := newMockHealthWaitConfigProvider(t)
					waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Timeout: time.Second}, nil)
					return waitConfigMock
				},
			},
			want: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
					ecosystem.UnavailableHealthStatus: {"nginx-ingress"},
				}},
				ComponentHealth: ecosystem.ComponentHealthResult{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i)
			},
		},
		{
			name: "should succeed",
			fields: fields{
				doguUseCaseFn: func(t *testing.T) doguInstallationUseCase {
					doguMock := newMockDoguInstallationUseCase(t)
					doguMock.EXPECT().WaitForHealthyDogus(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.DoguHealthResult, error) {
							return ecosystem.DoguHealthResult{DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
								ecosystem.UnavailableHealthStatus: {"nginx-ingress"},
							}}, nil
						})
					return doguMock
				},
				componentUseCaseFn: func(t *testing.T) componentInstallationUseCase {
					componentMock := newMockComponentInstallationUseCase(t)
					componentMock.EXPECT().WaitForHealthyComponents(mock.Anything).
						RunAndReturn(func(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
							return ecosystem.ComponentHealthResult{ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
								ecosystem.AvailableHealthStatus: {"k8s-dogu-operator"},
							}}, nil
						})
					return componentMock
				},
				waitConfigProviderFn: func(t *testing.T) healthWaitConfigProvider {
					waitConfigMock := newMockHealthWaitConfigProvider(t)
					waitConfigMock.EXPECT().GetWaitConfig(testCtx).Return(ecosystem.WaitConfig{Timeout: time.Second}, nil)
					return waitConfigMock
				},
			},
			want: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
					ecosystem.UnavailableHealthStatus: {"nginx-ingress"},
				}},
				ComponentHealth: ecosystem.ComponentHealthResult{ComponentsByStatus: map[ecosystem.HealthStatus][]ecosystem.ComponentName{
					ecosystem.AvailableHealthStatus: {"k8s-dogu-operator"},
				}},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &EcosystemHealthUseCase{
				doguUseCase:        tt.fields.doguUseCaseFn(t),
				componentUseCase:   tt.fields.componentUseCaseFn(t),
				waitConfigProvider: tt.fields.waitConfigProviderFn(t),
			}
			got, err := useCase.WaitForHealthyEcosystem(testCtx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
