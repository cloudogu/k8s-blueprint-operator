package v1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// DoguDiff is the comparison of a Dogu's desired state vs. its cluster state.
// It contains the operation that needs to be done to achieve this desired state.
type DoguDiff struct {
	Actual       DoguDiffState `json:"actual"`
	Expected     DoguDiffState `json:"expected"`
	NeededAction DoguAction    `json:"neededAction"`
}

// DoguDiffState is either the actual or desired state of a dogu in the cluster.
type DoguDiffState struct {
	Namespace         string `json:"namespace,omitempty"`
	Version           string `json:"version,omitempty"`
	InstallationState string `json:"installationState"`
}

// DoguAction is the action that needs to be done for a dogu
// to achieve the desired state in the cluster.
type DoguAction string

func convertToDoguDiffDTO(domainModel domain.DoguDiff) DoguDiff {
	return DoguDiff{
		Actual: DoguDiffState{
			Namespace:         string(domainModel.Actual.Namespace),
			Version:           domainModel.Actual.Version.Raw,
			InstallationState: domainModel.Actual.InstallationState.String(),
		},
		Expected: DoguDiffState{
			Namespace:         string(domainModel.Expected.Namespace),
			Version:           domainModel.Expected.Version.Raw,
			InstallationState: domainModel.Expected.InstallationState.String(),
		},
		NeededAction: DoguAction(domainModel.NeededAction),
	}
}

func convertToDoguDiffDomain(doguName string, dto DoguDiff) (domain.DoguDiff, error) {
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
		DoguName: common.SimpleDoguName(doguName),
		Actual: domain.DoguDiffState{
			Namespace:         common.DoguNamespace(dto.Actual.Namespace),
			Version:           actualVersion,
			InstallationState: actualState,
		},
		Expected: domain.DoguDiffState{
			Namespace:         common.DoguNamespace(dto.Expected.Namespace),
			Version:           expectedVersion,
			InstallationState: expectedState,
		},
		NeededAction: domain.Action(dto.NeededAction),
	}, nil
}
