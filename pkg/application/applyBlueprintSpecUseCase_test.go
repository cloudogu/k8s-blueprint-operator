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

	sut := NewApplyBlueprintSpecUseCase(repoMock, installUseCaseMock, healthMock, componentInstallUseCaseMock)

	assert.Equal(t, installUseCaseMock, sut.doguInstallUseCase)
	assert.Equal(t, componentInstallUseCaseMock, sut.componentInstallUseCase)
	assert.Equal(t, repoMock, sut.repo)
	assert.Equal(t, healthMock, sut.healthUseCase)
}

func TestApplyBlueprintSpecUseCase_PreProcessBlueprintApplication(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

		err := useCase.PreProcessBlueprintApplication(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationPreProcessed, spec.Status)
	})
	t.Run("repo error while saving", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

		err := useCase.PreProcessBlueprintApplication(testCtx, spec)

		require.ErrorIs(t, err, assert.AnError)
	})
	t.Run("do nothing on dry run", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
			Config: domain.BlueprintConfiguration{DryRun: true},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

		err := useCase.PreProcessBlueprintApplication(testCtx, spec)

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
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		var counter = 0
		repoMock.EXPECT().Update(testCtx, blueprint).RunAndReturn(func(ctx context.Context, spec *domain.BlueprintSpec) error {
			counter++
			assert.Equal(t, statusTransitions[counter], spec.Status)
			return nil
		}).Times(2)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		installUseCaseMock.EXPECT().WaitForHealthyDogus(testCtx).Return(ecosystem.DoguHealthResult{}, nil)
		componentInstallUseCase := newMockComponentInstallationUseCase(t)
		componentInstallUseCase.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		componentInstallUseCase.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCase}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, blueprint.Status)
	})
	t.Run("error waiting for dogu health", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Times(2)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		installUseCaseMock.EXPECT().WaitForHealthyDogus(testCtx).Return(ecosystem.DoguHealthResult{}, assert.AnError)
		componentInstallUseCase := newMockComponentInstallationUseCase(t)
		componentInstallUseCase.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		componentInstallUseCase.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCase}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})

	t.Run("fail to mark in progress", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseInProgress, blueprint.Status)
	})

	t.Run("fail to apply component state", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Times(2)

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})

	t.Run("fail to wait for component health", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Times(2)

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		componentInstallUseCaseMock.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})

	t.Run("fail to apply dogu state", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Times(2)

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		componentInstallUseCaseMock.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})

	t.Run("fail to apply state and fail to mark execution failed", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		counter := 0
		repoMock.EXPECT().Update(testCtx, blueprint).RunAndReturn(func(ctx context.Context, spec *domain.BlueprintSpec) error {
			counter++
			if counter == 1 {
				return nil
			} else {
				return assert.AnError
			}
		})

		componentInstallUseCaseMock := newMockComponentInstallationUseCase(t)
		componentInstallUseCaseMock.EXPECT().ApplyComponentStates(testCtx, blueprint).Return(nil)
		componentInstallUseCaseMock.EXPECT().WaitForHealthyComponents(testCtx).Return(ecosystem.ComponentHealthResult{}, nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock, componentInstallUseCase: componentInstallUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})
}

func TestApplyBlueprintSpecUseCase_CheckEcosystemHealthUpfront(t *testing.T) {
	t.Run("should fail to get health result", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: blueprintId,
		}

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(testCtx, false, false).Return(ecosystem.HealthResult{}, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(nil, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot check ecosystem health upfront of applying the blueprint \"blueprint1\"")
	})
	t.Run("should fail to update blueprint spec", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: blueprintId,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(mock.Anything, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the ecosystem health")
	})
	t.Run("should succeed, ignoring dogu and component health", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth:      true,
			IgnoreComponentHealth: true,
		}}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(mock.Anything, true, true).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprint)

		// then
		require.NoError(t, err)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(mock.Anything, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprint)

		// then
		require.NoError(t, err)
	})
}

func TestApplyBlueprintSpecUseCase_CheckEcosystemHealthAfterwards(t *testing.T) {
	t.Run("should fail to get health result", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: blueprintId,
		}

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(testCtx, false, false).Return(ecosystem.HealthResult{}, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(nil, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot check ecosystem health after applying the blueprint \"blueprint1\"")
	})

	t.Run("should fail to update blueprint spec", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: blueprintId,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(testCtx, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprint)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the ecosystem health")
	})

	t.Run("should succeed", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(testCtx, false, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprint)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseEcosystemHealthyAfterwards, blueprint.Status)
	})
}

func TestApplyBlueprintSpecUseCase_PostProcessBlueprintApplication(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprint)

		require.NoError(t, err)

		assert.Equal(t, domain.StatusPhaseCompleted, blueprint.Status)
		assert.Len(t, blueprint.Events, 2)
		assert.Contains(t, blueprint.Events, domain.SensitiveConfigDataCensoredEvent{}, blueprint.Events)
		assert.Contains(t, blueprint.Events, domain.CompletedEvent{}, blueprint.Events)
	})
	t.Run("repo error while saving", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyAfterwards,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
	})
}
