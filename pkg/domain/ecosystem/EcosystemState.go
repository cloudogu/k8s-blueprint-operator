package ecosystem

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type EcosystemState struct {
	InstalledDogus        map[common.SimpleDoguName]*DoguInstallation
	InstalledComponents   map[common.SimpleComponentName]*ComponentInstallation
	GlobalConfig          config.GlobalConfig
	ConfigByDogu          map[common.SimpleDoguName]config.DoguConfig
	SensitiveConfigByDogu map[common.SimpleDoguName]config.DoguConfig
}
