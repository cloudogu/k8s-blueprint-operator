package v1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// StateDiff is the result of comparing the EffectiveBlueprint to the current cluster state.
// It describes what operations need to be done to achieve the desired state of the blueprint.
type StateDiff struct {
	// DoguDiffs maps simple dogu names to the determined diff.
	DoguDiffs map[string]DoguDiff `json:"doguDiffs,omitempty"`
	// DoguDiffs maps simple dogu names to the determined diff.
	ComponentDiffs map[string]ComponentDiff `json:"componentDiffs,omitempty"`
}

func ConvertToStateDiffDTO(domainModel domain.StateDiff) StateDiff {
	doguDiffs := make(map[string]DoguDiff, len(domainModel.DoguDiffs))
	for _, doguDiff := range domainModel.DoguDiffs {
		doguDiffs[doguDiff.DoguName] = convertToDoguDiffDTO(doguDiff)
	}

	componentDiffsV1 := make(map[string]ComponentDiff, len(domainModel.ComponentDiffs))
	for _, componentDiff := range domainModel.ComponentDiffs {
		componentDiffsV1[componentDiff.ComponentName] = convertToComponentDiffDTO(componentDiff)
	}

	return StateDiff{
		DoguDiffs:      doguDiffs,
		ComponentDiffs: componentDiffsV1,
		// in the future, this will also contain registry diffs
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

	return domain.StateDiff{
		DoguDiffs:      doguDiffs,
		ComponentDiffs: componentDiffs,
	}, nil
}
