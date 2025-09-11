package serializer

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

func convertToComponentDiffDTO(domainModel domain.ComponentDiff) crd.ComponentDiff {
	actualVersion := ""
	expectedVersion := ""

	if domainModel.Actual.Version != nil {
		actualVersion = domainModel.Actual.Version.String()
	}
	if domainModel.Expected.Version != nil {
		expectedVersion = domainModel.Expected.Version.String()
	}

	neededActions := domainModel.NeededActions
	componentActions := make([]crd.ComponentAction, 0, len(neededActions))
	for _, action := range neededActions {
		componentActions = append(componentActions, crd.ComponentAction(action))
	}

	return crd.ComponentDiff{
		Actual: crd.ComponentDiffState{
			Namespace:    string(domainModel.Actual.Namespace),
			Version:      &actualVersion,
			Absent:       domainModel.Actual.Absent,
			DeployConfig: crd.DeployConfig(domainModel.Actual.DeployConfig),
		},
		Expected: crd.ComponentDiffState{
			Namespace:    string(domainModel.Expected.Namespace),
			Version:      &expectedVersion,
			Absent:       domainModel.Expected.Absent,
			DeployConfig: crd.DeployConfig(domainModel.Expected.DeployConfig),
		},
		NeededActions: componentActions,
	}
}

func convertToComponentDiffDomain(componentName string, dto crd.ComponentDiff) (domain.ComponentDiff, error) {
	var actualVersion *semver.Version
	var actualVersionErr error
	if dto.Actual.Version != nil && *dto.Actual.Version != "" {
		actualVersion, actualVersionErr = semver.NewVersion(*dto.Actual.Version)
		if actualVersionErr != nil {
			actualVersionErr = fmt.Errorf("failed to parse actual version %q: %w", *dto.Actual.Version, actualVersionErr)
		}
	}

	var expectedVersion *semver.Version
	var expectedVersionErr error
	if dto.Expected.Version != nil && *dto.Expected.Version != "" {
		expectedVersion, expectedVersionErr = semver.NewVersion(*dto.Expected.Version)
		if expectedVersionErr != nil {
			expectedVersionErr = fmt.Errorf("failed to parse expected version %q: %w", *dto.Expected.Version, expectedVersionErr)
		}
	}

	neededActions := dto.NeededActions
	componentActions := make([]domain.Action, 0, len(neededActions))
	for _, action := range neededActions {
		componentActions = append(componentActions, domain.Action(action))
	}

	err := errors.Join(actualVersionErr, expectedVersionErr)
	if err != nil {
		return domain.ComponentDiff{}, fmt.Errorf("failed to convert component diff dto %q to domain model: %w", componentName, err)
	}

	return domain.ComponentDiff{
		Name: common.SimpleComponentName(componentName),
		Actual: domain.ComponentDiffState{
			Namespace:    common.ComponentNamespace(dto.Actual.Namespace),
			Version:      actualVersion,
			Absent:       dto.Actual.Absent,
			DeployConfig: ecosystem.DeployConfig(dto.Actual.DeployConfig),
		},
		Expected: domain.ComponentDiffState{
			Namespace:    common.ComponentNamespace(dto.Expected.Namespace),
			Version:      expectedVersion,
			Absent:       dto.Expected.Absent,
			DeployConfig: ecosystem.DeployConfig(dto.Expected.DeployConfig),
		},
		NeededActions: componentActions,
	}, nil
}
