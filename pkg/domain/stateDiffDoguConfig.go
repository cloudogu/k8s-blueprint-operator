package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type DoguConfigDiffs []DoguConfigEntryDiff
type SensitiveDoguConfigDiffs = DoguConfigDiffs

func (diffs DoguConfigDiffs) HasChanges() bool {
	for _, diff := range diffs {
		if diff.NeededAction != ConfigActionNone {
			return true
		}
	}
	return false
}

func (diffs SensitiveDoguConfigDiffs) CensorValues() SensitiveDoguConfigDiffs {
	var censoredEntries []DoguConfigEntryDiff
	for _, entry := range diffs {
		actual := entry.Actual
		expected := entry.Expected
		if len(entry.Actual.Value) > 0 {
			actual.Value = censorValue
		}
		if len(entry.Expected.Value) > 0 {
			expected.Value = censorValue
		}
		censoredEntries = append(censoredEntries, DoguConfigEntryDiff{
			Key:          entry.Key,
			Actual:       actual,
			Expected:     expected,
			NeededAction: entry.NeededAction,
		})
	}
	return censoredEntries
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
type SensitiveDoguConfigEntryDiff = DoguConfigEntryDiff

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
	actualConfig map[cescommons.SimpleName]config.DoguConfig,
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
