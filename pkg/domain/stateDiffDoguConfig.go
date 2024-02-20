package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type DoguConfigDiffs []DoguConfigEntryDiff
type DoguConfigValueState ConfigValueState

type ConfigValueState struct {
	Value  string
	Exists bool
}
type DoguConfigEntryDiff struct {
	Key      common.DoguConfigKey
	Actual   DoguConfigValueState
	Expected DoguConfigValueState
	Action   ConfigAction
}

func newDoguConfigEntryDiff(key common.DoguConfigKey, actualValue common.DoguConfigValue, actualExists bool, expectedValue common.DoguConfigValue, expectedExists bool) DoguConfigEntryDiff {
	actual := DoguConfigValueState{
		Value:  string(actualValue),
		Exists: actualExists,
	}
	expected := DoguConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return DoguConfigEntryDiff{
		Key:      key,
		Actual:   actual,
		Expected: expected,
		Action:   getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineDoguConfigDiffs(
	config DoguConfig,
	actualDoguConfig map[common.DoguConfigKey]ecosystem.DoguConfigEntry,
) DoguConfigDiffs {
	var doguConfigDiff []DoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualEntry.Value, actualExists, expectedValue, true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualEntry.Value, actualExists, "", false))
	}
	return doguConfigDiff
}
