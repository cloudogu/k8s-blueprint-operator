package ecosystem

type DoguInstallation struct {
	Name    string
	Version string
	Health  DoguHealth
}

type DoguHealth = string

const (
	Healhty   DoguHealth = "healthy"
	Unhealthy DoguHealth = "unhealthy"
)
