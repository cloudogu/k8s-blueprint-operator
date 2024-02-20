package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type SensitiveDoguConfigDiffs []SensitiveDoguConfigEntryDiff
type SensitiveDoguConfigEntryDiff struct {
	Key                  common.SensitiveDoguConfigKey
	Actual               DoguConfigValueState
	Expected             DoguConfigValueState
	doguAlreadyInstalled bool
	Action               ConfigAction
}

func newSensitiveDoguConfigEntryDiff(
	key common.SensitiveDoguConfigKey,
	actualValue common.EncryptedDoguConfigValue,
	actualExists bool,
	expectedValue common.SensitiveDoguConfigValue,
	expectedExists bool,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigEntryDiff {
	actual := DoguConfigValueState{
		Value:  string(actualValue),
		Exists: actualExists,
	}
	expected := DoguConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return SensitiveDoguConfigEntryDiff{
		Key:                  key,
		Actual:               actual,
		Expected:             expected,
		doguAlreadyInstalled: doguAlreadyInstalled,
		Action:               getNeededSensitiveConfigAction(ConfigValueState(expected), ConfigValueState(actual), doguAlreadyInstalled),
	}
}

func determineSensitiveDoguConfigDiffs(
	config SensitiveDoguConfig,
	actualDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigDiffs {
	var doguConfigDiff []SensitiveDoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualEntry.Value, actualExists, expectedValue, true, doguAlreadyInstalled))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry, actualExists := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualEntry.Value, actualExists, "", false, doguAlreadyInstalled))
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
