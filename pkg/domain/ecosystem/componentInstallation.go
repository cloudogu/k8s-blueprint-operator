package ecosystem

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

// ComponentInstallation represents an installed or to be installed component in the ecosystem.
type ComponentInstallation struct {
	// Name identifies the component by simple dogu name and namespace, e.g 'k8s/k8s-dogu-operator'.
	Name common.QualifiedComponentName
	// ExpectedVersion is the version of the component which should be installed
	ExpectedVersion *semver.Version
	// ActualVersion is the version of the component which is actually installed
	ActualVersion *semver.Version
	// Status is the installation status of the component in the ecosystem
	Status string
	// Health is the current health status of the component in the ecosystem
	Health HealthStatus
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
	DeployConfig       DeployConfig
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
func InstallComponent(
	componentName common.QualifiedComponentName,
	expectedVersion *semver.Version,
	deployConfig DeployConfig,
) *ComponentInstallation {
	return &ComponentInstallation{
		Name:            componentName,
		ExpectedVersion: expectedVersion,
		DeployConfig:    deployConfig,
	}
}

func (ci *ComponentInstallation) Upgrade(expectedVersion *semver.Version) {
	ci.ExpectedVersion = expectedVersion
}

func (ci *ComponentInstallation) UpdateDeployConfig(deployConfig DeployConfig) {
	ci.DeployConfig = deployConfig
}
