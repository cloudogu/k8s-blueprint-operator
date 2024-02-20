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
	NeededAction         ConfigAction
}

func NewDoguConfigValueStateFromSensitiveEntry(actualEntry *ecosystem.SensitiveDoguConfigEntry) DoguConfigValueState {
	actualValue := common.SensitiveDoguConfigValue("")
	if actualEntry != nil {
		actualValue = common.SensitiveDoguConfigValue(actualEntry.Value)
	}
	return DoguConfigValueState{
		Value:  string(actualValue),
		Exists: actualEntry != nil,
	}
}

func newSensitiveDoguConfigEntryDiff(
	key common.SensitiveDoguConfigKey,
	actualEntry *ecosystem.SensitiveDoguConfigEntry,
	expectedValue common.SensitiveDoguConfigValue,
	expectedExists bool,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigEntryDiff {
	actual := NewDoguConfigValueStateFromSensitiveEntry(actualEntry)
	expected := DoguConfigValueState{
		Value:  string(expectedValue),
		Exists: expectedExists,
	}
	return SensitiveDoguConfigEntryDiff{
		Key:                  key,
		Actual:               actual,
		Expected:             expected,
		doguAlreadyInstalled: doguAlreadyInstalled,
		NeededAction:         getNeededSensitiveConfigAction(ConfigValueState(expected), ConfigValueState(actual), doguAlreadyInstalled),
	}
}

func determineSensitiveDoguConfigDiffs(
	config SensitiveDoguConfig,
	actualDoguConfig map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry,
	doguAlreadyInstalled bool,
) SensitiveDoguConfigDiffs {
	var doguConfigDiff []SensitiveDoguConfigEntryDiff
	// present entries
	for key, expectedValue := range config.Present {
		actualEntry := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualEntry, expectedValue, true, doguAlreadyInstalled))
	}
	// absent entries
	for _, key := range config.Absent {
		actualEntry := actualDoguConfig[key]
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualEntry, "", false, doguAlreadyInstalled))
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
