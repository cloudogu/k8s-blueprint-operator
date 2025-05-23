package application

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewApplyBlueprintSpecUseCase(t *testing.T) {
	repoMock := newMockBlueprintSpecRepository(t)
	installUseCaseMock := newMockDoguInstallationUseCase(t)
	componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
	healthMock := newMockEcosystemHealthUseCase(t)
	maintenanceMock := newMockMaintenanceMode(t)

	sut := NewApplyBlueprintSpecUseCase(repoMock, installUseCaseMock, healthMock, componentInstallUseCaseMock, maintenanceMock)

	assert.Equal(t, installUseCaseMock, sut.doguInstallUseCase)
	assert.Equal(t, componentInstallUseCaseMock, sut.componentInstallUseCase)
	assert.Equal(t, repoMock, sut.repo)
	assert.Equal(t, healthMock, sut.healthUseCase)
	assert.Equal(t, maintenanceMock, sut.maintenanceModeAdapter)
}

func TestApplyBlueprintSpecUseCase_PreProcessBlueprintApplication(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		maintenanceMock.EXPECT().Activate(testCtx, maintenanceTitle, maintenanceText).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationPreProcessed, spec.Status)
	})
	t.Run("repo error while loading", func(t *testing.T) {
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, nil)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("repo error while saving", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		maintenanceMock.EXPECT().Activate(testCtx, maintenanceTitle, maintenanceText).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("error activating maintenance mode", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		maintenanceMock.EXPECT().Activate(testCtx, maintenanceTitle, maintenanceText).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("do nothing on dry run", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
			Config: domain.BlueprintConfiguration{DryRun: true},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.NoError(t, err)
		require.Equal(t, 1, len(spec.Events))
		assert.Equal(t, domain.BlueprintDryRunEvent{}, spec.Events[0])
	})
}

func TestApplyBlueprintSpecUseCase_markInProgress(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.startApplying(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
	})

	t.Run("repo error while saving", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		maintenanceMock.EXPECT().Activate(testCtx, maintenanceTitle, maintenanceText).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("error activating maintenance mode", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		maintenanceMock.EXPECT().Activate(testCtx, maintenanceTitle, maintenanceText).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("do nothing on dry run", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
			Config: domain.BlueprintConfiguration{DryRun: true},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		maintenanceMock := newMockMaintenanceMode(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PreProcessBlueprintApplication(testCtx, blueprintId)

		require.NoError(t, err)
		require.Equal(t, 1, len(spec.Events))
		assert.Equal(t, domain.BlueprintDryRunEvent{}, spec.Events[0])
	})
}

func TestApplyBlueprintSpecUseCase_markBlueprintApplicationFailed(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markBlueprintApplicationFailed(testCtx, spec, assert.AnError)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markBlueprintApplicationFailed(testCtx, spec, assert.AnError)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})
}

func TestApplyBlueprintSpecUseCase_markBlueprintApplied(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markBlueprintApplied(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markBlueprintApplied(testCtx, spec)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, spec.Status)
	})
}

