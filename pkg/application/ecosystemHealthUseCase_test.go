package application

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

var (
	healthyDogu = map[ecosystem.HealthStatus][]cescommons.SimpleName{
		ecosystem.AvailableHealthStatus: {"postgresql"},
	}
	mixedDoguHealth = map[ecosystem.HealthStatus][]cescommons.SimpleName{
		ecosystem.AvailableHealthStatus:   {"postgresql"},
		ecosystem.UnavailableHealthStatus: {"postfix"},
		ecosystem.PendingHealthStatus:     {"scm"},
	}
)

func TestNewEcosystemHealthUseCase(t *testing.T) {
	doguUseCase := newMockDoguInstallationUseCase(t)
	blueprintRepo := newMockBlueprintSpecRepository(t)
	useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

	assert.Same(t, doguUseCase, useCase.doguUseCase)
}

func TestEcosystemHealthUseCase_CheckEcosystemHealth(t *testing.T) {
	t.Run("all healthy", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth: false,
			},
		}

		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: healthyDogu,
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		blueprintRepo := newMockBlueprintSpecRepository(t)
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		health, err := useCase.CheckEcosystemHealth(testCtx, blueprint)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth}, health)
		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionEcosystemHealthy))
	})

	t.Run("unhealthy", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth: false,
			},
		}

		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: mixedDoguHealth,
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		blueprintRepo := newMockBlueprintSpecRepository(t)
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		health, err := useCase.CheckEcosystemHealth(testCtx, blueprint)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "ecosystem is unhealthy")
		assert.ErrorContains(t, err, "2 dogu(s) are unhealthy: postfix, scm")
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth}, health)
		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionEcosystemHealthy))
	})

	t.Run("error updating blueprint", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth: false,
			},
		}

		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: mixedDoguHealth,
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		blueprintRepo := newMockBlueprintSpecRepository(t)
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		health, err := useCase.CheckEcosystemHealth(testCtx, blueprint)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not update health condition after health check")
		assert.Equal(t, ecosystem.HealthResult{}, health)
	})

	t.Run("error getting health", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth: false,
			},
		}

		doguHealth := ecosystem.DoguHealthResult{}

		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, assert.AnError)
		blueprintRepo := newMockBlueprintSpecRepository(t)
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil)
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		_, err := useCase.CheckEcosystemHealth(testCtx, blueprint)

		assert.ErrorIs(t, err, assert.AnError)
		assert.True(t, meta.IsStatusConditionPresentAndEqual(
			blueprint.Conditions, domain.ConditionEcosystemHealthy, metav1.ConditionUnknown,
		))
	})

	t.Run("no update without health change", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			Config: domain.BlueprintConfiguration{
				IgnoreDoguHealth: false,
			},
		}

		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: mixedDoguHealth,
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil).Twice()
		blueprintRepo := newMockBlueprintSpecRepository(t)
		blueprintRepo.EXPECT().Update(testCtx, blueprint).Return(nil).Once()
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		_, err := useCase.CheckEcosystemHealth(testCtx, blueprint)
		assert.ErrorContains(t, err, "ecosystem is unhealthy")
		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionEcosystemHealthy))
		_, err = useCase.CheckEcosystemHealth(testCtx, blueprint) //no repo.Update called again
		assert.ErrorContains(t, err, "ecosystem is unhealthy")
		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionEcosystemHealthy))
	})
}

func TestEcosystemHealthUseCase_getEcosystemHealth(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		doguHealth := ecosystem.DoguHealthResult{
			DogusByStatus: mixedDoguHealth,
		}
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(doguHealth, nil)
		blueprintRepo := newMockBlueprintSpecRepository(t)
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		health, err := useCase.getEcosystemHealth(testCtx, false)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{DoguHealth: doguHealth}, health)
	})

	t.Run("ok, ignore dogu health", func(t *testing.T) {
		blueprintRepo := newMockBlueprintSpecRepository(t)
		useCase := NewEcosystemHealthUseCase(nil, blueprintRepo)

		health, err := useCase.getEcosystemHealth(testCtx, true)

		require.NoError(t, err)
		assert.Equal(t, ecosystem.HealthResult{}, health)
	})

	t.Run("error checking dogu health", func(t *testing.T) {
		doguUseCase := newMockDoguInstallationUseCase(t)
		doguUseCase.EXPECT().CheckDoguHealth(mock.Anything).Return(ecosystem.DoguHealthResult{}, assert.AnError)
		blueprintRepo := newMockBlueprintSpecRepository(t)
		useCase := NewEcosystemHealthUseCase(doguUseCase, blueprintRepo)

		_, err := useCase.getEcosystemHealth(testCtx, false)

		require.ErrorIs(t, err, assert.AnError)
	})
}
