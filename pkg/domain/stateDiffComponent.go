package domain

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/blueprint-lib/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"golang.org/x/exp/maps"
	"reflect"
)

// ComponentDiffs contains the differences for all expected Components to the current ecosystem.ComponentInstallations.
type ComponentDiffs []ComponentDiff

// GetComponentDiffByName returns the diff for the given component name or an empty struct if it was not found.
func (diffs ComponentDiffs) GetComponentDiffByName(name v2.SimpleComponentName) ComponentDiff {
	for _, diff := range diffs {
		if diff.Name == name {
			return diff
		}
	}
	return ComponentDiff{}
}

// ComponentDiff represents the Diff for a single expected Component to the current ecosystem.ComponentInstallation.
type ComponentDiff struct {
	// Name contains the component's name.
	Name v2.SimpleComponentName
	// Actual contains that state of a component how it is currently found in the system.
	Actual ComponentDiffState
	// Expected contains that desired state of a component how it is supposed to be.
	Expected ComponentDiffState
	// NeededActions hints how the component should be handled by the application change automaton in order to reconcile
	// differences between Actual and Expected in the current system.
	NeededActions []Action
}

// ComponentDiffState contains all fields to make a diff for components in respect to another ComponentDiffState.
type ComponentDiffState struct {
	// Namespace is part of the address under which the component will be obtained. This namespace must NOT
	// to be confused with the K8s cluster namespace.
	Namespace v2.ComponentNamespace
	// Version contains the component's version.
	Version *semver.Version
	// InstallationState contains the component's target state.
	InstallationState v2.TargetState
	// DeployConfig contains generic properties for the component.
	DeployConfig ecosystem.DeployConfig
}

// IsExpectedVersion checks if the given version es equal to the expected version
func (diff ComponentDiff) IsExpectedVersion(actualVersion *semver.Version) bool {
	// expected is nil if the component is not in the blueprint, therefore no upgrade needs to happen
	if diff.Expected.Version == nil {
		return true
	}
	// actualVersion is nil if there is no component or no actual version in it yet.
	if actualVersion == nil {
		return false
	}
	return diff.Expected.Version.Equal(actualVersion)
}

func (diff ComponentDiff) HasChanges() bool {
	return len(diff.NeededActions) != 0
}

// String returns a string representation of the ComponentDiff.
func (diff *ComponentDiff) String() string {
	return fmt.Sprintf(
		"{Name: %q, Actual: %s, Expected: %s, NeededActions: %q}",
		diff.Name,
		diff.Actual.String(),
		diff.Expected.String(),
		diff.NeededActions,
	)
}

