package domain

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/pkg/errors"
)

// Component represents a CES component (e.g. operators), its version, and the installation state in which it is supposed to be
// after a blueprint was applied.
type Component struct {
	// Name defines the name of the component. Must not be empty.
	Name string
	// Version defines the version of the package that is to be installed. Must not be empty if the targetState is
	// "present"; otherwise it is optional and is not going to be interpreted.
	Version core.Version
	// TargetState defines a state of installation of this package. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState
}

// Validate checks if the component is semantically correct.
func (component *Component) Validate() error {
	if component.Name == "" {
		return errors.Errorf("component name must not be empty: %+v", component)
	}
	emptyVersion := core.Version{}
	if component.TargetState != TargetStateAbsent && component.Version == emptyVersion {
		return errors.Errorf("component version must not be empty: %s", component.Version.Raw)
	}

	return nil
}
