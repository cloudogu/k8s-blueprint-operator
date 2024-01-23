package stateDiffV1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// StateDiffV1 is the result of comparing the EffectiveBlueprint to the current cluster state.
// It describes what operations need to be done to achieve the desired state of the blueprint.
type StateDiffV1 struct {
	// DoguDiffs maps simple dogu names to the determined diff.
	DoguDiffs map[string]DoguDiffV1 `json:"doguDiffs,omitempty"`
	// DoguDiffs maps simple dogu names to the determined diff.
	ComponentDiffs map[string]ComponentDiffV1 `json:"componentDiffs,omitempty"`
}

func ConvertToDTO(domainModel domain.StateDiff) StateDiffV1 {
	doguDiffsV1 := make(map[string]DoguDiffV1, len(domainModel.DoguDiffs))
	for _, doguDiff := range domainModel.DoguDiffs {
		doguDiffsV1[doguDiff.DoguName] = convertToDoguDiffDTO(doguDiff)
	}

	componentDiffsV1 := make(map[string]ComponentDiffV1, len(domainModel.ComponentDiffs))
	for _, componentDiff := range domainModel.ComponentDiffs {
		componentDiffsV1[componentDiff.Name] = convertToComponentDiffDTO(componentDiff)
	}

	return StateDiffV1{
		DoguDiffs:      doguDiffsV1,
		ComponentDiffs: componentDiffsV1,
		// in the future, this will also contain registry diffs
	}
}

func ConvertToDomainModel(dto StateDiffV1) (domain.StateDiff, error) {
	var errs []error

	doguDiffs := make([]domain.DoguDiff, 0)
	for doguName, doguDiff := range dto.DoguDiffs {
		doguDiffDomainModel, err := convertToDoguDiffDomainModel(doguName, doguDiff)
		errs = append(errs, err)
		doguDiffs = append(doguDiffs, doguDiffDomainModel)
	}

	componentDiffs := make([]domain.ComponentDiff, 0)
	for componentName, componentDiff := range dto.ComponentDiffs {
		componentDiffDomainModel, err := convertToComponentDiffDomainModel(componentName, componentDiff)
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
