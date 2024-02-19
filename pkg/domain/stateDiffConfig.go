package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type CombinedDoguConfigDiff struct {
	DoguConfigDiff          DoguConfigDiff
	SensitiveDoguConfigDiff SensitiveDoguConfigDiff
}

type DoguConfigDiff []DoguConfigEntryDiff
type SensitiveDoguConfigDiff []SensitiveDoguConfigEntryDiff
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

func determineConfigDiff(
	blueprintConfig Config,
	actualGlobalConfig map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry,
	actualDoguConfig map[common.DoguConfigKey]ecosystem.DoguConfigEntry,
	actualSensitiveDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
) (CombinedDoguConfigDiff, GlobalConfigDiff, error) {
	return CombinedDoguConfigDiff{
			DoguConfigDiff:          determineDoguConfigDiffs(blueprintConfig.Dogus, actualDoguConfig),
			SensitiveDoguConfigDiff: determineSensitiveDoguConfigDiffs(blueprintConfig.Dogus, actualSensitiveDoguConfig),
		},
		determineGlobalConfigDiffs(blueprintConfig.Global, actualGlobalConfig),
		nil
}

func getNeededConfigAction(expected ConfigValueState, actual ConfigValueState) ConfigAction {
	if expected == actual {
		return ConfigActionNone
	}
	if !expected.Exists {
		return ConfigActionRemove
	}

	return ConfigActionSet
}
