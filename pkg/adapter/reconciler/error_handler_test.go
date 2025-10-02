package reconciler

import (
	"fmt"
	"testing"
	"time"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/application"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ctrl "sigs.k8s.io/controller-runtime"
)

func Test_decideRequeueForError(t *testing.T) {
	t.Run("should catch wrapped InternalError, issue a log line and requeue with error", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := domainservice.NewInternalError(assert.AnError, "a generic oh-noez")
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

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
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 1 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: A concurrent update happened in conflict to the processing of the blueprint spec. A retry could fix this issue")
	})
	t.Run("should catch wrapped NotFoundError, issue a log line and do not requeue timely", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domainservice.NotFoundError{
			WrappedError: assert.AnError,
			Message:      "a generic oh-noez",
		}
		errorChain := fmt.Errorf("could not do the thing: %w", intermediateErr)

		// when
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
		assert.Contains(t, logSinkMock.output, "0: Blueprint was not found, so maybe it was deleted in the meantime. No further evaluation will happen")
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
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 10 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: Ecosystem contains multiple blueprints - delete all but one. Retry later")
	})
	t.Run("NotFoundError, should retry if referenced config is missing", func(t *testing.T) {
		// given
		logSinkMock := newTrivialTestLogSink()
		testLogger := logr.New(logSinkMock)

		intermediateErr := &domainservice.NotFoundError{
			WrappedError: assert.AnError,
			Message:      "secret xyz does not exist",
		}
		errorChain := fmt.Errorf("%s: %w", application.REFERENCED_CONFIG_NOT_FOUND, intermediateErr)

		// when
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: 10 * time.Second}, actual)
		assert.Contains(t, logSinkMock.output, "0: Referenced config not found. Retry later")
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
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

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
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

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
		sut := NewErrorHandler()
		actual, err := sut.handleError(testLogger, errorChain)

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
