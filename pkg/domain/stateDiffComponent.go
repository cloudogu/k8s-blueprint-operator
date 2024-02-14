package domain

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"golang.org/x/exp/maps"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

// ComponentDiffs contains the differences for all expected Components to the current ecosystem.ComponentInstallations.
type ComponentDiffs []ComponentDiff

// Statistics aggregates various figures about the required actions of the ComponentDiffs.
func (cd ComponentDiffs) Statistics() (toInstall int, toUpgrade int, toUninstall int, other int) {
	for _, componentDiff := range cd {
		switch componentDiff.NeededAction {
		case ActionInstall:
			toInstall += 1
		case ActionUpgrade:
			toUpgrade += 1
		case ActionUninstall:
			toUninstall += 1
		default:
			other += 1
		}
	}
	return
}

// ComponentDiff represents the Diff for a single expected Component to the current ecosystem.ComponentInstallation.
type ComponentDiff struct {
	// Name contains the component's name.
	Name string
	// Actual contains that state of a component how it is currently found in the system.
	Actual ComponentDiffState
	// Expected contains that desired state of a component how it is supposed to be.
	Expected ComponentDiffState
	// NeededAction hints how the component should be handled by the application change automaton in order to reconcile
	// differences between Actual and Expected in the current system.
	NeededAction Action
}

// ComponentDiffState contains all fields to make a diff for components in respect to another ComponentDiffState.
type ComponentDiffState struct {
	// DistributionNamespace is part of the address under which the component will be obtained. This namespace must NOT
	// to be confused with the K8s cluster namespace.
	DistributionNamespace string
	// Version contains the component's version.
	Version *semver.Version
	// InstallationState contains the component's target state.
	InstallationState TargetState
}

// String returns a string representation of the ComponentDiff.
func (diff *ComponentDiff) String() string {
	return fmt.Sprintf(
		"{Name: %q, Actual: %s, Expected: %s, NeededAction: %q}",
		diff.Name,
		diff.Actual.String(),
		diff.Expected.String(),
		diff.NeededAction,
	)
}

// String returns a string representation of the ComponentDiffState.
func (diff *ComponentDiffState) String() string {
	return fmt.Sprintf(
		"{DistributionNamespace: %q, Version: %q, InstallationState: %q}",
		diff.DistributionNamespace,
		diff.getSafeVersionString(),
		diff.InstallationState,
	)
}

func (diff *ComponentDiffState) getSafeVersionString() string {
	if diff.Version != nil {
		return diff.Version.String()
	} else {
		return ""
	}
}

// determineComponentDiffs creates ComponentDiffs for all components in the blueprint and all installed components as well.
func determineComponentDiffs(blueprintComponents []Component, installedComponents map[string]*ecosystem.ComponentInstallation) ([]ComponentDiff, error) {
	var componentDiffs = map[string]ComponentDiff{}
	for _, blueprintComponent := range blueprintComponents {
		installedComponent := installedComponents[blueprintComponent.Name]
		compDiff, err := determineComponentDiff(&blueprintComponent, installedComponent)
		if err != nil {
			return nil, err
		}
		componentDiffs[blueprintComponent.Name] = compDiff
	}

	for _, installedComponent := range installedComponents {
		blueprintComponent, found := findComponentByName(blueprintComponents, installedComponent.Name)

		if !found {
			var notFoundInBlueprint *Component = nil
			compDiff, err := determineComponentDiff(notFoundInBlueprint, installedComponent)
			if err != nil {
				return nil, err
			}
			componentDiffs[installedComponent.Name] = compDiff
			continue
		}

		compDiff, err := determineComponentDiff(&blueprintComponent, installedComponent)
		if err != nil {
			return nil, err
		}
		componentDiffs[installedComponent.Name] = compDiff
	}
	return maps.Values(componentDiffs), nil
}

// determineComponentDiff creates a ComponentDiff out of a Component from the blueprint and the ecosystem.ComponentInstallation in the ecosystem.
// If the Component is nil (was not in the blueprint), the actual state is also the expected state.
// If the installedComponent is nil, it is considered to be not installed currently.
func determineComponentDiff(blueprintComponent *Component, installedComponent *ecosystem.ComponentInstallation) (ComponentDiff, error) {
	var expectedState, actualState ComponentDiffState
	componentName := "" // either blueprintComponent or installedComponent could be nil

	if installedComponent == nil {
		actualState = ComponentDiffState{
			InstallationState: TargetStateAbsent,
		}
	} else {
		componentName = installedComponent.Name
		actualState = ComponentDiffState{
			DistributionNamespace: installedComponent.DistributionNamespace,
			Version:               installedComponent.Version,
			InstallationState:     TargetStatePresent,
		}
	}

	if blueprintComponent == nil {
		expectedState = actualState
	} else {
		componentName = blueprintComponent.Name
		expectedState = ComponentDiffState{
			DistributionNamespace: blueprintComponent.DistributionNamespace,
			Version:               blueprintComponent.Version,
			InstallationState:     blueprintComponent.TargetState,
		}
	}

	nextAction, err := getNextComponentAction(expectedState, actualState)
	if err != nil {
		return ComponentDiff{}, fmt.Errorf("failed to determine diff for component %q : %w", componentName, err)
	}

	return ComponentDiff{
		Name:         componentName,
		Expected:     expectedState,
		Actual:       actualState,
		NeededAction: nextAction,
	}, nil
}

func findComponentByName(components []Component, name string) (Component, bool) {
	for _, component := range components {
		if component.Name == name {
			return component, true
		}
	}
	return Component{}, false
}

func getNextComponentAction(expected ComponentDiffState, actual ComponentDiffState) (Action, error) {
	if expected.InstallationState == actual.InstallationState {
		return decideOnEqualState(expected, actual)
	}

	return decideOnDifferentState(expected)
}

func decideOnEqualState(expected ComponentDiffState, actual ComponentDiffState) (Action, error) {
	switch expected.InstallationState {
	case TargetStatePresent:
		if expected.DistributionNamespace != actual.DistributionNamespace {
			return ActionSwitchComponentDistributionNamespace, nil
		}

		if expected.Version.GreaterThan(actual.Version) {
			return ActionUpgrade, nil
		}
		if expected.Version.Equal(actual.Version) {
			return ActionNone, nil
		}
		return ActionDowngrade, nil
	case TargetStateAbsent:
		return ActionNone, nil
	default:
		return ActionNone, fmt.Errorf("component has unexpected target state %q", expected.InstallationState)
	}
}

func decideOnDifferentState(expected ComponentDiffState) (Action, error) {
	// at this place, the actual state is always the opposite to the expected state so just follow the expected state.
	switch expected.InstallationState {
	case TargetStatePresent:
		return ActionInstall, nil
	case TargetStateAbsent:
		return ActionUninstall, nil
	default:
		return ActionNone, fmt.Errorf("component has unexpected installation state %q", expected.InstallationState)
	}
}
