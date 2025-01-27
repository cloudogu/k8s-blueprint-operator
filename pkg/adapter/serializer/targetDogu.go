package serializer

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	bpentities "github.com/cloudogu/k8s-blueprint-lib/json/entities"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

func ConvertDogus(dogus []bpentities.TargetDogu) ([]domain.Dogu, error) {
	var convertedDogus []domain.Dogu
	var errorList []error

	for _, dogu := range dogus {
		name, err := cescommons.QualifiedNameFromString(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		newState, err := ToDomainTargetState(dogu.TargetState)
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
		minVolumeSize, minVolumeSizeErr := ecosystem.GetQuantityReference(minVolumeSizeStr)
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
			TargetState:   newState,
			MinVolumeSize: minVolumeSize,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      maxBodySize,
				RewriteTarget:    ecosystem.RewriteTarget(dogu.PlatformConfig.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: ecosystem.AdditionalConfig(dogu.PlatformConfig.ReverseProxyConfig.AdditionalConfig),
			},
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func ConvertToDoguDTOs(dogus []domain.Dogu) ([]bpentities.TargetDogu, error) {
	var errorList []error
	converted := util.Map(dogus, func(dogu domain.Dogu) bpentities.TargetDogu {
		newState, err := ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)

		return bpentities.TargetDogu{
			Name:           dogu.Name.String(),
			Version:        dogu.Version.Raw,
			TargetState:    newState,
			PlatformConfig: convertPlatformConfigDTO(dogu),
		}
	})
	return converted, errors.Join(errorList...)
}

func convertPlatformConfigDTO(dogu domain.Dogu) bpentities.PlatformConfig {
	config := bpentities.PlatformConfig{}
	config.ResourceConfig = convertResourceConfigDTO(dogu)
	config.ReverseProxyConfig = convertReverseProxyConfigDTO(dogu)

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
	config.MinVolumeSize = ecosystem.GetQuantityString(dogu.MinVolumeSize)

	return config
}
