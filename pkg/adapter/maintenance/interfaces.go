package maintenance

import (
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type globalConfig interface {
	registry.ConfigurationContext
}

// switcher provides ways to activate and deactivate the maintenance mode.
type switcher interface {
	// Activate enables the maintenance mode.
	activate(content domainservice.MaintenancePageModel) error
	// Deactivate disables the maintenance mode.
	deactivate() error
}

type lock interface {
	// isActiveAndOurs returns two bools that determine if the maintenance mode is active and ours.
	isActiveAndOurs() (bool, bool, error)
}
