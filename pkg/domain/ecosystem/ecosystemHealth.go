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
	return fmt.Sprintf("ecosystem is unhealthy: %s", result.DoguHealth)
}

func (res *HealthResult) AllHealthy() bool {
	return res.DoguHealth.AllHealthy()
}
