package logging

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfigureLogger(t *testing.T) {
	t.Run("create logger with no log level set in env -> should use default", func(t *testing.T) {
		// given
		_ = os.Unsetenv(logLevelEnvVar)

		// when
		err := ConfigureLogger()

		// then
		assert.NoError(t, err)
	})
	t.Run("should not fail with empty string log level and return info level", func(t *testing.T) {
		// given
		t.Setenv(logLevelEnvVar, "")

		// when
		err := ConfigureLogger()

		// then
		assert.NoError(t, err)
		assert.Equal(t, logrus.InfoLevel, CurrentLogLevel)
	})

	t.Run("should not fail with DEBUG log level", func(t *testing.T) {
		// given
		t.Setenv(logLevelEnvVar, "DEBUG")

		// when
		err := ConfigureLogger()

		// then
		assert.NoError(t, err)
		assert.Equal(t, logrus.DebugLevel, CurrentLogLevel)
	})

	t.Run("create logger with invalid log level TEST_LEVEL", func(t *testing.T) {
		// given
		_ = os.Setenv(logLevelEnvVar, "TEST_LEVEL")

		// when
		err := ConfigureLogger()

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "value of log environment variable [LOG_LEVEL] is not a valid log level")
	})
}

func Test_libraryLogger_Debug(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", debugLevel, "[testLogger] test debug call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	logger.Debug("test debug call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Debugf(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", debugLevel, "[testLogger] myText - test debug call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	text := "myText"
	logger.Debugf("%s - %s", text, "test debug call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Error(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", errorLevel, "[testLogger] test error call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	logger.Error("test error call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Errorf(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", errorLevel, "[testLogger] myText - test error call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	text := "myText"
	logger.Errorf("%s - %s", text, "test error call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Info(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", infoLevel, "[testLogger] test info call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	logger.Info("test info call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Infof(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", infoLevel, "[testLogger] myText - test info call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	text := "myText"
	logger.Infof("%s - %s", text, "test info call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Warning(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", warningLevel, "[testLogger] test warning call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	logger.Warning("test warning call")

	// then
	// expectations asserted by mocks
}

func Test_libraryLogger_Warningf(t *testing.T) {
	// given
	loggerSink := newMockLogSink(t)
	loggerSink.On("Info", warningLevel, "[testLogger] myText - test warning call")
	logger := libraryLogger{name: "testLogger", logger: loggerSink}

	// when
	text := "myText"
	logger.Warningf("%s - %s", text, "test warning call")

	// then
	// expectations asserted by mocks
}
