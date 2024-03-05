package serializer

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"
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

func ConvertDogus(dogus []TargetDogu) ([]domain.Dogu, error) {
	var convertedDogus []domain.Dogu
	var errorList []error

	for _, dogu := range dogus {
		name, err := common.QualifiedDoguNameFromString(dogu.Name)
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

		minVolumeSize := dogu.PlatformConfig.ResourceConfig.MinVolumeSize
		var volumeQuantity ecosystem.VolumeSize
		var volumeQuantityErr error
		if minVolumeSize != "" {
			volumeQuantity, volumeQuantityErr = resource.ParseQuantity(minVolumeSize)
		}
		if volumeQuantityErr != nil {
			errorList = append(errorList, fmt.Errorf("could not parse minimum volume size %q for dogu %q", minVolumeSize, dogu.Name))
		}

		var maxBodySize *ecosystem.BodySize
		var parse ecosystem.BodySize
		var maxBodyQuantityErr error
		maxProxyBodySize := dogu.PlatformConfig.ReverseProxyConfig.MaxBodySize
		if maxProxyBodySize != "" {
			parse, maxBodyQuantityErr = resource.ParseQuantity(maxProxyBodySize)
			if maxBodyQuantityErr != nil {
				errorList = append(errorList, fmt.Errorf("could not parse maximum proxy body size %q for dogu %q", maxProxyBodySize, dogu.Name))
			}
			maxBodySize = &parse
		}

		convertedDogus = append(convertedDogus, domain.Dogu{
			Name:          name,
			Version:       version,
			TargetState:   newState,
			MinVolumeSize: volumeQuantity,
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

func ConvertToDoguDTOs(dogus []domain.Dogu) ([]TargetDogu, error) {
	var errorList []error
	converted := util.Map(dogus, func(dogu domain.Dogu) TargetDogu {
		newState, err := ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)

		maxBodySize := ""
		domainSize := dogu.ReverseProxyConfig.MaxBodySize
		if domainSize != nil {
			maxBodySize = domainSize.String()
		}
		// TODO
		minVolumeSize := dogu.MinVolumeSize.String()
		if minVolumeSize == "0" {
			minVolumeSize = ""
		}
		return TargetDogu{
			Name:        dogu.Name.String(),
			Version:     dogu.Version.Raw,
			TargetState: newState,
			PlatformConfig: PlatformConfig{
				ResourceConfig: ResourceConfig{
					MinVolumeSize: minVolumeSize,
				},
				ReverseProxyConfig: ReverseProxyConfig{
					MaxBodySize:      maxBodySize,
					RewriteTarget:    string(dogu.ReverseProxyConfig.RewriteTarget),
					AdditionalConfig: string(dogu.ReverseProxyConfig.AdditionalConfig),
				},
			},
		}
	})
	return converted, errors.Join(errorList...)
}
