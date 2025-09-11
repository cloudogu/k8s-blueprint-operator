package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

func ConvertToConfigDTO(config *domain.Config) *v2.Config {
	if config == nil {
		return nil
	}

	var dogus map[string][]v2.ConfigEntry
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[string][]v2.ConfigEntry, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[string(doguName)] = convertToCombinedDoguConfigDTO(doguConfig)
		}
	}

	return &v2.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config *v2.Config) *domain.Config {
	if config == nil {
		return nil
	}
	var dogus map[cescommons.SimpleName]domain.CombinedDoguConfig
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[cescommons.SimpleName]domain.CombinedDoguConfig, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[cescommons.SimpleName(doguName)] = convertToCombinedDoguConfigDomain(doguName, doguConfig)
		}
	}

	return &domain.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDomain(config.Global),
	}
}

func convertToCombinedDoguConfigDTO(config domain.CombinedDoguConfig) []v2.ConfigEntry {
	var result []v2.ConfigEntry
	result = append(result, convertToDoguConfigDTO(config.Config)...)
	result = append(result, convertToSensitiveDoguConfigDTO(config.SensitiveConfig)...)

	return result
}

func convertToCombinedDoguConfigDomain(doguName string, config []v2.ConfigEntry) domain.CombinedDoguConfig {
	return domain.CombinedDoguConfig{
		DoguName:        cescommons.SimpleName(doguName),
		Config:          convertToDoguConfigDomain(doguName, config),
		SensitiveConfig: convertToSensitiveDoguConfigDomain(doguName, config),
	}
}

func convertToDoguConfigDTO(config domain.DoguConfig) []v2.ConfigEntry {
	// empty struct -> nil
	if len(config.Present) == 0 && len(config.Absent) == 0 {
		return nil
	}

	result := make([]v2.ConfigEntry, len(config.Present)+len(config.Absent))
	// we check for empty values to make good use of default values
	// this makes testing easier
	for key, value := range config.Present {
		valueString := string(value)
		result = append(result, v2.ConfigEntry{
			Key:   string(key.Key),
			Value: &valueString,
		})
	}

	for _, key := range config.Absent {
		truePtr := true
		result = append(result, v2.ConfigEntry{
			Key:    string(key.Key),
			Absent: &truePtr,
		})
	}

	return result
}

func convertToDoguConfigDomain(doguName string, config []v2.ConfigEntry) domain.DoguConfig {
	if config == nil {
		return domain.DoguConfig{}
	}

	present := make(map[common.DoguConfigKey]common.DoguConfigValue, len(config))
	absent := make([]common.DoguConfigKey, len(config))

	absentIndex := 0
	for _, configEntry := range config {
		if configEntry.Sensitive != nil && *configEntry.Sensitive == true {
			continue
		}

		if configEntry.Absent == nil || *configEntry.Absent == false && configEntry.Value != nil {
			present[convertToDoguConfigKeyDomain(doguName, configEntry.Key)] = common.DoguConfigValue(*configEntry.Value)
		} else {
			absent[absentIndex] = convertToDoguConfigKeyDomain(doguName, configEntry.Key)
			absentIndex++
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

func convertToGlobalConfigDTO(config domain.GlobalConfig) []v2.ConfigEntry {
	// empty struct -> nil
	if len(config.Present) == 0 && len(config.Absent) == 0 {
		return nil
	}

	result := make([]v2.ConfigEntry, len(config.Present)+len(config.Absent))
	// we check for empty values to make good use of default values
	// this makes testing easier
	for key, value := range config.Present {
		valueString := string(value)
		result = append(result, v2.ConfigEntry{
			Key:   string(key),
			Value: &valueString,
		})
	}

	for _, key := range config.Absent {
		truePtr := true
		result = append(result, v2.ConfigEntry{
			Key:    string(key),
			Absent: &truePtr,
		})
	}

	return result
}

func convertToGlobalConfigDomain(config []v2.ConfigEntry) domain.GlobalConfig {
	if config == nil {
		return domain.GlobalConfig{}
	}

	present := make(map[common.GlobalConfigKey]common.GlobalConfigValue, len(config))
	absent := make([]common.GlobalConfigKey, len(config))

	absentIndex := 0
	for _, configEntry := range config {
		if configEntry.Absent != nil && *configEntry.Absent == false && configEntry.Value != nil {
			present[common.GlobalConfigKey(configEntry.Key)] = common.GlobalConfigValue(*configEntry.Value)
		} else {
			absent[absentIndex] = common.GlobalConfigKey(configEntry.Key)
			absentIndex++
		}
	}

	return domain.GlobalConfig{
		Present: present,
		Absent:  absent,
	}
}
