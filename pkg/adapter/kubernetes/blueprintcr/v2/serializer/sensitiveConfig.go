package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

func convertToSensitiveDoguConfigDTO(config domain.SensitiveDoguConfig) *v2.SensitiveDoguConfig {
	var present []string
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

	return &v2.SensitiveDoguConfig{
		Present: present,
		Absent:  absent,
	}
}

func convertToSensitiveDoguConfigDomain(doguName string, doguConfig v2.SensitiveDoguConfig) domain.SensitiveDoguConfig {
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
