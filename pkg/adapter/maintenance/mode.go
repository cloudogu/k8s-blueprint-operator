package maintenance

import (
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

// Mode contains methods to Activate and Deactivate the switcher.
// When it is active, a page is displayed to the user, telling them that there is maintenance going on.
type Mode struct {
	lock
	switcher
}

func New(globalConfig registry.ConfigurationContext) *Mode {
	return &Mode{
		lock:     &defaultLock{globalConfig: globalConfig},
		switcher: &defaultSwitcher{globalConfig: globalConfig},
	}
}

// Activate enables the maintenance mode, setting the given MaintenancePageModel
func (m *Mode) Activate(content domainservice.MaintenancePageModel) error {
	isActive, isOurs, err := m.lock.isActiveAndOurs()
	if err != nil {
		return domainservice.NewInternalError(err, "failed to check if maintenance mode is already active and ours")
	}

	if isActive && !isOurs {
		return domainservice.NewConflictError(nil, "cannot activate maintenance mode as it was already activated by another party")
	}

	err = m.switcher.activate(content)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to activate maintenance mode")
	}

	return nil
}

// Deactivate disables the maintenance mode.
func (m *Mode) Deactivate() error {
	isActive, isOurs, err := m.lock.isActiveAndOurs()
	if err != nil {
		return domainservice.NewInternalError(err, "failed to check if maintenance mode is already active and ours")
	}

	if !isActive {
		// do nothing
		return nil
	}

	if !isOurs {
		return domainservice.NewConflictError(nil, "cannot deactivate maintenance mode as it was activated by another party")
	}

	err = m.switcher.deactivate()
	if err != nil {
		return domainservice.NewInternalError(err, "failed to deactivate maintenance mode")
	}

	return nil
}
