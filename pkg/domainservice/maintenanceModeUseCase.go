package domainservice

import "fmt"

type MaintenancePageModel struct {
	Title string
	Text  string
}

type MaintenanceModeUseCase struct {
	maintenanceMode MaintenanceMode
}

func (m *MaintenanceModeUseCase) Activate(content MaintenancePageModel) error {
	lock, err := m.maintenanceMode.GetLock()
	if err != nil {
		return fmt.Errorf("failed to check if maintenance mode is already active: %w", err)
	}

	if lock.IsActive() {
		if lock.IsOurs() {
			// do nothing
			return nil
		}

		return fmt.Errorf("cannot activate maintenance mode as someone else already activated it")
	}

	err = m.maintenanceMode.Activate(content)
	if err != nil {
		return fmt.Errorf("failed to activate maintenance mode: %w", err)
	}

	return nil
}

func (m *MaintenanceModeUseCase) Deactivate() error {
	lock, err := m.maintenanceMode.GetLock()
	if err != nil {
		return fmt.Errorf("failed to check if maintenance mode is already active: %w", err)
	}

	if lock.IsActive() {
		if lock.IsOurs() {
			err := m.maintenanceMode.Deactivate()
			if err != nil {
				return fmt.Errorf("failed to deactivate maintenance mode: %w", err)
			}
		}

		return fmt.Errorf("cannot deactivate maintenance mode as it was activated by another application")
	}

	// do nothing
	return nil
}
