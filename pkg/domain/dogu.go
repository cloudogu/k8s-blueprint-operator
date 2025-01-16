package domain

import (
	"errors"
	"fmt"
	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"k8s.io/apimachinery/pkg/api/resource"
	"slices"
)

type DoguValidator struct {
	dogu bpv2.Dogu
}

func NewDoguValidator(dogu bpv2.Dogu) *DoguValidator {
	return &DoguValidator{dogu}
}

// validate checks if the Dogu is semantically correct.
func (doguValidator *DoguValidator) validate() error {
	var errorList []error
	if !slices.Contains(PossibleTargetStates, doguValidator.dogu.TargetState) {
		errorList = append(errorList, fmt.Errorf("dogu target state is invalid: %s", dogu.Name))
	}
	emptyVersion := core.Version{}
	if doguValidator.dogu.TargetState != TargetStateAbsent && doguValidator.dogu.Version == emptyVersion {
		errorList = append(errorList, fmt.Errorf("dogu version must not be empty: %s", doguValidator.dogu.Name))
	}

	// Storage is usually expressed in Binary SI. Using Decimal SI can cause problems because sizes will be
	// rounded up (longhorn does this in volume resize).
	minVolumeSize := doguValidator.dogu.MinVolumeSize
	if minVolumeSize != nil && !minVolumeSize.IsZero() && minVolumeSize.Format != resource.BinarySI {
		errorList = append(errorList, fmt.Errorf("dogu minimum volume size is not in Binary SI (\"Mi\" or \"Gi\"): %s", dogu.Name))
	}

	// Nginx only supports quantities in Decimal SI. This check can be removed if the dogu-operator implements an abstraction for the body size.
	maxBodySize := doguValidator.dogu.ReverseProxyConfig.MaxBodySize
	if maxBodySize != nil && !maxBodySize.IsZero() && maxBodySize.Format != resource.DecimalSI {
		errorList = append(errorList, fmt.Errorf("dogu proxy body size is not in Decimal SI (\"M\" or \"G\"): %s", doguValidator.dogu.Name))
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("dogu is invalid: %w", err)
	}
	return err
}

func FindDoguByName(dogus []bpv2.Dogu, name cescommons.SimpleName) (Dogu, bool) {
	for _, dogu := range dogus {
		if dogu.Name.SimpleName == name {
			return dogu, true
		}
	}
	return Dogu{}, false
}
