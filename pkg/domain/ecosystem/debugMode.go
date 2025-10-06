package ecosystem

const (
	DebugModeStatusComplete string = "Completed"
	DebugModeStatusFailed   string = "Failed"
)

// DebugMode represents an object holding information about the debug mode in the ecosystem.
type DebugMode struct {
	// Phase defines the current general state the resource is in.
	Phase string
}

func (d *DebugMode) IsActive() bool {
	return d.Phase != DebugModeStatusComplete && d.Phase != DebugModeStatusFailed
}
