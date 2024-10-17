package maintenance

import (
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

// Mode contains methods to Activate and Deactivate the switcher.
// When it is active, a page is displayed to the user, telling them that there is maintenance going on.
type Mode struct {
}

func New(globalConfig registry.ConfigurationContext) *Mode {
	return &Mode{}
}

// Activate enables the maintenance mode, setting the given MaintenancePageModel
func (m *Mode) Activate(content domainservice.MaintenancePageModel) error {
	// TODO: replace this completely with maintenance mode from k8s-registry-lib
	return nil
}

// Deactivate disables the maintenance mode.
func (m *Mode) Deactivate() error {
	// TODO: replace this completely with maintenance mode from k8s-registry-lib
	return nil
}
