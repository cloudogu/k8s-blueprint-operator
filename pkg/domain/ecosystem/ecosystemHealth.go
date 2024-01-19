package ecosystem

import (
	"fmt"
)

// HealthResult is a snapshot of the health states of all relevant parts of the running ecosystem.
type HealthResult struct {
	DoguHealth DoguHealthResult
}

func (result HealthResult) String() string {
	return fmt.Sprintf("ecosystem is unhealthy: %s", result.DoguHealth)
}

func (res *HealthResult) AllHealthy() bool {
	return res.DoguHealth.AllHealthy()
}
