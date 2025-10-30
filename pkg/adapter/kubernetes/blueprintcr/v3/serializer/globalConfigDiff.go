package serializer

import (
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func convertToGlobalConfigDiffDTO(domainModel domain.GlobalConfigDiffs) bpv3.GlobalConfigDiff {
	if len(domainModel) == 0 {
		return nil
	}

	globalConfigDiff := make(bpv3.GlobalConfigDiff, len(domainModel))
	for i, entryDiff := range domainModel {
		globalConfigDiff[i] = convertToGlobalConfigEntryDiffDTO(entryDiff)
	}
	return globalConfigDiff
}

func convertToGlobalConfigEntryDiffDTO(domainModel domain.GlobalConfigEntryDiff) bpv3.ConfigEntryDiff {
	return bpv3.ConfigEntryDiff{
		Key: string(domainModel.Key),
		Actual: bpv3.ConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: bpv3.ConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: bpv3.ConfigAction(domainModel.NeededAction),
	}
}
