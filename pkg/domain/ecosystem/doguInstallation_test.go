package ecosystem

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var version1_2_3_1, _ = core.ParseVersion("1.2.3-1")
var version1_2_3_2, _ = core.ParseVersion("1.2.3-2")

func TestDoguInstallation_GetQualifiedName(t *testing.T) {
	dogu := DoguInstallation{
		Namespace: "official",
		Name:      "postgresql",
	}

	name := dogu.GetQualifiedName()

	assert.Equal(t, "official/postgresql", name)
}

func TestInstallDogu(t *testing.T) {
	assert.Equal(t, &DoguInstallation{
		Namespace:     "official",
		Name:          "postgresql",
		Version:       version1_2_3_1,
		UpgradeConfig: UpgradeConfig{AllowNamespaceSwitch: false},
	}, InstallDogu("official", "postgresql", version1_2_3_1))
}

func TestDoguInstallation_IsHealthy(t *testing.T) {
	t.Run("is healthy", func(t *testing.T) {
		dogu := &DoguInstallation{
			Name:   "postgresql",
			Health: AvailableHealthStatus,
		}

		isHealthy := dogu.IsHealthy()

		assert.True(t, isHealthy)
	})

	t.Run("is unhealthy", func(t *testing.T) {
		dogu := &DoguInstallation{
			Name:   "postgresql",
			Health: UnavailableHealthStatus,
		}

		isHealthy := dogu.IsHealthy()

		assert.False(t, isHealthy)
	})
}

func TestDoguInstallation_Upgrade(t *testing.T) {
	dogu := &DoguInstallation{
		Namespace: "official",
		Name:      "postgresql",
		Version:   version1_2_3_1,
	}

	dogu.Upgrade(version1_2_3_2)

	assert.Equal(t, &DoguInstallation{
		Namespace: "official",
		Name:      "postgresql",
		Version:   version1_2_3_2,
	}, dogu)
}

func TestDoguInstallation_SwitchNamespace(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		dogu := &DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version1_2_3_1,
		}

		err := dogu.SwitchNamespace("premium", version1_2_3_2, true)

		require.NoError(t, err)
		assert.Equal(t, &DoguInstallation{
			Namespace: "premium",
			Name:      "postgresql",
			Version:   version1_2_3_2,
			UpgradeConfig: UpgradeConfig{
				AllowNamespaceSwitch: true,
			},
		}, dogu)
	})

	t.Run("namespace switch not allowed", func(t *testing.T) {
		dogu := &DoguInstallation{
			Namespace: "official",
			Name:      "postgresql",
			Version:   version1_2_3_1,
		}

		err := dogu.SwitchNamespace("premium", version1_2_3_2, false)

		require.ErrorContains(t, err, "not allowed to switch dogu namespace")
	})
}
