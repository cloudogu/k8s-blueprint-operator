package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type GlobalConfigDiffs []GlobalConfigEntryDiff

func (diffs GlobalConfigDiffs) HasChanges() bool {
	for _, globalConfigDiff := range diffs {
		if globalConfigDiff.NeededAction != ConfigActionNone {
			return true
		}
	}
	return false
}

func (diffs GlobalConfigDiffs) GetGlobalConfigDiffsByAction() map[ConfigAction][]GlobalConfigEntryDiff {
	return util.GroupBy(diffs, func(diff GlobalConfigEntryDiff) ConfigAction {
		return diff.NeededAction
	})
}

type GlobalConfigValueState ConfigValueState
type GlobalConfigEntryDiff struct {
	Key          common.GlobalConfigKey
	Actual       GlobalConfigValueState
	Expected     GlobalConfigValueState
	NeededAction ConfigAction
}

func NewGlobalConfigValueStateFromEntry(actualEntry *ecosystem.GlobalConfigEntry) GlobalConfigValueState {
	actualValue := common.GlobalConfigValue("")
	if actualEntry != nil {
		actualValue = actualEntry.Value
	}
	return GlobalConfigValueState{
		Value:  string(actualValue),
		Exists: actualEntry != nil,
	}
}

func (diffs GlobalConfigDiffs) countByAction() map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, diff := range diffs {
		countByAction[diff.NeededAction]++
	}
	return countByAction
}

func newGlobalConfigEntryDiff(
	key common.GlobalConfigKey,
	actualEntry *ecosystem.GlobalConfigEntry,
	expectedValue common.GlobalConfigValue,
	expectedExists bool,
) GlobalConfigEntryDiff {
	actual := NewGlobalConfigValueStateFromEntry(actualEntry)
	expected := GlobalConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return GlobalConfigEntryDiff{
		Key:          key,
		Actual:       actual,
		Expected:     expected,
		NeededAction: getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineGlobalConfigDiffs(
	config GlobalConfig,
	actualDoguConfig map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry,
) GlobalConfigDiffs {
	var configDiffs []GlobalConfigEntryDiff

	// present entries
	for key, expectedValue := range config.Present {
		actualEntry := actualDoguConfig[key]
		configDiffs = append(configDiffs, newGlobalConfigEntryDiff(key, actualEntry, expectedValue, true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry := actualDoguConfig[key]
		configDiffs = append(configDiffs, newGlobalConfigEntryDiff(key, actualEntry, "", false))
	}
	return configDiffs
}
