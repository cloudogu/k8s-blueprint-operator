package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

func determineGlobalConfigDiffs(
	config GlobalConfig,
	actualDoguConfig map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry,
) GlobalConfigDiff {
	var configDiffs []GlobalConfigEntryDiff

	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualDoguConfig[key]
		configDiffs = append(configDiffs, determineGlobalConfigDiff(key, string(actualEntry.Value), actualExists, string(expectedValue), true))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualDoguConfig[key]
		configDiffs = append(configDiffs, determineGlobalConfigDiff(key, string(actualEntry.Value), actualExists, string(actualEntry.Value), false))
	}
	return configDiffs
}

func determineGlobalConfigDiff(key common.GlobalConfigKey, actualValue string, actualExists bool, expectedValue string, expectedExists bool) GlobalConfigEntryDiff {
	actual := GlobalConfigValueState{
		Value:  actualValue,
		Exists: actualExists,
	}
	expected := GlobalConfigValueState{
		Value:  expectedValue,
		Exists: expectedExists,
	}
	return GlobalConfigEntryDiff{
		Key:      key,
		Actual:   actual,
		Expected: expected,
		Action:   getNeededConfigAction(ConfigValueState(actual), ConfigValueState(expected)),
	}
}
