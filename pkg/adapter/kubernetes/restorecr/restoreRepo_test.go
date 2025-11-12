package restorecr

import (
	"context"
	"testing"

	restorev1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var testCtx = context.Background()

func TestNewRestoreRepo(t *testing.T) {
	t.Run("should create new RestoreRepo", func(t *testing.T) {
		mRestoreClient := NewMockRestoreInterface(t)

		repo := NewRestoreRepo(mRestoreClient)

		assert.NotNil(t, repo)
		assert.Equal(t, mRestoreClient, repo.(*restoreRepo).restoreClient)
	})
}

func Test_restoreRepo_IsRestoreInProgress(t *testing.T) {
	t.Run("should return true if a restore is in progress", func(t *testing.T) {
		mRestoreClient := NewMockRestoreInterface(t)
		mRestoreClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&restorev1.RestoreList{
			Items: []restorev1.Restore{
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-1"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusNew}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-2"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusCompleted}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-3"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusDeleting}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-4"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusFailed}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-5"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusInProgress}},
			},
		}, nil)

		repo := &restoreRepo{restoreClient: mRestoreClient}

		result, err := repo.IsRestoreInProgress(testCtx)

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if no restore is in progress", func(t *testing.T) {
		mRestoreClient := NewMockRestoreInterface(t)
		mRestoreClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&restorev1.RestoreList{
			Items: []restorev1.Restore{
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-1"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusNew}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-2"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusCompleted}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-3"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusDeleting}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-4"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusFailed}},
				{ObjectMeta: metav1.ObjectMeta{Name: "restore-5"}, Status: restorev1.RestoreStatus{Status: restorev1.RestoreStatusCompleted}},
			},
		}, nil)

		repo := &restoreRepo{restoreClient: mRestoreClient}

		result, err := repo.IsRestoreInProgress(testCtx)

		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false if no restore exists", func(t *testing.T) {
		mRestoreClient := NewMockRestoreInterface(t)
		mRestoreClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&restorev1.RestoreList{
			Items: []restorev1.Restore{},
		}, nil)

		repo := &restoreRepo{restoreClient: mRestoreClient}

		result, err := repo.IsRestoreInProgress(testCtx)

		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should fail if there is an error listing restores", func(t *testing.T) {
		mRestoreClient := NewMockRestoreInterface(t)
		mRestoreClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, assert.AnError)

		repo := &restoreRepo{restoreClient: mRestoreClient}

		_, err := repo.IsRestoreInProgress(testCtx)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while listing restore CRs")
	})

	t.Run("should not fail if no restores could be found", func(t *testing.T) {
		mRestoreClient := NewMockRestoreInterface(t)
		mRestoreClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, k8serrors.NewNotFound(schema.GroupResource{}, "restores not found"))

		repo := &restoreRepo{restoreClient: mRestoreClient}

		result, err := repo.IsRestoreInProgress(testCtx)

		require.NoError(t, err)
		assert.False(t, result)
	})
}
