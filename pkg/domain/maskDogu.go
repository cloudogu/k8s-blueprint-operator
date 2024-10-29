package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"slices"
)

// MaskDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied for a blueprintMask.
type MaskDogu struct {
	// Name is the qualified name of the dogu.
	Name common.QualifiedDoguName
	// Version defines the version of the dogu that is to be installed. This version is optional and overrides
	// the version of the dogu from the blueprint.
	Version core.Version
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState
}

func (dogu MaskDogu) validate() error {
	var errorList []error
	errorList = append(errorList, dogu.Name.Validate())

	if !slices.Contains(PossibleTargetStates, dogu.TargetState) {
		errorList = append(errorList, fmt.Errorf("dogu mask is invalid: dogu target state is invalid: %s", dogu.Name))
	}

	return errors.Join(errorList...)
}
