package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"golang.org/x/exp/maps"
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
	Dogus []Dogu
	// Components contains a set of exact components versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []Component
	// Config contains all config entries to set via blueprint.
	Config Config
}

// Validate checks the structure and data of the blueprint statically and returns an error if there are any problems
func (blueprint *Blueprint) Validate() error {
	errorList := []error{
		blueprint.validateDogus(),
		blueprint.validateDoguUniqueness(),
		blueprint.validateComponents(),
		blueprint.validateComponentUniqueness(),
		blueprint.Config.Global.validate(),
	}
	doguConfigErrors := util.Map(maps.Values(blueprint.Config.Dogus), CombinedDoguConfig.validate)
	errorList = append(errorList, doguConfigErrors...)

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("blueprint is invalid: %w", err)
	}
	return err
}

func (blueprint *Blueprint) validateDogus() error {
	errorList := util.Map(blueprint.Dogus, func(dogu Dogu) error { return dogu.validate() })
	return errors.Join(errorList...)
}

// validateDoguUniqueness checks if dogus exist twice in the blueprint and returns an error if it's so.
func (blueprint *Blueprint) validateDoguUniqueness() error {
	doguNames := util.Map(blueprint.Dogus, func(dogu Dogu) common.SimpleDoguName { return dogu.Name.Name })
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
	componentNames := util.Map(blueprint.Components, func(component Component) common.SimpleComponentName { return component.Name.Name })
	duplicates := util.GetDuplicates(componentNames)
	if len(duplicates) != 0 {
		return fmt.Errorf("there are duplicate components: %v", duplicates)
	}
	return nil
}
