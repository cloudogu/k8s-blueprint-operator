package ecosystem

type HealthStatus = string

const (
	PendingHealthStatus     HealthStatus = ""
	AvailableHealthStatus   HealthStatus = "available"
	UnavailableHealthStatus HealthStatus = "unavailable"
)

// HealthResult is a snapshot of the health states of all relevant parts of the running ecosystem.
type HealthResult struct {
	DoguHealth DoguHealthResult
}

func (result HealthResult) String() string {
	return result.DoguHealth.String()
}

func (result HealthResult) AllHealthy() bool {
	return result.DoguHealth.AllHealthy()
}
