package ecosystem

import (
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"slices"
	"strings"
)

// DoguHealthResult is a snapshot of the health states of all dogus.
type DoguHealthResult struct {
	DogusByStatus map[HealthStatus][]cescommons.SimpleName
}

func (result DoguHealthResult) getUnhealthyDogus() []cescommons.SimpleName {
	var unhealthyDogus []cescommons.SimpleName
	for healthState, doguNames := range result.DogusByStatus {
		if healthState != AvailableHealthStatus {
			unhealthyDogus = append(unhealthyDogus, doguNames...)
		}
	}
	return unhealthyDogus
}

func (result DoguHealthResult) String() string {
	unhealthyDogus := util.Map(result.getUnhealthyDogus(), func(dogu cescommons.SimpleName) string { return string(dogu) })
	slices.Sort(unhealthyDogus)
	return fmt.Sprintf("%d dogu(s) are unhealthy: %s", len(unhealthyDogus), strings.Join(unhealthyDogus, ", "))
}

// CalculateDoguHealthResult collects the health states from DoguInstallation and creates a DoguHealthResult.
func CalculateDoguHealthResult(dogus []*DoguInstallation) DoguHealthResult {
	result := DoguHealthResult{
		DogusByStatus: map[HealthStatus][]cescommons.SimpleName{},
	}
	for _, dogu := range dogus {
		result.DogusByStatus[dogu.Health] = append(result.DogusByStatus[dogu.Health], dogu.Name.SimpleName)
	}
	return result
}

func (result DoguHealthResult) AllHealthy() bool {
	for healthState, doguNames := range result.DogusByStatus {
		if healthState != AvailableHealthStatus && len(doguNames) != 0 {
			return false
		}
	}
	return true
}
