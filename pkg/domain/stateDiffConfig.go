package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type CombinedDoguConfigDiffs struct {
	DoguConfigDiff          DoguConfigDiffs
	SensitiveDoguConfigDiff SensitiveDoguConfigDiffs
}

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
func (combinedDiff CombinedDoguConfigDiffs) censorValues() CombinedDoguConfigDiffs {
	for i, entry := range combinedDiff.SensitiveDoguConfigDiff {
		if len(entry.Actual.Value) > 0 {
			combinedDiff.SensitiveDoguConfigDiff[i].Actual.Value = censorValue
		}
		if len(entry.Expected.Value) > 0 {
			combinedDiff.SensitiveDoguConfigDiff[i].Expected.Value = censorValue
		}
	}
	return combinedDiff
}

func countByAction(combinedDogusConfigDiffs map[cescommons.SimpleDoguName]CombinedDoguConfigDiffs) map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, doguDiffs := range combinedDogusConfigDiffs {
		for _, diff := range doguDiffs.DoguConfigDiff {
			countByAction[diff.NeededAction]++
		}
		for _, diff := range doguDiffs.SensitiveDoguConfigDiff {
			countByAction[diff.NeededAction]++
		}
	}
	return countByAction
}

func determineConfigDiffs(
	blueprintConfig Config,
	globalConfig config.GlobalConfig,
	configByDogu map[cescommons.SimpleDoguName]config.DoguConfig,
	SensitiveConfigByDogu map[cescommons.SimpleDoguName]config.DoguConfig,
) (map[cescommons.SimpleDoguName]CombinedDoguConfigDiffs, GlobalConfigDiffs) {
	return determineDogusConfigDiffs(
			blueprintConfig.Dogus,
			configByDogu,
			SensitiveConfigByDogu,
		),
		determineGlobalConfigDiffs(blueprintConfig.Global, globalConfig)
}

func determineDogusConfigDiffs(
	combinedDoguConfigs map[cescommons.SimpleDoguName]CombinedDoguConfig,
	configByDogu map[cescommons.SimpleDoguName]config.DoguConfig,
	sensitiveConfigByDogu map[cescommons.SimpleDoguName]config.DoguConfig,
) map[cescommons.SimpleDoguName]CombinedDoguConfigDiffs {
	diffsPerDogu := map[cescommons.SimpleDoguName]CombinedDoguConfigDiffs{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		diffsPerDogu[doguName] = CombinedDoguConfigDiffs{
			DoguConfigDiff:          determineDoguConfigDiffs(combinedDoguConfig.Config, configByDogu),
			SensitiveDoguConfigDiff: determineSensitiveDoguConfigDiffs(combinedDoguConfig.SensitiveConfig, sensitiveConfigByDogu),
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

func (combinedDiff CombinedDoguConfigDiffs) HasChanges() bool {
	for _, diff := range combinedDiff.DoguConfigDiff {
		if diff.NeededAction != ConfigActionNone {
			return true
		}
	}
	for _, diff := range combinedDiff.SensitiveDoguConfigDiff {
		if diff.NeededAction != ConfigActionNone {
			return true
		}
	}
	return false
}
