package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

func convertToSensitiveDoguConfigDTO(config domain.SensitiveDoguConfig) []v2.ConfigEntry {
	// empty struct -> nil
	if len(config.Absent) == 0 && len(config.Present) == 0 {
		return nil
	}

	result := make([]v2.ConfigEntry, len(config.Present)+len(config.Absent))
	// we check for empty values to make good use of default values
	// this makes testing easier
	for key, value := range config.Present {
		result = append(result, v2.ConfigEntry{
			Key: string(key.Key),
			SecretRef: &v2.SecretReference{
				Name: value.SecretName,
				Key:  value.SecretKey,
			},
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

func convertToSensitiveDoguConfigDomain(doguName string, doguConfig []v2.ConfigEntry) domain.SensitiveDoguConfig {
	if doguConfig == nil {
		return domain.SensitiveDoguConfig{}
	}

	present := make(map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef, len(doguConfig))
	absent := make([]common.SensitiveDoguConfigKey, len(doguConfig))

	absentIndex := 0
	for _, configEntry := range doguConfig {
		if configEntry.Absent == nil || !*configEntry.Absent {
			if configEntry.Sensitive != nil && !*configEntry.Sensitive || configEntry.SecretRef == nil {
				continue
			}
			present[convertToSensitiveDoguConfigKeyDomain(doguName, configEntry.Key)] = domain.SensitiveValueRef{
				SecretName: configEntry.SecretRef.Name,
				SecretKey:  configEntry.SecretRef.Key,
			}
		} else {
			absent[absentIndex] = common.SensitiveDoguConfigKey{
				DoguName: cescommons.SimpleName(doguName),
				Key:      config.Key(configEntry.Key),
			}
			absentIndex++
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
