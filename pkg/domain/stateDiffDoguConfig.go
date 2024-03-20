package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type DoguConfigDiffs []DoguConfigEntryDiff

func (diffs DoguConfigDiffs) GetDoguConfigDiffByAction() map[ConfigAction][]DoguConfigEntryDiff {
	return util.GroupBy(diffs, func(diff DoguConfigEntryDiff) ConfigAction {
		return diff.NeededAction
	})
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

func NewDoguConfigValueStateFromEntry(actualEntry *ecosystem.DoguConfigEntry) DoguConfigValueState {
	actualValue := common.DoguConfigValue("")
	if actualEntry != nil {
		actualValue = actualEntry.Value
	}
	return DoguConfigValueState{
		Value:  string(actualValue),
		Exists: actualEntry != nil,
	}
}

func newDoguConfigEntryDiff(
	key common.DoguConfigKey,
	actualEntry *ecosystem.DoguConfigEntry,
	expectedValue common.DoguConfigValue,
	expectedExists bool,
) DoguConfigEntryDiff {
	actual := NewDoguConfigValueStateFromEntry(actualEntry)
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
	config DoguConfig,
	actualDoguConfig map[common.DoguConfigKey]*ecosystem.DoguConfigEntry,
) DoguConfigDiffs {
	var doguConfigDiff []DoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualEntry := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualEntry, expectedValue, true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualEntry, "", false))
	}
	return doguConfigDiff
}
