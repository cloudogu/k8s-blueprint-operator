package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

type SensitiveDoguConfigDiffs []SensitiveDoguConfigEntryDiff
type SensitiveDoguConfigEntryDiff struct {
	Key                  common.SensitiveDoguConfigKey
	Actual               DoguConfigValueState
	Expected             DoguConfigValueState
	DoguAlreadyInstalled bool
	NeededAction         ConfigAction
}

func NewDoguConfigValueStateFromSensitiveValue(actualValue common.SensitiveDoguConfigValue, actualExists bool) DoguConfigValueState {
	if !actualExists {
		actualValue = ""
	}
	return DoguConfigValueState{
		Value:  string(actualValue),
		Exists: actualExists,
	}
}

func newSensitiveDoguConfigEntryDiff(
	key common.SensitiveDoguConfigKey,
	actualValue common.SensitiveDoguConfigValue,
	actualExists bool,
	expectedValue common.SensitiveDoguConfigValue,
	expectedExists bool,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigEntryDiff {
	actual := NewDoguConfigValueStateFromSensitiveValue(actualValue, actualExists)
	expected := DoguConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return SensitiveDoguConfigEntryDiff{
		Key:                  key,
		Actual:               actual,
		Expected:             expected,
		DoguAlreadyInstalled: doguAlreadyInstalled,
		NeededAction:         getNeededSensitiveConfigAction(ConfigValueState(expected), ConfigValueState(actual), doguAlreadyInstalled),
	}
}

func determineSensitiveDoguConfigDiffs(
	config SensitiveDoguConfig,
	actualDoguConfig map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigDiffs {
	var doguConfigDiff []SensitiveDoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualValue, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualValue, actualExists, expectedValue, true, doguAlreadyInstalled))
	}
	// absent entries
	for _, key := range config.Absent {
		actualValue, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualValue, actualExists, "", false, doguAlreadyInstalled))
	}
	return doguConfigDiff
}

func getNeededSensitiveConfigAction(expected ConfigValueState, actual ConfigValueState, doguAlreadyInstalled bool) ConfigAction {
	action := getNeededConfigAction(expected, actual)
	if action == ConfigActionSet && !doguAlreadyInstalled {
		return ConfigActionSetToEncrypt
	}
	return action
}
