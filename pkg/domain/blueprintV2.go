package domain

import (
	"errors"
	"fmt"
)

// BlueprintV2 describes an abstraction of CES components that should be absent or present within one or more CES
// instances. When the same Blueprint is applied to two different CES instances it is required to leave two equal
// instances in terms of the components.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type BlueprintV2 struct {
	GeneralBlueprint
	// ID is the unique name of the set over all components. This blueprint ID should be used to distinguish from similar
	// blueprints between humans in an easy way. Must not be empty.
	ID string `json:"blueprintId"`
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu `json:"dogus,omitempty"`
	// Components contains a set of exact components versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []Component `json:"components,omitempty"`
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig RegistryConfig `json:"registryConfig,omitempty"`
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string `json:"registryConfigAbsent,omitempty"`
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted RegistryConfig `json:"registryConfigEncrypted,omitempty"`
}

type RegistryConfig map[string]map[string]interface{}

// Validate checks the structure and data of the blueprint statically and returns an error if there are any problems
func (blueprint *BlueprintV2) Validate() error {
	errorList := []error{
		blueprint.validateDogus(),
		blueprint.validateDoguUniqueness(),
		blueprint.validateComponents(),
		blueprint.validateComponentUniqueness(),
		blueprint.validateRegistryConfig(),
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("could not Validate blueprint: %w", err)
	}
	return err
}

func (blueprint *BlueprintV2) validateDogus() error {
	errorList := Map(blueprint.Dogus, func(dogu TargetDogu) error { return dogu.Validate() })
	return errors.Join(errorList...)
}

// validateDoguUniqueness checks if dogus exist twice in the blueprint and returns an error if it's so.
func (blueprint *BlueprintV2) validateDoguUniqueness() error {
	doguNames := Map(blueprint.Dogus, func(dogu TargetDogu) string { return dogu.Name })
	duplicates := getDuplicates(doguNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate dogus: %v", duplicates)
	}
	return nil
}

func (blueprint *BlueprintV2) validateComponents() error {
	errorList := Map(blueprint.Components, func(component Component) error { return component.Validate() })
	return errors.Join(errorList...)
}

// validateComponentUniqueness checks if components exist twice in the blueprint and returns an error if it's so.
func (blueprint *BlueprintV2) validateComponentUniqueness() error {
	componentNames := Map(blueprint.Components, func(component Component) string { return component.Name })
	duplicates := getDuplicates(componentNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate components: %v", duplicates)
	}
	return nil
}

func (blueprint *BlueprintV2) validateRegistryConfig() error {
	for key, value := range blueprint.RegistryConfig {
		if len(key) == 0 {
			return fmt.Errorf("a config key is empty")
		}

		err := validateKeysNotEmpty(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateKeysNotEmpty(config map[string]interface{}) error {
	for key, value := range config {
		if len(key) == 0 {
			return fmt.Errorf("a config key is empty")
		}

		switch vTyped := value.(type) {
		case map[string]interface{}:
			err := validateKeysNotEmpty(vTyped)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getDuplicates(list []string) []string {
	elementCount := make(map[string]int)

	// countByName
	for _, name := range list {
		elementCount[name] += 1
	}

	// get list of names with count != 1
	var duplicates []string
	for name, count := range elementCount {
		if count != 1 {
			duplicates = append(duplicates, name)
		}
	}
	return duplicates
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
