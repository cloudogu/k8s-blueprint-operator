package ecosystem

import (
	"github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type EcosystemState struct {
	InstalledDogus        map[cescommons.SimpleName]*DoguInstallation
	InstalledComponents   map[v2.SimpleComponentName]*ComponentInstallation
	GlobalConfig          config.GlobalConfig
	ConfigByDogu          map[cescommons.SimpleName]config.DoguConfig
	SensitiveConfigByDogu map[cescommons.SimpleName]config.DoguConfig
}
