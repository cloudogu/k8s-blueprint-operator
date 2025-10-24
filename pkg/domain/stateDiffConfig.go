package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type ConfigAction = bpv3.ConfigAction

const (
	// ConfigActionNone means that nothing is to do for this config key
	ConfigActionNone = bpv3.ConfigActionNone
	// ConfigActionSet means that the config key needs to be set as given
	ConfigActionSet = bpv3.ConfigActionSet
	// ConfigActionRemove means that the config key needs to be deleted
	ConfigActionRemove = bpv3.ConfigActionRemove
)

func countByAction(configActions []ConfigAction) map[ConfigAction]int {
	result := map[ConfigAction]int{}
	for _, action := range configActions {
		result[action]++
	}
	return result
}

func determineConfigDiffs(
	blueprintConfig Config,
	globalConfig config.GlobalConfig,
	configByDogu map[cescommons.SimpleName]config.DoguConfig,
	SensitiveConfigByDogu map[cescommons.SimpleName]config.DoguConfig,
	referencedSensitiveConfig map[common.DoguConfigKey]common.SensitiveDoguConfigValue,
	referencedConfig map[common.DoguConfigKey]common.DoguConfigValue,
) (
	map[cescommons.SimpleName]DoguConfigDiffs,
	map[cescommons.SimpleName]SensitiveDoguConfigDiffs,
	GlobalConfigDiffs,
) {
	return determineDogusConfigDiffs(blueprintConfig.Dogus, configByDogu, referencedConfig),
		determineSensitiveDogusConfigDiffs(blueprintConfig.Dogus, SensitiveConfigByDogu, referencedSensitiveConfig),
		determineGlobalConfigDiffs(blueprintConfig.Global, globalConfig)
}

func determineDogusConfigDiffs(
	blueprintDoguConfigs map[cescommons.SimpleName]DoguConfigEntries,
	configByDogu map[cescommons.SimpleName]config.DoguConfig,
	referencedValues map[common.DoguConfigKey]common.DoguConfigValue,
) map[cescommons.SimpleName]DoguConfigDiffs {
	var diffsPerDogu map[cescommons.SimpleName]DoguConfigDiffs
	for doguName, blueprintDoguConfig := range blueprintDoguConfigs {
		setConfigValues(doguName, blueprintDoguConfig, referencedValues)
		configDiffs := determineDoguConfigDiffs(doguName, blueprintDoguConfig, configByDogu, false)
		if len(configDiffs) > 0 {
			if diffsPerDogu == nil {
				diffsPerDogu = make(map[cescommons.SimpleName]DoguConfigDiffs)
			}
			diffsPerDogu[doguName] = configDiffs
		}
	}
	return diffsPerDogu
}

func determineSensitiveDogusConfigDiffs(
	blueprintDoguConfigs map[cescommons.SimpleName]DoguConfigEntries,
	configByDogu map[cescommons.SimpleName]config.DoguConfig,
	referencedValues map[common.DoguConfigKey]common.SensitiveDoguConfigValue,
) map[cescommons.SimpleName]DoguConfigDiffs {
	var diffsPerDogu map[cescommons.SimpleName]DoguConfigDiffs
	for doguName, blueprintDoguConfig := range blueprintDoguConfigs {
		setSensitiveConfigValues(doguName, blueprintDoguConfig, referencedValues)
		configDiffs := determineDoguConfigDiffs(doguName, blueprintDoguConfig, configByDogu, true)
		if len(configDiffs) > 0 {
			if diffsPerDogu == nil {
				diffsPerDogu = make(map[cescommons.SimpleName]DoguConfigDiffs)
			}
			diffsPerDogu[doguName] = configDiffs
		}
	}
	return diffsPerDogu
}

// setSensitiveConfigValues maps the referenced values to normal dogu config,
// so that the stateDiff can handle sensitive dogu config the same way as normal dogu config
func setSensitiveConfigValues(doguName cescommons.SimpleName, configEntries DoguConfigEntries, referencedValues map[common.DoguConfigKey]common.SensitiveDoguConfigValue) {
	for i, entry := range configEntries {
		if !entry.Absent && entry.Sensitive {
			key := common.DoguConfigKey{
				DoguName: doguName,
				Key:      entry.Key,
			}
			// we checked previously that all referenced values exist. Therefore, we need no error here.
			// in case of a bug, this will cause, that no expected value gets set while applying the blueprint.
			sensitiveValue := referencedValues[key]
			configEntries[i].Value = &sensitiveValue
		}
	}
}

// setConfigValues maps the referenced values to normal dogu config,
// so that the stateDiff can handle referenced dogu config the same way as normal dogu config
func setConfigValues(doguName cescommons.SimpleName, configEntries DoguConfigEntries, referencedValues map[common.DoguConfigKey]common.DoguConfigValue) {
	for i, entry := range configEntries {
		if !entry.Absent {
			key := common.DoguConfigKey{
				DoguName: doguName,
				Key:      entry.Key,
			}
			// we checked previously that all referenced values exist. Therefore, we need no error here.
			// in case of a bug, this will cause, that no expected value gets set while applying the blueprint.
			value := referencedValues[key]
			configEntries[i].Value = &value
		}
	}
}

func getNeededConfigAction(expected ConfigValueState, actual ConfigValueState) ConfigAction {
	if expected.Equal(actual) {
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
