package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type GlobalConfigDiffs []GlobalConfigEntryDiff

func (diffs GlobalConfigDiffs) HasChanges() bool {
	for _, globalConfigDiff := range diffs {
		if globalConfigDiff.NeededAction != ConfigActionNone {
			return true
		}
	}
	return false
}

func (diffs GlobalConfigDiffs) GetGlobalConfigDiffsByAction() map[ConfigAction][]GlobalConfigEntryDiff {
	return util.GroupBy(diffs, func(diff GlobalConfigEntryDiff) ConfigAction {
		return diff.NeededAction
	})
}

type GlobalConfigValueState ConfigValueState
type GlobalConfigEntryDiff struct {
	Key          common.GlobalConfigKey
	Actual       GlobalConfigValueState
	Expected     GlobalConfigValueState
	NeededAction ConfigAction
}

func (diffs GlobalConfigDiffs) countByAction() map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, diff := range diffs {
		countByAction[diff.NeededAction]++
	}
	return countByAction
}

func newGlobalConfigEntryDiff(
	key common.GlobalConfigKey,
	actualValue *common.GlobalConfigValue,
	actualExists bool,
	expectedValue *common.GlobalConfigValue,
	expectedExists bool,
) GlobalConfigEntryDiff {
	actual := GlobalConfigValueState{
		Value:  (*string)(actualValue),
		Exists: actualExists,
	}
	expected := GlobalConfigValueState{
		Value:  (*string)(expectedValue),
		Exists: expectedExists,
	}
	return GlobalConfigEntryDiff{
		Key:          key,
		Actual:       actual,
		Expected:     expected,
		NeededAction: getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineGlobalConfigDiffs(
	config GlobalConfigEntries,
	actualConfig config.GlobalConfig,
	referencedSensitiveGlobalConfig map[common.GlobalConfigKey]common.GlobalConfigValue,
	referencedGlobalConfig map[common.GlobalConfigKey]common.GlobalConfigValue,
) GlobalConfigDiffs {
	var configDiffs []GlobalConfigEntryDiff

	for _, expectedConfig := range config {
		println(fmt.Sprintf("Key: %s", expectedConfig.Key))
		var actualValue *common.GlobalConfigValue
		actualEntry, actualExists := actualConfig.Get(expectedConfig.Key)
		if actualExists {
			actualValue = &actualEntry
		}
		referencedValue := getReferencedGlobalConfigValue(expectedConfig.Key, referencedSensitiveGlobalConfig, referencedGlobalConfig)
		if referencedValue != nil {
			expectedConfig.Value = referencedValue
		}
		diff := newGlobalConfigEntryDiff(expectedConfig.Key, actualValue, actualExists, expectedConfig.Value, !expectedConfig.Absent)
		// only add diff if there are changes
		if diff.NeededAction != ConfigActionNone {
			configDiffs = append(configDiffs, diff)
		}
	}
	return configDiffs
}

func getReferencedGlobalConfigValue(
	key common.GlobalConfigKey,
	referencedSensitiveGlobalConfig map[common.GlobalConfigKey]common.GlobalConfigValue,
	referencedGlobalConfig map[common.GlobalConfigKey]common.GlobalConfigValue,
) *common.GlobalConfigValue {
	value, exists := referencedSensitiveGlobalConfig[key]
	if exists {
		value = common.GlobalConfigValue(base64.StdEncoding.EncodeToString([]byte(value)))
		return &value
	}
	value, exists = referencedGlobalConfig[key]
	if exists {
		return &value
	}
	return nil
}
