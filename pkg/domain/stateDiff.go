package domain

// StateDiff represents the diff between the defined state in the effective blueprint and the actual state in the ecosystem.
// If there is a state in the ecosystem, which is not represented in the effective blueprint, then the expected state is the actual state.
type StateDiff struct {
	DoguDiffs        DoguDiffs
	ComponentDiffs   ComponentDiffs
	DoguConfigDiff   map[string]DoguConfigDiff
	GlobalConfigDiff GlobalConfigDiff
}

type DoguConfigDiff struct {
	SensibleDoguConfigDiff SensibleDoguConfigDiff
	NormalDoguConfigDiff   NormalDoguConfigDiff
}

type SensibleDoguConfigDiff []ConfigKeyDiff

type NormalDoguConfigDiff []ConfigKeyDiff

type ConfigKeyDiff struct {
	Key      string
	Actual   ConfigValue
	Expected ConfigValue
	Action   ConfigAction
}

type ConfigValue struct {
	Value  string
	Exists bool
}

type GlobalConfigDiff []ConfigKeyDiff

type ConfigAction string

const (
	ConfigActionNone   ConfigAction = "none"
	ConfigActionSet    ConfigAction = "set"
	ConfigActionRemove ConfigAction = "remove"
)

func getNeededConfigAction(expected ConfigValue, actual ConfigValue) ConfigAction {
	if expected == actual {
		return ConfigActionNone
	}
	if !expected.Exists {
		return ConfigActionRemove
	}

	return ConfigActionSet
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
