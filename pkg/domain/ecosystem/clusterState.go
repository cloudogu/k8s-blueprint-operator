package ecosystem

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"golang.org/x/exp/maps"
)

type ClusterState struct {
	InstalledDogus               map[common.SimpleDoguName]*DoguInstallation
	InstalledComponents          map[common.SimpleComponentName]*ComponentInstallation
	GlobalConfig                 map[common.GlobalConfigKey]*GlobalConfigEntry
	DoguConfig                   map[common.DoguConfigKey]*DoguConfigEntry
	EncryptedDoguConfig          map[common.SensitiveDoguConfigKey]*SensitiveDoguConfigEntry
	DecryptedSensitiveDoguConfig map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue
}

func (state ClusterState) GetInstalledDoguNames() []common.SimpleDoguName {
	return maps.Keys(state.InstalledDogus)
}
