package application

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
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
