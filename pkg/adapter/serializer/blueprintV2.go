package serializer

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"strings"
)

type BlueprintApi string

const (
	V1 BlueprintApi = "v1"
	V2 BlueprintApi = "v2"
)

// GeneralBlueprint defines the minimum set to parse the blueprint API version string in order to select the right
// blueprint handling strategy. This is necessary in order to accommodate maximal changes in different blueprint API
// versions.
type GeneralBlueprint struct {
	// API is used to distinguish between different versions of the used API and impacts directly the interpretation of
	// this blueprint. Must not be empty.
	//
	// This field MUST NOT be MODIFIED or REMOVED because the API is paramount for distinguishing between different
	// blueprint version implementations.
	API BlueprintApi `json:"blueprintApi"`
}

// BlueprintV2 describes an abstraction of CES components that should be absent or present within one or more CES
// instances. When the same Blueprint is applied to two different CES instances it is required to leave two equal
// instances in terms of the components.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type BlueprintV2 struct {
	GeneralBlueprint
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu `json:"dogus,omitempty"`
	// Packages contains a set of exact package versions which should be present or absent in the CES instance after which
	// this blueprint was applied. The packages must correspond to the used operation system package manager. Optional.
	Components []TargetComponent `json:"components,omitempty"`
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig RegistryConfig `json:"registryConfig,omitempty"`
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string `json:"registryConfigAbsent,omitempty"`
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted RegistryConfig `json:"registryConfigEncrypted,omitempty"`
}

type RegistryConfig map[string]map[string]interface{}

// TargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type TargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
}

type TargetComponent struct {
	// Name defines the name of the component including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this component. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
}

func ConvertToBlueprintV2(blueprint domain.Blueprint) (BlueprintV2, error) {
	var errorList []error
	convertedDogus := util.Map(blueprint.Dogus, func(dogu domain.TargetDogu) TargetDogu {
		newState, err := toSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)
		return TargetDogu{
			Name:        dogu.GetQualifiedName(),
			Version:     dogu.Version,
			TargetState: newState,
		}
	})
	convertedComponents := util.Map(blueprint.Components, func(component domain.Component) TargetComponent {
		newState, err := toSerializerTargetState(component.TargetState)
		errorList = append(errorList, err)
		return TargetComponent{
			Name:        component.Name,
			Version:     component.Version,
			TargetState: newState,
		}
	})

	err := errors.Join(errorList...)
	if err != nil {
		return BlueprintV2{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return BlueprintV2{
		GeneralBlueprint:        GeneralBlueprint{V2},
		Dogus:                   convertedDogus,
		Components:              convertedComponents,
		RegistryConfig:          RegistryConfig(blueprint.RegistryConfig),
		RegistryConfigAbsent:    blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: RegistryConfig(blueprint.RegistryConfigEncrypted),
	}, nil
}

func convertToBlueprint(blueprint BlueprintV2) (domain.Blueprint, error) {
	convertedDogus, doguErr := convertDogus(blueprint.Dogus)
	convertedComponents, compErr := convertComponents(blueprint.Components)
	err := errors.Join(doguErr, compErr)
	if err != nil {
		return domain.Blueprint{}, fmt.Errorf("syntax of blueprintV2 is not correct: %w", err)
	}
	return domain.Blueprint{
		Dogus:                   convertedDogus,
		Components:              convertedComponents,
		RegistryConfig:          domain.RegistryConfig(blueprint.RegistryConfig),
		RegistryConfigAbsent:    blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: domain.RegistryConfig(blueprint.RegistryConfigEncrypted),
	}, nil
}

func convertDogus(dogus []TargetDogu) ([]domain.TargetDogu, error) {
	var convertedDogus []domain.TargetDogu
	var errorList []error

	for _, dogu := range dogus {
		doguNamespace, doguName, err := splitDoguName(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		newState, err := toDomainTargetState(dogu.TargetState)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		convertedDogus = append(convertedDogus, domain.TargetDogu{
			Namespace:   doguNamespace,
			Name:        doguName,
			Version:     dogu.Version,
			TargetState: newState,
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func convertComponents(components []TargetComponent) ([]domain.Component, error) {
	var convertedDogus []domain.Component
	var errorList []error

	for _, component := range components {
		newState, err := toDomainTargetState(component.TargetState)
		errorList = append(errorList, err)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}

		convertedDogus = append(convertedDogus, domain.Component{
			Name:        component.Name,
			Version:     component.Version,
			TargetState: newState,
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint components: %w", err)
	}

	return convertedDogus, err
}

func splitDoguName(doguName string) (string, string, error) {
	splitName := strings.Split(doguName, "/")
	if len(splitName) != 2 {
		return "", "", fmt.Errorf("dogu name needs to be in the form 'namespace/dogu' but is '%s'", doguName)
	}
	return splitName[0], splitName[1], nil
}
