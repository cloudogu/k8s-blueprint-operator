package serializer

import (
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

func convertToDoguDiffDTO(domainModel domain.DoguDiff) bpv3.DoguDiff {
	neededActions := domainModel.NeededActions
	doguActions := make([]bpv3.DoguAction, 0, len(neededActions))
	for _, action := range neededActions {
		doguActions = append(doguActions, bpv3.DoguAction(action))
	}

	return bpv3.DoguDiff{
		Actual:        convertToDoguDiffStateDTO(domainModel.Actual),
		Expected:      convertToDoguDiffStateDTO(domainModel.Expected),
		NeededActions: doguActions,
	}
}

func convertToDoguDiffStateDTO(domainModel domain.DoguDiffState) bpv3.DoguDiffState {
	var version *string
	if domainModel.Version != nil {
		version = &domainModel.Version.Raw
	}

	var reverseProxyConfig *bpv3.ReverseProxyConfig
	if !domainModel.ReverseProxyConfig.IsEmpty() {
		var rewriteTarget, additionalConfig *string
		if domainModel.ReverseProxyConfig.RewriteTarget != "" {
			rewriteTarget = (*string)(&domainModel.ReverseProxyConfig.RewriteTarget)
		}
		if domainModel.ReverseProxyConfig.AdditionalConfig != "" {
			additionalConfig = (*string)(&domainModel.ReverseProxyConfig.AdditionalConfig)
		}
		reverseProxyConfig = &bpv3.ReverseProxyConfig{
			RewriteTarget:    rewriteTarget,
			AdditionalConfig: additionalConfig,
			MaxBodySize:      ecosystem.GetQuantityString(domainModel.ReverseProxyConfig.MaxBodySize),
		}
	}

	var resourceConfig *bpv3.ResourceConfig
	if domainModel.MinVolumeSize != nil || domainModel.StorageClassName != nil {
		resourceConfig = &bpv3.ResourceConfig{
			MinVolumeSize:    ecosystem.GetQuantityString(domainModel.MinVolumeSize),
			StorageClassName: domainModel.StorageClassName,
		}
	}
	return bpv3.DoguDiffState{
		Namespace:          string(domainModel.Namespace),
		Version:            version,
		Absent:             domainModel.Absent,
		ResourceConfig:     resourceConfig,
		ReverseProxyConfig: reverseProxyConfig,
		AdditionalMounts:   convertAdditionalMountsToDoguDiffDTO(domainModel.AdditionalMounts),
	}
}

func convertAdditionalMountsToDoguDiffDTO(mounts []ecosystem.AdditionalMount) []bpv3.AdditionalMount {
	if len(mounts) == 0 {
		// an empty slice and nil are serialized differently
		// we want no entry instead of an empty json list if there are no mounts given
		return nil
	}
	result := make([]bpv3.AdditionalMount, len(mounts))
	for index, mount := range mounts {
		var subfolder *string
		if mount.Subfolder != "" {
			subfolder = &mount.Subfolder
		}
		result[index] = bpv3.AdditionalMount{
			SourceType: bpv3.DataSourceType(mount.SourceType),
			Name:       mount.Name,
			Volume:     mount.Volume,
			Subfolder:  subfolder,
		}
	}
	return result
}
