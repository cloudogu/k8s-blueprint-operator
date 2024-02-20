package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type GlobalConfigDiffs []GlobalConfigEntryDiff

func (diffs GlobalConfigDiffs) countByAction() map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, diff := range diffs {
		countByAction[diff.Action]++
	}
	return countByAction
}

type GlobalConfigValueState ConfigValueState
type GlobalConfigEntryDiff struct {
	Key      common.GlobalConfigKey
	Actual   GlobalConfigValueState
	Expected GlobalConfigValueState
	Action   ConfigAction
}

func newGlobalConfigEntryDiff(
	key common.GlobalConfigKey,
	actualValue common.GlobalConfigValue,
	actualExists bool,
	expectedValue common.GlobalConfigValue,
	expectedExists bool,
) GlobalConfigEntryDiff {
	actual := GlobalConfigValueState{
		Value:  string(actualValue),
		Exists: actualExists,
	}
	expected := GlobalConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return GlobalConfigEntryDiff{
		Key:      key,
		Actual:   actual,
		Expected: expected,
		Action:   getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineGlobalConfigDiffs(
	config GlobalConfig,
	actualDoguConfig map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry,
) GlobalConfigDiffs {
	var configDiffs []GlobalConfigEntryDiff

	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualDoguConfig[key]
		configDiffs = append(configDiffs, newGlobalConfigEntryDiff(key, actualEntry.Value, actualExists, expectedValue, true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualDoguConfig[key]
		configDiffs = append(configDiffs, newGlobalConfigEntryDiff(key, actualEntry.Value, actualExists, "", false))
	}
	return configDiffs
}
