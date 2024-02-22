package v1

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

type GlobalConfigDiff []GlobalConfigEntryDiff

type GlobalConfigValueState ConfigValueState
type GlobalConfigEntryDiff struct {
	Key          string                 `json:"key,omitempty"`
	Actual       GlobalConfigValueState `json:"actual,omitempty"`
	Expected     GlobalConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction           `json:"neededAction,omitempty"`
}

func convertToGlobalConfigDiffDomain(dto GlobalConfigDiff) domain.GlobalConfigDiffs {
	if len(dto) == 0 {
		return nil
	}

	globalConfigDiff := make(domain.GlobalConfigDiffs, len(dto))
	for i, entryDiff := range dto {
		globalConfigDiff[i] = convertToGlobalConfigEntryDiffDomain(entryDiff)
	}
	return globalConfigDiff
}

func convertToGlobalConfigEntryDiffDomain(dto GlobalConfigEntryDiff) domain.GlobalConfigEntryDiff {
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

func convertToGlobalConfigDiffDTO(domainModel domain.GlobalConfigDiffs) GlobalConfigDiff {
	if len(domainModel) == 0 {
		return nil
	}

	globalConfigDiff := make(GlobalConfigDiff, len(domainModel))
	for i, entryDiff := range domainModel {
		globalConfigDiff[i] = convertToGlobalConfigEntryDiffDTO(entryDiff)
	}
	return globalConfigDiff
}

func convertToGlobalConfigEntryDiffDTO(domainModel domain.GlobalConfigEntryDiff) GlobalConfigEntryDiff {
	return GlobalConfigEntryDiff{
		Key: string(domainModel.Key),
		Actual: GlobalConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: GlobalConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: ConfigAction(domainModel.NeededAction),
	}
}
