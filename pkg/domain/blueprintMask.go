package domain

import (
	"errors"
	"fmt"
)

// BlueprintMask describes an abstraction of CES components that should alter a blueprint definition before
// applying it to a CES system via a blueprint upgrade. The blueprint mask should not change the blueprint JSON file
// itself, but is applied to the information in it to generate a new, effective blueprint.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint mask API version.
type BlueprintMask struct {
	// Dogus contains a set of dogus which alters the states of the dogus in the blueprint this mask is applied on.
	// The names and target states of all dogus must not be empty.
	Dogus []MaskTargetDogu `json:"dogus"`
}

// Validate checks the structure and data of a blueprint mask and returns an error if there are any problems
func (blueprintMask *BlueprintMask) Validate() error {
	errorList := []error{
		blueprintMask.validateDogus(),
		blueprintMask.validateDoguUniqueness(),
	}
	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint mask is invalid: %w", err)
	}
	return err
}

func (blueprintMask *BlueprintMask) validateDogus() error {
	for _, dogu := range blueprintMask.Dogus {

		err := dogu.validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// validateDoguUniqueness checks if the same dogu exists in the blueprint mask with different versions or
// target states and returns an error if it's so.
func (blueprintMask *BlueprintMask) validateDoguUniqueness() error {
	seenDogu := make(map[string]bool)

	for _, dogu := range blueprintMask.Dogus {
		_, seen := seenDogu[dogu.Name]
		if seen {
			return fmt.Errorf("could not Validate blueprint mask, there is at least one duplicate for this dogu: %s", dogu.Name)
		}
		seenDogu[dogu.Name] = true
	}

	return nil
}

func (blueprintMask *BlueprintMask) FindDoguByName(name string) (MaskTargetDogu, error) {
	for doguIndex, dogu := range blueprintMask.Dogus {
		if dogu.Name == name {
			return blueprintMask.Dogus[doguIndex], nil
		}
	}
	return MaskTargetDogu{}, fmt.Errorf("could not find dogu name %s in blueprint", name)
}
