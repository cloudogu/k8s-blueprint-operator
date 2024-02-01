package v1

import (
	"errors"
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// ComponentDiff is the comparison of a Component's desired state vs. its cluster state.
// It contains the operation that needs to be done to achieve this desired state.
type ComponentDiff struct {
	Actual       ComponentDiffState `json:"actual"`
	Expected     ComponentDiffState `json:"expected"`
	NeededAction ComponentAction    `json:"neededAction"`
}

// ComponentDiffState is either the actual or desired state of a component in the cluster.
type ComponentDiffState struct {
	Version           string `json:"version,omitempty"`
	InstallationState string `json:"installationState"`
}

// ComponentAction is the action that needs to be done for a component
// to achieve the desired state in the cluster.
type ComponentAction string

func convertToComponentDiffDTO(domainModel domain.ComponentDiff) ComponentDiff {
	return ComponentDiff{
		Actual: ComponentDiffState{
			Version:           domainModel.Actual.Version.Raw,
			InstallationState: domainModel.Actual.InstallationState.String(),
		},
		Expected: ComponentDiffState{
			Version:           domainModel.Expected.Version.Raw,
			InstallationState: domainModel.Expected.InstallationState.String(),
		},
		NeededAction: ComponentAction(domainModel.NeededAction),
	}
}

func convertToComponentDiffDomain(componentName string, dto ComponentDiff) (domain.ComponentDiff, error) {
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
		ComponentName: componentName,
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
