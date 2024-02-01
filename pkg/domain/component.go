package domain

import (
	"errors"
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
)

// Component represents a CES component (e.g. operators), its version, and the installation state in which it is supposed to be
// after a blueprint was applied.
type Component struct {
	// Name defines the name of the component. Must not be empty.
	Name string
	// DistributionNamespace is part of the address under which the component will be obtained. This namespace does NOT
	// to be confused with the K8s cluster namespace.
	DistributionNamespace string
	// Version defines the version of the package that is to be installed. Must not be empty if the targetState is
	// "present"; otherwise it is optional and is not going to be interpreted.
	Version core.Version
	// TargetState defines a state of installation of this package. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState
}

// Validate checks if the component is semantically correct.
func (component *Component) Validate() error {
	if component.Name == "" {
		return fmt.Errorf("component name must not be empty: %+v", component)
	}

	emptyVersion := core.Version{}
	if component.TargetState == TargetStatePresent {
		var versionErr error
		if component.Version == emptyVersion {
			versionErr = fmt.Errorf("version of component %q must not be empty: %s", component.Name, component.Version.Raw)
		}
		var namespaceErr error
		if component.DistributionNamespace == "" {
			namespaceErr = fmt.Errorf("distribution namespace of component %q must not be empty", component.Name)
		}
		return errors.Join(versionErr, namespaceErr)
	}

	return nil
}
