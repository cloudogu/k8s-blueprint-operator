package blueprintV2

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

type Config struct {
	Dogus  map[string]CombinedDoguConfig `json:"dogus,omitempty"`
	Global GlobalConfig                  `json:"global,omitempty"`
}

type CombinedDoguConfig struct {
	Config          DoguConfig          `json:"config,omitempty"`
	SensitiveConfig SensitiveDoguConfig `json:"sensitiveConfig,omitempty"`
}

type DoguConfig presentAbsentConfig

type SensitiveDoguConfig presentAbsentConfig

type GlobalConfig presentAbsentConfig

type presentAbsentConfig struct {
	Present map[string]string `json:"present,omitempty"`
	Absent  []string          `json:"absent,omitempty"`
}

func ConvertToConfigDTO(config domain.Config) Config {
	dogus := make(map[string]CombinedDoguConfig, len(config.Dogus))
	for doguName, doguConfig := range config.Dogus {
		dogus[string(doguName)] = convertToCombinedDoguConfigDTO(doguConfig)
	}
	return Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config Config) domain.Config {
	dogus := make(map[common.SimpleDoguName]domain.CombinedDoguConfig, len(config.Dogus))
	for doguName, doguConfig := range config.Dogus {
		dogus[common.SimpleDoguName(doguName)] = convertToCombinedDoguConfigDomain(doguName, doguConfig)
	}
	return domain.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDomain(config.Global),
	}
}

func convertToCombinedDoguConfigDTO(config domain.CombinedDoguConfig) CombinedDoguConfig {
	return CombinedDoguConfig{
		Config:          convertToDoguConfigDTO(config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDTO(config.SensitiveConfig),
	}
}

func convertToCombinedDoguConfigDomain(doguName string, config CombinedDoguConfig) domain.CombinedDoguConfig {
	return domain.CombinedDoguConfig{
		DoguName:        common.SimpleDoguName(doguName),
		Config:          convertToDoguConfigDomain(doguName, config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDomain(doguName, config.SensitiveConfig),
	}
}

func convertToDoguConfigDTO(config domain.DoguConfig) DoguConfig {
	present := make(map[string]string, len(config.Present))
	for key, value := range config.Present {
		present[key.Key] = string(value)
	}

	absent := make([]string, len(config.Absent))
	for i, key := range config.Absent {
		absent[i] = key.Key
	}

	return DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigDomain(doguName string, config DoguConfig) domain.DoguConfig {
	present := make(map[common.DoguConfigKey]common.DoguConfigValue, len(config.Present))
	for key, value := range config.Present {
		present[convertToDoguConfigKeyDomain(doguName, key)] = common.DoguConfigValue(value)
	}

	absent := make([]common.DoguConfigKey, len(config.Absent))
	for i, key := range config.Absent {
		absent[i] = convertToDoguConfigKeyDomain(doguName, key)
	}

	return domain.DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigKeyDomain(doguName, key string) common.DoguConfigKey {
	return common.DoguConfigKey{
		DoguName: common.SimpleDoguName(doguName),
		Key:      key,
	}
}

func convertToSensitiveDoguConfigDTO(config domain.SensitiveDoguConfig) SensitiveDoguConfig {
	present := make(map[string]string, len(config.Present))
	for key, value := range config.Present {
		present[key.Key] = string(value)
	}

	absent := make([]string, len(config.Absent))
	for i, key := range config.Absent {
		absent[i] = key.Key
	}

	return SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigDomain(doguName string, config SensitiveDoguConfig) domain.SensitiveDoguConfig {
	present := make(map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, len(config.Present))
	for key, value := range config.Present {
		present[convertToSensitiveDoguConfigKeyDomain(doguName, key)] = common.SensitiveDoguConfigValue(value)
	}

	absent := make([]common.SensitiveDoguConfigKey, len(config.Absent))
	for i, key := range config.Absent {
		absent[i] = convertToSensitiveDoguConfigKeyDomain(doguName, key)
	}

	return domain.SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigKeyDomain(doguName, key string) common.SensitiveDoguConfigKey {
	return common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{
		DoguName: common.SimpleDoguName(doguName),
		Key:      key,
	},
	}
}

func convertToGlobalConfigDTO(config domain.GlobalConfig) GlobalConfig {
	present := make(map[string]string, len(config.Present))
	for key, value := range config.Present {
		present[string(key)] = string(value)
	}

	absent := make([]string, len(config.Absent))
	for i, key := range config.Absent {
		absent[i] = string(key)
	}

	return GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToGlobalConfigDomain(config GlobalConfig) domain.GlobalConfig {
	present := make(map[common.GlobalConfigKey]common.GlobalConfigValue, len(config.Present))
	for key, value := range config.Present {
		present[common.GlobalConfigKey(key)] = common.GlobalConfigValue(value)
	}

	absent := make([]common.GlobalConfigKey, len(config.Absent))
	for i, key := range config.Absent {
		absent[i] = common.GlobalConfigKey(key)
	}

	return domain.GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}
