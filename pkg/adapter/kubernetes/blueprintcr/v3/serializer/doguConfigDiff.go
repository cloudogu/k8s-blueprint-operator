package serializer

import (
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func convertToDoguConfigEntryDiffsDTO(domainDiffs domain.DoguConfigDiffs, isSensitive bool) bpv3.DoguConfigDiff {
	var dtoDiffs []bpv3.ConfigEntryDiff
	for _, domainDiff := range domainDiffs {
		dtoDiffs = append(dtoDiffs, convertToDoguConfigEntryDiffDTO(domainDiff, isSensitive))
	}
	return dtoDiffs
}

func convertToDoguConfigEntryDiffDTO(domainModel domain.DoguConfigEntryDiff, isSensitive bool) bpv3.ConfigEntryDiff {
	actual := bpv3.ConfigValueState{
		Exists: domainModel.Actual.Exists,
	}
	expected := bpv3.ConfigValueState{
		Exists: domainModel.Expected.Exists,
	}
	if !isSensitive {
		actual.Value = domainModel.Actual.Value
		expected.Value = domainModel.Expected.Value
	}
	return bpv3.ConfigEntryDiff{
		Key:          string(domainModel.Key.Key),
		Actual:       actual,
		Expected:     expected,
		NeededAction: bpv3.ConfigAction(domainModel.NeededAction),
	}
}
