package ecosystem

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
)

type DoguHealthResult struct {
	UnhealthyDogus []UnhealthyDogu
}

type UnhealthyDogu struct {
	Namespace string
	Name      string
	Version   core.Version
	Health    HealthStatus
}

func (ud UnhealthyDogu) String() string {
	return fmt.Sprintf("%s/%s:%s is %s", ud.Namespace, ud.Name, ud.Version.Raw, ud.Health)
}
