package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	crd "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-registry-lib/config"
)

func convertToDoguConfigDiffsDomain(doguName string, dtoDiffs crd.DoguConfigDiff) domain.DoguConfigDiffs {
	var doguConfigDiff domain.DoguConfigDiffs
	for _, entryDiff := range dtoDiffs {
		doguConfigDiff = append(doguConfigDiff, convertToDoguConfigEntryDiffDomain(doguName, entryDiff))
	}
	return doguConfigDiff
}

func convertToDoguConfigEntryDiffDomain(doguName string, dto crd.DoguConfigEntryDiff) domain.DoguConfigEntryDiff {
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

func convertToDoguConfigEntryDiffsDTO(domainDiffs domain.DoguConfigDiffs) []crd.DoguConfigEntryDiff {
	var dtoDiffs []crd.DoguConfigEntryDiff
	for _, domainDiff := range domainDiffs {
		dtoDiffs = append(dtoDiffs, convertToDoguConfigEntryDiffDTO(domainDiff))
	}
	return dtoDiffs
}

func convertToDoguConfigEntryDiffDTO(domainModel domain.DoguConfigEntryDiff) crd.DoguConfigEntryDiff {
	return crd.DoguConfigEntryDiff{
		Key: string(domainModel.Key.Key),
		Actual: crd.DoguConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: crd.DoguConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: crd.ConfigAction(domainModel.NeededAction),
	}
}
