package ecosystem

import (
	"fmt"
)

type HealthStatus = string

const (
	PendingHealthStatus      HealthStatus = ""
	AvailableHealthStatus    HealthStatus = "available"
	UnavailableHealthStatus  HealthStatus = "unavailable"
	NotInstalledHealthStatus HealthStatus = "not installed"
)

// HealthResult is a snapshot of the health states of all relevant parts of the running ecosystem.
type HealthResult struct {
	DoguHealth      DoguHealthResult
	ComponentHealth ComponentHealthResult
}

func (result HealthResult) String() string {
	return fmt.Sprintf("ecosystem health:\n  %s\n  %s", result.DoguHealth, result.ComponentHealth)
}

func (result HealthResult) AllHealthy() bool {
	return result.DoguHealth.AllHealthy() &&
		result.ComponentHealth.AllHealthy()
}
