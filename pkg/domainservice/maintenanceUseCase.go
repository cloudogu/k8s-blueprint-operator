package domainservice

// MaintenancePageModel contains data that gets displayed when the maintenance mode is active.
type MaintenancePageModel struct {
	Title string
	Text  string
}

type MaintenanceUseCase struct {
	MaintenanceMode
}

func NewMaintenanceUseCase(maintenanceMode MaintenanceMode) *MaintenanceUseCase {
	return &MaintenanceUseCase{MaintenanceMode: maintenanceMode}
}
