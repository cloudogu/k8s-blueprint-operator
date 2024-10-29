package ecosystem

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

var version1231, _ = core.ParseVersion("1.2.3-1")
var version1232, _ = core.ParseVersion("1.2.3-2")

func TestInstallDogu(t *testing.T) {
	volumeSize := resource.MustParse("1Gi")
	proxyBodySize := resource.MustParse("1G")
	dogu := InstallDogu(postgresqlQualifiedName, version1231, &volumeSize, ReverseProxyConfig{MaxBodySize: &proxyBodySize, RewriteTarget: "/", AdditionalConfig: "additional"})
	assert.Equal(t, &DoguInstallation{
		Name:          postgresqlQualifiedName,
		Version:       version1231,
		UpgradeConfig: UpgradeConfig{AllowNamespaceSwitch: false},
		MinVolumeSize: &volumeSize,
		ReverseProxyConfig: ReverseProxyConfig{
			MaxBodySize:      &proxyBodySize,
			RewriteTarget:    "/",
			AdditionalConfig: "additional",
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

	dogu.Upgrade(version1232)

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
			Name: common.QualifiedDoguName{
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
		dogu.UpdateProxyRewriteTarget("/")

		// then
		assert.Equal(t, RewriteTarget("/"), dogu.ReverseProxyConfig.RewriteTarget)
	})
}

func TestDoguInstallation_UpdateProxyAdditionalConfig(t *testing.T) {
	t.Run("should set property", func(t *testing.T) {
		// given
		dogu := DoguInstallation{}

		// when
		dogu.UpdateProxyAdditionalConfig("config")

		// then
		assert.Equal(t, AdditionalConfig("config"), dogu.ReverseProxyConfig.AdditionalConfig)
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