func TestApplyBlueprintSpecUseCase_ApplyBlueprintSpec(t *testing.T) {
	statusTransitions := map[int]domain.StatusPhase{
		1: domain.StatusPhaseInProgress,
		2: domain.StatusPhaseBlueprintApplied,
	}
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		var counter = 0
		repoMock.EXPECT().Update(testCtx, spec).RunAndReturn(func(ctx context.Context, spec *domain.BlueprintSpec) error {
			counter++
			assert.Equal(t, statusTransitions[counter], spec.Status)
			return nil
		}).Times(2)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(nil)
		installUseCaseMock.EXPECT().WaitForHealthyDogus(testCtx).Return(ecosystem.DoguHealthResult{}, nil)
		componentInstallUseCase := newMockComponentInstallationUseCase(t)
		componentInstallUseCase.EXPECT().ApplyComponentStates(testCtx, "blueprintId").Return(nil)
		componentInstallUseCase.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCase}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, spec.Status)
	})
	t.Run("error waiting for dogu health", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil).Times(2)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(nil)
		installUseCaseMock.EXPECT().WaitForHealthyDogus(testCtx).Return(ecosystem.DoguHealthResult{}, assert.AnError)
		componentInstallUseCase := newMockComponentInstallationUseCase(t)
		componentInstallUseCase.EXPECT().ApplyComponentStates(testCtx, "blueprintId").Return(nil)
		componentInstallUseCase.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCase}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})

	t.Run("cannot load spec", func(t *testing.T) {
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(nil, assert.AnError)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint to apply blueprint spec")
	})

	t.Run("fail to mark in progress", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
	})

	t.Run("fail to apply component state", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil).Times(2)

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, "blueprintId").Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})

	t.Run("fail to wait for component health", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil).Times(2)

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, "blueprintId").Return(nil)
		componentInstallUseCaseMock.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})

	t.Run("fail to apply dogu state", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil).Times(2)

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, "blueprintId").Return(nil)
		componentInstallUseCaseMock.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})

	t.Run("fail to apply state and fail to mark execution failed", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		counter := 0
		repoMock.EXPECT().Update(testCtx, spec).RunAndReturn(func(ctx context.Context, spec *domain.BlueprintSpec) error {
			counter++
			if counter == 1 {
				return nil
			} else {
				return assert.AnError
			}
		})

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, "blueprintId").Return(nil)
		componentInstallUseCaseMock.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})
}

func TestApplyBlueprintSpecUseCase_CheckEcosystemHealthUpfront(t *testing.T) {
	t.Run("should fail to get blueprint spec", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"blueprint1\" to check ecosystem health")
	})
	t.Run("should fail to get health result", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{}, nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx, false, false).Return(ecosystem.HealthResult{}, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot check ecosystem health upfront of applying the blueprint \"blueprint1\"")
	})
	t.Run("should fail to update blueprint spec", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		repoMock.EXPECT().Update(testCtx, blueprintSpec).Return(assert.AnError)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(mock.Anything, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the ecosystem health")
	})
	t.Run("should succeed, ignoring dogu and component health", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:      true,
			IgnoreComponentHealth: true,
		}}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		repoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(mock.Anything, true, true).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		repoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(mock.Anything, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})
}

func TestApplyBlueprintSpecUseCase_CheckEcosystemHealthAfterwards(t *testing.T) {
	t.Run("should fail to get blueprint spec", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"blueprint1\" to check ecosystem health")
	})

	t.Run("should fail to get health result", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{}, nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx, false, false).Return(ecosystem.HealthResult{}, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot check ecosystem health after applying the blueprint \"blueprint1\"")
	})

	t.Run("should fail to update blueprint spec", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		repoMock.EXPECT().Update(testCtx, blueprintSpec).Return(assert.AnError)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the ecosystem health")
	})

	t.Run("should succeed", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		repoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprintId)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseEcosystemHealthyAfterwards, blueprintSpec.Status)
	})
}

func TestApplyBlueprintSpecUseCase_PostProcessBlueprintApplication(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		maintenanceMock := newMockMaintenanceMode(t)
		maintenanceMock.EXPECT().Deactivate(testCtx).Return(nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprintId)

		require.NoError(t, err)

		assert.Equal(t, domain.StatusPhaseCompleted, spec.Status)
		assert.Len(t, spec.Events, 2)
		assert.Contains(t, spec.Events, domain.SensitiveConfigDataCensoredEvent{}, spec.Events)
		assert.Contains(t, spec.Events, domain.CompletedEvent{}, spec.Events)
	})

	t.Run("repo error while loading", func(t *testing.T) {
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, nil)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("repo error while saving", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		maintenanceMock := newMockMaintenanceMode(t)
		maintenanceMock.EXPECT().Deactivate(testCtx).Return(nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("error deactivating maintenance mode", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(spec, nil)
		maintenanceMock := newMockMaintenanceMode(t)
		maintenanceMock.EXPECT().Deactivate(testCtx).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil, maintenanceMock)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprintId)

		require.ErrorIs(t, err, assert.AnError)
	})
}
