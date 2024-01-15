package ecosystem

import (
	"fmt"
)

type HealthResult struct {
	DoguHealth DoguHealthResult
}

func (result HealthResult) String() string {
	return fmt.Sprintf("ecosystem is unhealthy: %s", result.DoguHealth)
}

func (res *HealthResult) AllHealthy() bool {
	return res.DoguHealth.AllHealthy()
}
