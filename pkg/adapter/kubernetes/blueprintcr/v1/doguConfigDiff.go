package v1

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type CombinedDoguConfigDiff struct {
	DoguConfigDiff          DoguConfigDiff          `json:"doguConfigDiff,omitempty"`
	SensitiveDoguConfigDiff SensitiveDoguConfigDiff `json:"sensitiveDoguConfigDiff,omitempty"`
}

type DoguConfigValueState ConfigValueState

type DoguConfigDiff []DoguConfigEntryDiff
type DoguConfigEntryDiff struct {
	Key          string               `json:"key"`
	Actual       DoguConfigValueState `json:"actual"`
	Expected     DoguConfigValueState `json:"expected"`
	NeededAction ConfigAction         `json:"neededAction"`
}

type SensitiveDoguConfigDiff []SensitiveDoguConfigEntryDiff
type SensitiveDoguConfigEntryDiff struct {
	Key          string               `json:"key"`
	Actual       DoguConfigValueState `json:"actual"`
	Expected     DoguConfigValueState `json:"expected"`
	NeededAction ConfigAction         `json:"neededAction"`
}

func convertToCombinedDoguConfigDiffDomain(doguName string, dto CombinedDoguConfigDiff) domain.CombinedDoguConfigDiffs {
	var doguConfigDiff domain.DoguConfigDiffs
	if len(dto.DoguConfigDiff) != 0 {
		doguConfigDiff = make(domain.DoguConfigDiffs, len(dto.DoguConfigDiff))
		for i, entryDiff := range dto.DoguConfigDiff {
			doguConfigDiff[i] = convertToDoguConfigEntryDiffDomain(doguName, entryDiff)
		}
	}

	var sensitiveDoguConfigDiff domain.SensitiveDoguConfigDiffs
	if len(dto.SensitiveDoguConfigDiff) != 0 {
		sensitiveDoguConfigDiff = make(domain.SensitiveDoguConfigDiffs, len(dto.SensitiveDoguConfigDiff))
		for i, entryDiff := range dto.SensitiveDoguConfigDiff {
			sensitiveDoguConfigDiff[i] = convertToSensitiveDoguConfigEntryDiffDomain(doguName, entryDiff)
		}
	}

	return domain.CombinedDoguConfigDiffs{
		DoguConfigDiff:          doguConfigDiff,
		SensitiveDoguConfigDiff: sensitiveDoguConfigDiff,
	}
}

func convertToDoguConfigEntryDiffDomain(doguName string, dto DoguConfigEntryDiff) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			DoguName: common.SimpleDoguName(doguName),
			Key:      config.Key(dto.Key),
		},
		Actual: domain.DoguConfigValueState{
			Value:  dto.Actual.Value,
			Exists: dto.Actual.Exists,
		},
		Expected: domain.DoguConfigValueState{
			Value:  dto.Expected.Value,
			Exists: dto.Expected.Exists,
		},
		NeededAction: domain.ConfigAction(dto.NeededAction),
	}
}

func convertToSensitiveDoguConfigEntryDiffDomain(doguName string, dto SensitiveDoguConfigEntryDiff) domain.SensitiveDoguConfigEntryDiff {
	return domain.SensitiveDoguConfigEntryDiff{
		Key: common.SensitiveDoguConfigKey{
			DoguName: common.SimpleDoguName(doguName),
			Key:      config.Key(dto.Key),
		},
		Actual: domain.DoguConfigValueState{
			Value:  dto.Actual.Value,
			Exists: dto.Actual.Exists,
		},
		Expected: domain.DoguConfigValueState{
			Value:  dto.Expected.Value,
			Exists: dto.Expected.Exists,
		},
		NeededAction: domain.ConfigAction(dto.NeededAction),
	}
}

func convertToCombinedDoguConfigDiffDTO(domainModel domain.CombinedDoguConfigDiffs) CombinedDoguConfigDiff {
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
		Key: string(domainModel.Key.Key),
		Actual: DoguConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: DoguConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: ConfigAction(domainModel.NeededAction),
	}
}

func convertToSensitiveDoguConfigEntryDiffDTO(domainModel domain.SensitiveDoguConfigEntryDiff) SensitiveDoguConfigEntryDiff {
	return SensitiveDoguConfigEntryDiff{
		Key: string(domainModel.Key.Key),
		Actual: DoguConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: DoguConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: ConfigAction(domainModel.NeededAction),
	}
}
