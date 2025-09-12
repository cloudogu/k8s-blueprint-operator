package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type DoguConfigDiffs []DoguConfigEntryDiff
type SensitiveDoguConfigDiffs = DoguConfigDiffs

func (diffs DoguConfigDiffs) HasChanges() bool {
	for _, diff := range diffs {
		if diff.NeededAction != ConfigActionNone {
			return true
		}
	}
	return false
}

type DoguConfigValueState ConfigValueState

type ConfigValueState struct {
	Value  *string
	Exists bool
}

func (a ConfigValueState) Equal(b ConfigValueState) bool {
	if a.Exists != b.Exists {
		return false
	}
	if a.Value == b.Value { // covers both nil and same address
		return true
	}
	if a.Value == nil || b.Value == nil {
		return false
	}
	return *a.Value == *b.Value
}

type DoguConfigEntryDiff struct {
	Key          common.DoguConfigKey
	Actual       DoguConfigValueState
	Expected     DoguConfigValueState
	NeededAction ConfigAction
}
type SensitiveDoguConfigEntryDiff = DoguConfigEntryDiff

func newDoguConfigEntryDiff(
	key common.DoguConfigKey,
	actualValue *config.Value,
	actualExists bool,
	expectedValue *common.DoguConfigValue,
	expectedExists bool,
) DoguConfigEntryDiff {
	actual := DoguConfigValueState{
		Value:  (*string)(actualValue),
		Exists: actualExists,
	}
	expected := DoguConfigValueState{
		Value:  (*string)(expectedValue),
		Exists: expectedExists,
	}
	return DoguConfigEntryDiff{
		Key:          key,
		Actual:       actual,
		Expected:     expected,
		NeededAction: getNeededConfigAction(ConfigValueState(expected), ConfigValueState(actual)),
	}
}

func determineDoguConfigDiffs(
	wantedConfig DoguConfig,
	actualConfig map[cescommons.SimpleName]config.DoguConfig,
) DoguConfigDiffs {
	var doguConfigDiff []DoguConfigEntryDiff
	// present entries
	for key, expectedValue := range wantedConfig.Present {
		var actualValue *config.Value
		actualEntry, exists := actualConfig[key.DoguName].Get(key.Key)
		if exists {
			actualValue = &actualEntry
		}
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualValue, exists, &expectedValue, true))
	}
	// absent entries
	for _, key := range wantedConfig.Absent {
		var actualValue *config.Value
		actualEntry, exists := actualConfig[key.DoguName].Get(key.Key)
		if exists {
			actualValue = &actualEntry
		}
		doguConfigDiff = append(doguConfigDiff, newDoguConfigEntryDiff(key, actualValue, exists, nil, false))
	}
	return doguConfigDiff
}
