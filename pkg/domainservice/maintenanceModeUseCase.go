package domainservice

import "fmt"

// MaintenancePageModel contains data that gets displayed when the maintenance mode is active.
type MaintenancePageModel struct {
	Title string
	Text  string
}

// MaintenanceModeUseCase contains methods to Activate and Deactivate the MaintenanceMode.
// When it is active, a page is displayed to the user, telling them that there is maintenance going on.
type MaintenanceModeUseCase struct {
	maintenanceMode MaintenanceMode
}

func NewMaintenanceModeUseCase(maintenanceMode MaintenanceMode) *MaintenanceModeUseCase {
	return &MaintenanceModeUseCase{maintenanceMode: maintenanceMode}
}

// Activate enables the maintenance mode, setting the given MaintenancePageModel
func (m *MaintenanceModeUseCase) Activate(content MaintenancePageModel) error {
	lock, err := m.maintenanceMode.GetLock()
	if err != nil {
		return fmt.Errorf("failed to check if maintenance mode is already active: %w", err)
	}

	if lock.IsActive() && !lock.IsOurs() {
		return fmt.Errorf("cannot activate maintenance mode as someone else already activated it")
	}

	err = m.maintenanceMode.Activate(content)
	if err != nil {
		return fmt.Errorf("failed to activate maintenance mode: %w", err)
	}

	return nil
}

// Deactivate disables the maintenance mode.
func (m *MaintenanceModeUseCase) Deactivate() error {
	lock, err := m.maintenanceMode.GetLock()
	if err != nil {
		return fmt.Errorf("failed to check if maintenance mode is already active: %w", err)
	}

	if !lock.IsActive() {
		// do nothing
		return nil
	}

	if !lock.IsOurs() {
		return fmt.Errorf("cannot deactivate maintenance mode as it was activated by another application")
	}

	err = m.maintenanceMode.Deactivate()
	if err != nil {
		return fmt.Errorf("failed to deactivate maintenance mode: %w", err)
	}

	return nil
}
