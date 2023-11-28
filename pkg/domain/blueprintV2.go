package domain

import "github.com/pkg/errors"

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

// Validate checks the structure and data of the blueprint and returns an error if there are any problems
func (blueprint *BlueprintV2) Validate() error {
	err := blueprint.validateDogus()
	if err != nil {
		return err
	}

	err = blueprint.validateDoguUniqueness()
	if err != nil {
		return err
	}

	err = blueprint.validateComponents()
	if err != nil {
		return err
	}

	err = blueprint.validateRegistryConfig()
	if err != nil {
		return err
	}

	return nil
}

func (blueprint *BlueprintV2) validateDogus() error {
	for _, dogu := range blueprint.Dogus {

		err := dogu.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// validateDoguUniqueness checks if the same dogu exists in the blueprint with different versions or target states and
// returns an error if it's so.
func (blueprint *BlueprintV2) validateDoguUniqueness() error {
	seenDogu := make(map[string]bool)

	for _, dogu := range blueprint.Dogus {
		_, seen := seenDogu[dogu.Name]
		if seen {
			return errors.Errorf("could not validate blueprint, there is at least one duplicate for this dogu: %s", dogu.Name)
		}
		seenDogu[dogu.Name] = true
	}

	return nil
}

func (blueprint *BlueprintV2) validateComponents() error {
	for _, component := range blueprint.Components {

		err := component.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (blueprint *BlueprintV2) validateRegistryConfig() error {
	for key, value := range blueprint.RegistryConfig {
		if len(key) == 0 {
			return errors.Errorf("could not validate blueprint, a config key is empty")
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
			return errors.Errorf("could not validate blueprint, a config key is empty")
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
