package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
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

func convertToDoguConfigEntryDiffsDTO(domainDiffs domain.DoguConfigDiffs, isSensitive bool) []crd.DoguConfigEntryDiff {
	var dtoDiffs []crd.DoguConfigEntryDiff
	for _, domainDiff := range domainDiffs {
		dtoDiffs = append(dtoDiffs, convertToDoguConfigEntryDiffDTO(domainDiff, isSensitive))
	}
	return dtoDiffs
}

func convertToDoguConfigEntryDiffDTO(domainModel domain.DoguConfigEntryDiff, isSensitive bool) crd.DoguConfigEntryDiff {
	actual := crd.DoguConfigValueState{
		Exists: domainModel.Actual.Exists,
	}
	expected := crd.DoguConfigValueState{
		Exists: domainModel.Expected.Exists,
	}
	if !isSensitive {
		actual.Value = domainModel.Actual.Value
		expected.Value = domainModel.Expected.Value
	}
	return crd.DoguConfigEntryDiff{
		Key:          string(domainModel.Key.Key),
		Actual:       actual,
		Expected:     expected,
		NeededAction: crd.ConfigAction(domainModel.NeededAction),
	}
}
