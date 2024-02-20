package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"slices"
)

type CombinedDoguConfigDiff struct {
	DoguConfigDiff          DoguConfigDiffs
	SensitiveDoguConfigDiff SensitiveDoguConfigDiffs
}

type ConfigAction string

const (
	ConfigActionNone         ConfigAction = "none"
	ConfigActionSet          ConfigAction = "set"
	ConfigActionSetToEncrypt ConfigAction = "setToEncrypt"
	ConfigActionRemove       ConfigAction = "remove"
)

func determineConfigDiffs(
	blueprintConfig Config,
	actualGlobalConfig map[common.GlobalConfigKey]ecosystem.GlobalConfigEntry,
	actualDoguConfig map[common.DoguConfigKey]ecosystem.DoguConfigEntry,
	actualSensitiveDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
	alreadyInstalledDogus []common.SimpleDoguName,
) (map[common.SimpleDoguName]CombinedDoguConfigDiff, GlobalConfigDiffs) {
	return determineDogusConfigDiffs(blueprintConfig.Dogus, actualDoguConfig, actualSensitiveDoguConfig, alreadyInstalledDogus),
		determineGlobalConfigDiffs(blueprintConfig.Global, actualGlobalConfig)
}

func determineDogusConfigDiffs(
	combinedDoguConfigs map[common.SimpleDoguName]CombinedDoguConfig,
	actualDoguConfig map[common.DoguConfigKey]ecosystem.DoguConfigEntry,
	actualSensitiveDoguConfig map[common.SensitiveDoguConfigKey]ecosystem.SensitiveDoguConfigEntry,
	installedDogus []common.SimpleDoguName,
) map[common.SimpleDoguName]CombinedDoguConfigDiff {
	diffsPerDogu := map[common.SimpleDoguName]CombinedDoguConfigDiff{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		diffsPerDogu[doguName] = CombinedDoguConfigDiff{
			DoguConfigDiff:          determineDoguConfigDiffs(combinedDoguConfig.Config, actualDoguConfig),
			SensitiveDoguConfigDiff: determineSensitiveDoguConfigDiffs(combinedDoguConfig.SensitiveConfig, actualSensitiveDoguConfig, slices.Contains(installedDogus, doguName)),
		}
	}
	return diffsPerDogu
}

func getNeededConfigAction(expected ConfigValueState, actual ConfigValueState) ConfigAction {
	if expected == actual {
		return ConfigActionNone
	}
	if !expected.Exists {
		if actual.Exists {
			return ConfigActionRemove
		} else {
			return ConfigActionNone
		}
	}
	return ConfigActionSet
}
