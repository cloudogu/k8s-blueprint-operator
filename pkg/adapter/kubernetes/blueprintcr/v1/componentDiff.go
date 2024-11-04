package v1

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

// ComponentDiff is the comparison of a Component's desired state vs. its cluster state.
// It contains the operation that needs to be done to achieve this desired state.
type ComponentDiff struct {
	// Actual contains the component's state in the current system.
	Actual ComponentDiffState `json:"actual"`
	// Expected contains the desired component's target state.
	Expected ComponentDiffState `json:"expected"`
	// NeededActions contains the refined actions as decided by the application's state determination automaton.
	NeededActions []ComponentAction `json:"neededActions"`
}

// ComponentDiffState is either the actual or desired state of a component in the cluster. The fields will be used to
// determine the kind of changed if there is a drift between actual or desired state.
type ComponentDiffState struct {
	// Namespace is part of the address under which the component will be obtained. This namespace must NOT
	// to be confused with the K8s cluster namespace.
	Namespace string `json:"distributionNamespace,omitempty"`
	// Version contains the component's version.
	Version string `json:"version,omitempty"`
	// InstallationState contains the component's installation state. Such a state correlate with the domain Actions:
	//
	//  - domain.ActionInstall
	//  - domain.ActionUninstall
	//  - and so on
	InstallationState string `json:"installationState"`
	// DeployConfig contains generic properties for the component.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	DeployConfig serializer.DeployConfig `json:"deployConfig,omitempty"`
}

// ComponentAction is the action that needs to be done for a component
// to achieve the desired state in the cluster.
type ComponentAction string

func convertToComponentDiffDTO(domainModel domain.ComponentDiff) ComponentDiff {
	actualVersion := ""
	expectedVersion := ""

	if domainModel.Actual.Version != nil {
		actualVersion = domainModel.Actual.Version.String()
	}
	if domainModel.Expected.Version != nil {
		expectedVersion = domainModel.Expected.Version.String()
	}

	neededActions := domainModel.NeededActions
	componentActions := make([]ComponentAction, 0, len(neededActions))
	for _, action := range neededActions {
		componentActions = append(componentActions, ComponentAction(action))
	}

	return ComponentDiff{
		Actual: ComponentDiffState{
			Namespace:         string(domainModel.Actual.Namespace),
			Version:           actualVersion,
			InstallationState: domainModel.Actual.InstallationState.String(),
			DeployConfig:      serializer.DeployConfig(domainModel.Actual.DeployConfig),
		},
		Expected: ComponentDiffState{
			Namespace:         string(domainModel.Expected.Namespace),
			Version:           expectedVersion,
			InstallationState: domainModel.Expected.InstallationState.String(),
			DeployConfig:      serializer.DeployConfig(domainModel.Expected.DeployConfig),
		},
		NeededActions: componentActions,
	}
}

func convertToComponentDiffDomain(componentName string, dto ComponentDiff) (domain.ComponentDiff, error) {
	var actualVersion *semver.Version
	var actualVersionErr error
	if dto.Actual.Version != "" {
		actualVersion, actualVersionErr = semver.NewVersion(dto.Actual.Version)
		if actualVersionErr != nil {
			actualVersionErr = fmt.Errorf("failed to parse actual version %q: %w", dto.Actual.Version, actualVersionErr)
		}
	}

	var expectedVersion *semver.Version
	var expectedVersionErr error
	if dto.Expected.Version != "" {
		expectedVersion, expectedVersionErr = semver.NewVersion(dto.Expected.Version)
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

	actualDistributionNamespace := dto.Actual.Namespace
	expectedDistributionNamespace := dto.Expected.Namespace

	neededActions := dto.NeededActions
	componentActions := make([]domain.Action, 0, len(neededActions))
	for _, action := range neededActions {
		componentActions = append(componentActions, domain.Action(action))
	}

	err := errors.Join(actualVersionErr, expectedVersionErr, actualStateErr, expectedStateErr)
	if err != nil {
		return domain.ComponentDiff{}, fmt.Errorf("failed to convert component diff dto %q to domain model: %w", componentName, err)
	}

	return domain.ComponentDiff{
		Name: common.SimpleComponentName(componentName),
		Actual: domain.ComponentDiffState{
			Namespace:         common.ComponentNamespace(actualDistributionNamespace),
			Version:           actualVersion,
			InstallationState: actualState,
			DeployConfig:      ecosystem.DeployConfig(dto.Actual.DeployConfig),
		},
		Expected: domain.ComponentDiffState{
			Namespace:         common.ComponentNamespace(expectedDistributionNamespace),
			Version:           expectedVersion,
			InstallationState: expectedState,
			DeployConfig:      ecosystem.DeployConfig(dto.Expected.DeployConfig),
		},
		NeededActions: componentActions,
	}, nil
}
