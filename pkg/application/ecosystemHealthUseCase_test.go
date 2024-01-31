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
	useCase := NewEcosystemHealthUseCase(doguUseCase, time.Minute)

	assert.Equal(t, doguUseCase, useCase.doguUseCase)
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
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, time.Minute)

		health, err := useCase.CheckEcosystemHealth(testCtx, false)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth}, health)
	})

	t.Run("ok, ignore dogu health", func(t *testing.T) {
		useCase := NewEcosystemHealthUseCase(nil, time.Minute)

		health, err := useCase.CheckEcosystemHealth(testCtx, true)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{}, health)
	})

	t.Run("error checking dogu health", func(t *testing.T) {
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(ecosystem.DoguHealthResult{}, assert.AnError)
		useCase := NewEcosystemHealthUseCase(doguUseCase, time.Minute)

		_, err := useCase.CheckEcosystemHealth(testCtx, false)

		require.ErrorIs(t, err, assert.AnError)
	})
}
