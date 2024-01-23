package stateDiffV1

import (
	"errors"
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// ComponentDiffV1 is the comparison of a Component's desired state vs. its cluster state.
// It contains the operation that needs to be done to achieve this desired state.
type ComponentDiffV1 struct {
	Actual       ComponentDiffV1State `json:"actual"`
	Expected     ComponentDiffV1State `json:"expected"`
	NeededAction ComponentActionV1    `json:"neededAction"`
}

// ComponentDiffV1State is either the actual or desired state of a component in the cluster.
type ComponentDiffV1State struct {
	Version           string `json:"version,omitempty"`
	InstallationState string `json:"installationState"`
}

// ComponentActionV1 is the action that needs to be done for a component
// to achieve the desired state in the cluster.
type ComponentActionV1 string

func convertToComponentDiffDTO(domainModel domain.ComponentDiff) ComponentDiffV1 {
	return ComponentDiffV1{
		Actual: ComponentDiffV1State{
			Version:           domainModel.Actual.Version.Raw,
			InstallationState: domainModel.Actual.InstallationState.String(),
		},
		Expected: ComponentDiffV1State{
			Version:           domainModel.Expected.Version.Raw,
			InstallationState: domainModel.Expected.InstallationState.String(),
		},
		NeededAction: ComponentActionV1(domainModel.NeededAction),
	}
}

func convertToComponentDiffDomainModel(componentName string, dto ComponentDiffV1) (domain.ComponentDiff, error) {
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
		return domain.ComponentDiff{}, fmt.Errorf("failed to convert component diff dto %q to domain model: %w", componentName, err)
	}

	return domain.ComponentDiff{
		Name: componentName,
		Actual: domain.ComponentDiffState{
			Version:           actualVersion,
			InstallationState: actualState,
		},
		Expected: domain.ComponentDiffState{
			Version:           expectedVersion,
			InstallationState: expectedState,
		},
		NeededAction: domain.Action(dto.NeededAction),
	}, nil
}
