package config

import (
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewOperatorConfig(t *testing.T) {
	t.Run("should fail to read stage and parse version", func(t *testing.T) {
		// given
		t.Setenv(StageEnvVar, "")
		err := os.Unsetenv(StageEnvVar)
		require.NoError(t, err)

		oldLog := log
		defer func() { log = oldLog }()
		logMock := newMockLogSink(t)
		logMock.EXPECT().Init(mock.Anything).Return()
		logMock.EXPECT().Error(mock.Anything, "Error reading stage environment variable. Use stage production").Return()
		log = logr.New(logMock)

		// when
		actual, err := NewOperatorConfig("not-semver")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse version")
		assert.Nil(t, actual)
	})
	t.Run("should use development stage and fail to get namespace", func(t *testing.T) {
		// given
		t.Setenv(StageEnvVar, StageDevelopment)
		t.Setenv(namespaceEnvVar, "")
		err := os.Unsetenv(namespaceEnvVar)
		require.NoError(t, err)

		oldLog := log
		defer func() { log = oldLog }()
		logMock := newMockLogSink(t)
		logMock.EXPECT().Init(mock.Anything).Return()
		logMock.EXPECT().Enabled(0).Return(true)
		logMock.EXPECT().Info(0, "Version: [0.1.0]").Return()
		logMock.EXPECT().Info(0, "Starting in development mode! This is not recommended for production!").Return()
		log = logr.New(logMock)

		// when
		actual, err := NewOperatorConfig("0.1.0")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read namespace: failed to get env var [NAMESPACE]: environment variable NAMESPACE must be set")
		assert.Nil(t, actual)
	})
	t.Run("should use development stage and succeed", func(t *testing.T) {
		// given
		t.Setenv(StageEnvVar, StageDevelopment)
		t.Setenv(namespaceEnvVar, "ecosystem")

		oldLog := log
		defer func() { log = oldLog }()
		logMock := newMockLogSink(t)
		logMock.EXPECT().Init(mock.Anything).Return()
		logMock.EXPECT().Enabled(0).Return(true)
		logMock.EXPECT().Info(0, "Version: [0.1.0]").Return()
		logMock.EXPECT().Info(0, "Starting in development mode! This is not recommended for production!").Return()
		logMock.EXPECT().Info(0, "Deploying the k8s dogu operator in namespace ecosystem").Return()
		log = logr.New(logMock)

		// when
		actual, err := NewOperatorConfig("0.1.0")

		// then
		require.NoError(t, err)
		expected := &OperatorConfig{
			Version:   semver.MustParse("0.1.0"),
			Namespace: "ecosystem",
		}
		assert.Equal(t, expected, actual)
	})
}
