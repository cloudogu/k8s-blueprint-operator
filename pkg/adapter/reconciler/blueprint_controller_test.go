package reconciler

import (
	"context"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var testCtx = context.Background()

const testBlueprint = "test-blueprint"

func TestNewBlueprintReconciler(t *testing.T) {
	reconciler := NewBlueprintReconciler(nil, nil)
	assert.NotNil(t, reconciler)
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

	scheme.AddKnownTypes(gv, &v1.Blueprint{})
	return scheme
}

func TestBlueprintReconciler_Reconcile(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		// given
		request := ctrl.Request{NamespacedName: types.NamespacedName{Name: testBlueprint}}
		sut := &BlueprintReconciler{}

		// when
		actual, err := sut.Reconcile(testCtx, request)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
	})
}
