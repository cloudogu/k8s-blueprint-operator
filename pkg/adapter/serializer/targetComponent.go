package serializer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/blueprint-lib/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
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
	DeployConfig DeployConfig `json:"deployConfig,omitempty"`
}

func (in TargetComponent) DeepCopyInto(out *TargetComponent) {
	if out != nil {
		out.Name = in.Name
		out.Version = in.Version
		out.TargetState = in.TargetState
		out.DeployConfig = *in.DeployConfig.DeepCopy()
	}
}

type DeployConfig map[string]interface{}

func (in *DeployConfig) DeepCopy() *DeployConfig {
	out := new(DeployConfig)
	in.DeepCopyInto(out)
	return out
}

func (in *DeployConfig) DeepCopyInto(out *DeployConfig) {
	if out != nil {
		jsonStr, err := json.Marshal(in)
		if err != nil {
			return
		}
		err = json.Unmarshal(jsonStr, in)
		if err != nil {
			return
		}
	}
}

// ConvertComponents takes a slice of TargetComponent and returns a new slice with their DTO equivalent.
func ConvertComponents(components []TargetComponent) ([]v2.Component, error) {
	var convertedComponents []v2.Component
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

		name, err := v2.QualifiedComponentNameFromString(component.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		convertedComponents = append(convertedComponents, v2.Component{
			Name:         name,
			Version:      version,
			TargetState:  newState,
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
func ConvertToComponentDTOs(components []v2.Component) ([]TargetComponent, error) {
	var errorList []error
	converted := util.Map(components, func(component v2.Component) TargetComponent {
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
			DeployConfig: DeployConfig(component.DeployConfig),
		}
	})
	return converted, errors.Join(errorList...)
}
