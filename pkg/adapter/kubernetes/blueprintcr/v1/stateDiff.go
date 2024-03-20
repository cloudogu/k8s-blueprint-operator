package v1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
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

	var doguConfigDiffs map[string]CombinedDoguConfigDiff
	if len(domainModel.DoguConfigDiffs) != 0 {
		doguConfigDiffs = make(map[string]CombinedDoguConfigDiff)
		for doguName, doguConfigDiff := range domainModel.DoguConfigDiffs {
			doguConfigDiffs[string(doguName)] = convertToCombinedDoguConfigDiffDTO(doguConfigDiff)
		}
	}

	return StateDiff{
		DoguDiffs:        doguDiffs,
		ComponentDiffs:   componentDiffsV1,
		DoguConfigDiffs:  doguConfigDiffs,
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

	var doguConfigDiffs map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs
	if len(dto.DoguConfigDiffs) != 0 {
		doguConfigDiffs = make(map[common.SimpleDoguName]domain.CombinedDoguConfigDiffs)
		for doguName, doguConfigDiff := range dto.DoguConfigDiffs {
			doguConfigDiffs[common.SimpleDoguName(doguName)] = convertToCombinedDoguConfigDiffDomain(doguName, doguConfigDiff)
		}
	}

	return domain.StateDiff{
		DoguDiffs:         doguDiffs,
		ComponentDiffs:    componentDiffs,
		DoguConfigDiffs:   doguConfigDiffs,
		GlobalConfigDiffs: convertToGlobalConfigDiffDomain(dto.GlobalConfigDiff),
	}, nil
}
