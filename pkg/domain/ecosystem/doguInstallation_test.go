package ecosystem

import (
	"testing"
	"time"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	version1231, _   = core.ParseVersion("1.2.3-1")
	version1232, _   = core.ParseVersion("1.2.3-2")
	rewriteTarget    = "/"
	additionalConfig = "additional"
	subfolder        = "different_subfolder"
)

func TestInstallDogu(t *testing.T) {
	volumeSize := resource.MustParse("1Gi")
	proxyBodySize := resource.MustParse("1G")
	dogu := InstallDogu(
		postgresqlQualifiedName,
		&version1231,
		&volumeSize,
		ReverseProxyConfig{MaxBodySize: &proxyBodySize, RewriteTarget: RewriteTarget(rewriteTarget), AdditionalConfig: AdditionalConfig(additionalConfig)},
		[]AdditionalMount{
			{
				SourceType: DataSourceConfigMap,
				Name:       "configmap",
				Volume:     "volume",
				Subfolder:  subfolder,
			},
		},
	)
	assert.Equal(t, &DoguInstallation{
		Name:          postgresqlQualifiedName,
		Version:       version1231,
		UpgradeConfig: UpgradeConfig{AllowNamespaceSwitch: false},
		MinVolumeSize: &volumeSize,
		ReverseProxyConfig: ReverseProxyConfig{
			MaxBodySize:      &proxyBodySize,
			RewriteTarget:    RewriteTarget(rewriteTarget),
			AdditionalConfig: AdditionalConfig(additionalConfig),
		},
		AdditionalMounts: []AdditionalMount{
			{
				SourceType: DataSourceConfigMap,
				Name:       "configmap",
				Volume:     "volume",
				Subfolder:  subfolder,
			},
		},
	}, dogu)
}

func TestDoguInstallation_IsHealthy(t *testing.T) {
	t.Run("is healthy", func(t *testing.T) {
		dogu := &DoguInstallation{
			Name:   postgresqlQualifiedName,
			Health: AvailableHealthStatus,
		}

		isHealthy := dogu.IsHealthy()

		assert.True(t, isHealthy)
	})

	t.Run("is unhealthy", func(t *testing.T) {
		dogu := &DoguInstallation{
			Name:   postgresqlQualifiedName,
			Health: UnavailableHealthStatus,
		}

		isHealthy := dogu.IsHealthy()

		assert.False(t, isHealthy)
	})
}

func TestDoguInstallation_Upgrade(t *testing.T) {
	dogu := &DoguInstallation{
		Name:    postgresqlQualifiedName,
		Version: version1231,
	}

	dogu.Upgrade(&version1232)

	assert.Equal(t, &DoguInstallation{
		Name:    postgresqlQualifiedName,
		Version: version1232,
	}, dogu)
}

func TestDoguInstallation_SwitchNamespace(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		dogu := &DoguInstallation{
			Name: postgresqlQualifiedName,
		}

		err := dogu.SwitchNamespace("premium", true)

		require.NoError(t, err)
		assert.Equal(t, &DoguInstallation{
			Name: cescommons.QualifiedName{
				Namespace:  "premium",
				SimpleName: "postgresql",
			},
			UpgradeConfig: UpgradeConfig{
				AllowNamespaceSwitch: true,
			},
		}, dogu)
	})

	t.Run("namespace switch not allowed", func(t *testing.T) {
		dogu := &DoguInstallation{
			Name: postgresqlQualifiedName,
		}

		err := dogu.SwitchNamespace("premium", false)

		require.ErrorContains(t, err, "not allowed to switch dogu namespace")
	})
}

func TestDoguInstallation_UpdateProxyBodySize(t *testing.T) {
	t.Run("should set property", func(t *testing.T) {
		// given
		bodySize := resource.MustParse("1G")
		dogu := DoguInstallation{}

		// when
		dogu.UpdateProxyBodySize(&bodySize)

		// then
		assert.Equal(t, &bodySize, dogu.ReverseProxyConfig.MaxBodySize)
	})
}

func TestDoguInstallation_UpdateProxyRewriteTarget(t *testing.T) {
	t.Run("should set property", func(t *testing.T) {
		// given
		dogu := DoguInstallation{}

		// when
		dogu.UpdateProxyRewriteTarget(RewriteTarget(rewriteTarget))

		// then
		assert.Equal(t, RewriteTarget(rewriteTarget), dogu.ReverseProxyConfig.RewriteTarget)
	})
}

