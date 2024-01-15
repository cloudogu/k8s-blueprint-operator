package ecosystem

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"slices"
	"strings"
)

type DoguName string

type DoguHealthResult struct {
	DogusByStatus map[HealthStatus][]DoguName
}

func (result DoguHealthResult) GetUnhealthyDogus() []DoguName {
	var unhealthyDogus []DoguName
	for healthState, doguNames := range result.DogusByStatus {
		if healthState != AvailableHealthStatus {
			unhealthyDogus = append(unhealthyDogus, doguNames...)
		}
	}
	return unhealthyDogus
}

func (result DoguHealthResult) String() string {
	unhealthyDogus := util.Map(result.GetUnhealthyDogus(), func(dogu DoguName) string { return string(dogu) })
	slices.Sort(unhealthyDogus)
	return fmt.Sprintf("%d dogus are unhealthy: %s", len(unhealthyDogus), strings.Join(unhealthyDogus, ", "))
}

func CalculateDoguHealthResult(dogus []*DoguInstallation) DoguHealthResult {
	result := DoguHealthResult{
		DogusByStatus: map[HealthStatus][]DoguName{},
	}
	for _, dogu := range dogus {
		result.DogusByStatus[dogu.Health] = append(result.DogusByStatus[dogu.Health], DoguName(dogu.Name))
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
