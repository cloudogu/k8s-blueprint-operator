package serializer

import (
	"errors"
	"fmt"
	bpentities "github.com/cloudogu/k8s-blueprint-lib/json/entities"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

func ConvertDogus(dogus []bpv2.Dogu) ([]domain.Dogu, error) {
	var convertedDogus []domain.Dogu
	var errorList []error

	for _, dogu := range dogus {
		name, err := cescommons.QualifiedNameFromString(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		var version core.Version
		if dogu.Version != "" {
			version, err = core.ParseVersion(dogu.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of target dogu %q: %w", dogu.Name, err))
				continue
			}
		}

		minVolumeSizeStr := dogu.PlatformConfig.ResourceConfig.MinVolumeSize
		minVolumeSize, minVolumeSizeErr := ecosystem.GetNonNilQuantityRef(minVolumeSizeStr)
		if minVolumeSizeErr != nil {
			errorList = append(errorList, fmt.Errorf("could not parse minimum volume size %q for dogu %q", minVolumeSizeStr, dogu.Name))
			continue
		}

		maxBodySizeStr := dogu.PlatformConfig.ReverseProxyConfig.MaxBodySize
		maxBodySize, maxBodySizeErr := ecosystem.GetQuantityReference(maxBodySizeStr)
		if maxBodySizeErr != nil {
			errorList = append(errorList, fmt.Errorf("could not parse maximum proxy body size %q for dogu %q", maxBodySizeStr, dogu.Name))
			continue
		}

		convertedDogus = append(convertedDogus, domain.Dogu{
			Name:          name,
			Version:       version,
			TargetState:   ToDomainTargetState(dogu.Absent),
			MinVolumeSize: *minVolumeSize,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      maxBodySize,
				RewriteTarget:    ecosystem.RewriteTarget(dogu.PlatformConfig.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: ecosystem.AdditionalConfig(dogu.PlatformConfig.ReverseProxyConfig.AdditionalConfig),
			},
			AdditionalMounts: convertAdditionalMountsFromDTOToDomain(dogu.PlatformConfig.AdditionalMountsConfig),
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func convertAdditionalMountsFromDTOToDomain(mounts []bpv2.AdditionalMount) []ecosystem.AdditionalMount {
	var result []ecosystem.AdditionalMount
	for _, m := range mounts {
		result = append(result, ecosystem.AdditionalMount{
			SourceType: ecosystem.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		})
	}

	return result
}

func ConvertToDoguDTOs(dogus []domain.Dogu) ([]bpv2.Dogu, error) {
	var errorList []error
	converted := util.Map(dogus, func(dogu domain.Dogu) bpv2.Dogu {
		return bpv2.Dogu{
			Name:           dogu.Name.String(),
			Version:        dogu.Version.Raw,
			Absent:         ToSerializerAbsentState(dogu.TargetState),
			PlatformConfig: convertPlatformConfigDTO(dogu),
		}
	})
	//TODO: remove error return value if there is no error needed after refactoring
	return converted, errors.Join(errorList...)
}

func convertPlatformConfigDTO(dogu domain.Dogu) bpv2.PlatformConfig {
	config := bpentities.PlatformConfig{}
	config.ResourceConfig = convertResourceConfigDTO(dogu)
	config.ReverseProxyConfig = convertReverseProxyConfigDTO(dogu)
	config.AdditionalMountsConfig = convertAdditionalMountsConfig(dogu)

	return config
}

func convertReverseProxyConfigDTO(dogu domain.Dogu) bpentities.ReverseProxyConfig {
	config := bpentities.ReverseProxyConfig{}
	config.RewriteTarget = string(dogu.ReverseProxyConfig.RewriteTarget)
	config.AdditionalConfig = string(dogu.ReverseProxyConfig.AdditionalConfig)
	config.MaxBodySize = ecosystem.GetQuantityString(dogu.ReverseProxyConfig.MaxBodySize)

	return config
}

func convertResourceConfigDTO(dogu domain.Dogu) bpentities.ResourceConfig {
	config := bpentities.ResourceConfig{}
	config.MinVolumeSize = ecosystem.GetQuantityString(&dogu.MinVolumeSize)

	return config
}

func convertAdditionalMountsConfig(dogu domain.Dogu) []bpentities.AdditionalMount {
	var config []bpentities.AdditionalMount
	for _, m := range dogu.AdditionalMounts {
		config = append(config, bpentities.AdditionalMount{
			SourceType: bpentities.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		})
	}
	return config
}
