package serializer

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

func ConvertDogus(dogus []bpv3.Dogu) ([]domain.Dogu, error) {
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
		result.Absent = ptr.Deref(dogu.Absent, false)

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

func convertPlatformConfigFromDTOToDomain(dtoDogu *bpv3.Dogu, domainDogu *domain.Dogu) error {
	if dtoDogu.PlatformConfig == nil {
		return nil
	}

	var minVolumeSize, maxBodySize *resource.Quantity
	var additionalMounts []ecosystem.AdditionalMount
	var err error
	if dtoDogu.PlatformConfig.ResourceConfig != nil {
		domainDogu.StorageClassName = dtoDogu.PlatformConfig.ResourceConfig.StorageClassName

		if dtoDogu.PlatformConfig.ResourceConfig.MinVolumeSize != nil {
			minVolumeSizeStr := dtoDogu.PlatformConfig.ResourceConfig.MinVolumeSize
			minVolumeSize, err = ecosystem.GetNonNilQuantityRef(*minVolumeSizeStr)
			if err != nil {
				return fmt.Errorf("could not parse minimum volume size %q for dogu %q", *minVolumeSizeStr, dtoDogu.Name)
			}
			domainDogu.MinVolumeSize = minVolumeSize
		}
	}

	if dtoDogu.PlatformConfig.ReverseProxyConfig != nil {
		if dtoDogu.PlatformConfig.ReverseProxyConfig.MaxBodySize != nil {
			maxBodySizeStr := dtoDogu.PlatformConfig.ReverseProxyConfig.MaxBodySize
			maxBodySize, err = ecosystem.GetQuantityReference(*maxBodySizeStr)
			if err != nil {
				return fmt.Errorf("could not parse maximum proxy body size %q for dogu %q", *maxBodySizeStr, dtoDogu.Name)
			}
		}

		domainDogu.ReverseProxyConfig = ecosystem.ReverseProxyConfig{
			MaxBodySize:      maxBodySize,
			RewriteTarget:    ecosystem.RewriteTarget(ptr.Deref(dtoDogu.PlatformConfig.ReverseProxyConfig.RewriteTarget, "")),
			AdditionalConfig: ecosystem.AdditionalConfig(ptr.Deref(dtoDogu.PlatformConfig.ReverseProxyConfig.AdditionalConfig, "")),
		}
	}

	if dtoDogu.PlatformConfig.AdditionalMountsConfig != nil {
		additionalMounts = convertAdditionalMountsFromDTOToDomain(dtoDogu.PlatformConfig.AdditionalMountsConfig)
		domainDogu.AdditionalMounts = additionalMounts
	}

	return nil
}

func ConvertMaskDogus(dogus []bpv3.MaskDogu) ([]domain.MaskDogu, error) {
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

		convertedDogus = append(convertedDogus, domain.MaskDogu{
			Name:    name,
			Version: version,
			Absent:  ptr.Deref(dogu.Absent, false),
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func convertAdditionalMountsFromDTOToDomain(mounts []bpv3.AdditionalMount) []ecosystem.AdditionalMount {
	var result []ecosystem.AdditionalMount
	for _, m := range mounts {
		result = append(result, ecosystem.AdditionalMount{
			SourceType: ecosystem.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  ptr.Deref(m.Subfolder, ""),
		})
	}

	return result
}

func ConvertToDoguDTOs(dogus []domain.Dogu) []bpv3.Dogu {
	converted := util.Map(dogus, func(dogu domain.Dogu) bpv3.Dogu {
		var version *string
		if dogu.Version != nil {
			version = &dogu.Version.Raw
		}
		return bpv3.Dogu{
			Name:           dogu.Name.String(),
			Version:        version,
			Absent:         &dogu.Absent,
			PlatformConfig: convertPlatformConfigDTO(dogu),
		}
	})
	return converted
}

func convertPlatformConfigDTO(dogu domain.Dogu) *bpv3.PlatformConfig {
	if dogu.ReverseProxyConfig.IsEmpty() && dogu.MinVolumeSize == nil && len(dogu.AdditionalMounts) == 0 {
		return nil
	}

	config := bpv3.PlatformConfig{}
	config.ResourceConfig = convertResourceConfigDTO(dogu)
	config.ReverseProxyConfig = convertReverseProxyConfigDTO(dogu)
	config.AdditionalMountsConfig = convertAdditionalMountsConfigDTO(dogu)

	return &config
}

func convertReverseProxyConfigDTO(dogu domain.Dogu) *bpv3.ReverseProxyConfig {
	var rewriteTarget, additionalConfig *string
	if dogu.ReverseProxyConfig.RewriteTarget != "" {
		rewriteTarget = (*string)(&dogu.ReverseProxyConfig.RewriteTarget)
	}
	if dogu.ReverseProxyConfig.AdditionalConfig != "" {
		additionalConfig = (*string)(&dogu.ReverseProxyConfig.AdditionalConfig)
	}
	return &bpv3.ReverseProxyConfig{
		RewriteTarget:    rewriteTarget,
		AdditionalConfig: additionalConfig,
		MaxBodySize:      ecosystem.GetQuantityString(dogu.ReverseProxyConfig.MaxBodySize),
	}
}

func convertResourceConfigDTO(dogu domain.Dogu) *bpv3.ResourceConfig {
	config := bpv3.ResourceConfig{}
	config.MinVolumeSize = ecosystem.GetQuantityString(dogu.MinVolumeSize)
	config.StorageClassName = dogu.StorageClassName

	return &config
}

func convertAdditionalMountsConfigDTO(dogu domain.Dogu) []bpv3.AdditionalMount {
	var config []bpv3.AdditionalMount
	for _, m := range dogu.AdditionalMounts {
		var subfolder *string
		if m.Subfolder != "" {
			subfolder = &m.Subfolder
		}
		config = append(config, bpv3.AdditionalMount{
			SourceType: bpv3.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  subfolder,
		})
	}
	return config
}
