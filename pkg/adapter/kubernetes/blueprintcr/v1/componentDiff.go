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
	// Actual contains the component's state in the current system.
	Actual ComponentDiffState `json:"actual"`
	// Expected contains the desired component's target state.
	Expected ComponentDiffState `json:"expected"`
	// NeededAction contains the refined action as decided by the application's state determination automaton.
	NeededAction ComponentAction `json:"neededAction"`
}

// ComponentDiffState is either the actual or desired state of a component in the cluster. The fields will be used to
// determine the kind of changed if there is a drift between actual or desired state.
type ComponentDiffState struct {
	// DistributionNamespace is part of the address under which the component will be obtained. This namespace must NOT
	// to be confused with the K8s cluster namespace.
	DistributionNamespace string `json:"distributionNamespace"`
	// Version contains the component's version.
	Version string `json:"version,omitempty"`
	// InstallationState contains the component's installation state. Such a state correlate with the domain Actions:
	//
	//  - domain.ActionInstall
	//  - domain.ActionUninstall
	//  - and so on
	InstallationState string `json:"installationState"`
}

// ComponentAction is the action that needs to be done for a component
// to achieve the desired state in the cluster.
type ComponentAction string

func convertToComponentDiffDTO(domainModel domain.ComponentDiff) ComponentDiff {
	return ComponentDiff{
		Actual: ComponentDiffState{
			DistributionNamespace: domainModel.Actual.DistributionNamespace,
			Version:               domainModel.Actual.Version.Raw,
			InstallationState:     domainModel.Actual.InstallationState.String(),
		},
		Expected: ComponentDiffState{
			DistributionNamespace: domainModel.Expected.DistributionNamespace,
			Version:               domainModel.Expected.Version.Raw,
			InstallationState:     domainModel.Expected.InstallationState.String(),
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

	actualDistributionNamespace := dto.Actual.DistributionNamespace
	expectedDistributionNamespace := dto.Expected.DistributionNamespace

	err := errors.Join(actualVersionErr, expectedVersionErr, actualStateErr, expectedStateErr)
	if err != nil {
		return domain.ComponentDiff{}, fmt.Errorf("failed to convert component diff dto %q to domain model: %w", componentName, err)
	}

	return domain.ComponentDiff{
		ComponentName: componentName,
		Actual: domain.ComponentDiffState{
			DistributionNamespace: actualDistributionNamespace,
			Version:               actualVersion,
			InstallationState:     actualState,
		},
		Expected: domain.ComponentDiffState{
			DistributionNamespace: expectedDistributionNamespace,
			Version:               expectedVersion,
			InstallationState:     expectedState,
		},
		NeededAction: domain.Action(dto.NeededAction),
	}, nil
}
