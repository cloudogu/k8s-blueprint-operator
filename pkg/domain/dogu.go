package domain

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/api/resource"
	"slices"
)

type DataSourceType string

//goland:noinspection GoUnusedConst
const (
	// DataSourceConfigMap mounts a config map as a data source.
	DataSourceConfigMap DataSourceType = "ConfigMap"
	// DataSourceSecret mounts a secret as a data source.
	DataSourceSecret DataSourceType = "Secret"
)

// Dogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type Dogu struct {
	// Name defines the name of the dogu, e.g. "official/postgresql"
	Name cescommons.QualifiedName
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
	// AdditionalMounts provides the possibility to mount additional data into the dogu.
	AdditionalMounts []AdditionalMount
}

// AdditionalMount is a description of what data should be mounted to a specific Dogu volume (already defined in dogu.json).
type AdditionalMount struct {
	// SourceType defines where the data is coming from.
	// Valid options are:
	//   ConfigMap - data stored in a kubernetes ConfigMap.
	//   Secret - data stored in a kubernetes Secret.
	SourceType DataSourceType
	// Name is the name of the data source.
	Name string
	// Volume is the name of the volume to which the data should be mounted. It is defined in the respective dogu.json.
	Volume string
	// Subfolder defines a subfolder in which the data should be put within the volume.
	// +optional
	Subfolder string
}

// validate checks if the Dogu is semantically correct.
func (dogu Dogu) validate() error {
	var errorList []error
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

func FindDoguByName(dogus []Dogu, name cescommons.SimpleName) (Dogu, bool) {
	for _, dogu := range dogus {
		if dogu.Name.SimpleName == name {
			return dogu, true
		}
	}
	return Dogu{}, false
}
