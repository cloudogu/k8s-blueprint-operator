package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

func determineDoguConfigDiffs(
	config DoguConfig,
	actualDoguConfig map[common.DoguConfigKey]ecosystem.DoguConfigEntry,
) DoguConfigDiffs {
	var doguConfigDiff []DoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, determineDoguConfigDiff(key, string(actualEntry.Value), actualExists, string(expectedValue), true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, determineDoguConfigDiff(key, string(actualEntry.Value), actualExists, "", false))
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
		Action:   getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}
