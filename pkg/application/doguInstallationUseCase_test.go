package application

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const blueprintId = "blueprint1"

func TestDoguInstallationUseCase_CheckDoguHealth(t *testing.T) {
	t.Run("should fail to get blueprint spec", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(nil, assert.AnError)

		doguRepoMock := newMockDoguInstallationRepository(t)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot load blueprint spec \"blueprint1\" to check dogu health")
	})
	t.Run("should fail to get dogus", func(t *testing.T) {
		// given
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(&domain.BlueprintSpec{}, nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(nil, assert.AnError)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot evaluate dogu health states for blueprint spec \"blueprint1\"")
	})
	t.Run("should fail to update blueprint spec", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		blueprintSpecRepoMock.EXPECT().Update(testCtx, blueprintSpec).Return(assert.AnError)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "cannot save blueprint spec \"blueprint1\" after checking the dogu health")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		blueprintSpec := &domain.BlueprintSpec{}
		blueprintSpecRepoMock := newMockBlueprintSpecRepository(t)
		blueprintSpecRepoMock.EXPECT().GetById(testCtx, blueprintId).Return(blueprintSpec, nil)
		blueprintSpecRepoMock.EXPECT().Update(testCtx, blueprintSpec).Return(nil)

		doguRepoMock := newMockDoguInstallationRepository(t)
		doguRepoMock.EXPECT().GetAll(testCtx).Return(map[string]*ecosystem.DoguInstallation{}, nil)

		sut := NewDoguInstallationUseCase(blueprintSpecRepoMock, doguRepoMock)

		// when
		err := sut.CheckDoguHealth(testCtx, blueprintId)

		// then
		require.NoError(t, err)
	})
}
