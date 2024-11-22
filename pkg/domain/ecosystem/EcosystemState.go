package ecosystem

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type EcosystemState struct {
	InstalledDogus        map[cescommons.SimpleName]*DoguInstallation
	InstalledComponents   map[common.SimpleComponentName]*ComponentInstallation
	GlobalConfig          config.GlobalConfig
	ConfigByDogu          map[cescommons.SimpleName]config.DoguConfig
	SensitiveConfigByDogu map[cescommons.SimpleName]config.DoguConfig
}
