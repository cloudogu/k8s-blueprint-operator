package serializer

import (
	crd "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

func convertToGlobalConfigDiffDomain(dto crd.GlobalConfigDiff) domain.GlobalConfigDiffs {
	if len(dto) == 0 {
		return nil
	}

	globalConfigDiff := make(domain.GlobalConfigDiffs, len(dto))
	for i, entryDiff := range dto {
		globalConfigDiff[i] = convertToGlobalConfigEntryDiffDomain(entryDiff)
	}
	return globalConfigDiff
}

func convertToGlobalConfigEntryDiffDomain(dto crd.GlobalConfigEntryDiff) domain.GlobalConfigEntryDiff {
	return domain.GlobalConfigEntryDiff{
		Key: common.GlobalConfigKey(dto.Key),
		Actual: domain.GlobalConfigValueState{
			Value:  dto.Actual.Value,
			Exists: dto.Actual.Exists,
		},
		Expected: domain.GlobalConfigValueState{
			Value:  dto.Expected.Value,
			Exists: dto.Expected.Exists,
		},
		NeededAction: domain.ConfigAction(dto.NeededAction),
	}
}

func convertToGlobalConfigDiffDTO(domainModel domain.GlobalConfigDiffs) crd.GlobalConfigDiff {
	if len(domainModel) == 0 {
		return nil
	}

	globalConfigDiff := make(crd.GlobalConfigDiff, len(domainModel))
	for i, entryDiff := range domainModel {
		globalConfigDiff[i] = convertToGlobalConfigEntryDiffDTO(entryDiff)
	}
	return globalConfigDiff
}

func convertToGlobalConfigEntryDiffDTO(domainModel domain.GlobalConfigEntryDiff) crd.GlobalConfigEntryDiff {
	return crd.GlobalConfigEntryDiff{
		Key: string(domainModel.Key),
		Actual: crd.GlobalConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: crd.GlobalConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: crd.ConfigAction(domainModel.NeededAction),
	}
}
