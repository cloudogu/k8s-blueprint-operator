package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"slices"
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
	// ConfigActionSetEncrypted means that the config key needs to be encrypted
	ConfigActionSetEncrypted ConfigAction = "setEncrypted"
	// ConfigActionSetToEncrypt means that the config key needs to be encrypted but another service needs to do this.
	// This can happen if a dogu is not yet installed and therefore no encryption key pair is available.
	ConfigActionSetToEncrypt ConfigAction = "setToEncrypt"
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

func countByAction(combinedDogusConfigDiffs map[common.SimpleDoguName]CombinedDoguConfigDiffs) map[ConfigAction]int {
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
	clusterState ecosystem.EcosystemState,
) (map[common.SimpleDoguName]CombinedDoguConfigDiffs, GlobalConfigDiffs) {
	return determineDogusConfigDiffs(blueprintConfig.Dogus, clusterState.DoguConfig, clusterState.DecryptedSensitiveDoguConfig, clusterState.GetInstalledDoguNames()),
		determineGlobalConfigDiffs(blueprintConfig.Global, clusterState.GlobalConfig)
}

func determineDogusConfigDiffs(
	combinedDoguConfigs map[common.SimpleDoguName]CombinedDoguConfig,
	actualDoguConfig map[common.DoguConfigKey]*ecosystem.DoguConfigEntry,
	actualSensitiveDoguConfig map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
	installedDogus []common.SimpleDoguName,
) map[common.SimpleDoguName]CombinedDoguConfigDiffs {
	diffsPerDogu := map[common.SimpleDoguName]CombinedDoguConfigDiffs{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		diffsPerDogu[doguName] = CombinedDoguConfigDiffs{
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
