package ecosystem

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
	"golang.org/x/exp/maps"
)

type EcosystemState struct {
	InstalledDogus      map[common.SimpleDoguName]*DoguInstallation
	InstalledComponents map[common.SimpleComponentName]*ComponentInstallation
	GlobalConfig        config.GlobalConfig
	DoguConfig          map[common.DoguConfigKey]*DoguConfigEntry
	SensitiveDoguConfig map[common.SensitiveDoguConfigKey]*SensitiveDoguConfigEntry
}

func (state EcosystemState) GetInstalledDoguNames() []common.SimpleDoguName {
	return maps.Keys(state.InstalledDogus)
}
