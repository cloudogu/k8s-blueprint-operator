package ecosystem

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
)

// DoguInstallation represents an installed or to be installed dogu in the ecosystem.
type DoguInstallation struct {
	// Namespace is the namespace of the dogu, e.g. 'official' like in 'official/postgresql'
	Namespace string
	// Name is the simple name of the dogu, e.g. 'postgresql' like in 'official/postgresql'.
	// the name is also the id of the dogu in the ecosystem as only one dogu with this name can be installed.
	Name string
	// Version is the version of the dogu
	Version core.Version
	// Status is the installation status of the dogu in the ecosystem
	Status string
	// Health is the current health status of the dogu in the ecosystem
	Health HealthStatus
	// UpgradeConfig contains configuration for dogu upgrades
	UpgradeConfig UpgradeConfig
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
}

const (
	DoguStatusNotInstalled = ""
	DoguStatusInstalling   = "installing"
	DoguStatusUpgrading    = "upgrading"
	DoguStatusDeleting     = "deleting"
	DoguStatusInstalled    = "installed"
	DoguStatusPVCResizing  = "resizing PVC"
)

type HealthStatus = string

const (
	PendingHealthStatus     HealthStatus = ""
	AvailableHealthStatus   HealthStatus = "available"
	UnavailableHealthStatus HealthStatus = "unavailable"
)

// UpgradeConfig contains configuration hints regarding aspects during the upgrade of dogus.
type UpgradeConfig struct {
	// AllowNamespaceSwitch lets a dogu switch its dogu namespace during an upgrade. The dogu must be technically the
	// same dogu which did reside in a different namespace. The remote dogu's version must be equal to or greater than
	// the version of the local dogu.
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch,omitempty"`
}

// GetQualifiedName returns the dogu name with namespace, e.g. official/postgresql
func (dogu *DoguInstallation) GetQualifiedName() string {
	return fmt.Sprintf("%s/%s", dogu.Namespace, dogu.Name)
}

// InstallDogu is a factory for new DoguInstallation's.
func InstallDogu(namespace string, doguName string, version core.Version) *DoguInstallation {
	return &DoguInstallation{
		Namespace:     namespace,
		Name:          doguName,
		Version:       version,
		UpgradeConfig: UpgradeConfig{AllowNamespaceSwitch: false},
	}
}

func (dogu *DoguInstallation) IsHealthy() bool {
	return dogu.Health == AvailableHealthStatus
}

func (dogu *DoguInstallation) Upgrade(newVersion core.Version) {
	dogu.Version = newVersion
	dogu.UpgradeConfig.AllowNamespaceSwitch = false
}

func (dogu *DoguInstallation) SwitchNamespace(newNamespace string, newVersion core.Version, isNamespaceSwitchAllowed bool) error {
	if !isNamespaceSwitchAllowed {
		return fmt.Errorf("not allowed to switch dogu namespace from %q to %q", dogu.Namespace, newNamespace)
	}
	dogu.Namespace = newNamespace
	dogu.Version = newVersion
	dogu.UpgradeConfig.AllowNamespaceSwitch = true
	return nil
}
