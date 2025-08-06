package application

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

func TestNewEcosystemHealthUseCase(t *testing.T) {
	doguUseCase := newMockDoguInstallationUseCase(t)
	componentUseCase := newMockComponentInstallationUseCase(t)
	useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase)

	assert.Same(t, doguUseCase, useCase.doguUseCase)
	assert.Same(t, componentUseCase, useCase.componentUseCase)
}

func TestEcosystemHealthUseCase_CheckEcosystemHealth(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
				ecosystem.AvailableHealthStatus:   {"postgresql"},
				ecosystem.UnavailableHealthStatus: {"postfix"},
				ecosystem.PendingHealthStatus:     {"scm"},
			},
		}
		componentHealth := ecosystem.ComponentHealthResult{
			ComponentsByStatus: map[ecosystem.HealthStatus][]common.SimpleComponentName{
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
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase)

		health, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth, ComponentHealth: componentHealth}, health)
	})

	t.Run("ok, ignore dogu health", func(t *testing.T) {
		componentHealth := ecosystem.ComponentHealthResult{
			ComponentsByStatus: map[ecosystem.HealthStatus][]common.SimpleComponentName{
				ecosystem.NotInstalledHealthStatus: {"k8s-dogu-operator"},
				ecosystem.UnavailableHealthStatus:  {"k8s-etcd"},
				ecosystem.PendingHealthStatus:      {"k8s-service-discovery"},
				ecosystem.AvailableHealthStatus:    {"k8s-component-operator"},
			},
		}
		componentUseCase := newMockComponentInstallationUseCase(t)
		componentUseCase.EXPECT().CheckComponentHealth(testCtx).Return(componentHealth, nil)
		useCase := NewEcosystemHealthUseCase(nil, componentUseCase)

		health, err := useCase.CheckEcosystemHealth(testCtx, true, false)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{ComponentHealth: componentHealth}, health)
	})

	t.Run("ok, ignore component health", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
				ecosystem.AvailableHealthStatus:   {"postgresql"},
				ecosystem.UnavailableHealthStatus: {"postfix"},
				ecosystem.PendingHealthStatus:     {"scm"},
			},
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, nil)

		health, err := useCase.CheckEcosystemHealth(testCtx, false, true)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth}, health)
	})

	t.Run("error checking dogu health", func(t *testing.T) {
		componentHealth := ecosystem.ComponentHealthResult{
			ComponentsByStatus: map[ecosystem.HealthStatus][]common.SimpleComponentName{
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
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase)

		_, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("error checking component health", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
				ecosystem.AvailableHealthStatus:   {"postgresql"},
				ecosystem.UnavailableHealthStatus: {"postfix"},
				ecosystem.PendingHealthStatus:     {"scm"},
			},
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		componentUseCase := newMockComponentInstallationUseCase(t)
		componentUseCase.EXPECT().CheckComponentHealth(testCtx).Return(ecosystem.ComponentHealthResult{}, assert.AnError)
		useCase := NewEcosystemHealthUseCase(doguUseCase, componentUseCase)

		_, err := useCase.CheckEcosystemHealth(testCtx, false, false)

		require.ErrorIs(t, err, assert.AnError)
	})
}
