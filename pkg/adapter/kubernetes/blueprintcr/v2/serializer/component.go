package serializer

import (
	"errors"
	"fmt"

	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"k8s.io/utils/ptr"

	"github.com/Masterminds/semver/v3"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

// ConvertComponents takes a slice of TargetComponent and returns a new slice with their DTO equivalent.
func ConvertComponents(components []v2.Component) ([]domain.Component, error) {
	var convertedComponents []domain.Component
	var errorList []error

	for _, component := range components {
		var version *semver.Version
		if component.Version != nil && *component.Version != "" {
			var err error
			version, err = semver.NewVersion(*component.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of target component %q: %w", component.Name, err))
				continue
			}
		}

		name, err := common.QualifiedComponentNameFromString(component.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		convertedComponents = append(convertedComponents, domain.Component{
			Name:         name,
			Version:      version,
			Absent:       ptr.Deref(component.Absent, false),
			DeployConfig: ecosystem.DeployConfig(component.DeployConfig),
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedComponents, fmt.Errorf("cannot convert blueprint components: %w", err)
	}

	return convertedComponents, err
}

// ConvertToComponentDTOs takes a slice of Component DTOs and returns a new slice with their domain equivalent.
func ConvertToComponentDTOs(components []domain.Component) []v2.Component {
	converted := util.Map(components, func(component domain.Component) v2.Component {
		// convert the distribution namespace back into the name field so the EffectiveBlueprint has the same syntax
		// as the original blueprint json from the Blueprint resource.
		joinedComponentName := component.Name.String()
		var version *string
		if !component.Absent && component.Version != nil {
			version = ptr.To(component.Version.String())
		}

		return v2.Component{
			Name:         joinedComponentName,
			Version:      version,
			Absent:       &component.Absent,
			DeployConfig: map[string]interface{}(component.DeployConfig),
		}
	})
	return converted
}
