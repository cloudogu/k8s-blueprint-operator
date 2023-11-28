package domain

import "github.com/pkg/errors"

// TargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type TargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState `json:"targetState"`
}

func (dogu TargetDogu) validate() error {
	if dogu.Name == "" {
		return errors.Errorf("could not validate blueprint, dogu field Name must not be empty: %s", dogu)
	}
	if dogu.TargetState != TargetStateAbsent && dogu.Version == "" {
		return errors.Errorf("could not validate blueprint, dogu field Version must not be empty: %s", dogu)
	}

	return nil
}
