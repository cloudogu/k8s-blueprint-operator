package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type SensitiveDoguConfigDiffs []SensitiveDoguConfigEntryDiff

func (diffs SensitiveDoguConfigDiffs) GetSensitiveDoguConfigDiffByAction() map[ConfigAction][]SensitiveDoguConfigEntryDiff {
	return util.GroupBy(diffs, func(diff SensitiveDoguConfigEntryDiff) ConfigAction {
		return diff.NeededAction
	})
}

type SensitiveDoguConfigEntryDiff struct {
	Key          common.SensitiveDoguConfigKey
	Actual       DoguConfigValueState
	Expected     DoguConfigValueState
	NeededAction ConfigAction
}

func newSensitiveDoguConfigEntryDiff(
	key common.SensitiveDoguConfigKey,
	actualValue common.SensitiveDoguConfigValue,
	actualExists bool,
	expectedValue common.SensitiveDoguConfigValue,
	expectedExists bool,
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
		Key:          key,
		Actual:       actual,
		Expected:     expected,
		NeededAction: getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineSensitiveDoguConfigDiffs(
	wantedConfig SensitiveDoguConfig,
	actualConfig map[common.SimpleDoguName]config.DoguConfig,
) SensitiveDoguConfigDiffs {
	var doguConfigDiff []SensitiveDoguConfigEntryDiff
	// present entries
	for key, expectedValue := range wantedConfig.Present {
		actualEntry, exists := actualConfig[key.DoguName].Get(key.Key)
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualEntry, exists, expectedValue, true))
	}
	// absent entries
	for _, key := range wantedConfig.Absent {
		actualEntry, exists := actualConfig[key.DoguName].Get(key.Key)
		doguConfigDiff = append(doguConfigDiff, newSensitiveDoguConfigEntryDiff(key, actualEntry, exists, "", false))
	}
	return doguConfigDiff
}
