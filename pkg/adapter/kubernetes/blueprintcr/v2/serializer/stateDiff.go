package serializer

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToStateDiffDTO(domainModel domain.StateDiff) *crd.StateDiff {
	doguDiffs := make(map[string]crd.DoguDiff, len(domainModel.DoguDiffs))
	for _, doguDiff := range domainModel.DoguDiffs {
		doguDiffs[string(doguDiff.DoguName)] = convertToDoguDiffDTO(doguDiff)
	}

	componentDiffs := make(map[string]crd.ComponentDiff, len(domainModel.ComponentDiffs))
	for _, componentDiff := range domainModel.ComponentDiffs {
		componentDiffs[string(componentDiff.Name)] = convertToComponentDiffDTO(componentDiff)
	}

	var doguConfigDiffs map[string]crd.ConfigDiff

	if len(domainModel.DoguConfigDiffs) != 0 || len(domainModel.SensitiveDoguConfigDiffs) != 0 {
		doguConfigDiffs = make(map[string]crd.ConfigDiff)
		for doguName, doguConfigDiff := range domainModel.DoguConfigDiffs {
			doguConfigDiffs[doguName.String()] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff, false)
		}
		for doguName, doguConfigDiff := range domainModel.SensitiveDoguConfigDiffs {
			doguConfigDiffs[doguName.String()] = append(doguConfigDiffs[doguName.String()], convertToDoguConfigEntryDiffsDTO(doguConfigDiff, true)...)
		}
	}

	return &crd.StateDiff{
		DoguDiffs:        doguDiffs,
		ComponentDiffs:   componentDiffs,
		DoguConfigDiffs:  doguConfigDiffs,
		GlobalConfigDiff: convertToGlobalConfigDiffDTO(domainModel.GlobalConfigDiffs),
	}
}

func ConvertToStateDiffDomain(dto *crd.StateDiff) (domain.StateDiff, error) {
	if dto == nil {
		return domain.StateDiff{}, nil
	}

	var errs []error

	var doguDiffs []domain.DoguDiff
	for doguName, doguDiff := range dto.DoguDiffs {
		doguDiffDomainModel, err := convertToDoguDiffDomain(doguName, doguDiff)
		errs = append(errs, err)
		doguDiffs = append(doguDiffs, doguDiffDomainModel)
	}

	var componentDiffs []domain.ComponentDiff
	for componentName, componentDiff := range dto.ComponentDiffs {
		componentDiffDomainModel, err := convertToComponentDiffDomain(componentName, componentDiff)
		errs = append(errs, err)
		componentDiffs = append(componentDiffs, componentDiffDomainModel)
	}

	err := errors.Join(errs...)
	if err != nil {
		return domain.StateDiff{}, fmt.Errorf("failed to convert state diff DTO to domain model: %w", err)
	}

	var doguConfigDiffs map[cescommons.SimpleName]domain.DoguConfigDiffs
	var sensitiveDoguConfigDiffs map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs
	if len(dto.DoguConfigDiffs) != 0 {
		doguConfigDiffs = map[cescommons.SimpleName]domain.DoguConfigDiffs{}
		sensitiveDoguConfigDiffs = map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{}
		for doguName, doguConfigDiff := range dto.DoguConfigDiffs {
			// TODO: remove sensitive config diff from domain
			doguConfigDiffs[cescommons.SimpleName(doguName)] = convertToDoguConfigDiffsDomain(doguName, doguConfigDiff)
			sensitiveDoguConfigDiffs[cescommons.SimpleName(doguName)] = convertToDoguConfigDiffsDomain(doguName, doguConfigDiff)
		}
	}

	return domain.StateDiff{
		DoguDiffs:         doguDiffs,
		ComponentDiffs:    componentDiffs,
		DoguConfigDiffs:   doguConfigDiffs,
		GlobalConfigDiffs: convertToGlobalConfigDiffDomain(dto.GlobalConfigDiff),
	}, nil
}
