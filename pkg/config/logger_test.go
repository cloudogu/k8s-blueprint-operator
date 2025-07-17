package config

import (
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	uberzap "go.uber.org/zap"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestConfigureLogger(t *testing.T) {
	originalSetLogger := ctrl.SetLogger
	defer func() {
		ctrl.SetLogger = originalSetLogger
	}()
	t.Setenv(logLevelEnvVar, "debug")
	Stage = StageDevelopment

	ctrl.SetLogger = func(l logr.Logger) {
		require.NotNil(t, l)
		assert.NotNil(t, l.GetSink())
	}

	ConfigureLogger()
}

func Test_getZapOptions(t *testing.T) {
	t.Run("should configure devMode and level", func(t *testing.T) {
		t.Setenv(logLevelEnvVar, "debug")
		Stage = StageDevelopment

		options := getZapOptions()

		assert.Equal(t, options.Level, uberzap.NewAtomicLevelAt(uberzap.DebugLevel))
		assert.True(t, options.Development)
	})

	t.Run("should use default level when not configured", func(t *testing.T) {
		err := os.Unsetenv(logLevelEnvVar)
		require.NoError(t, err)
		Stage = StageProduction

		options := getZapOptions()

		assert.Equal(t, options.Level, uberzap.NewAtomicLevelAt(uberzap.InfoLevel))
		assert.False(t, options.Development)
	})

	t.Run("should use default level when invalid", func(t *testing.T) {
		t.Setenv(logLevelEnvVar, "invalid")
		Stage = StageDevelopment

		options := getZapOptions()

		assert.Equal(t, options.Level, uberzap.NewAtomicLevelAt(uberzap.InfoLevel))
		assert.True(t, options.Development)
	})
}
