package v1

import (
	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
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

func ConvertToConfigDTO(config bpv2.Config) Config {
	var dogus map[string]CombinedDoguConfig
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[string]CombinedDoguConfig, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[string(doguName)] = convertToCombinedDoguConfigDTO(doguConfig)
		}
	}

	return Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config Config) bpv2.Config {
	var dogus map[cescommons.SimpleName]bpv2.CombinedDoguConfig
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[cescommons.SimpleName]bpv2.CombinedDoguConfig, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[cescommons.SimpleName(doguName)] = convertToCombinedDoguConfigDomain(doguName, doguConfig)
		}
	}

	return bpv2.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDomain(config.Global),
	}
}

func convertToCombinedDoguConfigDTO(config bpv2.CombinedDoguConfig) CombinedDoguConfig {
	return CombinedDoguConfig{
		Config:          convertToDoguConfigDTO(config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDTO(config.SensitiveConfig),
	}
}

func convertToCombinedDoguConfigDomain(doguName string, config CombinedDoguConfig) bpv2.CombinedDoguConfig {
	return bpv2.CombinedDoguConfig{
		DoguName:        cescommons.SimpleName(doguName),
		Config:          convertToDoguConfigDomain(doguName, config.Config),
		SensitiveConfig: convertToSensitiveDoguConfigDomain(doguName, config.SensitiveConfig),
	}
}

func convertToDoguConfigDTO(config bpv2.DoguConfig) DoguConfig {
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

	return DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigDomain(doguName string, config DoguConfig) bpv2.DoguConfig {
	var present map[bpv2.DoguConfigKey]bpv2.DoguConfigValue
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Present) != 0 {
		present = make(map[bpv2.DoguConfigKey]bpv2.DoguConfigValue, len(config.Present))
		for key, value := range config.Present {
			present[convertToDoguConfigKeyDomain(doguName, key)] = bpv2.DoguConfigValue(value)
		}
	}

	var absent []bpv2.DoguConfigKey
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Absent) != 0 {
		absent = make([]bpv2.DoguConfigKey, len(config.Absent))
		for i, key := range config.Absent {
			absent[i] = convertToDoguConfigKeyDomain(doguName, key)
		}
	}

	return bpv2.DoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToDoguConfigKeyDomain(doguName, key string) bpv2.DoguConfigKey {
	return bpv2.DoguConfigKey{
		DoguName: cescommons.SimpleName(doguName),
		Key:      config.Key(key),
	}
}

func convertToSensitiveDoguConfigDTO(config bpv2.SensitiveDoguConfig) SensitiveDoguConfig {
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

	return SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigDomain(doguName string, doguConfig SensitiveDoguConfig) bpv2.SensitiveDoguConfig {
	var present map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(doguConfig.Present) != 0 {
		present = make(map[bpv2.SensitiveDoguConfigKey]bpv2.SensitiveDoguConfigValue, len(doguConfig.Present))
		for key, value := range doguConfig.Present {
			present[convertToSensitiveDoguConfigKeyDomain(doguName, key)] = bpv2.SensitiveDoguConfigValue(value)
		}
	}

	var absent []bpv2.SensitiveDoguConfigKey
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(doguConfig.Absent) != 0 {
		absent = make([]bpv2.SensitiveDoguConfigKey, len(doguConfig.Absent))
		for i, key := range doguConfig.Absent {
			absent[i] = bpv2.SensitiveDoguConfigKey{
				DoguName: cescommons.SimpleName(doguName),
				Key:      config.Key(key),
			}
		}
	}

	return bpv2.SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigKeyDomain(doguName, key string) bpv2.SensitiveDoguConfigKey {
	return bpv2.SensitiveDoguConfigKey{
		DoguName: cescommons.SimpleName(doguName),
		Key:      config.Key(key),
	}
}

func convertToGlobalConfigDTO(config bpv2.GlobalConfig) GlobalConfig {
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

	return GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToGlobalConfigDomain(config GlobalConfig) bpv2.GlobalConfig {
	var present map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Present) != 0 {
		present = make(map[bpv2.GlobalConfigKey]bpv2.GlobalConfigValue, len(config.Present))
		for key, value := range config.Present {
			present[bpv2.GlobalConfigKey(key)] = bpv2.GlobalConfigValue(value)
		}
	}

	var absent []bpv2.GlobalConfigKey
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Absent) != 0 {
		absent = make([]bpv2.GlobalConfigKey, len(config.Absent))
		for i, key := range config.Absent {
			absent[i] = bpv2.GlobalConfigKey(key)
		}
	}

	return bpv2.GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}
