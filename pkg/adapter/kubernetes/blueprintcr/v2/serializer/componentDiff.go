package serializer

import (
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"k8s.io/utils/ptr"
)

func convertToComponentDiffDTO(domainModel domain.ComponentDiff) crd.ComponentDiff {
	var actualVersion, expectedVersion *string

	if domainModel.Actual.Version != nil {
		actualVersion = ptr.To(domainModel.Actual.Version.String())
	}
	if domainModel.Expected.Version != nil {
		expectedVersion = ptr.To(domainModel.Expected.Version.String())
	}

	neededActions := domainModel.NeededActions
	componentActions := make([]crd.ComponentAction, 0, len(neededActions))
	for _, action := range neededActions {
		componentActions = append(componentActions, crd.ComponentAction(action))
	}

	return crd.ComponentDiff{
		Actual: crd.ComponentDiffState{
			Namespace:    string(domainModel.Actual.Namespace),
			Version:      actualVersion,
			Absent:       domainModel.Actual.Absent,
			DeployConfig: crd.DeployConfig(domainModel.Actual.DeployConfig),
		},
		Expected: crd.ComponentDiffState{
			Namespace:    string(domainModel.Expected.Namespace),
			Version:      expectedVersion,
			Absent:       domainModel.Expected.Absent,
			DeployConfig: crd.DeployConfig(domainModel.Expected.DeployConfig),
		},
		NeededActions: componentActions,
	}
}
