package effectiveBlueprintV1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

var configKeySeparator = "/"

// EffectiveBlueprintV1 describes an abstraction of CES components that should be absent or present within one or more CES
// instances after combining the blueprint with the blueprint mask.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type EffectiveBlueprintV1 struct {
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu `json:"dogus,omitempty"`
	// Packages contains a set of exact package versions which should be present or absent in the CES instance after which
	// this blueprint was applied. The packages must correspond to the used operation system package manager. Optional.
	Components []TargetComponent `json:"components,omitempty"`
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig map[string]string `json:"registryConfig,omitempty"`
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string `json:"registryConfigAbsent,omitempty"`
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted map[string]string `json:"registryConfigEncrypted,omitempty"`
}

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

func ConvertToEffectiveBlueprintV1(blueprint domain.EffectiveBlueprint) (EffectiveBlueprintV1, error) {
	var errorList []error
	convertedDogus := util.Map(blueprint.Dogus, func(dogu domain.Dogu) TargetDogu {
		newState, err := serializer.ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)
		return TargetDogu{
			Name:        dogu.GetQualifiedName(),
			Version:     dogu.Version.Raw,
			TargetState: newState,
		}
	})
	convertedComponents := util.Map(blueprint.Components, func(component domain.Component) TargetComponent {
		newState, err := serializer.ToSerializerTargetState(component.TargetState)
		errorList = append(errorList, err)
		return TargetComponent{
			Name:        component.Name,
			Version:     component.Version.Raw,
			TargetState: newState,
		}
	})

	err := errors.Join(errorList...)
	if err != nil {
		return EffectiveBlueprintV1{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return EffectiveBlueprintV1{
		Dogus:                   convertedDogus,
		Components:              convertedComponents,
		RegistryConfig:          flattenRegistryConfig(blueprint.RegistryConfig),
		RegistryConfigAbsent:    blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: flattenRegistryConfig(blueprint.RegistryConfigEncrypted),
	}, nil
}

func ConvertToEffectiveBlueprint(blueprint EffectiveBlueprintV1) (domain.EffectiveBlueprint, error) {
	convertedDogus, doguErr := convertDogus(blueprint.Dogus)
	convertedComponents, compErr := convertComponents(blueprint.Components)
	err := errors.Join(doguErr, compErr)
	if err != nil {
		return domain.EffectiveBlueprint{}, fmt.Errorf("syntax of blueprintV2 is not correct: %w", err)
	}
	return domain.EffectiveBlueprint{
		Dogus:                   convertedDogus,
		Components:              convertedComponents,
		RegistryConfig:          convertToRegistryConfig(blueprint.RegistryConfig),
		RegistryConfigAbsent:    blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: convertToRegistryConfig(blueprint.RegistryConfigEncrypted),
	}, nil
}

func convertDogus(dogus []TargetDogu) ([]domain.Dogu, error) {
	var convertedDogus []domain.Dogu
	var errorList []error

	for _, dogu := range dogus {
		doguNamespace, doguName, err := serializer.SplitDoguName(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		newState, err := serializer.ToDomainTargetState(dogu.TargetState)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		var version core.Version
		if dogu.Version != "" {
			version, err = core.ParseVersion(dogu.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of target dogu %q: %w", dogu.Name, err))
				continue
			}
		}
		convertedDogus = append(convertedDogus, domain.Dogu{
			Namespace:   doguNamespace,
			Name:        doguName,
			Version:     version,
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
	var convertedComponents []domain.Component
	var errorList []error

	for _, component := range components {
		newState, err := serializer.ToDomainTargetState(component.TargetState)
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
		convertedComponents = append(convertedComponents, domain.Component{
			Name:        component.Name,
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

func convertToRegistryConfig(flattenedConfig map[string]string) domain.RegistryConfig {
	//TODO: implement this
	return nil
}

func flattenRegistryConfig(config domain.RegistryConfig) map[string]string {
	intermediateResult := make(map[string]interface{})
	for key, value := range config {
		for key2, value2 := range value {
			intermediateResult[key+configKeySeparator+key2] = value2
		}
	}
	keyToValueConfig := make(map[string]string)
	flattenMap("", intermediateResult, keyToValueConfig)

	return keyToValueConfig
}

func flattenMap(prefix string, src map[string]interface{}, dest map[string]string) {
	if len(prefix) > 0 {
		prefix += configKeySeparator
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flattenMap(prefix+k, child, dest)
		//case []interface{}: there should be no arrays in the config
		default:
			dest[prefix+k] = fmt.Sprintf("%s", v)
		}
	}
}
