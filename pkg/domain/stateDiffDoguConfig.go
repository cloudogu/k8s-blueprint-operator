package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type DoguConfigDiffs []DoguConfigEntryDiff

func (diffs DoguConfigDiffs) GetDoguConfigDiffByAction() map[ConfigAction][]DoguConfigEntryDiff {
	return util.GroupBy(diffs, func(diff DoguConfigEntryDiff) ConfigAction {
		return diff.NeededAction
	})
}

func (diffs DoguConfigDiffs) HasChangesForDogu(dogu common.SimpleDoguName) bool {
	for _, diff := range diffs {
		if diff.Key.DoguName == dogu && diff.NeededAction != ConfigActionNone {
			return true
		}
	}
	return false
}

type DoguConfigValueState ConfigValueState

type ConfigValueState struct {
	Value  string
	Exists bool
}
type DoguConfigEntryDiff struct {
	Key          common.DoguConfigKey
	Actual       DoguConfigValueState
	Expected     DoguConfigValueState
	NeededAction ConfigAction
}

func newDoguConfigEntryDiff(
	key common.DoguConfigKey,
	actualValue config.Value,
	actualExists bool,
	expectedValue common.DoguConfigValue,
	expectedExists bool,
) DoguConfigEntryDiff {
	actual := DoguConfigValueState{
		Value:  string(actualValue),
		Exists: actualExists,
	}
	expected := DoguConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return DoguConfigEntryDiff{
		Key:          key,
		Actual:       actual,
		Expected:     expected,
		NeededAction: getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineDoguConfigDiffs(
	wantedConfig DoguConfig,
	actualConfig map[common.SimpleDoguName]config.DoguConfig,
) DoguConfigDiffs {
	var doguConfigDiff []DoguConfigEntryDiff
	// present entries
	for key, expectedValue := range wantedConfig.Present {
		actualEntry, exists := actualConfig[key.DoguName].Get(key.Key)
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualEntry, exists, expectedValue, true))
	}
	// absent entries
	for _, key := range wantedConfig.Absent {
		actualEntry, exists := actualConfig[key.DoguName].Get(key.Key)
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualEntry, exists, "", false))
	}
	return doguConfigDiff
}
