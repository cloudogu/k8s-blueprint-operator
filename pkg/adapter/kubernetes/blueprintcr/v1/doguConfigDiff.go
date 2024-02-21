package v1

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain"

type CombinedDoguConfigDiff struct {
	DoguConfigDiff          DoguConfigDiff          `json:"doguConfigDiff,omitempty"`
	SensitiveDoguConfigDiff SensitiveDoguConfigDiff `json:"sensitiveDoguConfigDiff,omitempty"`
}

type DoguConfigValueState ConfigValueState

type DoguConfigDiff []DoguConfigEntryDiff
type DoguConfigEntryDiff struct {
	Key          string               `json:"key,omitempty"`
	Actual       DoguConfigValueState `json:"actual,omitempty"`
	Expected     DoguConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction         `json:"neededAction,omitempty"`
}

type SensitiveDoguConfigDiff []SensitiveDoguConfigEntryDiff
type SensitiveDoguConfigEntryDiff struct {
	Key          string               `json:"key,omitempty"`
	Actual       DoguConfigValueState `json:"actual,omitempty"`
	Expected     DoguConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction         `json:"neededAction,omitempty"`
}

func convertToCombinedDoguConfigDiffDTO(domainModel domain.CombinedDoguConfigDiff) CombinedDoguConfigDiff {
	var doguConfigDiff DoguConfigDiff
	if len(domainModel.DoguConfigDiff) != 0 {
		doguConfigDiff = make(DoguConfigDiff, len(domainModel.DoguConfigDiff))
		for i, entryDiff := range domainModel.DoguConfigDiff {
			doguConfigDiff[i] = convertToDoguConfigEntryDiffDTO(entryDiff)
		}
	}

	var sensitiveDoguConfigDiff SensitiveDoguConfigDiff
	if len(domainModel.SensitiveDoguConfigDiff) != 0 {
		sensitiveDoguConfigDiff = make(SensitiveDoguConfigDiff, len(domainModel.SensitiveDoguConfigDiff))
		for i, entryDiff := range domainModel.SensitiveDoguConfigDiff {
			sensitiveDoguConfigDiff[i] = convertToSensitiveDoguConfigEntryDiffDTO(entryDiff)
		}
	}

	return CombinedDoguConfigDiff{
		DoguConfigDiff:          doguConfigDiff,
		SensitiveDoguConfigDiff: sensitiveDoguConfigDiff,
	}
}

func convertToDoguConfigEntryDiffDTO(domainModel domain.DoguConfigEntryDiff) DoguConfigEntryDiff {
	return DoguConfigEntryDiff{
		Key: domainModel.Key.Key,
		Actual: DoguConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: DoguConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: ConfigAction(domainModel.Action),
	}
}

func convertToSensitiveDoguConfigEntryDiffDTO(domainModel domain.SensitiveDoguConfigEntryDiff) SensitiveDoguConfigEntryDiff {
	return SensitiveDoguConfigEntryDiff{
		Key: domainModel.Key.Key,
		Actual: DoguConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: DoguConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: ConfigAction(domainModel.Action),
	}
}
