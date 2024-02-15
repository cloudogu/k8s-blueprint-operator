package domain

type DoguConfigDiff struct {
	NormalDoguConfigDiff    NormalDoguConfigDiff
	SensitiveDoguConfigDiff SensitiveDoguConfigDiff
}

type SensitiveDoguConfigDiff []ConfigKeyDiff

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