// String returns a string representation of the ComponentDiffState.
func (diff *ComponentDiffState) String() string {
	return fmt.Sprintf(
		"{Namespace: %q, Version: %q, InstallationState: %q}",
		diff.Namespace,
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
func determineComponentDiffs(blueprintComponents []v2.Component, installedComponents map[v2.SimpleComponentName]*ecosystem.ComponentInstallation) ([]ComponentDiff, error) {
	var componentDiffs = map[v2.SimpleComponentName]ComponentDiff{}
	for _, blueprintComponent := range blueprintComponents {
		installedComponent := installedComponents[blueprintComponent.Name.SimpleName]
		compDiff, err := determineComponentDiff(&blueprintComponent, installedComponent)
		if err != nil {
			return nil, err
		}
		componentDiffs[blueprintComponent.Name.SimpleName] = compDiff
	}

	for _, installedComponent := range installedComponents {
		_, found := findComponentByName(blueprintComponents, installedComponent.Name.SimpleName)
		// Only create ComponentDiff if the installed component is not found in the blueprint.
		// If the installed component is in blueprint the ComponentDiff was already determined above.
		if !found {
			compDiff, err := determineComponentDiff(nil, installedComponent)
			if err != nil {
				return nil, err
			}
			componentDiffs[installedComponent.Name.SimpleName] = compDiff
		}
	}
	return maps.Values(componentDiffs), nil
}

// determineComponentDiff creates a ComponentDiff out of a Component from the blueprint and the ecosystem.ComponentInstallation in the ecosystem.
// If the Component is nil (was not in the blueprint), the actual state is also the expected state.
// If the installedComponent is nil, it is considered to be not installed currently.
func determineComponentDiff(blueprintComponent *v2.Component, installedComponent *ecosystem.ComponentInstallation) (ComponentDiff, error) {
	var expectedState, actualState ComponentDiffState
	componentName := v2.SimpleComponentName("") // either blueprintComponent or installedComponent could be nil

	if installedComponent == nil {
		actualState = ComponentDiffState{
			InstallationState: v2.TargetStateAbsent,
		}
	} else {
		componentName = installedComponent.Name.SimpleName
		actualState = ComponentDiffState{
			Namespace:         installedComponent.Name.Namespace,
			Version:           installedComponent.ExpectedVersion,
			InstallationState: v2.TargetStatePresent,
			DeployConfig:      installedComponent.DeployConfig,
		}
	}

	if blueprintComponent == nil {
		expectedState = actualState
	} else {
		componentName = blueprintComponent.Name.SimpleName
		expectedState = ComponentDiffState{
			Namespace:         blueprintComponent.Name.Namespace,
			Version:           blueprintComponent.Version,
			InstallationState: blueprintComponent.TargetState,
			DeployConfig:      blueprintComponent.DeployConfig,
		}
	}

	nextActions, err := getComponentActions(expectedState, actualState)
	if err != nil {
		return ComponentDiff{}, fmt.Errorf("failed to determine diff for component %q : %w", componentName, err)
	}

	return ComponentDiff{
		Name:          componentName,
		Expected:      expectedState,
		Actual:        actualState,
		NeededActions: nextActions,
	}, nil
}

func findComponentByName(components []v2.Component, name v2.SimpleComponentName) (v2.Component, bool) {
	for _, component := range components {
		if component.Name.SimpleName == name {
			return component, true
		}
	}
	return v2.Component{}, false
}

func getComponentActions(expected ComponentDiffState, actual ComponentDiffState) ([]Action, error) {
	if expected.InstallationState == actual.InstallationState {
		return decideOnEqualState(expected, actual)
	}

	return decideOnDifferentState(expected)
}

func decideOnEqualState(expected ComponentDiffState, actual ComponentDiffState) ([]Action, error) {
	var neededActions []Action

	switch expected.InstallationState {
	case v2.TargetStatePresent:
		return getActionsForEqualPresentState(expected, actual), nil
	case v2.TargetStateAbsent:
		return neededActions, nil
	default:
		return nil, fmt.Errorf("component has unexpected target state %q", expected.InstallationState)
	}
}

func getActionsForEqualPresentState(expected ComponentDiffState, actual ComponentDiffState) []Action {
	var neededActions []Action

	if expected.Namespace != actual.Namespace {
		neededActions = append(neededActions, ActionSwitchComponentNamespace)
	}

	if !reflect.DeepEqual(expected.DeployConfig, actual.DeployConfig) {
		// Do update only if any DeployConfig contains data.
		// A nil DeployConfig and an empty DeployConfig are not deeply equal. But in this case we do not want to update the DeployConfig.
		if len(expected.DeployConfig) != 0 || len(actual.DeployConfig) != 0 {
			neededActions = append(neededActions, ActionUpdateComponentDeployConfig)
		}
	}

	if expected.Version.GreaterThan(actual.Version) {
		neededActions = append(neededActions, ActionUpgrade)
	} else if actual.Version.GreaterThan(expected.Version) {
		neededActions = append(neededActions, ActionDowngrade)
	}

	return neededActions
}

func decideOnDifferentState(expected ComponentDiffState) ([]Action, error) {
	// at this place, the actual state is always the opposite to the expected state so just follow the expected state.
	switch expected.InstallationState {
	case v2.TargetStatePresent:
		return []Action{ActionInstall}, nil
	case v2.TargetStateAbsent:
		return []Action{ActionUninstall}, nil
	default:
		return nil, fmt.Errorf("component has unexpected installation state %q", expected.InstallationState)
	}
}
