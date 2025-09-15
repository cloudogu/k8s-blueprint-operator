package serializer

import (
	"errors"
	"fmt"
	"slices"

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

	var dogus []cescommons.SimpleName
	var combinedConfigDiffs map[string]crd.CombinedDoguConfigDiff
	var doguConfigDiffByDogu map[cescommons.SimpleName]crd.DoguConfigDiff
	var sensitiveDoguConfigDiff map[cescommons.SimpleName]crd.SensitiveDoguConfigDiff

	if len(domainModel.DoguConfigDiffs) != 0 || len(domainModel.SensitiveDoguConfigDiffs) != 0 {
		combinedConfigDiffs = make(map[string]crd.CombinedDoguConfigDiff)
		doguConfigDiffByDogu = make(map[cescommons.SimpleName]crd.DoguConfigDiff)
		for doguName, doguConfigDiff := range domainModel.DoguConfigDiffs {
			doguConfigDiffByDogu[doguName] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff, false)
			dogus = append(dogus, doguName)
		}
		sensitiveDoguConfigDiff = make(map[cescommons.SimpleName]crd.SensitiveDoguConfigDiff)
		for doguName, doguConfigDiff := range domainModel.SensitiveDoguConfigDiffs {
			sensitiveDoguConfigDiff[doguName] = convertToDoguConfigEntryDiffsDTO(doguConfigDiff, true)
			dogus = append(dogus, doguName)
		}

		// remove duplicates, so we have a complete list of all dogus with config
		dogus = slices.Compact(dogus)
		for _, doguName := range dogus {
			combinedConfigDiffs[string(doguName)] = crd.CombinedDoguConfigDiff{
				DoguConfigDiff:          doguConfigDiffByDogu[doguName],
				SensitiveDoguConfigDiff: sensitiveDoguConfigDiff[doguName],
			}
		}
	}

	return &crd.StateDiff{
		DoguDiffs:        doguDiffs,
		ComponentDiffs:   componentDiffs,
		DoguConfigDiffs:  combinedConfigDiffs,
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
			doguConfigDiffs[cescommons.SimpleName(doguName)] = convertToDoguConfigDiffsDomain(doguName, crd.ConfigDiff(doguConfigDiff.DoguConfigDiff))
			sensitiveDoguConfigDiffs[cescommons.SimpleName(doguName)] = convertToDoguConfigDiffsDomain(doguName, crd.ConfigDiff(doguConfigDiff.SensitiveDoguConfigDiff))
		}
	}

	return domain.StateDiff{
		DoguDiffs:                doguDiffs,
		ComponentDiffs:           componentDiffs,
		DoguConfigDiffs:          doguConfigDiffs,
		SensitiveDoguConfigDiffs: sensitiveDoguConfigDiffs,
		GlobalConfigDiffs:        convertToGlobalConfigDiffDomain(dto.GlobalConfigDiff),
	}, nil
}
