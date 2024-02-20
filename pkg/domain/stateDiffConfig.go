package domain

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

type CombinedDoguConfigDiff struct {
	DoguConfigDiff          DoguConfigDiff
	SensitiveDoguConfigDiff SensitiveDoguConfigDiff
}

type DoguConfigDiff []DoguConfigEntryDiff
type SensitiveDoguConfigDiff []SensitiveDoguConfigEntryDiff

func (combinedDiff *CombinedDoguConfigDiff) censorValues() {
	for i, entry := range combinedDiff.SensitiveDoguConfigDiff {
		if len(entry.Actual.Value) > 0 {
			combinedDiff.SensitiveDoguConfigDiff[i].Actual.Value = censorValue
		}
		if len(entry.Expected.Value) > 0 {
			combinedDiff.SensitiveDoguConfigDiff[i].Expected.Value = censorValue
		}
	}
}

type GlobalConfigDiff []GlobalConfigEntryDiff

type DoguConfigValueState ConfigValueState
type EncryptedDoguConfigValueState ConfigValueState
type GlobalConfigValueState ConfigValueState

type DoguConfigEntryDiff struct {
	Key      common.DoguConfigKey
	Actual   DoguConfigValueState
	Expected DoguConfigValueState
	Action   ConfigAction
}

type SensitiveDoguConfigEntryDiff struct {
	Key      common.SensitiveDoguConfigKey
	Actual   EncryptedDoguConfigValueState
	Expected EncryptedDoguConfigValueState
	Action   ConfigAction
}

type GlobalConfigEntryDiff struct {
	Key      common.GlobalConfigKey
	Actual   GlobalConfigValueState
	Expected GlobalConfigValueState
	Action   ConfigAction
}

type ConfigValueState struct {
	Value  string
	Exists bool
}

type ConfigAction string

const (
	ConfigActionNone   ConfigAction = "none"
	ConfigActionSet    ConfigAction = "set"
	ConfigActionRemove ConfigAction = "remove"
)

func getNeededConfigAction(expected ConfigValueState, actual ConfigValueState) ConfigAction {
	if expected == actual {
		return ConfigActionNone
	}
	if !expected.Exists {
		return ConfigActionRemove
	}

	return ConfigActionSet
}
