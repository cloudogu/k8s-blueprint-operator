package serializer

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	crd "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

func convertToDoguDiffDTO(domainModel domain.DoguDiff) crd.DoguDiff {
	neededActions := domainModel.NeededActions
	doguActions := make([]crd.DoguAction, 0, len(neededActions))
	for _, action := range neededActions {
		doguActions = append(doguActions, crd.DoguAction(action))
	}

	return crd.DoguDiff{
		Actual: crd.DoguDiffState{
			Namespace:         string(domainModel.Actual.Namespace),
			Version:           domainModel.Actual.Version.Raw,
			InstallationState: domainModel.Actual.InstallationState.String(),
			ResourceConfig: crd.ResourceConfig{
				MinVolumeSize: convertMinimumVolumeSizeToDTO(domainModel.Actual.MinVolumeSize),
			},
			ReverseProxyConfig: crd.ReverseProxyConfig{
				MaxBodySize:      ecosystem.GetQuantityString(domainModel.Actual.ReverseProxyConfig.MaxBodySize),
				RewriteTarget:    string(domainModel.Actual.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: string(domainModel.Actual.ReverseProxyConfig.AdditionalConfig),
			},
			AdditionalMounts: convertAdditionalMountsToDoguDiffDTO(domainModel.Actual.AdditionalMounts),
		},
		Expected: crd.DoguDiffState{
			Namespace:         string(domainModel.Expected.Namespace),
			Version:           domainModel.Expected.Version.Raw,
			InstallationState: domainModel.Expected.InstallationState.String(),
			ResourceConfig: crd.ResourceConfig{
				MinVolumeSize: convertMinimumVolumeSizeToDTO(domainModel.Expected.MinVolumeSize),
			},
			ReverseProxyConfig: crd.ReverseProxyConfig{
				MaxBodySize:      ecosystem.GetQuantityString(domainModel.Expected.ReverseProxyConfig.MaxBodySize),
				RewriteTarget:    string(domainModel.Expected.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: string(domainModel.Expected.ReverseProxyConfig.AdditionalConfig),
			},
			AdditionalMounts: convertAdditionalMountsToDoguDiffDTO(domainModel.Expected.AdditionalMounts),
		},
		NeededActions: doguActions,
	}
}

func convertMinimumVolumeSizeToDTO(minVolSize ecosystem.VolumeSize) string {
	if minVolSize.IsZero() {
		return ""
	} else {
		return minVolSize.String()
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
	if actualStateErr != nil {
		actualStateErr = fmt.Errorf("failed to convert expected dogu diff state: %w", actualStateErr)
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

	var version core.Version
	var versionErr error
	if dto.Version != "" {
		version, versionErr = core.ParseVersion(dto.Version)
		if versionErr != nil {
			versionErr = fmt.Errorf("failed to parse version %q: %w", dto.Version, versionErr)
		}
	}
	errorList = append(errorList, versionErr)

	state, stateErr := serializer.ToDomainTargetState(dto.InstallationState)
	if stateErr != nil {
		errorList = append(errorList, fmt.Errorf("failed to parse installation state %q: %w", dto.InstallationState, stateErr))
	}

	minVolumeSize, volumeSizeErr := ecosystem.GetNonNilQuantityRef(dto.ResourceConfig.MinVolumeSize)
	if volumeSizeErr != nil {
		errorList = append(errorList, fmt.Errorf("failed to parse minimum volume size %q: %w", dto.ResourceConfig.MinVolumeSize, volumeSizeErr))
	}

	maxBodySize, bodySizeErr := ecosystem.GetQuantityReference(dto.ReverseProxyConfig.MaxBodySize)
	if bodySizeErr != nil {
		errorList = append(errorList, fmt.Errorf("failed to parse maximum proxy body size %q: %w", dto.ReverseProxyConfig.MaxBodySize, bodySizeErr))
	}

	return domain.DoguDiffState{
		Namespace:         cescommons.Namespace(dto.Namespace),
		Version:           version,
		InstallationState: state,
		MinVolumeSize:     *minVolumeSize,
		ReverseProxyConfig: ecosystem.ReverseProxyConfig{
			MaxBodySize:      maxBodySize,
			RewriteTarget:    ecosystem.RewriteTarget(dto.ReverseProxyConfig.RewriteTarget),
			AdditionalConfig: ecosystem.AdditionalConfig(dto.ReverseProxyConfig.AdditionalConfig),
		},
		AdditionalMounts: convertAdditionalMountsToDoguDiffDomain(dto.AdditionalMounts),
	}, errors.Join(errorList...)
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
