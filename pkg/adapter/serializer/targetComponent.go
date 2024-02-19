package serializer

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type TargetComponent struct {
	// Name defines the name of the component including its distribution namespace, f. i. "k8s/k8s-dogu-operator". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the component that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this component. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
	// DeployNamespace defines the namespace where the component should be installed to. Actually this is only used for
	// the component `k8s-longhorn` because it requires the `longhorn-system` namespace.
	DeployNamespace string `json:"deployNamespace"`
}

// ConvertComponents takes a slice of TargetComponent and returns a new slice with their DTO equivalent.
func ConvertComponents(components []TargetComponent) ([]domain.Component, error) {
	var convertedComponents []domain.Component
	var errorList []error

	for _, component := range components {
		newState, err := ToDomainTargetState(component.TargetState)
		errorList = append(errorList, err)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		var version *semver.Version
		if component.Version != "" {
			version, err = semver.NewVersion(component.Version)
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
			Name:        name,
			Version:     version,
			TargetState: newState,
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedComponents, fmt.Errorf("cannot convert blueprint components: %w", err)
	}

	return convertedComponents, err
}

// ConvertToComponentDTOs takes a slice of Component DTOs and returns a new slice with their domain equivalent.
func ConvertToComponentDTOs(components []domain.Component) ([]TargetComponent, error) {
	var errorList []error
	converted := util.Map(components, func(component domain.Component) TargetComponent {
		newState, err := ToSerializerTargetState(component.TargetState)
		errorList = append(errorList, err)

		// convert the distribution namespace back into the name field so the EffectiveBlueprint has the same syntax
		// as the original blueprint json from the Blueprint resource.
		joinedComponentName := component.Name.String()
		version := ""
		if newState == "present" {
			version = component.Version.String()
		}

		// TODO Delete this if the blueprint can handle a component configuration.
		// This section would contain the deployNamespace in a generic Map.
		var deployNamespace string
		if component.Name == common.K8sK8sLonghornName {
			deployNamespace = "longhorn-system"
		}
		return TargetComponent{
			Name:            joinedComponentName,
			Version:         version,
			TargetState:     newState,
			DeployNamespace: deployNamespace,
		}
	})
	return converted, errors.Join(errorList...)
}
