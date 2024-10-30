package v1

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

// DoguDiff is the comparison of a Dogu's desired state vs. its cluster state.
// It contains the operation that needs to be done to achieve this desired state.
type DoguDiff struct {
	Actual        DoguDiffState `json:"actual"`
	Expected      DoguDiffState `json:"expected"`
	NeededActions []DoguAction  `json:"neededActions"`
}

// DoguDiffState is either the actual or desired state of a dogu in the cluster.
type DoguDiffState struct {
	Namespace          string             `json:"namespace,omitempty"`
	Version            string             `json:"version,omitempty"`
	InstallationState  string             `json:"installationState"`
	ResourceConfig     ResourceConfig     `json:"resourceConfig,omitempty"`
	ReverseProxyConfig ReverseProxyConfig `json:"reverseProxyConfig,omitempty"`
}

type ResourceConfig struct {
	MinVolumeSize string `json:"minVolumeSize,omitempty"`
}

type ReverseProxyConfig struct {
	MaxBodySize      string `json:"maxBodySize,omitempty"`
	RewriteTarget    string `json:"rewriteTarget,omitempty"`
	AdditionalConfig string `json:"additionalConfig,omitempty"`
}

// DoguAction is the action that needs to be done for a dogu
// to achieve the desired state in the cluster.
type DoguAction string

func convertToDoguDiffDTO(domainModel domain.DoguDiff) DoguDiff {
	neededActions := domainModel.NeededActions
	doguActions := make([]DoguAction, 0, len(neededActions))
	for _, action := range neededActions {
		doguActions = append(doguActions, DoguAction(action))
	}

	return DoguDiff{
		Actual: DoguDiffState{
			Namespace:         string(domainModel.Actual.Namespace),
			Version:           domainModel.Actual.Version.Raw,
			InstallationState: domainModel.Actual.InstallationState.String(),
			ResourceConfig: ResourceConfig{
				MinVolumeSize: ecosystem.GetQuantityString(domainModel.Actual.MinVolumeSize),
			},
			ReverseProxyConfig: ReverseProxyConfig{
				MaxBodySize:      ecosystem.GetQuantityString(domainModel.Actual.ReverseProxyConfig.MaxBodySize),
				RewriteTarget:    string(domainModel.Actual.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: string(domainModel.Actual.ReverseProxyConfig.AdditionalConfig),
			},
		},
		Expected: DoguDiffState{
			Namespace:         string(domainModel.Expected.Namespace),
			Version:           domainModel.Expected.Version.Raw,
			InstallationState: domainModel.Expected.InstallationState.String(),
			ResourceConfig: ResourceConfig{
				MinVolumeSize: ecosystem.GetQuantityString(domainModel.Expected.MinVolumeSize),
			},
			ReverseProxyConfig: ReverseProxyConfig{
				MaxBodySize:      ecosystem.GetQuantityString(domainModel.Expected.ReverseProxyConfig.MaxBodySize),
				RewriteTarget:    string(domainModel.Expected.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: string(domainModel.Expected.ReverseProxyConfig.AdditionalConfig),
			},
		},
		NeededActions: doguActions,
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

	actualMinVolumeSize, actualVolumeSizeErr := ecosystem.GetQuantityReference(dto.Actual.ResourceConfig.MinVolumeSize)
	if actualVolumeSizeErr != nil {
		actualVolumeSizeErr = fmt.Errorf("failed to parse actual minimum volume size %q: %w", dto.Actual.ResourceConfig.MinVolumeSize, actualVolumeSizeErr)
	}
	expectedMinVolumeSize, expectedVolumeSizeErr := ecosystem.GetQuantityReference(dto.Expected.ResourceConfig.MinVolumeSize)
	if expectedVolumeSizeErr != nil {
		expectedVolumeSizeErr = fmt.Errorf("failed to parse expected minimum volume size %q: %w", dto.Expected.ResourceConfig.MinVolumeSize, expectedVolumeSizeErr)
	}

	actualMaxBodySize, actualBodySizeErr := ecosystem.GetQuantityReference(dto.Actual.ReverseProxyConfig.MaxBodySize)
	if actualBodySizeErr != nil {
		actualBodySizeErr = fmt.Errorf("failed to parse actual maximum proxy body size %q: %w", dto.Actual.ReverseProxyConfig.MaxBodySize, actualBodySizeErr)
	}
	expectedMaxBodySize, expectedBodySizeErr := ecosystem.GetQuantityReference(dto.Expected.ReverseProxyConfig.MaxBodySize)
	if expectedBodySizeErr != nil {
		expectedBodySizeErr = fmt.Errorf("failed to parse expected maximum proxy body size %q: %w", dto.Expected.ReverseProxyConfig.MaxBodySize, expectedBodySizeErr)
	}

	err := errors.Join(actualVersionErr, expectedVersionErr, actualStateErr, expectedStateErr, actualVolumeSizeErr, expectedVolumeSizeErr, actualBodySizeErr, expectedBodySizeErr)
	if err != nil {
		return domain.DoguDiff{}, fmt.Errorf("failed to convert dogu diff dto %q to domain model: %w", doguName, err)
	}

	neededActions := dto.NeededActions
	doguActions := make([]domain.Action, 0, len(neededActions))
	for _, action := range neededActions {
		doguActions = append(doguActions, domain.Action(action))
	}

	return domain.DoguDiff{
		DoguName: cescommons.SimpleDoguName(doguName),
		Actual: domain.DoguDiffState{
			Namespace:         cescommons.DoguNamespace(dto.Actual.Namespace),
			Version:           actualVersion,
			InstallationState: actualState,
			MinVolumeSize:     actualMinVolumeSize,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      actualMaxBodySize,
				RewriteTarget:    ecosystem.RewriteTarget(dto.Actual.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: ecosystem.AdditionalConfig(dto.Actual.ReverseProxyConfig.AdditionalConfig),
			},
		},
		Expected: domain.DoguDiffState{
			Namespace:         cescommons.DoguNamespace(dto.Expected.Namespace),
			Version:           expectedVersion,
			InstallationState: expectedState,
			MinVolumeSize:     expectedMinVolumeSize,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      expectedMaxBodySize,
				RewriteTarget:    ecosystem.RewriteTarget(dto.Expected.ReverseProxyConfig.RewriteTarget),
				AdditionalConfig: ecosystem.AdditionalConfig(dto.Expected.ReverseProxyConfig.AdditionalConfig),
			},
		},
		NeededActions: doguActions,
	}, nil
}
