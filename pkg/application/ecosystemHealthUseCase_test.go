package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewEcosystemHealthUseCase(t *testing.T) {
	doguUseCase := newMockDoguInstallationUseCase(t)
	componentUseCase := newMockComponentInstallationUseCase(t)
	useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, time.Minute)

	assert.Equal(t, doguUseCase, useCase.doguUseCase)
	assert.Equal(t, componentUseCase, useCase.componentUseCase)
	assert.Equal(t, time.Minute, useCase.healthCheckTimeOut)
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
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, time.Minute)

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
		useCase := NewEcosystemHealthUseCase(nil, componentUseCase, time.Minute)

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
		useCase := NewEcosystemHealthUseCase(doguUseCase, nil, time.Minute)

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
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, time.Minute)

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
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase, time.Minute)

		_, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.ErrorIs(t, err, assert.AnError)
	})
}
