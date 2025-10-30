package application

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRestoreInProgressUseCase(t *testing.T) {
	t.Run("should create new RestoreInProgressUseCase", func(t *testing.T) {
		mRestoreRepo := newMockRestoreRepository(t)

		sut := NewRestoreInProgressUseCase(mRestoreRepo)

		require.NotNil(t, sut)
		assert.Equal(t, mRestoreRepo, sut.restoreRepo)
	})
}

func TestRestoreInProgressUseCase_CheckRestoreInProgress(t *testing.T) {
	t.Run("should return RestoreInProgressError if restore is in progress", func(t *testing.T) {
		mRestoreRepo := newMockRestoreRepository(t)
		mRestoreRepo.EXPECT().IsRestoreInProgress(testCtx).Return(true, nil)

		sut := &RestoreInProgressUseCase{
			restoreRepo: mRestoreRepo,
		}

		err := sut.CheckRestoreInProgress(testCtx)

		require.Error(t, err)
		var restoreErr *domain.RestoreInProgressError
		assert.ErrorAs(t, err, &restoreErr)
		assert.Equal(t, "cannot apply blueprint because a restore is in progress", restoreErr.Message)
	})

	t.Run("should return nil if no restore is in progress", func(t *testing.T) {
		mRestoreRepo := newMockRestoreRepository(t)
		mRestoreRepo.EXPECT().IsRestoreInProgress(testCtx).Return(false, nil)

		sut := &RestoreInProgressUseCase{
			restoreRepo: mRestoreRepo,
		}

		err := sut.CheckRestoreInProgress(testCtx)

		require.NoError(t, err)
	})

	t.Run("should fail if no restore is in progress", func(t *testing.T) {
		mRestoreRepo := newMockRestoreRepository(t)
		mRestoreRepo.EXPECT().IsRestoreInProgress(testCtx).Return(false, assert.AnError)

		sut := &RestoreInProgressUseCase{
			restoreRepo: mRestoreRepo,
		}

		err := sut.CheckRestoreInProgress(testCtx)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		var restoreErr *domainservice.InternalError
		assert.ErrorAs(t, err, &restoreErr)
		assert.Equal(t, "error while checking if a restore is in progress", restoreErr.Message)
	})
}