func TestDoguInstallation_UpdateProxyAdditionalConfig(t *testing.T) {
	t.Run("should set property", func(t *testing.T) {
		// given
		dogu := DoguInstallation{}

		// when
		dogu.UpdateProxyAdditionalConfig(AdditionalConfig(additionalConfig))

		// then
		assert.Equal(t, AdditionalConfig(additionalConfig), dogu.ReverseProxyConfig.AdditionalConfig)
	})
}

func TestDoguInstallation_UpdateMinVolumeSize(t *testing.T) {
	t.Run("should set property", func(t *testing.T) {
		// given
		volumeSize := resource.MustParse("1Gi")
		dogu := DoguInstallation{}

		// when
		dogu.UpdateMinVolumeSize(&volumeSize)

		// then
		assert.Equal(t, &volumeSize, dogu.MinVolumeSize)
	})
}

func TestDoguInstallation_IsVersionUpToDate(t *testing.T) {
	type fields struct {
		Version          core.Version
		InstalledVersion core.Version
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "equal Versions is up to date",
			fields: fields{
				Version:          version1231,
				InstalledVersion: version1231,
			},
			want: true,
		},
		{
			name: "Version newer than installed version is not up to date",
			fields: fields{
				Version:          version1232,
				InstalledVersion: version1231,
			},
			want: false,
		},
		{
			name: "Installed version empty is not up to date",
			fields: fields{
				Version:          version1232,
				InstalledVersion: core.Version{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dogu := &DoguInstallation{
				Version:          tt.fields.Version,
				InstalledVersion: tt.fields.InstalledVersion,
			}
			assert.Equalf(t, tt.want, dogu.IsVersionUpToDate(), "IsVersionUpToDate()")
		})
	}
}

func TestDoguInstallation_IsConfigUpToDate(t *testing.T) {
	timeMay := v1.NewTime(time.Date(2024, time.May, 23, 10, 0, 0, 0, time.UTC))
	timeJune := v1.NewTime(time.Date(2024, time.June, 23, 10, 0, 0, 0, time.UTC))
	timeJuly := v1.NewTime(time.Date(2024, time.July, 23, 10, 0, 0, 0, time.UTC))
	type args struct {
		globalConfigUpdateTime *v1.Time
		doguConfigUpdateTime   *v1.Time
	}
	tests := []struct {
		name      string
		StartedAt v1.Time
		args      args
		want      bool
	}{
		{
			name:      "all equal is up to date",
			StartedAt: timeMay,
			args: args{
				globalConfigUpdateTime: &timeMay,
				doguConfigUpdateTime:   &timeMay,
			},
			want: true,
		},
		{
			name:      "StartedAt newest date is up to date",
			StartedAt: timeJuly,
			args: args{
				globalConfigUpdateTime: &timeMay,
				doguConfigUpdateTime:   &timeJune,
			},
			want: true,
		},
		{
			name:      "globalConfigUpdateTime newest date is not up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: &timeJuly,
				doguConfigUpdateTime:   &timeMay,
			},
			want: false,
		},
		{
			name:      "doguConfigUpdateTime newest date is not up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: &timeMay,
				doguConfigUpdateTime:   &timeJuly,
			},
			want: false,
		},
		{
			name:      "doguConfigUpdateTime nil and StartedAt newest is up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: &timeMay,
				doguConfigUpdateTime:   nil,
			},
			want: true,
		},
		{
			name:      "doguConfigUpdateTime nil and StartedAt not newest is not up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: &timeJuly,
				doguConfigUpdateTime:   nil,
			},
			want: false,
		},
		{
			name:      "globalConfigUpdateTime nil and StartedAt newest is up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: nil,
				doguConfigUpdateTime:   &timeMay,
			},
			want: true,
		},
		{
			name:      "globalConfigUpdateTime nil and StartedAt not newest is not up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: nil,
				doguConfigUpdateTime:   &timeJuly,
			},
			want: false,
		},
		{
			name:      "globalConfigUpdateTime nil and doguConfigUpdateTime nil is up to date",
			StartedAt: timeJune,
			args: args{
				globalConfigUpdateTime: nil,
				doguConfigUpdateTime:   nil,
			},
			want: true,
		},
		{
			name:      "StartedAt empty is not up to date",
			StartedAt: v1.Time{},
			args: args{
				globalConfigUpdateTime: &timeMay,
				doguConfigUpdateTime:   nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dogu := &DoguInstallation{
				StartedAt: tt.StartedAt,
			}
			assert.Equalf(t, tt.want, dogu.IsConfigUpToDate(tt.args.globalConfigUpdateTime, tt.args.doguConfigUpdateTime), "IsConfigUpToDate(%v, %v)", tt.args.globalConfigUpdateTime, tt.args.doguConfigUpdateTime)
		})
	}
}
