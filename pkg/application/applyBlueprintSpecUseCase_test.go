package application

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewApplyBlueprintSpecUseCase(t *testing.T) {
	repoMock := newMockBlueprintSpecRepository(t)
	installUseCaseMock := newMockDoguInstallationUseCase(t)

	sut := NewApplyBlueprintSpecUseCase(repoMock, installUseCaseMock)

	assert.Equal(t, installUseCaseMock, sut.doguInstallUseCase)
	assert.Equal(t, repoMock, sut.repo)
}

func TestApplyBlueprintSpecUseCase_markInProgress(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseDogusHealthy,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markInProgress(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseDogusHealthy,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markInProgress(testCtx, spec)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
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

		err := useCase.markWaitingForHealthyEcosystem(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseWaitForHealthyEcosystem, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseInProgress,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markWaitingForHealthyEcosystem(testCtx, spec)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseWaitForHealthyEcosystem, spec.Status)
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
		2: domain.StatusPhaseWaitForHealthyEcosystem,
		3: domain.StatusPhaseCompleted,
	}
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseDogusHealthy,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		var counter = 0
		repoMock.EXPECT().Update(testCtx, spec).RunAndReturn(func(ctx context.Context, spec *domain.BlueprintSpec) error {
			counter++
			assert.Equal(t, spec.Status, statusTransitions[counter])
			return nil
		}).Times(3)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(nil)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseCompleted, spec.Status)
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
			Status: domain.StatusPhaseDogusHealthy,
		}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().GetById(testCtx, "blueprintId").Return(spec, nil)
		repoMock.EXPECT().Update(testCtx, spec).Return(assert.AnError)

		//installUseCaseMock := newMockDoguInstallationUseCase(t)
		//installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, "blueprintId").Return(nil)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: nil}

		err := useCase.ApplyBlueprintSpec(testCtx, "blueprintId")

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseInProgress, spec.Status)
	})

	t.Run("fail to apply state", func(t *testing.T) {
		spec := &domain.BlueprintSpec{
			Status: domain.StatusPhaseDogusHealthy,
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
			Status: domain.StatusPhaseDogusHealthy,
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
