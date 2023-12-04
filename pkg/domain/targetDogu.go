package domain

import (
	"errors"
	"fmt"
	"slices"
)

// TargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type TargetDogu struct {
	// Namespace defines the namespace of the dogu, e.g. "official". Must not be empty.
	Namespace string
	// Name defines the name of the dogu excluding the namespace, e.g. "nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState `json:"targetState"`
}

func (dogu TargetDogu) GetQualifiedName() string {
	return fmt.Sprintf("%s/%s", dogu.Namespace, dogu.Name)
}

// validate checks if the TargetDogu is semantically correct.
func (dogu TargetDogu) validate() error {
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
	if dogu.TargetState != TargetStateAbsent && dogu.Version == "" {
		errorList = append(errorList, fmt.Errorf("dogu field Version must not be empty: %s", dogu.GetQualifiedName()))
	}
	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("dogu is invalid: %w", err)
	}
	return err
}
