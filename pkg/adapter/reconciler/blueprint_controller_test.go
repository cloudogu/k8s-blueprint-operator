package reconciler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/log"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
)

var testCtx = context.Background()

const testBlueprint = "test-blueprint"

func TestNewBlueprintReconciler(t *testing.T) {
	reconciler := NewBlueprintReconciler(nil, nil, "", time.Duration(0))
	assert.NotNil(t, reconciler)
	assert.NotNil(t, reconciler.errorHandler)
}

func TestBlueprintReconciler_SetupWithManager(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		// given
		sut := &BlueprintReconciler{}

		// when
		err := sut.SetupWithManager(nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "must provide a non-nil Manager")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		ctrlManMock := newMockControllerManager(t)
		ctrlManMock.EXPECT().GetControllerOptions().Return(config.Controller{})
		ctrlManMock.EXPECT().GetScheme().Return(createScheme(t))
		logger := log.FromContext(testCtx)
		ctrlManMock.EXPECT().GetLogger().Return(logger)
		ctrlManMock.EXPECT().Add(mock.Anything).Return(nil)
		ctrlManMock.EXPECT().GetCache().Return(nil)

		sut := &BlueprintReconciler{}

		// when
		err := sut.SetupWithManager(ctrlManMock)

		// then
		require.NoError(t, err)
	})
}

func createScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()
	gv, err := schema.ParseGroupVersion("k8s.cloudogu.com/v1")
	assert.NoError(t, err)

	scheme.AddKnownTypes(gv, &bpv2.Blueprint{})
	return scheme
}

func TestBlueprintReconciler_Reconcile(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		// given
		request := ctrl.Request{NamespacedName: types.NamespacedName{Name: testBlueprint}}
		changeHandlerMock := NewMockBlueprintChangeHandler(t)
		sut := &BlueprintReconciler{blueprintChangeHandler: changeHandlerMock}

		changeHandlerMock.EXPECT().CheckForMultipleBlueprintResources(testCtx).Return(nil)
		changeHandlerMock.EXPECT().HandleUntilApplied(testCtx, testBlueprint).Return(nil)
		// when
		actual, err := sut.Reconcile(testCtx, request)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
	})

	t.Run("should fail on multiple blueprint resource error", func(t *testing.T) {
		// given
		request := ctrl.Request{NamespacedName: types.NamespacedName{Name: testBlueprint}}
		changeHandlerMock := NewMockBlueprintChangeHandler(t)
		sut := &BlueprintReconciler{blueprintChangeHandler: changeHandlerMock}

		changeHandlerMock.EXPECT().CheckForMultipleBlueprintResources(testCtx).Return(assert.AnError)
		// when
		_, err := sut.Reconcile(testCtx, request)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should fail on HandleUntilApplied error", func(t *testing.T) {
		// given
		request := ctrl.Request{NamespacedName: types.NamespacedName{Name: testBlueprint}}
		changeHandlerMock := NewMockBlueprintChangeHandler(t)
		sut := &BlueprintReconciler{blueprintChangeHandler: changeHandlerMock}

		changeHandlerMock.EXPECT().CheckForMultipleBlueprintResources(testCtx).Return(nil)
		changeHandlerMock.EXPECT().HandleUntilApplied(testCtx, testBlueprint).Return(errors.New("test"))
		// when
		_, err := sut.Reconcile(testCtx, request)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "test")
	})

	t.Run("should handle error with requeue", func(t *testing.T) {
		mockHandler := NewMockBlueprintChangeHandler(t)
		mockRepo := NewMockBlueprintSpecRepository(t)

		reconciler := NewBlueprintReconciler(mockHandler, mockRepo, "test-namespace", 5*time.Second)

		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-blueprint",
				Namespace: "test-namespace",
			},
		}

		ctx := context.TODO()
		testErr := &domainservice.ConflictError{Message: "conflict error"}

		mockHandler.EXPECT().CheckForMultipleBlueprintResources(ctx).Return(nil)
		mockHandler.EXPECT().HandleUntilApplied(ctx, "test-blueprint").Return(testErr)

		result, err := reconciler.Reconcile(ctx, req)

		assert.NoError(t, err) // Error should be handled by ErrorHandler
		assert.Equal(t, ctrl.Result{RequeueAfter: 1 * time.Second}, result)
	})

	t.Run("should reconcile on pending change", func(t *testing.T) {
		mockHandler := NewMockBlueprintChangeHandler(t)
		mockRepo := NewMockBlueprintSpecRepository(t)

		reconciler := NewBlueprintReconciler(mockHandler, mockRepo, "test-namespace", 5*time.Second)

		// Set up debounce to have pending request
		reconciler.debounce.AllowOrMark(1 * time.Second)
		reconciler.debounce.AllowOrMark(1 * time.Second) // This marks as pending
		assert.True(t, reconciler.debounce.pending)

		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-blueprint",
				Namespace: "test-namespace",
			},
		}

		ctx := context.TODO()

		mockHandler.EXPECT().CheckForMultipleBlueprintResources(ctx).Return(nil)
		mockHandler.EXPECT().HandleUntilApplied(ctx, "test-blueprint").Return(nil)

		result, err := reconciler.Reconcile(ctx, req)

		assert.NoError(t, err)
		assert.True(t, result.RequeueAfter > 0)
	})
}

func TestBlueprintReconciler_getBlueprintRequest(t *testing.T) {
	ctx := context.TODO()

	t.Run("one blueprint gets successful request", func(t *testing.T) {
		list := &bpv2.BlueprintList{
			Items: []bpv2.Blueprint{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-blueprint",
						Namespace: "test-namespace",
					},
				},
			},
		}

		mockRepo := NewMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().List(ctx).Return(list, nil)

		reconciler := &BlueprintReconciler{blueprintRepo: mockRepo}
		result := reconciler.getBlueprintRequest(ctx)

		expected := []reconcile.Request{{
			NamespacedName: types.NamespacedName{
				Name:      "test-blueprint",
				Namespace: "test-namespace",
			},
		}}

		assert.Equal(t, expected, result)
	})

	t.Run("no reconcile request on error from repository", func(t *testing.T) {
		mockRepo := NewMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().List(ctx).Return(nil, errors.New("repo error"))

		reconciler := &BlueprintReconciler{blueprintRepo: mockRepo}
		result := reconciler.getBlueprintRequest(ctx)

		assert.Nil(t, result)
	})

	t.Run("no reconcile request when no blueprints", func(t *testing.T) {
		list := &bpv2.BlueprintList{Items: []bpv2.Blueprint{}}

		mockRepo := NewMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().List(ctx).Return(list, nil)

		reconciler := &BlueprintReconciler{blueprintRepo: mockRepo}
		result := reconciler.getBlueprintRequest(ctx)

		assert.Nil(t, result)
	})

	t.Run("no reconcile request when multiple blueprints", func(t *testing.T) {
		list := &bpv2.BlueprintList{
			Items: []bpv2.Blueprint{
				{ObjectMeta: metav1.ObjectMeta{Name: "bp1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "bp2"}},
			},
		}

		mockRepo := NewMockBlueprintSpecRepository(t)
		mockRepo.EXPECT().List(ctx).Return(list, nil)

		reconciler := &BlueprintReconciler{blueprintRepo: mockRepo}
		result := reconciler.getBlueprintRequest(ctx)

		assert.Nil(t, result)
	})
}
