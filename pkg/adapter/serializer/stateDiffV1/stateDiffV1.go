package stateDiffV1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// StateDiffV1 is the result of comparing the EffectiveBlueprint to the current cluster state.
// It describes what operations need to be done to achieve the desired state of the blueprint.
type StateDiffV1 struct {
	// DoguDiffs maps simple dogu names to the determined diff.
	DoguDiffs map[string]DoguDiffV1 `json:"doguDiffs,omitempty"`
	// DoguDiffs maps simple dogu names to the determined diff.
	ComponentDiffs map[string]DoguDiffV1 `json:"componentDiffs,omitempty"`
}

// DoguDiffV1 is the comparison of a Dogu's desired state vs. its cluster state.
// It contains the operation that needs to be done to achieve this desired state.
type DoguDiffV1 struct {
	Actual       DoguDiffV1State `json:"actual"`
	Expected     DoguDiffV1State `json:"expected"`
	NeededAction DoguActionV1    `json:"neededAction"`
}

// DoguDiffV1State is either the actual or desired state of a dogu in the cluster.
type DoguDiffV1State struct {
	Namespace         string `json:"namespace,omitempty"`
	Version           string `json:"version,omitempty"`
	InstallationState string `json:"installationState"`
}

// DoguActionV1 is the action that needs to be done for a dogu
// to achieve the desired state in the cluster.
type DoguActionV1 string

func ConvertToDTO(domainModel domain.StateDiff) StateDiffV1 {
	doguDiffsV1 := make(map[string]DoguDiffV1, len(domainModel.DoguDiffs))
	for _, doguDiff := range domainModel.DoguDiffs {
		doguDiffsV1[doguDiff.DoguName] = convertToDoguDiffDTO(doguDiff)
	}

	return StateDiffV1{
		DoguDiffs: doguDiffsV1,
		// in the future, this will also contain component and registry diffs
	}
}

func convertToDoguDiffDTO(domainModel domain.DoguDiff) DoguDiffV1 {
	return DoguDiffV1{
		Actual: DoguDiffV1State{
			Namespace:         domainModel.Actual.Namespace,
			Version:           domainModel.Actual.Version.Raw,
			InstallationState: domainModel.Actual.InstallationState.String(),
		},
		Expected: DoguDiffV1State{
			Namespace:         domainModel.Expected.Namespace,
			Version:           domainModel.Expected.Version.Raw,
			InstallationState: domainModel.Expected.InstallationState.String(),
		},
		NeededAction: DoguActionV1(domainModel.NeededAction),
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

	err := errors.Join(errs...)
	if err != nil {
		return domain.StateDiff{}, fmt.Errorf("failed to convert state diff dto to domain model: %w", err)
	}

	return domain.StateDiff{
		DoguDiffs: doguDiffs,
	}, nil
}

func convertToDoguDiffDomainModel(doguName string, dto DoguDiffV1) (domain.DoguDiff, error) {
	var actualVersion core.Version
	var actualVersionErr error
	if dto.Actual.Version != "" {
		actualVersion, actualVersionErr = core.ParseVersion(dto.Actual.Version)
		if actualVersionErr != nil {
			actualVersionErr = fmt.Errorf("failed to parse actual version %q: %w", dto.Actual.Version, actualVersionErr)
		}
	}

	var expectedVersion core.Version
	var expectedVersionErr error
	if dto.Expected.Version != "" {
		expectedVersion, expectedVersionErr = core.ParseVersion(dto.Expected.Version)
		if expectedVersionErr != nil {
			expectedVersionErr = fmt.Errorf("failed to parse expected version %q: %w", dto.Expected.Version, expectedVersionErr)
		}
	}

	actualState, actualStateErr := serializer.ToDomainTargetState(dto.Actual.InstallationState)
	if actualStateErr != nil {
		actualStateErr = fmt.Errorf("failed to parse actual installation state %q: %w", dto.Actual.InstallationState, actualStateErr)
	}

	expectedState, expectedStateErr := serializer.ToDomainTargetState(dto.Expected.InstallationState)
	if expectedStateErr != nil {
		expectedStateErr = fmt.Errorf("failed to parse expected installation state %q: %w", dto.Expected.InstallationState, expectedStateErr)
	}

	err := errors.Join(actualVersionErr, expectedVersionErr, actualStateErr, expectedStateErr)
	if err != nil {
		return domain.DoguDiff{}, fmt.Errorf("failed to convert dogu diff dto %q to domain model: %w", doguName, err)
	}

	return domain.DoguDiff{
		DoguName: doguName,
		Actual: domain.DoguDiffState{
			Namespace:         dto.Actual.Namespace,
			Version:           actualVersion,
			InstallationState: actualState,
		},
		Expected: domain.DoguDiffState{
			Namespace:         dto.Expected.Namespace,
			Version:           expectedVersion,
			InstallationState: expectedState,
		},
		NeededAction: domain.Action(dto.NeededAction),
	}, nil
}
