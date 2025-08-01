package config

import (
	"github.com/cloudogu/cesapp-lib/core"
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

func TestGetRemoteConfiguration(t *testing.T) {
	type args struct {
		endpoint  string
		urlSchema string
	}
	tests := []struct {
		name    string
		args    args
		want    *core.Remote
		wantErr assert.ErrorAssertionFunc
		setEnv  func(t *testing.T)
	}{
		{
			name:    "test default url schema",
			want:    &core.Remote{Endpoint: "https://example.com/", URLSchema: "default", CacheDir: "/tmp/dogu-registry-cache"},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "default")
			},
		},
		{
			name:    "test dogu registry endpoint environment variable not set",
			want:    nil,
			wantErr: assert.Error,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "default")
			},
		},
		{
			name:    "test default url schema with 'dogus' suffix",
			want:    &core.Remote{Endpoint: "https://example.com/", URLSchema: "default", CacheDir: "/tmp/dogu-registry-cache"},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "default")
			},
		},
		{
			name:    "test default url schema with 'dogus/' suffix",
			want:    &core.Remote{Endpoint: "https://example.com/", URLSchema: "default", CacheDir: "/tmp/dogu-registry-cache"},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "default")
			},
		},
		{
			name:    "test non-default url schema",
			want:    &core.Remote{Endpoint: "https://example.com/", URLSchema: "index", CacheDir: "/tmp/dogu-registry-cache"},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "index")
			},
		},
		{
			name:    "test non-default url schema with 'dogus' suffix",
			want:    &core.Remote{Endpoint: "https://example.com/dogus", URLSchema: "index", CacheDir: "/tmp/dogu-registry-cache"},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "index")
			},
		},
		{
			name:    "test non-default url schema with 'dogus/' suffix",
			want:    &core.Remote{Endpoint: "https://example.com/dogus/", URLSchema: "index", CacheDir: "/tmp/dogu-registry-cache"},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "index")
			},
		},
		{
			name: "test with proxy",
			want: &core.Remote{Endpoint: "https://example.com/dogus/", URLSchema: "index", CacheDir: "/tmp/dogu-registry-cache", ProxySettings: core.ProxySettings{
				Enabled:  true,
				Server:   "host",
				Port:     3128,
				Username: "user",
				Password: "password",
			}},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "index")
				t.Setenv("PROXY_URL", "https://user:password@host:3128")
			},
		},
		{
			name:    "test proxy invalid url",
			args:    args{endpoint: "https://example.com/dogus/", urlSchema: "index"},
			want:    nil,
			wantErr: assert.Error,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "index")
				t.Setenv("PROXY_URL", "://f")
			},
		},
		{
			name:    "test proxy invalid port",
			args:    args{endpoint: "https://example.com/dogus/", urlSchema: "index"},
			want:    nil,
			wantErr: assert.Error,
			setEnv: func(t *testing.T) {
				t.Setenv("DOGU_REGISTRY_ENDPOINT", "https://example.com/dogus/")
				t.Setenv("DOGU_REGISTRY_URLSCHEMA", "index")
				t.Setenv("PROXY_URL", "https://user:password@host:invalid")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv != nil {
				tt.setEnv(t)
			}

			config, err := GetRemoteConfiguration()

			tt.wantErr(t, err)

			assert.Equalf(t, tt.want, config, "getRemoteConfig(%v, %v)", tt.args.endpoint, tt.args.urlSchema)
		})
	}
}

func TestGetRemoteCredentials(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		t.Setenv(doguRegistryUsernameEnvVar, "user")
		t.Setenv(doguRegistryPasswordEnvVar, "pass")
		config, err := GetRemoteCredentials()

		require.NoError(t, err)
		assert.Equal(t, "user", config.Username)
		assert.Equal(t, "pass", config.Password)
	})
	t.Run("no user", func(t *testing.T) {
		t.Setenv(doguRegistryPasswordEnvVar, "pass")
		_, err := GetRemoteCredentials()

		require.Error(t, err)
	})
	t.Run("no pass", func(t *testing.T) {
		t.Setenv(doguRegistryUsernameEnvVar, "user")
		_, err := GetRemoteCredentials()

		require.Error(t, err)
	})
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "test log level not set",
			wantErr: assert.Error,
		},
		{
			name:    "test log level set to debug",
			want:    "debug",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != "" {
				t.Setenv(logLevelEnvVar, tt.want)
			} else {
				// first set it so it got rolled back afterward
				t.Setenv(logLevelEnvVar, "")
				// then unset it, so environments with this envVar also work with this test
				err := os.Unsetenv(logLevelEnvVar)
				if err != nil {
					require.NoError(t, err)
				}
			}
			got, err := GetLogLevel()
			if !tt.wantErr(t, err, "GetLogLevel()") {
				return
			}
			assert.Equalf(t, tt.want, got, "GetLogLevel()")
		})
	}
}
