package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

func determineSensitiveDoguConfigDiffs(
	config SensitiveDoguConfig,
	actualDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigDiffs {
	var doguConfigDiff []SensitiveDoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, determineSensitiveDoguConfigDiff(key, string(actualEntry.Value), actualExists, string(expectedValue), true, doguAlreadyInstalled))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, determineSensitiveDoguConfigDiff(key, string(actualEntry.Value), actualExists, "", false, doguAlreadyInstalled))
	}
	return doguConfigDiff
}

func determineSensitiveDoguConfigDiff(key common.SensitiveDoguConfigKey, actualValue string, actualExists bool, expectedValue string, expectedExists bool, doguAlreadyInstalled bool) SensitiveDoguConfigEntryDiff {
	actual := DoguConfigValueState{
		Value:  actualValue,
		Exists: actualExists,
	}
	expected := DoguConfigValueState{
		Value:  expectedValue,
		Exists: expectedExists,
	}
	return SensitiveDoguConfigEntryDiff{
		Key:        key,
		Actual:     actual,
		Expected:   expected,
		doguExists: doguAlreadyInstalled,
		Action:     getNeededSensitiveConfigAction(ConfigValueState(expected), ConfigValueState(actual), doguAlreadyInstalled),
	}
}

func getNeededSensitiveConfigAction(expected ConfigValueState, actual ConfigValueState, doguAlreadyInstalled bool) ConfigAction {
	action := getNeededConfigAction(expected, actual)
	if action == ConfigActionSet && !doguAlreadyInstalled {
		return ConfigActionSetToEncrypt
	}
	return action
}
