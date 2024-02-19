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
			SensitiveDoguConfigDiff: determineSensitiveDoguConfigDiff(blueprintConfig.Dogus, actualSensitiveDoguConfig),
		},
		determineGlobalConfigDiff(blueprintConfig.Global, actualGlobalConfig),
		nil
}

func determineDoguConfigDiffs(
	config map[common.SimpleDoguName]CombinedDoguConfig,
	actualDoguConfig map[common.DoguConfigKey]ecosystem.DoguConfigEntry,
) DoguConfigDiff {
	var doguConfigDiff []DoguConfigEntryDiff

	for _, doguConfig := range config {
		// present entries
		for key, expectedValue := range doguConfig.Config.Present {
			actualEntry, actualExists := actualDoguConfig[key]
			doguConfigDiff = append(doguConfigDiff, determineDoguConfigDiff(key, string(actualEntry.Value), actualExists, string(expectedValue), true))
		}
		// absent entries
		for _, key := range doguConfig.Config.Absent {
			actualEntry, actualExists := actualDoguConfig[key]
			doguConfigDiff = append(doguConfigDiff, determineDoguConfigDiff(key, string(actualEntry.Value), actualExists, string(actualEntry.Value), false))
		}
	}
	return doguConfigDiff
}

func determineDoguConfigDiff(key common.DoguConfigKey, actualValue string, actualExists bool, expectedValue string, expectedExists bool) DoguConfigEntryDiff {
	actual := DoguConfigValueState{
		Value:  actualValue,
		Exists: actualExists,
	}
	expected := DoguConfigValueState{
		Value:  expectedValue,
		Exists: expectedExists,
	}
	return DoguConfigEntryDiff{
		Key:      key,
		Actual:   actual,
		Expected: expected,
		Action:   getNeededConfigAction(ConfigValueState(actual), ConfigValueState(expected)),
	}
}

func determineSensitiveDoguConfigDiff(
	config map[common.SimpleDoguName]CombinedDoguConfig,
	actualDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
) SensitiveDoguConfigDiff {
	return nil
}

func determineGlobalConfigDiff(
	config GlobalConfig,
	actualDoguConfig map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry,
) GlobalConfigDiff {
	return nil
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
