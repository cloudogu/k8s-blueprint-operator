package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

// Blueprint describes an abstraction of CES components that should be absent or present within one or more CES
// instances. When the same Blueprint is applied to two different CES instances it is required to leave two equal
// instances in terms of the components.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type Blueprint struct {
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu
	// Components contains a set of exact components versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []Component
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig RegistryConfig
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted RegistryConfig
}

type RegistryConfig map[string]map[string]interface{}

// Validate checks the structure and data of the blueprint statically and returns an error if there are any problems
func (blueprint *Blueprint) Validate() error {
	errorList := []error{
		blueprint.validateDogus(),
		blueprint.validateDoguUniqueness(),
		blueprint.validateComponents(),
		blueprint.validateComponentUniqueness(),
		blueprint.validateRegistryConfig(),
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint is invalid: %w", err)
	}
	return err
}

func (blueprint *Blueprint) validateDogus() error {
	errorList := util.Map(blueprint.Dogus, func(dogu TargetDogu) error { return dogu.validate() })
	return errors.Join(errorList...)
}

// validateDoguUniqueness checks if dogus exist twice in the blueprint and returns an error if it's so.
func (blueprint *Blueprint) validateDoguUniqueness() error {
	doguNames := util.Map(blueprint.Dogus, func(dogu TargetDogu) string { return dogu.Name })
	duplicates := util.GetDuplicates(doguNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate dogus: %v", duplicates)
	}
	return nil
}

func (blueprint *Blueprint) validateComponents() error {
	errorList := util.Map(blueprint.Components, func(component Component) error { return component.Validate() })
	return errors.Join(errorList...)
}

// validateComponentUniqueness checks if components exist twice in the blueprint and returns an error if it's so.
func (blueprint *Blueprint) validateComponentUniqueness() error {
	componentNames := util.Map(blueprint.Components, func(component Component) string { return component.Name })
	duplicates := util.GetDuplicates(componentNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate components: %v", duplicates)
	}
	return nil
}

func (blueprint *Blueprint) validateRegistryConfig() error {
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

func (blueprint *Blueprint) FindDoguByName(name string) (TargetDogu, error) {
	for doguIndex, dogu := range blueprint.Dogus {
		if dogu.Name == name {
			return blueprint.Dogus[doguIndex], nil
		}
	}
	return TargetDogu{}, fmt.Errorf("could not find dogu name %s in blueprint", name)
}

// GetWantedDogus returns a list of all dogus which should be installed
func (blueprint *Blueprint) GetWantedDogus() []TargetDogu {
	var wantedDogus []TargetDogu
	for _, dogu := range blueprint.Dogus {
		if dogu.TargetState == TargetStatePresent {
			wantedDogus = append(wantedDogus, dogu)
		}
	}
	return wantedDogus
}
