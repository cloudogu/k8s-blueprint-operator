package serializer

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

type TargetComponent struct {
	// Name defines the name of the component including its distribution namespace, f. i. "k8s/k8s-dogu-operator". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the component that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this component. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
	// DeployConfig defines a generic property map for the component configuration. This field is optional.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	DeployConfig map[string]interface{} `json:"deployConfig,omitempty"`
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
			Name:         name,
			Version:      version,
			TargetState:  newState,
			DeployConfig: component.DeployConfig,
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

		return TargetComponent{
			Name:         joinedComponentName,
			Version:      version,
			TargetState:  newState,
			DeployConfig: component.DeployConfig,
		}
	})
	return converted, errors.Join(errorList...)
}
