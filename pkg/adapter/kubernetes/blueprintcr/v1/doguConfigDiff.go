package v1

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
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

// +kubebuilder:object:generate:=false
type SensitiveDoguConfigDiff = DoguConfigDiff

// +kubebuilder:object:generate:=false
type SensitiveDoguConfigEntryDiff = DoguConfigEntryDiff

func convertToDoguConfigDiffsDomain(doguName string, dtoDiffs DoguConfigDiff) domain.DoguConfigDiffs {
	var doguConfigDiff domain.DoguConfigDiffs
	for _, entryDiff := range dtoDiffs {
		doguConfigDiff = append(doguConfigDiff, convertToDoguConfigEntryDiffDomain(doguName, entryDiff))
	}
	return doguConfigDiff
}

func convertToDoguConfigEntryDiffDomain(doguName string, dto DoguConfigEntryDiff) domain.DoguConfigEntryDiff {
	return domain.DoguConfigEntryDiff{
		Key: common.DoguConfigKey{
			DoguName: cescommons.SimpleName(doguName),
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

func convertToDoguConfigEntryDiffsDTO(domainDiffs domain.DoguConfigDiffs) []DoguConfigEntryDiff {
	var dtoDiffs []DoguConfigEntryDiff
	for _, domainDiff := range domainDiffs {
		dtoDiffs = append(dtoDiffs, convertToDoguConfigEntryDiffDTO(domainDiff))
	}
	return dtoDiffs
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
