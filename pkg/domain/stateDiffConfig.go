package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
	"maps"
)

type ConfigAction string

const (
	// ConfigActionNone means that nothing is to do for this config key
	ConfigActionNone ConfigAction = "none"
	// ConfigActionSet means that the config key needs to be set as given
	ConfigActionSet ConfigAction = "set"
	// ConfigActionRemove means that the config key needs to be deleted
	ConfigActionRemove ConfigAction = "remove"
)

// censorValues censors all sensitive configuration data to make them unrecognisable.
func censorValues(sensitiveConfigByDogu map[common.SimpleDoguName]SensitiveDoguConfigDiffs) map[common.SimpleDoguName]SensitiveDoguConfigDiffs {
	//TODO: add/fix tests for this
	censoredByDogu := maps.Clone(sensitiveConfigByDogu)
	for dogu, entryDiffs := range sensitiveConfigByDogu {
		censoredByDogu[dogu] = entryDiffs.CensorValues()
	}
	return censoredByDogu
}

func countByAction(diffsByDogu map[common.SimpleDoguName]DoguConfigDiffs) map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, doguDiffs := range diffsByDogu {
		for _, diff := range doguDiffs {
			countByAction[diff.NeededAction]++
		}
	}
	return countByAction
}

func determineConfigDiffs(
	blueprintConfig Config,
	globalConfig config.GlobalConfig,
	configByDogu map[common.SimpleDoguName]config.DoguConfig,
	SensitiveConfigByDogu map[common.SimpleDoguName]config.DoguConfig,
) (
	map[common.SimpleDoguName]DoguConfigDiffs,
	map[common.SimpleDoguName]SensitiveDoguConfigDiffs,
	GlobalConfigDiffs,
) {
	return determineDogusConfigDiffs(blueprintConfig.Dogus, configByDogu),
		determineSensitiveDogusConfigDiffs(blueprintConfig.Dogus, SensitiveConfigByDogu),
		determineGlobalConfigDiffs(blueprintConfig.Global, globalConfig)
}

func determineDogusConfigDiffs(
	combinedDoguConfigs map[common.SimpleDoguName]CombinedDoguConfig,
	configByDogu map[common.SimpleDoguName]config.DoguConfig,
) map[common.SimpleDoguName]DoguConfigDiffs {
	diffsPerDogu := map[common.SimpleDoguName]DoguConfigDiffs{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		diffsPerDogu[doguName] = determineDoguConfigDiffs(combinedDoguConfig.Config, configByDogu)
	}
	return diffsPerDogu
}

func determineSensitiveDogusConfigDiffs(
	combinedDoguConfigs map[common.SimpleDoguName]CombinedDoguConfig,
	configByDogu map[common.SimpleDoguName]config.DoguConfig,
) map[common.SimpleDoguName]DoguConfigDiffs {
	diffsPerDogu := map[common.SimpleDoguName]DoguConfigDiffs{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		diffsPerDogu[doguName] = determineDoguConfigDiffs(combinedDoguConfig.SensitiveConfig, configByDogu)
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
