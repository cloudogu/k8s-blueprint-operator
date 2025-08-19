package application

import (
	"context"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApplyBlueprintSpecUseCase(t *testing.T) {
	repoMock := newMockBlueprintSpecRepository(t)
	installUseCaseMock := newMockDoguInstallationUseCase(t)
	healthMock := newMockEcosystemHealthUseCase(t)

	sut := NewApplyBlueprintSpecUseCase(repoMock, installUseCaseMock, healthMock)

	assert.Equal(t, installUseCaseMock, sut.doguInstallUseCase)
	assert.Equal(t, repoMock, sut.repo)
	assert.Equal(t, healthMock, sut.healthUseCase)
}

func TestApplyBlueprintSpecUseCase_markBlueprintApplicationFailed(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		spec := &domain.BlueprintSpec{}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markBlueprintApplicationFailed(testCtx, spec, assert.AnError)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{}

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
		spec := &domain.BlueprintSpec{}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, spec).Return(nil)
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.markBlueprintApplied(testCtx, spec)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, spec.Status)
	})

	t.Run("repo error", func(t *testing.T) {
		spec := &domain.BlueprintSpec{}

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
	t.Run("ok", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		installUseCaseMock.EXPECT().WaitForHealthyDogus(testCtx).Return(ecosystem.DoguHealthResult{}, nil)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.NoError(t, err)
		assert.Equal(t, domain.StatusPhaseBlueprintApplied, blueprint.Status)
	})
	t.Run("error waiting for dogu health", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Once()

		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(nil)
		installUseCaseMock.EXPECT().WaitForHealthyDogus(testCtx).Return(ecosystem.DoguHealthResult{}, assert.AnError)

		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})

	t.Run("fail to apply dogu state", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{}
		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil).Once()
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})

	t.Run("fail to apply state and fail to mark execution failed", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{}
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
		installUseCaseMock := newMockDoguInstallationUseCase(t)
		installUseCaseMock.EXPECT().ApplyDoguStates(testCtx, blueprint).Return(assert.AnError)
		useCase := ApplyBlueprintSpecUseCase{repo: repoMock, doguInstallUseCase: installUseCaseMock}

		err := useCase.ApplyBlueprintSpec(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, domain.StatusPhaseBlueprintApplicationFailed, blueprint.Status)
	})
}

func TestApplyBlueprintSpecUseCase_PostProcessBlueprintApplication(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprint)

		require.NoError(t, err)

		assert.Equal(t, domain.StatusPhaseCompleted, blueprint.Status)
		assert.Len(t, blueprint.Events, 1)
		assert.Contains(t, blueprint.Events, domain.CompletedEvent{}, blueprint.Events)
	})
	t.Run("repo error while saving", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		useCase := NewApplyBlueprintSpecUseCase(repoMock, nil, nil)

		err := useCase.PostProcessBlueprintApplication(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
	})
}
