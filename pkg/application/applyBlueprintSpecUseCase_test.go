package application

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewApplyBlueprintSpecUseCase(t *testing.T) {
	repoMock := newMockBlueprintSpecRepository(t)
	installUseCaseMock := newMockDoguInstallationUseCase(t)
	healthMock := newMockEcosystemHealthUseCase(t)

	sut := NewApplyBlueprintSpecUseCase(repoMock, installUseCaseMock, healthMock, nil)

	assert.Equal(t, installUseCaseMock, sut.doguInstallUseCase)
	assert.Equal(t, repoMock, sut.repo)
	assert.Equal(t, healthMock, sut.healthUseCase)
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

		shouldApply, err := useCase.startApplying(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
		assert.True(t, shouldApply)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		shouldApply, err := useCase.startApplying(testCtx, spec)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
		assert.False(t, shouldApply)
	})
}

func TestApplyBlueprintSpecUseCase_MarkFailed(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.MarkFailed(testCtx, spec, assert.AnError)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseFailed, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.MarkFailed(testCtx, spec, assert.AnError)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseFailed, spec.Status)
	})
}

func TestApplyBlueprintSpecUseCase_markWaitingForHealthyEcosystem(t *testing.T) {
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

func TestApplyBlueprintSpecUseCase_markCompleted(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markCompleted(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseCompleted, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markCompleted(testCtx, spec)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseCompleted, spec.Status)
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
			assert.Equal(t, spec.Status, statusTransitions[counter])
			return nil
		}).Times(2)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(nil)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, spec.Status)
	})

	t.Run("should do nothing and return nil on dry run", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Config: domain.BlueprintConfiguration{
				DryRun: true,
			},
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil).Times(1)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil).Times(1).Run(func(args mock.Arguments) {
			spec := args.Get(1).(*domain.BlueprintSpec)
			assert.Equal(t, domain.BlueprintDryRunEvent{}, spec.Events[0])
		})

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.NoError(t, err)
	})

	t.Run("should return error on error updating blueprint spec with dry run event", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Config: domain.BlueprintConfiguration{
				DryRun: true,
			},
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil).Times(1)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError).Times(1).Run(func(args mock.Arguments) {
			spec := args.Get(1).(*domain.BlueprintSpec)
			assert.Equal(t, domain.BlueprintDryRunEvent{}, spec.Events[0])
		})

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot mark blueprint as in progress")
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

		// installUseCaseMock := newMockDoguInstallationUseCase(t)
		// installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(nil)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
	})

	t.Run("fail to apply state", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseEcosystemHealthyUpfront,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil).Times(2)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseFailed, spec.Status)
	})

	t.Run("fail to apply state and fail to mark failed", func(t *testing.T) {
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

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseFailed, spec.Status)
	})
}

func TestApplyBlueprintSpecUseCase_CheckEcosystemHealthUpfront(t *testing.T) {
	t.Run("should fail to get blueprint spec", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

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
		healthMock.EXPECT().CheckEcosystemHealth(testCtx, false).Return(ecosystem.HealthResult{}, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

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
		healthMock.EXPECT().CheckEcosystemHealth(mock.Anything, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthUpfront(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the ecosystem health")
	})
	t.Run("should succeed, ignoring dogu health", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{Config: domain.BlueprintConfiguration{
			IgnoreDoguHealth: true,
		}}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		repoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		healthMock := newMockEcosystemHealthUseCase(t)
		healthMock.EXPECT().CheckEcosystemHealth(mock.Anything, true).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

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
		healthMock.EXPECT().CheckEcosystemHealth(mock.Anything, false).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

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

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, nil, nil)

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
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx).Return(ecosystem.HealthResult{}, assert.AnError)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

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
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

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
		healthMock.EXPECT().WaitForHealthyEcosystem(testCtx).Return(ecosystem.HealthResult{}, nil)

		sut := NewApplyBlueprintSpecUseCase(repoMock, nil, healthMock, nil)

		// when
		err := sut.CheckEcosystemHealthAfterwards(testCtx, blueprintId)

		// then
		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseEcosystemHealthyAfterwards, blueprintSpec.Status)
	})
}
