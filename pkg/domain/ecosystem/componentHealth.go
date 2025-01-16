package ecosystem

import (
	"fmt"
	"github.com/cloudogu/blueprint-lib/v2"
	"slices"
	"strings"
	"time"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

type RequiredComponent struct {
	Name v2.SimpleComponentName
}

type WaitConfig struct {
	Timeout  time.Duration
	Interval time.Duration
}

// ComponentHealthResult is a snapshot of all components' health states.
type ComponentHealthResult struct {
	ComponentsByStatus map[HealthStatus][]v2.SimpleComponentName
}

func (result ComponentHealthResult) getUnhealthyComponents() []v2.SimpleComponentName {
	var unhealthyComponents []v2.SimpleComponentName
	for healthState, componentNames := range result.ComponentsByStatus {
		if healthState != AvailableHealthStatus {
			unhealthyComponents = append(unhealthyComponents, componentNames...)
		}
	}
	return unhealthyComponents
}

func (result ComponentHealthResult) String() string {
	unhealthyComponents := util.Map(result.getUnhealthyComponents(), func(dogu v2.SimpleComponentName) string { return string(dogu) })
	slices.Sort(unhealthyComponents)
	return fmt.Sprintf("%d component(s) are unhealthy: %s", len(unhealthyComponents), strings.Join(unhealthyComponents, ", "))
}

// CalculateComponentHealthResult checks if all required components are installed,
// collects the health states from ComponentInstallation and creates a ComponentHealthResult.
func CalculateComponentHealthResult(installedComponents map[v2.SimpleComponentName]*ComponentInstallation, requiredComponents []RequiredComponent) ComponentHealthResult {
	result := ComponentHealthResult{
		ComponentsByStatus: map[HealthStatus][]v2.SimpleComponentName{},
	}
	for _, required := range requiredComponents {
		_, installed := installedComponents[required.Name]
		if !installed {
			result.ComponentsByStatus[NotInstalledHealthStatus] = append(result.ComponentsByStatus[NotInstalledHealthStatus], v2.SimpleComponentName(required.Name))
		}
	}
	for _, component := range installedComponents {
		result.ComponentsByStatus[component.Health] = append(result.ComponentsByStatus[component.Health], component.Name.SimpleName)
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
