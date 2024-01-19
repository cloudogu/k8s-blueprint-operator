package domain

import (
	"fmt"
	"github.com/go-logr/logr"
	"golang.org/x/exp/maps"

	"github.com/cloudogu/cesapp-lib/core"
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
	Name         string
	Actual       ComponentDiffState
	Expected     ComponentDiffState
	NeededAction Action
}

// ComponentDiffState contains all fields to make a diff for components in respect to another ComponentDiffState.
type ComponentDiffState struct {
	Version           core.Version
	InstallationState TargetState
}

// determineComponentDiffs creates ComponentDiffs for all components in the blueprint and all installed components as well.
func determineComponentDiffs(logger logr.Logger, blueprintComponents []Component, installedComponents map[string]*ecosystem.ComponentInstallation) []ComponentDiff {
	var componentDiffs = map[string]ComponentDiff{}
	for _, blueprintComponent := range blueprintComponents {
		installedComponent := installedComponents[blueprintComponent.Name]
		componentDiffs[blueprintComponent.Name] = determineComponentDiff(logger, &blueprintComponent, installedComponent)
	}

	for _, installedComponent := range installedComponents {
		blueprintComponent, notFound := findComponentByName(blueprintComponents, installedComponent.Name)

		if notFound == nil {
			componentDiffs[installedComponent.Name] = determineComponentDiff(logger, &blueprintComponent, installedComponent)
			continue
		}

		var notFoundInBlueprint *Component = nil
		componentDiffs[installedComponent.Name] = determineComponentDiff(logger, notFoundInBlueprint, installedComponent)
	}
	return maps.Values(componentDiffs)
}

// determineComponentDiff creates a ComponentDiff out of a Component from the blueprint and the ecosystem.ComponentInstallation in the ecosystem.
// If the Component is nil (was not in the blueprint), the actual state is also the expected state.
// If the installedComponent is nil, it is considered to be not installed currently.
func determineComponentDiff(logger logr.Logger, blueprintComponent *Component, installedComponent *ecosystem.ComponentInstallation) ComponentDiff {
	var expectedState, actualState ComponentDiffState
	componentName := "" // either blueprintComponent or installedComponent could be nil

	if installedComponent == nil {
		actualState = ComponentDiffState{
			InstallationState: TargetStateAbsent,
		}
	} else {
		componentName = installedComponent.Name
		actualState = ComponentDiffState{
			Version:           installedComponent.Version,
			InstallationState: TargetStatePresent,
		}
	}

	if blueprintComponent == nil {
		expectedState = actualState
	} else {
		componentName = blueprintComponent.Name
		expectedState = ComponentDiffState{
			Version:           blueprintComponent.Version,
			InstallationState: blueprintComponent.TargetState,
		}
	}

	return ComponentDiff{
		Name:         componentName,
		Expected:     expectedState,
		Actual:       actualState,
		NeededAction: getNextComponentAction(logger, expectedState, actualState),
	}
}

func findComponentByName(components []Component, name string) (Component, error) {
	for _, component := range components {
		if component.Name == name {
			return component, nil
		}
	}
	return Component{}, fmt.Errorf("could not find component '%s'", name)
}

func getNextComponentAction(logger logr.Logger, expected ComponentDiffState, actual ComponentDiffState) Action {
	if expected.InstallationState == actual.InstallationState {
		return decideOnEqualState(logger, expected, actual)
	}

	return decideOnDifferentState(logger, expected)
}

func decideOnEqualState(logger logr.Logger, expected ComponentDiffState, actual ComponentDiffState) Action {
	switch expected.InstallationState {
	case TargetStatePresent:
		if expected.Version.IsNewerThan(actual.Version) {
			return ActionUpgrade
		}
		if expected.Version.IsEqualTo(actual.Version) {
			return ActionNone
		}
		return ActionDowngrade
	case TargetStateAbsent:
		return ActionNone
	default:
		logger.Info("Warning: Component has unexpected target state, deciding for no action",
			"component", expected.InstallationState, "expected states", "")
		return ActionNone

	}
}

func decideOnDifferentState(logger logr.Logger, expected ComponentDiffState) Action {
	// at this place, the actual state is always the opposite to the expected state so just follow the expected state.
	switch expected.InstallationState {
	case TargetStatePresent:
		return ActionInstall
	case TargetStateAbsent:
		return ActionUninstall
	default:
		logger.Info("Warning: Component has unexpected installation state, deciding for no action",
			"component", expected.InstallationState, "expected states", "")
		return ActionNone
	}
}
