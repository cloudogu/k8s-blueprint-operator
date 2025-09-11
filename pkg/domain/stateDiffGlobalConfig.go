package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"github.com/cloudogu/k8s-registry-lib/config"
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

func (diffs GlobalConfigDiffs) countByAction() map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, diff := range diffs {
		countByAction[diff.NeededAction]++
	}
	return countByAction
}

func newGlobalConfigEntryDiff(
	key common.GlobalConfigKey,
	actualValue *common.GlobalConfigValue,
	actualExists bool,
	expectedValue *common.GlobalConfigValue,
	expectedExists bool,
) GlobalConfigEntryDiff {
	actual := GlobalConfigValueState{
		Value:  (*string)(actualValue),
		Exists: actualExists,
	}
	expected := GlobalConfigValueState{
		Value:  (*string)(expectedValue),
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
	actualConfig config.GlobalConfig,
) GlobalConfigDiffs {
	var configDiffs []GlobalConfigEntryDiff

	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualConfig.Get(key)
		configDiffs = append(configDiffs, newGlobalConfigEntryDiff(key, &actualEntry, actualExists, &expectedValue, true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualConfig.Get(key)
		configDiffs = append(configDiffs, newGlobalConfigEntryDiff(key, &actualEntry, actualExists, nil, false))
	}
	return configDiffs
}
