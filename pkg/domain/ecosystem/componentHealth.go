package ecosystem

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type ComponentName string

// ComponentHealthResult is a snapshot of all components' health states.
type ComponentHealthResult struct {
	ComponentsByStatus map[HealthStatus][]ComponentName
}

func (result ComponentHealthResult) getUnhealthyComponents() []ComponentName {
	var unhealthyComponents []ComponentName
	for healthState, componentNames := range result.ComponentsByStatus {
		if healthState != AvailableHealthStatus {
			unhealthyComponents = append(unhealthyComponents, componentNames...)
		}
	}
	return unhealthyComponents
}

func (result ComponentHealthResult) String() string {
	unhealthyComponents := util.Map(result.getUnhealthyComponents(), func(dogu ComponentName) string { return string(dogu) })
	slices.Sort(unhealthyComponents)
	return fmt.Sprintf("%d components are unhealthy: %s", len(unhealthyComponents), strings.Join(unhealthyComponents, ", "))
}

// CalculateComponentHealthResult checks if all required components are installed,
// collects the health states from ComponentInstallation and creates a ComponentHealthResult.
func CalculateComponentHealthResult(installedComponents map[string]*ComponentInstallation, requiredComponents []domain.RequiredComponent) ComponentHealthResult {
	result := ComponentHealthResult{
		ComponentsByStatus: map[HealthStatus][]ComponentName{},
	}
	for _, required := range requiredComponents {
		_, installed := installedComponents[required.Name]
		if !installed {
			result.ComponentsByStatus[NotInstalledHealthStatus] = append(result.ComponentsByStatus[NotInstalledHealthStatus], ComponentName(required.Name))
		}
	}
	for _, component := range installedComponents {
		result.ComponentsByStatus[component.Health] = append(result.ComponentsByStatus[component.Health], ComponentName(component.Name))
	}
	return result
}

func (result ComponentHealthResult) AllHealthy() bool {
	for healthState, componentNames := range result.ComponentsByStatus {
		if healthState != AvailableHealthStatus && len(componentNames) != 0 {
			return false
		}
	}
	return true
}
