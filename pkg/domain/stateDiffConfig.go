package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
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

func countByAction(diffsByDogu map[cescommons.SimpleName]DoguConfigDiffs) map[ConfigAction]int {
	countByAction := map[ConfigAction]int{}
	for _, doguDiffs := range diffsByDogu {
		for _, diff := range doguDiffs {
			countByAction[diff.NeededAction]++
		}
	}
	return countByAction
}

func determineConfigDiffs(
	blueprintConfig *Config,
	globalConfig config.GlobalConfig,
	configByDogu map[cescommons.SimpleName]config.DoguConfig,
	SensitiveConfigByDogu map[cescommons.SimpleName]config.DoguConfig,
	referencedSensitiveConfig map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue,
) (
	map[cescommons.SimpleName]DoguConfigDiffs,
	map[cescommons.SimpleName]SensitiveDoguConfigDiffs,
	GlobalConfigDiffs,
) {
	if blueprintConfig == nil {
		return nil, nil, nil
	}

	return determineDogusConfigDiffs(blueprintConfig.Dogus, configByDogu),
		determineSensitiveDogusConfigDiffs(blueprintConfig.Dogus, SensitiveConfigByDogu, referencedSensitiveConfig),
		determineGlobalConfigDiffs(blueprintConfig.Global, globalConfig)
}

func determineDogusConfigDiffs(
	combinedDoguConfigs map[cescommons.SimpleName]CombinedDoguConfig,
	configByDogu map[cescommons.SimpleName]config.DoguConfig,
) map[cescommons.SimpleName]DoguConfigDiffs {
	diffsPerDogu := map[cescommons.SimpleName]DoguConfigDiffs{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		diffsPerDogu[doguName] = determineDoguConfigDiffs(combinedDoguConfig.Config, configByDogu)
	}
	return diffsPerDogu
}

func determineSensitiveDogusConfigDiffs(
	combinedDoguConfigs map[cescommons.SimpleName]CombinedDoguConfig,
	configByDogu map[cescommons.SimpleName]config.DoguConfig,
	referencedValues map[common.DoguConfigKey]common.SensitiveDoguConfigValue,
) map[cescommons.SimpleName]DoguConfigDiffs {
	diffsPerDogu := map[cescommons.SimpleName]DoguConfigDiffs{}
	for doguName, combinedDoguConfig := range combinedDoguConfigs {
		expectedSensitiveConfig := createDoguConfigFromReferencedValues(combinedDoguConfig.SensitiveConfig, referencedValues)
		diffsPerDogu[doguName] = determineDoguConfigDiffs(expectedSensitiveConfig, configByDogu)
	}
	return diffsPerDogu
}

// createDoguConfigFromReferencedValues maps the referenced values to normal dogu config,
// so that the stateDiff can handle sensitive dogu config the same way as normal dogu config
func createDoguConfigFromReferencedValues(
	sensitiveConfig SensitiveDoguConfig,
	referencedValues map[common.DoguConfigKey]common.SensitiveDoguConfigValue,
) DoguConfig {
	configEntries := map[common.DoguConfigKey]common.DoguConfigValue{}
	for key := range sensitiveConfig.Present {
		// we checked previously that all referenced values exist. Therefore, we need no error here.
		// in case of a bug, this will cause, that no expected value gets set while applying the blueprint.
		configEntries[key] = referencedValues[key]
	}
	return DoguConfig{
		Present: configEntries,
		Absent:  sensitiveConfig.Absent,
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
