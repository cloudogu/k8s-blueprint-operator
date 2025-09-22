package serializer

import (
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func convertToGlobalConfigDiffDTO(domainModel domain.GlobalConfigDiffs) crd.GlobalConfigDiff {
	if len(domainModel) == 0 {
		return nil
	}

	globalConfigDiff := make(crd.GlobalConfigDiff, len(domainModel))
	for i, entryDiff := range domainModel {
		globalConfigDiff[i] = convertToGlobalConfigEntryDiffDTO(entryDiff)
	}
	return globalConfigDiff
}

func convertToGlobalConfigEntryDiffDTO(domainModel domain.GlobalConfigEntryDiff) crd.ConfigEntryDiff {
	return crd.ConfigEntryDiff{
		Key: string(domainModel.Key),
		Actual: crd.ConfigValueState{
			Value:  domainModel.Actual.Value,
			Exists: domainModel.Actual.Exists,
		},
		Expected: crd.ConfigValueState{
			Value:  domainModel.Expected.Value,
			Exists: domainModel.Expected.Exists,
		},
		NeededAction: crd.ConfigAction(domainModel.NeededAction),
	}
}
