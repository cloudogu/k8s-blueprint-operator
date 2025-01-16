package domain

import (
	"errors"
	"fmt"
)

type ComponentValidator struct {
	component *bpv2.Component
}

func NewComponentValidator(component bpv2.Component) *ComponentValidator {
	return &ComponentValidator{component: component}
}

// Validate checks if the component is semantically correct.
func (compValidator *ComponentValidator) validate() error {
	nameError := compValidator.component.Name.Validate()

	var versionErr error
	if compValidator.component.TargetState == TargetStatePresent {
		if compValidator.component.Version == nil {
			versionErr = fmt.Errorf("version of component %q must not be empty", compValidator.component.Name)
		}
	}

	return errors.Join(versionErr, nameError)
}
