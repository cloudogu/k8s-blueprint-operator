package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"slices"
)

// Dogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type Dogu struct {
	// Name defines the name of the dogu, e.g. "official/postgresql"
	Name common.QualifiedDoguName
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version core.Version
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState TargetState
}

// validate checks if the Dogu is semantically correct.
func (dogu Dogu) validate() error {
	var errorList []error
	errorList = append(errorList, dogu.Name.Validate())
	if !slices.Contains(PossibleTargetStates, dogu.TargetState) {
		errorList = append(errorList, fmt.Errorf("dogu target state is invalid: %s", dogu.Name))
	}
	emptyVersion := core.Version{}
	if dogu.TargetState != TargetStateAbsent && dogu.Version == emptyVersion {
		errorList = append(errorList, fmt.Errorf("dogu version must not be empty: %s", dogu.Name))
	}
	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("dogu is invalid: %w", err)
	}
	return err
}

func FindDoguByName(dogus []Dogu, name common.SimpleDoguName) (Dogu, error) {
	for _, dogu := range dogus {
		if dogu.Name.SimpleName == name {
			return dogu, nil
		}
	}
	return Dogu{}, fmt.Errorf("could not find dogu '%s'", name)
}
