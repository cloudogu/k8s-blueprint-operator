package ecosystem

import (
	"github.com/cloudogu/cesapp-lib/core"
)

// ComponentInstallation represents an installed or to be installed component in the ecosystem.
type ComponentInstallation struct {
	// Name is the name of the component, e.g. 'k8s-dogu-operator'.
	// The name is also the ID of the component in the ecosystem as only one component with this name can be installed.
	Name string
	// Version is the version of the component
	Version core.Version
	// Status is the installation status of the component in the ecosystem
	Status string
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
}

const (
	ComponentStatusNotInstalled = ""
	ComponentStatusInstalling   = "installing"
	ComponentStatusUpgrading    = "upgrading"
	ComponentStatusDeleting     = "deleting"
	ComponentStatusInstalled    = "installed"
	ComponentStatusIgnored      = "ignored"
)
