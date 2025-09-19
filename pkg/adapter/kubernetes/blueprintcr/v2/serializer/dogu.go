package serializer

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

func ConvertDogus(dogus []bpv2.Dogu) ([]domain.Dogu, error) {
	var convertedDogus []domain.Dogu
	var errorList []error

	for _, dogu := range dogus {
		result := domain.Dogu{}
		name, err := cescommons.QualifiedNameFromString(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		result.Name = name

		var version *core.Version
		if dogu.Version != nil && *dogu.Version != "" {
			coreVersion, err := core.ParseVersion(*dogu.Version)
			version = &coreVersion
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of target dogu %q: %w", dogu.Name, err))
				continue
			}
		}
		result.Version = version

		absent := false
		if dogu.Absent != nil {
			absent = *dogu.Absent
		}
		result.Absent = absent

		err = convertPlatformConfigFromDTOToDomain(&dogu, &result)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		convertedDogus = append(convertedDogus, result)
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func convertPlatformConfigFromDTOToDomain(dtoDogu *bpv2.Dogu, domainDogu *domain.Dogu) error {
	if dtoDogu.PlatformConfig == nil {
		return nil
	}

	var minVolumeSize, maxBodySize *resource.Quantity
	var additionalMounts []ecosystem.AdditionalMount
	var err error
	if dtoDogu.PlatformConfig.ResourceConfig != nil && dtoDogu.PlatformConfig.ResourceConfig.MinVolumeSize != nil {
		minVolumeSizeStr := dtoDogu.PlatformConfig.ResourceConfig.MinVolumeSize
		minVolumeSize, err = ecosystem.GetNonNilQuantityRef(*minVolumeSizeStr)
		if err != nil {
			return fmt.Errorf("could not parse minimum volume size %q for dogu %q", *minVolumeSizeStr, dtoDogu.Name)
		}
		domainDogu.MinVolumeSize = minVolumeSize
	}

	if dtoDogu.PlatformConfig.ReverseProxyConfig != nil {
		if dtoDogu.PlatformConfig.ReverseProxyConfig.MaxBodySize != nil {
			maxBodySizeStr := dtoDogu.PlatformConfig.ReverseProxyConfig.MaxBodySize
			maxBodySize, err = ecosystem.GetQuantityReference(*maxBodySizeStr)
			if err != nil {
				return fmt.Errorf("could not parse maximum proxy body size %q for dogu %q", *maxBodySizeStr, dtoDogu.Name)
			}
		}

		domainDogu.ReverseProxyConfig = &ecosystem.ReverseProxyConfig{
			MaxBodySize:      maxBodySize,
			RewriteTarget:    dtoDogu.PlatformConfig.ReverseProxyConfig.RewriteTarget,
			AdditionalConfig: dtoDogu.PlatformConfig.ReverseProxyConfig.AdditionalConfig,
		}
	}

	if dtoDogu.PlatformConfig.AdditionalMountsConfig != nil {
		additionalMounts = convertAdditionalMountsFromDTOToDomain(dtoDogu.PlatformConfig.AdditionalMountsConfig)
		domainDogu.AdditionalMounts = additionalMounts
	}

	return nil
}

func ConvertMaskDogus(dogus []bpv2.MaskDogu) ([]domain.MaskDogu, error) {
	var convertedDogus []domain.MaskDogu
	var errorList []error

	for _, dogu := range dogus {
		name, err := cescommons.QualifiedNameFromString(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		var version core.Version
		if dogu.Version != nil && *dogu.Version != "" {
			version, err = core.ParseVersion(*dogu.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of mask dogu %q: %w", dogu.Name, err))
				continue
			}
		}

		absent := false
		if dogu.Absent != nil {
			absent = *dogu.Absent
		}

		convertedDogus = append(convertedDogus, domain.MaskDogu{
			Name:    name,
			Version: version,
			Absent:  absent,
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

func ConvertToDoguDTOs(dogus []domain.Dogu) []bpv2.Dogu {
	converted := util.Map(dogus, func(dogu domain.Dogu) bpv2.Dogu {
		var version *string
		if dogu.Version != nil {
			version = &dogu.Version.Raw
		}
		return bpv2.Dogu{
			Name:           dogu.Name.String(),
			Version:        version,
			Absent:         &dogu.Absent,
			PlatformConfig: convertPlatformConfigDTO(dogu),
		}
	})
	return converted
}

func convertPlatformConfigDTO(dogu domain.Dogu) *bpv2.PlatformConfig {
	if dogu.ReverseProxyConfig == nil && dogu.MinVolumeSize == nil && len(dogu.AdditionalMounts) == 0 {
		return nil
	}

	config := bpv2.PlatformConfig{}
	config.ResourceConfig = convertResourceConfigDTO(dogu)
	config.ReverseProxyConfig = convertReverseProxyConfigDTO(dogu)
	config.AdditionalMountsConfig = convertAdditionalMountsConfigDTO(dogu)

	return &config
}

func convertReverseProxyConfigDTO(dogu domain.Dogu) *bpv2.ReverseProxyConfig {
	config := bpv2.ReverseProxyConfig{}
	if dogu.ReverseProxyConfig != nil {
		config.RewriteTarget = dogu.ReverseProxyConfig.RewriteTarget
		config.AdditionalConfig = dogu.ReverseProxyConfig.AdditionalConfig
		config.MaxBodySize = ecosystem.GetQuantityString(dogu.ReverseProxyConfig.MaxBodySize)
	}

	return &config
}

func convertResourceConfigDTO(dogu domain.Dogu) *bpv2.ResourceConfig {
	config := bpv2.ResourceConfig{}
	config.MinVolumeSize = ecosystem.GetQuantityString(dogu.MinVolumeSize)

	return &config
}

func convertAdditionalMountsConfigDTO(dogu domain.Dogu) []bpv2.AdditionalMount {
	var config []bpv2.AdditionalMount
	for _, m := range dogu.AdditionalMounts {
		config = append(config, bpv2.AdditionalMount{
			SourceType: bpv2.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		})
	}
	return config
}
