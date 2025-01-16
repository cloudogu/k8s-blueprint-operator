package serializer

import (
	"errors"
	"fmt"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

// TargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type TargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState    string         `json:"targetState"`
	PlatformConfig PlatformConfig `json:"platformConfig,omitempty"`
}

type ResourceConfig struct {
	MinVolumeSize string `json:"minVolumeSize,omitempty"`
}

type ReverseProxyConfig struct {
	MaxBodySize      string `json:"maxBodySize,omitempty"`
	RewriteTarget    string `json:"rewriteTarget,omitempty"`
	AdditionalConfig string `json:"additionalConfig,omitempty"`
}

type PlatformConfig struct {
	ResourceConfig     ResourceConfig     `json:"resource,omitempty"`
	ReverseProxyConfig ReverseProxyConfig `json:"reverseProxy,omitempty"`
}

func ConvertDogus(dogus []TargetDogu) ([]bpv2.Dogu, error) {
	var convertedDogus []bpv2.Dogu
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

		convertedDogus = append(convertedDogus, bpv2.Dogu{
			Name:          name,
			Version:       version,
			TargetState:   newState,
			MinVolumeSize: minVolumeSize,
			ReverseProxyConfig: bpv2.ReverseProxyConfig{
				MaxBodySize:      maxBodySize,
				RewriteTarget:    bpv2.RewriteTarget(dogu.PlatformConfig.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: bpv2.AdditionalConfig(dogu.PlatformConfig.ReverseProxyConfig.AdditionalConfig),
			},
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func ConvertToDoguDTOs(dogus []bpv2.Dogu) ([]TargetDogu, error) {
	var errorList []error
	converted := util.Map(dogus, func(dogu bpv2.Dogu) TargetDogu {
		newState, err := ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)

		return TargetDogu{
			Name:           dogu.Name.String(),
			Version:        dogu.Version.Raw,
			TargetState:    newState,
			PlatformConfig: convertPlatformConfigDTO(dogu),
		}
	})
	return converted, errors.Join(errorList...)
}

func convertPlatformConfigDTO(dogu bpv2.Dogu) PlatformConfig {
	config := PlatformConfig{}
	config.ResourceConfig = convertResourceConfigDTO(dogu)
	config.ReverseProxyConfig = convertReverseProxyConfigDTO(dogu)

	return config
}

func convertReverseProxyConfigDTO(dogu bpv2.Dogu) ReverseProxyConfig {
	config := ReverseProxyConfig{}
	config.RewriteTarget = string(dogu.ReverseProxyConfig.RewriteTarget)
	config.AdditionalConfig = string(dogu.ReverseProxyConfig.AdditionalConfig)
	config.MaxBodySize = ecosystem.GetQuantityString(dogu.ReverseProxyConfig.MaxBodySize)

	return config
}

func convertResourceConfigDTO(dogu bpv2.Dogu) ResourceConfig {
	config := ResourceConfig{}
	config.MinVolumeSize = ecosystem.GetQuantityString(dogu.MinVolumeSize)

	return config
}
