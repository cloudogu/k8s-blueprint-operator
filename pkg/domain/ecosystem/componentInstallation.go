package ecosystem

import (
	"github.com/cloudogu/cesapp-lib/core"
)

// ComponentInstallation represents an installed or to be installed component in the ecosystem.
type ComponentInstallation struct {
	// Namespace is the namespace of the component, e.g. 'k8s' like in 'k8s/k8s-dogu-operator'.
	Namespace string
	// Name is the simple name of the component, e.g. `k8s-dogu-operator`.
	// The name is also the id of the component in the ecosystem as only one component with this name can be installed.
	Name string
	// DeployNamespace is the cluster namespace where the component is deployed, e.g. `ecosystem` or `longhorn-system`
	// The default value is empty and indicated that the component should be deployed in the current namespace.
	DeployNamespace string
	// Version is the version of the component.
	Version core.Version
	// Status is the installation status of the component in the ecosystem.
	Status string
	// ValuesYamlOverwrite represents a helm configuration as string in yaml format.
	// Example:
	// ```
	// controller:
	//   env:
	//     logLevel: info
	// ```
	ValuesYamlOverwrite string
	// MappedValues represents also a helm configuration like ValuesYamlOverwrite.
	// The difference here is that these values will be mapped by the component-operator with a metadata file in the component's chart.
	MappedValues map[string]string
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
	// Health is the current health status of the component in the ecosystem
	Health HealthStatus
}

const (
	// ComponentStatusNotInstalled represents a status for a component that is not installed
	ComponentStatusNotInstalled = ""
	// ComponentStatusInstalling represents a status for a component that is currently being installed
	ComponentStatusInstalling = "installing"
	// ComponentStatusUpgrading represents a status for a component that is currently being upgraded
	ComponentStatusUpgrading = "upgrading"
	// ComponentStatusDeleting represents a status for a component that is currently being deleted
	ComponentStatusDeleting = "deleting"
	// ComponentStatusInstalled represents a status for a component that was successfully installed
	ComponentStatusInstalled = "installed"
	// ComponentStatusTryToInstall represents a status for a component that is not installed but its install process is in requeue loop.
	ComponentStatusTryToInstall = "tryToInstall"
	// ComponentStatusTryToUpgrade represents a status for a component that is installed but its actual upgrade process is in requeue loop.
	// In this state the component can be healthy but the version in the spec is not installed.
	ComponentStatusTryToUpgrade = "tryToUpgrade"
	// ComponentStatusTryToDelete represents a status for a component that is installed but its delete process is in requeue loop.
	// In this state the component can be healthy.
	ComponentStatusTryToDelete = "tryToDelete"
)

// InstallComponent is a factory for new ComponentInstallation's.
func InstallComponent(namespace, componentName string, version core.Version) *ComponentInstallation {
	return &ComponentInstallation{
		Namespace: namespace,
		Name:      componentName,
		Version:   version,
		// DeployNamespace:     deployNamespace,
		// ValuesYamlOverwrite: valuesYamlOverwrite,
	}
}

func (ci *ComponentInstallation) Upgrade(version core.Version) {
	ci.Version = version
}
