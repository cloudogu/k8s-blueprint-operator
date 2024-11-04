package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/api/resource"
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
	// MinVolumeSize is the minimum storage of the dogu. This field is optional and can be nil to indicate that no
	// storage is needed.
	MinVolumeSize *ecosystem.VolumeSize
	// ReverseProxyConfig defines configuration for the ecosystem reverse proxy. This field is optional.
	ReverseProxyConfig ecosystem.ReverseProxyConfig
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

	// Storage is usually expressed in Binary SI. Using Decimal SI can cause problems because sizes will be
	// rounded up (longhorn does this in volume resize).
	minVolumeSize := dogu.MinVolumeSize
	if minVolumeSize != nil && !minVolumeSize.IsZero() && minVolumeSize.Format != resource.BinarySI {
		errorList = append(errorList, fmt.Errorf("dogu minimum volume size is not in Binary SI (\"Mi\" or \"Gi\"): %s", dogu.Name))
	}

	// Nginx only supports quantities in Decimal SI. This check can be removed if the dogu-operator implements an abstraction for the body size.
	maxBodySize := dogu.ReverseProxyConfig.MaxBodySize
	if maxBodySize != nil && !maxBodySize.IsZero() && maxBodySize.Format != resource.DecimalSI {
		errorList = append(errorList, fmt.Errorf("dogu proxy body size is not in Decimal SI (\"M\" or \"G\"): %s", dogu.Name))
	}

	err := errors.Join(errorList...)
	if err != nil {
		err = fmt.Errorf("dogu is invalid: %w", err)
	}
	return err
}

func FindDoguByName(dogus []Dogu, name common.SimpleDoguName) (Dogu, bool) {
	for _, dogu := range dogus {
		if dogu.Name.SimpleName == name {
			return dogu, true
		}
	}
	return Dogu{}, false
}
