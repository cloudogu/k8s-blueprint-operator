package domain

import (
	"errors"
	"fmt"
	"slices"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
)

func validateMask(dogu bpv2.MaskDogu) error {
	var errorList []error
	errorList = append(errorList, dogu.Name.Validate())

	if !slices.Contains(bpv2.PossibleTargetStates, dogu.TargetState) {
		errorList = append(errorList, fmt.Errorf("dogu mask is invalid: dogu target state is invalid: %s", dogu.Name))
	}

	return errors.Join(errorList...)
}
