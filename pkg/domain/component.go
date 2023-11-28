package domain

import "github.com/pkg/errors"

// Component represents a CES component (e.g. operators), its version, and the installation state in which it is supposed to be
// after a blueprint was applied.
type Component struct {
	// Name defines the name of the component. Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the package that is to be installed. Must not be empty if the targetState is
	// "present"; otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this package. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState `json:"targetState"`
}

func (component *Component) validate() error {
	if component.Name == "" {
		return errors.Errorf("could not validate blueprint, component name must not be empty: %s", component)
	}
	if component.TargetState != TargetStateAbsent && component.Version == "" {
		return errors.Errorf("could not validate blueprint, component version must not be empty: %s", component)
	}

	return nil
}
