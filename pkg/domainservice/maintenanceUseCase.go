package domainservice

type MaintenanceUseCase struct {
	MaintenanceMode
}

func NewMaintenanceUseCase(maintenanceMode MaintenanceMode) *MaintenanceUseCase {
	return &MaintenanceUseCase{MaintenanceMode: maintenanceMode}
}
