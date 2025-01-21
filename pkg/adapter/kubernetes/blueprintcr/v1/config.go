package v1

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	bpentities "github.com/cloudogu/k8s-blueprint-lib/json/entities"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
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

func ConvertToConfigDTO(config domain.Config) bpentities.TargetConfig {
	var dogus map[string]bpentities.CombinedDoguConfig
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[string]bpentities.CombinedDoguConfig, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[string(doguName)] = convertToCombinedDoguConfigDTO(doguConfig)
		}
	}

	return bpentities.TargetConfig{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config bpentities.TargetConfig) domain.Config {
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

func convertToCombinedDoguConfigDTO(config domain.CombinedDoguConfig) bpentities.CombinedDoguConfig {
	return bpentities.CombinedDoguConfig{
		Config:          convertToDoguConfigDTO(config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDTO(config.SensitiveConfig),
	}
}

func convertToCombinedDoguConfigDomain(doguName string, config bpentities.CombinedDoguConfig) domain.CombinedDoguConfig {
	return domain.CombinedDoguConfig{
		DoguName:        cescommons.SimpleName(doguName),
		Config:          convertToDoguConfigDomain(doguName, config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDomain(doguName, config.SensitiveConfig),
	}
}

func convertToDoguConfigDTO(config domain.DoguConfig) bpentities.DoguConfig {
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

	return bpentities.DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigDomain(doguName string, config bpentities.DoguConfig) domain.DoguConfig {
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

func convertToSensitiveDoguConfigDTO(config domain.SensitiveDoguConfig) bpentities.SensitiveDoguConfig {
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

	return bpentities.SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigDomain(doguName string, doguConfig bpentities.SensitiveDoguConfig) domain.SensitiveDoguConfig {
	var present map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(doguConfig.Present) != 0 {
		present = make(map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue, len(doguConfig.Present))
		for key, value := range doguConfig.Present {
			present[convertToSensitiveDoguConfigKeyDomain(doguName, key)] = common.SensitiveDoguConfigValue(value)
		}
	}

	var absent []common.SensitiveDoguConfigKey
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(doguConfig.Absent) != 0 {
		absent = make([]common.SensitiveDoguConfigKey, len(doguConfig.Absent))
		for i, key := range doguConfig.Absent {
			absent[i] = common.SensitiveDoguConfigKey{
				DoguName: cescommons.SimpleName(doguName),
				Key:      config.Key(key),
			}
		}
	}

	return domain.SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigKeyDomain(doguName, key string) common.SensitiveDoguConfigKey {
	return common.SensitiveDoguConfigKey{
		DoguName: cescommons.SimpleName(doguName),
		Key:      config.Key(key),
	}
}

func convertToGlobalConfigDTO(config domain.GlobalConfig) bpentities.GlobalConfig {
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

	return bpentities.GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToGlobalConfigDomain(config bpentities.GlobalConfig) domain.GlobalConfig {
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
