package serializer

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

type TargetComponent struct {
	// Name defines the name of the component including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the component that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this component. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
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
		var version core.Version
		if component.Version != "" {
			version, err = core.ParseVersion(component.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of target component %q: %w", component.Name, err))
				continue
			}
		}

		namespace, name, err := SplitComponentName(component.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		convertedComponents = append(convertedComponents, domain.Component{
			Name:                  name,
			DistributionNamespace: namespace,
			Version:               version,
			TargetState:           newState,
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
		return TargetComponent{
			Name:        component.Name,
			Version:     component.Version.Raw,
			TargetState: newState,
		}
	})
	return converted, errors.Join(errorList...)
}
