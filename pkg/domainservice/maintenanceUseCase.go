package domainservice

type MaintenanceUseCase struct {
	maintenanceAdapter MaintenanceMode
}

func NewMaintenanceUseCase(maintenanceMode MaintenanceMode) *MaintenanceUseCase {
	return &MaintenanceUseCase{maintenanceAdapter: maintenanceMode}
}
