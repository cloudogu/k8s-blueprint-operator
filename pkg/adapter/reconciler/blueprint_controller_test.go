package reconciler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/go-logr/logr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
	reconciler := NewBlueprintReconciler(nil)
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

	t.Run("should succeed", func(t *testing.T) {
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
}

func Test_decideRequeueForError(t *testing.T) {
	t.Run("should catch wrapped InternalError, issue a log line and requeue with error", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := domainservice.NewInternalError(assert.AnError, "a generic oh-noez")
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.Error(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
		assert.Contains(t, logSinkMock.output, "0: An internal error occurred and can maybe be fixed by retrying it later")
	})
	t.Run("should catch wrapped ConflictError, issue a log line and requeue timely", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domainservice.ConflictError{
			WrappedError: assert.AnError,
			Message:      "a generic oh-noez",
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 1 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: A concurrent update happened in conflict to the processing of the blueprint spec. A retry could fix this issue")
	})
	t.Run("should catch wrapped NotFoundError, issue a log line and requeue", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domainservice.NotFoundError{
			WrappedError: assert.AnError,
			Message:      "a generic oh-noez",
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 10 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: Resource was not found, so maybe it was deleted in the meantime. Retry later")
	})
	t.Run("should catch wrapped MultipleBlueprintsError, issue a error log line and requeue", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domain.MultipleBlueprintsError{
			Message: "multiple blueprints found",
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 10 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: Ecosystem contains multiple blueprints - delete all except one. Retry later")
	})
	t.Run("NotFoundError, should not retry if DoNotRetry-Flag is set", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domainservice.NotFoundError{
			WrappedError: assert.AnError,
			Message:      "Blueprint does not exist",
			DoNotRetry:   true,
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
		assert.Contains(t, logSinkMock.output, "0: Did not find resource and a retry is not expected to fix this issue. There will be no further automatic evaluation.")
	})
	t.Run("should catch wrapped InvalidBlueprintError, issue a log line and do not requeue", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domain.InvalidBlueprintError{
			WrappedError: assert.AnError,
			Message:      "a generic oh-noez",
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
		assert.Contains(t, logSinkMock.output, "0: Blueprint is invalid, therefore there will be no further evaluation.")
	})
	t.Run("should catch wrapped StateDiffNotEmptyError, issue a log line and requeue timely", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domain.StateDiffNotEmptyError{
			Message: "a generic oh-noez",
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 1 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: requeue until state diff is empty")
	})
	t.Run("should catch general errors, issue a log line and return requeue with error", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		errorChain := fmt.Errorf("everything goes down the drain: %w", assert.AnError)

		// when
		actual, err := decideRequeueForError(testLogger, errorChain)

		// then
		require.Error(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
		assert.Contains(t, logSinkMock.output, "0: An unknown error type occurred. Retry with default backoff")
	})
}

type testLogSink struct {
	output []string
	r      logr.RuntimeInfo
}

func newTrivialTestLogSink() *testLogSink {
	var output []string
	return &testLogSink{output: output, r: logr.RuntimeInfo{CallDepth: 1}}
}

func (t *testLogSink) doLog(level int, msg string, _ ...interface{}) {
	t.output = append(t.output, fmt.Sprintf("%d: %s", level, msg))
}
func (t *testLogSink) Init(info logr.RuntimeInfo) { t.r = info }
func (t *testLogSink) Enabled(int) bool           { return true }
func (t *testLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	t.doLog(level, msg, keysAndValues...)
}
func (t *testLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	t.doLog(0, msg, append(keysAndValues, err)...)
}
func (t *testLogSink) WithValues(...interface{}) logr.LogSink { return t }
func (t *testLogSink) WithName(string) logr.LogSink           { return t }
