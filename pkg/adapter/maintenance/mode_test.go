package maintenance

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/errors"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testCtx   = context.TODO()
	testTitle = "title"
	testText  = "text"
)

func TestMode_Activate(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		mock.EXPECT().
			Activate(testCtx, repository.MaintenanceModeDescription{
				Title: testTitle,
				Text:  testText,
			}).Return(nil)

		err := mode.Activate(testCtx, testTitle, testText)
		require.NoError(t, err)
	})

	t.Run("conflict with maintenance lock", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		givenErr := errors.NewConflictError(assert.AnError)
		mock.EXPECT().
			Activate(testCtx, repository.MaintenanceModeDescription{
				Title: testTitle,
				Text:  testText,
			}).Return(givenErr)

		err := mode.Activate(testCtx, testTitle, testText)
		require.Error(t, err)
		assert.ErrorContains(t, err, givenErr.Error())
		assert.True(t, domainservice.IsConflictError(err))
	})

	t.Run("connection error", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		givenErr := errors.NewConnectionError(assert.AnError)
		mock.EXPECT().
			Activate(testCtx, repository.MaintenanceModeDescription{
				Title: testTitle,
				Text:  testText,
			}).Return(givenErr)

		err := mode.Activate(testCtx, testTitle, testText)
		require.Error(t, err)
		assert.ErrorContains(t, err, givenErr.Error())
		assert.True(t, domainservice.IsInternalError(err))
	})

	t.Run("unknown error", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		givenErr := assert.AnError
		mock.EXPECT().
			Activate(testCtx, repository.MaintenanceModeDescription{
				Title: testTitle,
				Text:  testText,
			}).Return(givenErr)

		err := mode.Activate(testCtx, testTitle, testText)
		require.Error(t, err)
		assert.ErrorContains(t, err, givenErr.Error())
		assert.True(t, domainservice.IsInternalError(err))
	})
}

func TestMode_Deactivate(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		mock.EXPECT().Deactivate(testCtx).Return(nil)

		err := mode.Deactivate(testCtx)
		require.NoError(t, err)
	})

	t.Run("conflict with maintenance lock", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		givenErr := errors.NewConflictError(assert.AnError)
		mock.EXPECT().Deactivate(testCtx).Return(givenErr)

		err := mode.Deactivate(testCtx)
		require.Error(t, err)
		assert.ErrorContains(t, err, givenErr.Error())
		assert.True(t, domainservice.IsConflictError(err))
	})

	t.Run("connection error", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		givenErr := errors.NewConnectionError(assert.AnError)
		mock.EXPECT().Deactivate(testCtx).Return(givenErr)

		err := mode.Deactivate(testCtx)
		require.Error(t, err)
		assert.ErrorContains(t, err, givenErr.Error())
		assert.True(t, domainservice.IsInternalError(err))
	})

	t.Run("unknown error", func(t *testing.T) {
		mock := newMockLibMaintenanceModeAdapter(t)
		mode := NewMaintenanceModeAdapter(mock)

		givenErr := assert.AnError
		mock.EXPECT().Deactivate(testCtx).Return(givenErr)

		err := mode.Deactivate(testCtx)
		require.Error(t, err)
		assert.ErrorContains(t, err, givenErr.Error())
		assert.True(t, domainservice.IsInternalError(err))
	})
}
