package ecosystem

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

// EcosystemState describes the actual state of the ecosystem, which is used to compare it with the expected state in the state diff.
type EcosystemState struct {
	InstalledDogus        map[cescommons.SimpleName]*DoguInstallation
	GlobalConfig          config.GlobalConfig
	ConfigByDogu          map[cescommons.SimpleName]config.DoguConfig
	SensitiveConfigByDogu map[cescommons.SimpleName]config.DoguConfig
}
