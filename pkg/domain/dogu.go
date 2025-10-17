package domain

import (
	"errors"
	"fmt"
	"strings"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Dogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type Dogu struct {
	// Name defines the name of the dogu, e.g. "official/postgresql"
	Name cescommons.QualifiedName
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version *core.Version
	// Absent defines if the dogu should be absent in the ecosystem. Defaults to false.
	Absent bool
	// MinVolumeSize is the minimum storage of the dogu. 0 indicates that the default size should be set.
	// Reducing this value below the actual volume size has no impact as we do not support downsizing.
	MinVolumeSize *ecosystem.VolumeSize
	// ReverseProxyConfig defines configuration for the ecosystem reverse proxy. This field is optional.
	ReverseProxyConfig ecosystem.ReverseProxyConfig
	// AdditionalMounts provides the possibility to mount additional data into the dogu.
	AdditionalMounts []ecosystem.AdditionalMount
}

// validate checks if the Dogu is semantically correct.
func (dogu Dogu) validate() error {
	var errorList []error

	emptyVersion := core.Version{}
	if !dogu.Absent && (dogu.Version == nil || *dogu.Version == emptyVersion) {
		errorList = append(errorList, fmt.Errorf("dogu version must not be empty: %s", dogu.Name))
	}
	// minVolumeSize is already checked while unmarshalling json/yaml

	// Nginx only supports quantities in Decimal SI. This check can be removed if the dogu-operator implements an abstraction for the body size.
	maxBodySize := dogu.ReverseProxyConfig.MaxBodySize
	if maxBodySize != nil && !maxBodySize.IsZero() && maxBodySize.Format != resource.DecimalSI {
		errorList = append(errorList, fmt.Errorf("dogu proxy body size is not in Decimal SI (\"M\" or \"G\"): %s", dogu.Name))
	}

	for _, mount := range dogu.AdditionalMounts {
		if mount.SourceType != ecosystem.DataSourceConfigMap && mount.SourceType != ecosystem.DataSourceSecret {
			errorList = append(errorList, fmt.Errorf(
				"dogu additional mounts sourceType must be one of '%s', '%s': %s",
				ecosystem.DataSourceConfigMap,
				ecosystem.DataSourceSecret,
				dogu.Name,
			))
		}
		if strings.HasPrefix(mount.Subfolder, "/") {
			errorList = append(errorList, fmt.Errorf("dogu additional mounts Subfolder must be a relative path : %s", dogu.Name))
		}
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
