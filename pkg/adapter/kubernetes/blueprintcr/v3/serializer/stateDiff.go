package serializer

import (
	"slices"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToStateDiffDTO(domainModel domain.StateDiff) *bpv3.StateDiff {
	doguDiffs := make(map[string]bpv3.DoguDiff, len(domainModel.DoguDiffs))
	for _, doguDiff := range domainModel.DoguDiffs {
		doguDiffs[string(doguDiff.DoguName)] = convertToDoguDiffDTO(doguDiff)
	}

	var dogus []cescommons.SimpleName
	var combinedConfigDiffs map[string]bpv3.CombinedDoguConfigDiff
	var doguConfigDiffByDogu map[cescommons.SimpleName]bpv3.DoguConfigDiff
	var sensitiveDoguConfigDiff map[cescommons.SimpleName]bpv3.SensitiveDoguConfigDiff

	if len(domainModel.DoguConfigDiffs) != 0 || len(domainModel.SensitiveDoguConfigDiffs) != 0 {
		combinedConfigDiffs = make(map[string]bpv3.CombinedDoguConfigDiff)
		doguConfigDiffByDogu = make(map[cescommons.SimpleName]bpv3.DoguConfigDiff)
		for doguName, doguConfigDiff := range domainModel.DoguConfigDiffs {
			doguConfigDiffByDogu[doguName] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff, false)
			dogus = append(dogus, doguName)
		}
		sensitiveDoguConfigDiff = make(map[cescommons.SimpleName]bpv3.SensitiveDoguConfigDiff)
		for doguName, doguConfigDiff := range domainModel.SensitiveDoguConfigDiffs {
			sensitiveDoguConfigDiff[doguName] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff, true)
			dogus = append(dogus, doguName)
		}

		// remove duplicates, so we have a complete list of all dogus with config
		dogus = slices.Compact(dogus)
		for _, doguName := range dogus {
			combinedConfigDiffs[string(doguName)] = bpv3.CombinedDoguConfigDiff{
				DoguConfigDiff:          doguConfigDiffByDogu[doguName],
				SensitiveDoguConfigDiff: sensitiveDoguConfigDiff[doguName],
			}
		}
	}

	return &bpv3.StateDiff{
		DoguDiffs:        doguDiffs,
		DoguConfigDiffs:  combinedConfigDiffs,
		GlobalConfigDiff: convertToGlobalConfigDiffDTO(domainModel.GlobalConfigDiffs),
	}
}
