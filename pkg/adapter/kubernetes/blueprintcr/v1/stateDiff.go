package v1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"slices"
)

// StateDiff is the result of comparing the EffectiveBlueprint to the current cluster state.
// It describes what operations need to be done to achieve the desired state of the blueprint.
type StateDiff struct {
	// DoguDiffs maps simple dogu names to the determined diff.
	DoguDiffs map[string]DoguDiff `json:"doguDiffs,omitempty"`
	// ComponentDiffs maps simple component names to the determined diff.
	ComponentDiffs map[string]ComponentDiff `json:"componentDiffs,omitempty"`
	// DoguConfigDiffs maps simple dogu names to the determined config diff.
	DoguConfigDiffs map[string]CombinedDoguConfigDiff `json:"doguConfigDiffs,omitempty"`
	// GlobalConfigDiff is the difference between the GlobalConfig in the EffectiveBlueprint and the cluster state.
	GlobalConfigDiff GlobalConfigDiff `json:"globalConfigDiff,omitempty"`
}

func ConvertToStateDiffDTO(domainModel domain.StateDiff) StateDiff {
	doguDiffs := make(map[string]DoguDiff, len(domainModel.DoguDiffs))
	for _, doguDiff := range domainModel.DoguDiffs {
		doguDiffs[string(doguDiff.DoguName)] = convertToDoguDiffDTO(doguDiff)
	}

	componentDiffsV1 := make(map[string]ComponentDiff, len(domainModel.ComponentDiffs))
	for _, componentDiff := range domainModel.ComponentDiffs {
		componentDiffsV1[string(componentDiff.Name)] = convertToComponentDiffDTO(componentDiff)
	}

	var dogus []common.SimpleDoguName
	var combinedConfigDiffs map[string]CombinedDoguConfigDiff
	var doguConfigDiffByDogu map[common.SimpleDoguName]DoguConfigDiff
	var sensitiveDoguConfigDiff map[common.SimpleDoguName]SensitiveDoguConfigDiff

	if len(domainModel.DoguConfigDiffs) != 0 || len(domainModel.SensitiveDoguConfigDiffs) != 0 {
		combinedConfigDiffs = make(map[string]CombinedDoguConfigDiff)
		for doguName, doguConfigDiff := range domainModel.DoguConfigDiffs {
			doguConfigDiffByDogu = make(map[common.SimpleDoguName]DoguConfigDiff)
			doguConfigDiffByDogu[doguName] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff)
			dogus = append(dogus, doguName)
		}
		for doguName, doguConfigDiff := range domainModel.SensitiveDoguConfigDiffs {
			sensitiveDoguConfigDiff = make(map[common.SimpleDoguName]SensitiveDoguConfigDiff)
			sensitiveDoguConfigDiff[doguName] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff)
			dogus = append(dogus, doguName)
		}

		//remove duplicates, so we have a complete list of all dogus with config
		dogus = slices.Compact(dogus)
		for _, doguName := range dogus {
			combinedConfigDiffs[string(doguName)] = CombinedDoguConfigDiff{
				DoguConfigDiff:          doguConfigDiffByDogu[doguName],
				SensitiveDoguConfigDiff: sensitiveDoguConfigDiff[doguName],
			}
		}
	}

	return StateDiff{
		DoguDiffs:        doguDiffs,
		ComponentDiffs:   componentDiffsV1,
		DoguConfigDiffs:  combinedConfigDiffs,
		GlobalConfigDiff: convertToGlobalConfigDiffDTO(domainModel.GlobalConfigDiffs),
	}
}

func ConvertToStateDiffDomain(dto StateDiff) (domain.StateDiff, error) {
	var errs []error

	doguDiffs := make([]domain.DoguDiff, 0)
	for doguName, doguDiff := range dto.DoguDiffs {
		doguDiffDomainModel, err := convertToDoguDiffDomain(doguName, doguDiff)
		errs = append(errs, err)
		doguDiffs = append(doguDiffs, doguDiffDomainModel)
	}

	componentDiffs := make([]domain.ComponentDiff, 0)
	for componentName, componentDiff := range dto.ComponentDiffs {
		componentDiffDomainModel, err := convertToComponentDiffDomain(componentName, componentDiff)
		errs = append(errs, err)
		componentDiffs = append(componentDiffs, componentDiffDomainModel)
	}

	err := errors.Join(errs...)
	if err != nil {
		return domain.StateDiff{}, fmt.Errorf("failed to convert state diff DTO to domain model: %w", err)
	}

	var doguConfigDiffs map[common.SimpleDoguName]domain.DoguConfigDiffs
	var sensitiveDoguConfigDiffs map[common.SimpleDoguName]domain.SensitiveDoguConfigDiffs
	if len(dto.DoguConfigDiffs) != 0 {
		doguConfigDiffs = map[common.SimpleDoguName]domain.DoguConfigDiffs{}
		sensitiveDoguConfigDiffs = map[common.SimpleDoguName]domain.SensitiveDoguConfigDiffs{}
		for doguName, combinedConfigDiff := range dto.DoguConfigDiffs {
			doguConfigDiffs[common.SimpleDoguName(doguName)] = convertToDoguConfigDiffsDomain(doguName, combinedConfigDiff.DoguConfigDiff)
			sensitiveDoguConfigDiffs[common.SimpleDoguName(doguName)] = convertToDoguConfigDiffsDomain(doguName, combinedConfigDiff.SensitiveDoguConfigDiff)
		}
	}

	return domain.StateDiff{
		DoguDiffs:                doguDiffs,
		ComponentDiffs:           componentDiffs,
		DoguConfigDiffs:          doguConfigDiffs,
		SensitiveDoguConfigDiffs: sensitiveDoguConfigDiffs,
		GlobalConfigDiffs:        convertToGlobalConfigDiffDomain(dto.GlobalConfigDiff),
	}, nil
}
