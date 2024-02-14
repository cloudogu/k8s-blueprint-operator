package ecosystem

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

// ComponentInstallation represents an installed or to be installed component in the ecosystem.
type ComponentInstallation struct {
	// Name identifies the component by simple dogu name and namespace, e.g 'k8s/k8s-dogu-operator'.
	Name common.QualifiedComponentName
	// DeployNamespace is the cluster namespace where the component is deployed, e.g. `ecosystem` or `longhorn-system`
	// The default value is empty and indicated that the component should be deployed in the current namespace.
	// TODO: this field breaks the abstraction of the domain against kubernetes. We should discuss if a generic property list is better.
	DeployNamespace string
	// Version is the version of the component
	Version *semver.Version
	// Status is the installation status of the component in the ecosystem
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
	// Health is the current health status of the component in the ecosystem
	Health HealthStatus
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
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
	ComponentStatusIgnored  = "ignored"
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
func InstallComponent(componentName common.QualifiedComponentName, version *semver.Version) *ComponentInstallation {
	// TODO Delete this if the blueprint can handle a component configuration.
	// This section would contain the deployNamespace in a generic Map.
	var deployNamespace string

	if componentName == common.K8sK8sLonghornName {
		deployNamespace = "longhorn-system"
	}

	return &ComponentInstallation{
		Name:            componentName,
		Version:         version,
		DeployNamespace: deployNamespace,
		// ValuesYamlOverwrite: valuesYamlOverwrite,
	}
}

func (ci *ComponentInstallation) Upgrade(version *semver.Version) {
	ci.Version = version
}
