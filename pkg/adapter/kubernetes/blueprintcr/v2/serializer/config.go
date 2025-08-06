package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

func ConvertToConfigDTO(config domain.Config) v2.Config {
	var dogus map[string]v2.CombinedDoguConfig
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[string]v2.CombinedDoguConfig, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[string(doguName)] = convertToCombinedDoguConfigDTO(doguConfig)
		}
	}

	return v2.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config v2.Config) domain.Config {
	var dogus map[cescommons.SimpleName]domain.CombinedDoguConfig
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[cescommons.SimpleName]domain.CombinedDoguConfig, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[cescommons.SimpleName(doguName)] = convertToCombinedDoguConfigDomain(doguName, doguConfig)
		}
	}

	return domain.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDomain(config.Global),
	}
}

func convertToCombinedDoguConfigDTO(config domain.CombinedDoguConfig) v2.CombinedDoguConfig {
	return v2.CombinedDoguConfig{
		Config:          convertToDoguConfigDTO(config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDTO(config.SensitiveConfig),
	}
}

func convertToCombinedDoguConfigDomain(doguName string, config v2.CombinedDoguConfig) domain.CombinedDoguConfig {
	return domain.CombinedDoguConfig{
		DoguName:        cescommons.SimpleName(doguName),
		Config:          convertToDoguConfigDomain(doguName, config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDomain(doguName, config.SensitiveConfig),
	}
}

func convertToDoguConfigDTO(config domain.DoguConfig) *v2.DoguConfig {
	// empty struct -> nil
	if len(config.Present) == 0 && len(config.Absent) == 0 {
		return nil
	}

	var present map[string]string
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Present) != 0 {
		present = make(map[string]string, len(config.Present))
		for key, value := range config.Present {
			present[string(key.Key)] = string(value)
		}
	}

	var absent []string
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Absent) != 0 {
		absent = make([]string, len(config.Absent))
		for i, key := range config.Absent {
			absent[i] = string(key.Key)
		}
	}

	return &v2.DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigDomain(doguName string, config *v2.DoguConfig) domain.DoguConfig {
	if config == nil {
		return domain.DoguConfig{}
	}

	var present map[common.DoguConfigKey]common.DoguConfigValue
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Present) != 0 {
		present = make(map[common.DoguConfigKey]common.DoguConfigValue, len(config.Present))
		for key, value := range config.Present {
			present[convertToDoguConfigKeyDomain(doguName, key)] = common.DoguConfigValue(value)
		}
	}

	var absent []common.DoguConfigKey
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Absent) != 0 {
		absent = make([]common.DoguConfigKey, len(config.Absent))
		for i, key := range config.Absent {
			absent[i] = convertToDoguConfigKeyDomain(doguName, key)
		}
	}

	return domain.DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigKeyDomain(doguName, key string) common.DoguConfigKey {
	return common.DoguConfigKey{
		DoguName: cescommons.SimpleName(doguName),
		Key:      config.Key(key),
	}
}

func convertToGlobalConfigDTO(config domain.GlobalConfig) v2.GlobalConfig {
	var present map[string]string
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Present) != 0 {
		present = make(map[string]string, len(config.Present))
		for key, value := range config.Present {
			present[string(key)] = string(value)
		}
	}

	var absent []string
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Absent) != 0 {
		absent = make([]string, len(config.Absent))
		for i, key := range config.Absent {
			absent[i] = string(key)
		}
	}

	return v2.GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToGlobalConfigDomain(config v2.GlobalConfig) domain.GlobalConfig {
	var present map[common.GlobalConfigKey]common.GlobalConfigValue
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Present) != 0 {
		present = make(map[common.GlobalConfigKey]common.GlobalConfigValue, len(config.Present))
		for key, value := range config.Present {
			present[common.GlobalConfigKey(key)] = common.GlobalConfigValue(value)
		}
	}

	var absent []common.GlobalConfigKey
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Absent) != 0 {
		absent = make([]common.GlobalConfigKey, len(config.Absent))
		for i, key := range config.Absent {
			absent[i] = common.GlobalConfigKey(key)
		}
	}

	return domain.GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}
