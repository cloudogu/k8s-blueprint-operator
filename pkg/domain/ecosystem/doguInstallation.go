package ecosystem

type DoguInstallation struct {
	Name   string
	Health DoguHealth
}

type DoguHealth = string

const (
	Healhty   DoguHealth = "healthy"
	Unhealthy DoguHealth = "unhealthy"
)
