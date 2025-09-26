package serializer

import (
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
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
