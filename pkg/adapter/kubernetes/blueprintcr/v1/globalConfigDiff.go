package v1

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain"

type GlobalConfigDiff []GlobalConfigEntryDiff

type GlobalConfigValueState ConfigValueState
type GlobalConfigEntryDiff struct {
	Key          string                 `json:"key,omitempty"`
	Actual       GlobalConfigValueState `json:"actual,omitempty"`
	Expected     GlobalConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction           `json:"neededAction,omitempty"`
}

func convertToGlobalConfigDiffDTO(domainModel domain.GlobalConfigDiff) GlobalConfigDiff {
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
		NeededAction: ConfigAction(domainModel.Action),
	}
}
