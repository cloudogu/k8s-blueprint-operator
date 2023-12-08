package domain

import (
	"errors"
	"fmt"
	"slices"
)

// MaskTargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied for a blueprintMask.
type MaskTargetDogu struct {
	// Namespace defines the namespace of the dogu, e.g. "official". Must not be empty.
	Namespace string
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty. If you set another namespace than in the normal blueprint, a
	Name string
	// Version defines the version of the dogu that is to be installed. This version is optional and overrides
	// the version of the dogu from the blueprint.
	Version string
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState
}

func (dogu MaskTargetDogu) GetQualifiedName() string {
	return fmt.Sprintf("%s/%s", dogu.Namespace, dogu.Name)
}

func (dogu MaskTargetDogu) validate() error {
	var errorList []error
	if dogu.Namespace == "" {
		errorList = append(errorList, fmt.Errorf("dogu field Namespace must not be empty: %s", dogu.GetQualifiedName()))
	}
	if dogu.Name == "" {
		errorList = append(errorList, fmt.Errorf("dogu field Name must not be empty: %s", dogu.GetQualifiedName()))
	}
	if !slices.Contains(PossibleTargetStates, dogu.TargetState) {
		errorList = append(errorList, fmt.Errorf("dogu target state is invalid: %s", dogu.GetQualifiedName()))
	}
	//TODO: parse version and validate the real format with cesapp-lib
	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("dogu mask is invalid: %w", err)
	}
	return err
}
