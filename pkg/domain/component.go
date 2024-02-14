package domain

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

// Component represents a CES component (e.g. operators), its version, and the installation state in which it is supposed to be
// after a blueprint was applied.
type Component struct {
	// TODO Add ComponentConfig for CRD fields like deployNamespace or helmValuesOverwrite.
	// Name defines the name and namespace of the component. Must not be empty.
	Name common.QualifiedComponentName
	// Version defines the version of the package that is to be installed. Must not be empty if the targetState is
	// "present"; otherwise it is optional and is not going to be interpreted.
	Version *semver.Version
	// TargetState defines a state of installation of this package. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState
}

// Validate checks if the component is semantically correct.
func (component *Component) Validate() error {
	nameError := component.Name.Validate()

	var versionErr error
	if component.TargetState == TargetStatePresent {
		if component.Version == nil {
			versionErr = fmt.Errorf("version of component %q must not be empty", component.Name)
		}
	}

	return errors.Join(versionErr, nameError)
}
