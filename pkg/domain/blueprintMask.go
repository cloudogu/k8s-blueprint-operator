package domain

import (
	"errors"
	"fmt"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

type blueprintMaskValidator struct {
	blueprintMask bpv2.BlueprintMask
}

func newBlueprintMaskValidator(blueprintMask bpv2.BlueprintMask) *blueprintMaskValidator {
	return &blueprintMaskValidator{blueprintMask: blueprintMask}
}

// Validate checks the structure and data of a blueprint mask and returns an error if there are any problems
func (validator *blueprintMaskValidator) validate() error {
	errorList := []error{
		validator.validateDogus(),
		validator.validateDoguUniqueness(),
	}
	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint mask is invalid: %w", err)
	}
	return err
}

func (validator *blueprintMaskValidator) validateDogus() error {
	errorList := util.Map(validator.blueprintMask.Dogus, func(maskDogu bpv2.MaskDogu) error {
		return validateMask(maskDogu)
	})
	return errors.Join(errorList...)
}

// validateDoguUniqueness checks if dogus exist twice in the blueprint and returns an error if it's so.
func (validator *blueprintMaskValidator) validateDoguUniqueness() error {
	doguNames := util.Map(validator.blueprintMask.Dogus, func(dogu bpv2.MaskDogu) cescommons.SimpleName { return dogu.Name.SimpleName })
	duplicates := util.GetDuplicates(doguNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate dogus: %v", duplicates)
	}
	return nil
}

func (validator *blueprintMaskValidator) FindDoguByName(name cescommons.SimpleName) (bpv2.MaskDogu, error) {
	for doguIndex, dogu := range validator.blueprintMask.Dogus {
		if dogu.Name.SimpleName == name {
			return validator.blueprintMask.Dogus[doguIndex], nil
		}
	}
	return bpv2.MaskDogu{}, fmt.Errorf("could not find dogu name %q in blueprint", name)
}
