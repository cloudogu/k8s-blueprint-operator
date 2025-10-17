package serializer

import (
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func convertToDoguConfigEntryDiffsDTO(domainDiffs domain.DoguConfigDiffs, isSensitive bool) crd.DoguConfigDiff {
	var dtoDiffs []crd.ConfigEntryDiff
	for _, domainDiff := range domainDiffs {
		dtoDiffs = append(dtoDiffs, convertToDoguConfigEntryDiffDTO(domainDiff, isSensitive))
	}
	return dtoDiffs
}

func convertToDoguConfigEntryDiffDTO(domainModel domain.DoguConfigEntryDiff, isSensitive bool) crd.ConfigEntryDiff {
	actual := crd.ConfigValueState{
		Exists: domainModel.Actual.Exists,
	}
	expected := crd.ConfigValueState{
		Exists: domainModel.Expected.Exists,
	}
	if !isSensitive {
		actual.Value = domainModel.Actual.Value
		expected.Value = domainModel.Expected.Value
	}
	return crd.ConfigEntryDiff{
		Key:          string(domainModel.Key.Key),
		Actual:       actual,
		Expected:     expected,
		NeededAction: crd.ConfigAction(domainModel.NeededAction),
	}
}
