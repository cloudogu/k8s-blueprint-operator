package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

// BlueprintMask describes an abstraction of CES components that should alter a blueprint definition before
// applying it to a CES system via a blueprint upgrade. The blueprint mask does not change the blueprint
// itself, but is applied to the information in it to generate a new, effective blueprint.
type BlueprintMask struct {
	// Dogus contains a set of dogus which alters the states of the dogus in the blueprint this mask is applied on.
	// The names and target states of all dogus must not be empty.
	Dogus []MaskDogu
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
	errorList := util.MapWithFunction(blueprintMask.Dogus, func(maskDogu MaskDogu) error { return maskDogu.validate() })
	return errors.Join(errorList...)
}

// validateDoguUniqueness checks if dogus exist twice in the blueprint and returns an error if it's so.
func (blueprintMask *BlueprintMask) validateDoguUniqueness() error {
	doguNames := util.MapWithFunction(blueprintMask.Dogus, func(dogu MaskDogu) string { return dogu.Name })
	duplicates := util.GetDuplicates(doguNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate dogus: %v", duplicates)
	}
	return nil
}

func (blueprintMask *BlueprintMask) FindDoguByName(name string) (MaskDogu, error) {
	for doguIndex, dogu := range blueprintMask.Dogus {
		if dogu.Name == name {
			return blueprintMask.Dogus[doguIndex], nil
		}
	}
	return MaskDogu{}, fmt.Errorf("could not find dogu name %q in blueprint", name)
}
