package domain

import "github.com/pkg/errors"

// MaskTargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied for a blueprintMask.
type MaskTargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty. If you set another namespace than in the normal blueprint, a
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. This version is optional and overrides
	// the version of the dogu from the blueprint.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState `json:"targetState"`
}

func (dogu MaskTargetDogu) validate() error {
	if dogu.Name == "" {
		return errors.Errorf("could not validate blueprint mask, dogu field Name must not be empty: %s", dogu)
	}
	if dogu.TargetState.String() == "" {
		return errors.Errorf("could not validate dogu, dogu target state must not be empty: %s", dogu)
	}

	return nil
}
