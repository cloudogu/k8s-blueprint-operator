package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type SensitiveDoguConfigDiffs []SensitiveDoguConfigEntryDiff

func (diffs SensitiveDoguConfigDiffs) GetSensitiveDoguConfigDiffByAction() map[ConfigAction][]SensitiveDoguConfigEntryDiff {
	return util.GroupBy(diffs, func(diff SensitiveDoguConfigEntryDiff) ConfigAction {
		return diff.NeededAction
	})
}

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
		NeededAction:         getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
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
