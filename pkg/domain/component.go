package domain

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

// Component represents a CES component (e.g. operators), its version, and the installation state in which it is supposed to be
// after a blueprint was applied.
type Component struct {
	// Name defines the name and namespace of the component. Must not be empty.
	Name common.QualifiedComponentName
	// Version defines the version of the package that is to be installed. Must not be empty if the targetState is
	// "present"; otherwise it is optional and is not going to be interpreted.
	Version *semver.Version
	// Absent defines if the dogu should be absent in the ecosystem. Defaults to false.
	Absent bool
	// DeployConfig defines generic properties for the component. This field is optional.
	DeployConfig ecosystem.DeployConfig
}

// Validate checks if the component is semantically correct.
func (component *Component) Validate() error {
	nameError := component.Name.Validate()

	var versionErr error
	if !component.Absent && component.Version == nil {
		versionErr = fmt.Errorf("version of component %q must not be empty", component.Name)
	}

	return errors.Join(versionErr, nameError)
}
