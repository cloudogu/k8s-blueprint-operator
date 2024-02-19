package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

func determineSensitiveDoguConfigDiffs(
	config map[common.SimpleDoguName]CombinedDoguConfig,
	actualDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
) SensitiveDoguConfigDiff {
	var doguConfigDiff []SensitiveDoguConfigEntryDiff

	for _, doguConfig := range config {
		// present entries
		for key, expectedValue := range doguConfig.SensitiveConfig.Present {
			actualEntry, actualExists := actualDoguConfig[key]
			doguConfigDiff = append(doguConfigDiff, determineSensitiveDoguConfigDiff(key, string(actualEntry.Value), actualExists, string(expectedValue), true))
		}
		// absent entries
		for _, key := range doguConfig.SensitiveConfig.Absent {
			actualEntry, actualExists := actualDoguConfig[key]
			doguConfigDiff = append(doguConfigDiff, determineSensitiveDoguConfigDiff(key, string(actualEntry.Value), actualExists, string(actualEntry.Value), false))
		}
	}
	return doguConfigDiff
}

func determineSensitiveDoguConfigDiff(key common.SensitiveDoguConfigKey, actualValue string, actualExists bool, expectedValue string, expectedExists bool) SensitiveDoguConfigEntryDiff {
	actual := EncryptedDoguConfigValueState{
		Value:  actualValue,
		Exists: actualExists,
	}
	expected := EncryptedDoguConfigValueState{
		Value:  expectedValue,
		Exists: expectedExists,
	}
	return SensitiveDoguConfigEntryDiff{
		Key:      key,
		Actual:   actual,
		Expected: expected,
		Action:   getNeededConfigAction(ConfigValueState(actual), ConfigValueState(expected)),
	}
}
