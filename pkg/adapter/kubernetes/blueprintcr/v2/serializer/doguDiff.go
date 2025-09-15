package serializer

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/api/resource"
)

func convertToDoguDiffDTO(domainModel domain.DoguDiff) crd.DoguDiff {
	neededActions := domainModel.NeededActions
	doguActions := make([]crd.DoguAction, 0, len(neededActions))
	for _, action := range neededActions {
		doguActions = append(doguActions, crd.DoguAction(action))
	}

	return crd.DoguDiff{
		Actual:        convertToDoguDiffStateDTO(domainModel.Actual),
		Expected:      convertToDoguDiffStateDTO(domainModel.Expected),
		NeededActions: doguActions,
	}
}

func convertToDoguDiffStateDTO(domainModel domain.DoguDiffState) crd.DoguDiffState {
	var version *string
	if domainModel.Version != nil {
		version = &domainModel.Version.Raw
	}

	var reverseProxyConfig *crd.ReverseProxyConfig
	if domainModel.ReverseProxyConfig != nil {
		reverseProxyConfig = &crd.ReverseProxyConfig{
			RewriteTarget:    domainModel.ReverseProxyConfig.RewriteTarget,
			AdditionalConfig: domainModel.ReverseProxyConfig.AdditionalConfig,
			MaxBodySize:      ecosystem.GetQuantityString(domainModel.ReverseProxyConfig.MaxBodySize),
		}
	}

	var resourceConfig *crd.ResourceConfig
	if domainModel.MinVolumeSize != nil {
		resourceConfig = &crd.ResourceConfig{
			MinVolumeSize: ecosystem.GetQuantityString(domainModel.MinVolumeSize),
		}
	}
	return crd.DoguDiffState{
		Namespace:          string(domainModel.Namespace),
		Version:            version,
		Absent:             domainModel.Absent,
		ResourceConfig:     resourceConfig,
		ReverseProxyConfig: reverseProxyConfig,
		AdditionalMounts:   convertAdditionalMountsToDoguDiffDTO(domainModel.AdditionalMounts),
	}
}

func convertAdditionalMountsToDoguDiffDTO(mounts []ecosystem.AdditionalMount) []crd.AdditionalMount {
	if len(mounts) == 0 {
		// an empty slice and nil are serialized differently
		// we want no entry instead of an empty json list if there are no mounts given
		return nil
	}
	result := make([]crd.AdditionalMount, len(mounts))
	for index, mount := range mounts {
		result[index] = crd.AdditionalMount{
			SourceType: crd.DataSourceType(mount.SourceType),
			Name:       mount.Name,
			Volume:     mount.Volume,
			Subfolder:  mount.Subfolder,
		}
	}
	return result
}

func convertToDoguDiffDomain(doguName string, dto crd.DoguDiff) (domain.DoguDiff, error) {

	actualState, actualStateErr := convertDoguDiffStateDomain(dto.Actual)
	if actualStateErr != nil {
		actualStateErr = fmt.Errorf("failed to convert actual dogu diff state: %w", actualStateErr)
	}
	expectedState, expectedStateErr := convertDoguDiffStateDomain(dto.Expected)
	if expectedStateErr != nil {
		expectedStateErr = fmt.Errorf("failed to convert expected dogu diff state: %w", expectedStateErr)
	}

	err := errors.Join(actualStateErr, expectedStateErr)
	if err != nil {
		return domain.DoguDiff{}, fmt.Errorf("failed to convert dogu diff dto %q to domain model: %w", doguName, err)
	}

	neededActions := dto.NeededActions
	doguActions := make([]domain.Action, 0, len(neededActions))
	for _, action := range neededActions {
		doguActions = append(doguActions, domain.Action(action))
	}

	return domain.DoguDiff{
		DoguName:      cescommons.SimpleName(doguName),
		Expected:      expectedState,
		Actual:        actualState,
		NeededActions: doguActions,
	}, nil
}

func convertDoguDiffStateDomain(dto crd.DoguDiffState) (domain.DoguDiffState, error) {
	var errorList []error

	var version *core.Version
	var err error
	if dto.Version != nil && *dto.Version != "" {
		var coreVersion core.Version
		coreVersion, err = core.ParseVersion(*dto.Version)
		version = &coreVersion
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to parse version %q: %w", *dto.Version, err))
		}
	}

	var minVolumeSize, maxBodySize *resource.Quantity
	if dto.ResourceConfig != nil && dto.ResourceConfig.MinVolumeSize != nil {
		minVolumeSizeStr := dto.ResourceConfig.MinVolumeSize
		minVolumeSize, err = ecosystem.GetNonNilQuantityRef(*minVolumeSizeStr)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to parse minimum volume size %q: %w", *minVolumeSizeStr, err))
		}
	}

	var reverseProxyConfig *ecosystem.ReverseProxyConfig
	if dto.ReverseProxyConfig != nil {
		if dto.ReverseProxyConfig.MaxBodySize != nil {
			maxBodySizeStr := dto.ReverseProxyConfig.MaxBodySize
			maxBodySize, err = ecosystem.GetQuantityReference(*maxBodySizeStr)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("failed to parse maximum proxy body size %q: %w", *maxBodySizeStr, err))
			}
		}
		reverseProxyConfig = &ecosystem.ReverseProxyConfig{
			MaxBodySize:      maxBodySize,
			RewriteTarget:    dto.ReverseProxyConfig.RewriteTarget,
			AdditionalConfig: dto.ReverseProxyConfig.AdditionalConfig,
		}
	}

	if len(errorList) != 0 {
		return domain.DoguDiffState{}, errors.Join(errorList...)
	}

	return domain.DoguDiffState{
		Namespace:          cescommons.Namespace(dto.Namespace),
		Version:            version,
		Absent:             dto.Absent,
		MinVolumeSize:      minVolumeSize,
		ReverseProxyConfig: reverseProxyConfig,
		AdditionalMounts:   convertAdditionalMountsToDoguDiffDomain(dto.AdditionalMounts),
	}, nil
}

func convertAdditionalMountsToDoguDiffDomain(mounts []crd.AdditionalMount) []ecosystem.AdditionalMount {
	if len(mounts) == 0 {
		// an empty slice and nil are serialized differently
		// we want no entry instead of an empty json list if there are no mounts given
		return nil
	}
	result := make([]ecosystem.AdditionalMount, len(mounts))

	for index, mount := range mounts {
		result[index] = ecosystem.AdditionalMount{
			SourceType: ecosystem.DataSourceType(mount.SourceType),
			Name:       mount.Name,
			Volume:     mount.Volume,
			Subfolder:  mount.Subfolder,
		}
	}

	return result
}
