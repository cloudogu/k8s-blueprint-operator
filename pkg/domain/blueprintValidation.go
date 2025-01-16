package domain

import (
	"errors"
	"fmt"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

type blueprintValidator struct {
	blueprint bpv2.Blueprint
}

func newBlueprintValidator(blueprint bpv2.Blueprint) *blueprintValidator {
	return &blueprintValidator{blueprint: blueprint}
}

// validate checks the structure and data of the blueprint statically and returns an error if there are any problems
func (validator *blueprintValidator) validate() error {
	errorList := []error{
		validator.validateDogus(),
		validator.validateDoguUniqueness(),
		validator.validateComponents(),
		validator.validateComponentUniqueness(),
		newConfigValidator(validator.blueprint.Config).validate(),
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint is invalid: %w", err)
	}
	return err
}

func (validator *blueprintValidator) validateDogus() error {
	errorList := util.Map(validator.blueprint.Dogus, func(dogu bpv2.Dogu) error { return NewDoguValidator(dogu).validate() })
	return errors.Join(errorList...)
}

// validateDoguUniqueness checks if dogus exist twice in the blueprint and returns an error if it's so.
func (validator *blueprintValidator) validateDoguUniqueness() error {
	doguNames := util.Map(validator.blueprint.Dogus, func(dogu bpv2.Dogu) cescommons.SimpleName { return dogu.Name.SimpleName })
	duplicates := util.GetDuplicates(doguNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate dogus: %v", duplicates)
	}
	return nil
}

func (validator *blueprintValidator) validateComponents() error {
	errorList := util.Map(validator.blueprint.Components, func(component bpv2.Component) error { return NewComponentValidator(component).validate() })
	return errors.Join(errorList...)
}

// validateComponentUniqueness checks if components exist twice in the blueprint and returns an error if it's so.
func (validator *blueprintValidator) validateComponentUniqueness() error {
	componentNames := util.Map(validator.blueprint.Components, func(component bpv2.Component) bpv2.SimpleComponentName { return component.Name.SimpleName })
	duplicates := util.GetDuplicates(componentNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate components: %v", duplicates)
	}
	return nil
}
