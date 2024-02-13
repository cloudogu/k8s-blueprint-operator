package domain

// StateDiff represents the diff between the defined state in the effective blueprint and the actual state in the ecosystem.
// If there is a state in the ecosystem, which is not represented in the effective blueprint, then the expected state is the actual state.
type StateDiff struct {
	DoguDiffs        DoguDiffs
	ComponentDiffs   ComponentDiffs
	DoguConfigDiff   map[string]DoguConfigDiff
	GlobalConfigDiff GlobalConfigDiff
}

// Action represents a needed Action for a dogu to reach the expected state.
type Action string

const (
	ActionNone                                 = "none"
	ActionInstall                              = "install"
	ActionUninstall                            = "uninstall"
	ActionUpgrade                              = "upgrade"
	ActionDowngrade                            = "downgrade"
	ActionSwitchDoguNamespace                  = "dogu namespace switch"
	ActionSwitchComponentDistributionNamespace = "component distribution namespace switch"
)
