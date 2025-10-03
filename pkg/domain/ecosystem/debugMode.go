package ecosystem

const (
	DebugModeStatusSet             string = "SetDebugMode"
	DebugModeStatusWaitForRollback string = "WaitForRollback"
	DebugModeStatusRollback        string = "Rollback"
)

// DebugMode represents an object holding information about the debug mode in the ecosystem.
type DebugMode struct {
	// Phase defines the current general state the resource is in.
	Phase string
}

func (d *DebugMode) IsActive() bool {
	return d.Phase == DebugModeStatusSet || d.Phase == DebugModeStatusWaitForRollback || d.Phase == DebugModeStatusRollback
}
