package domain

import "github.com/pkg/errors"

// BlueprintMaskV1 describes an abstraction of CES components that should alter a blueprint definition before
// applying it to a CES system via a blueprint upgrade. The blueprint mask should not change the blueprint JSON file
// itself, but is applied to the information in it to generate a new, effective blueprint.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint mask API version.
type BlueprintMaskV1 struct {
	GeneralBlueprintMask
	// ID is the unique name of the set over all components. This blueprint mask ID should be used to distinguish
	// from similar blueprint masks between humans in an easy way. Must not be empty.
	ID string `json:"blueprintMaskId"`
	// Dogus contains a set of dogus which alters the states of the dogus in the blueprint this mask is applied on.
	// The names and target states of all dogus must not be empty.
	Dogus []MaskTargetDogu `json:"dogus"`
}

// Validate checks the structure and data of a blueprint mask and returns an error if there are any problems
func (blueprintMask *BlueprintMaskV1) Validate() error {
	if blueprintMask.API == "" {
		return errors.Errorf("could not validate mask API, mask API must not be empty")
	}
	if blueprintMask.ID == "" {
		return errors.Errorf("could not validate mask ID, mask ID must not be empty")
	}

	err := blueprintMask.validateDogus()
	if err != nil {
		return err
	}

	err = blueprintMask.validateDoguUniqueness()
	if err != nil {
		return err
	}

	return nil
}

func (blueprintMask *BlueprintMaskV1) validateDogus() error {
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
func (blueprintMask *BlueprintMaskV1) validateDoguUniqueness() error {
	seenDogu := make(map[string]bool)

	for _, dogu := range blueprintMask.Dogus {
		_, seen := seenDogu[dogu.Name]
		if seen {
			return errors.Errorf("could not validate blueprint mask, there is at least one duplicate for this dogu: %s", dogu.Name)
		}
		seenDogu[dogu.Name] = true
	}

	return nil
}
