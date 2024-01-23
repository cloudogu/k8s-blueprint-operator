package serializer

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
)

// TargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied.
type TargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty.
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. Must not be empty if the targetState is "present";
	// otherwise it is optional and is not going to be interpreted.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
}

func ConvertDogus(dogus []TargetDogu) ([]domain.Dogu, error) {
	var convertedDogus []domain.Dogu
	var errorList []error

	for _, dogu := range dogus {
		doguNamespace, doguName, err := SplitDoguName(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		newState, err := ToDomainTargetState(dogu.TargetState)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		var version core.Version
		if dogu.Version != "" {
			version, err = core.ParseVersion(dogu.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of target dogu %q: %w", dogu.Name, err))
				continue
			}
		}
		convertedDogus = append(convertedDogus, domain.Dogu{
			Namespace:   doguNamespace,
			Name:        doguName,
			Version:     version,
			TargetState: newState,
		})
	}

	err := errors.Join(errorList...)
	if err != nil {
		return convertedDogus, fmt.Errorf("cannot convert blueprint dogus: %w", err)
	}

	return convertedDogus, err
}

func ConvertToDoguDTOs(dogus []domain.Dogu) ([]TargetDogu, error) {
	var errorList []error
	converted := util.MapWithFunction(dogus, func(dogu domain.Dogu) TargetDogu {
		newState, err := ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)
		return TargetDogu{
			Name:        dogu.GetQualifiedName(),
			Version:     dogu.Version.Raw,
			TargetState: newState,
		}
	})
	return converted, errors.Join(errorList...)
}
