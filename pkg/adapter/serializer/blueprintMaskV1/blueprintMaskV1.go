package blueprintMaskV1

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

// BlueprintMaskV1 describes an abstraction of CES components that should alter a blueprint definition before
// applying it to a CES system via a blueprint upgrade. The blueprint mask should not change the blueprint JSON file
// itself, but is applied to the information in it to generate a new effective blueprint.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint mask API version.
type BlueprintMaskV1 struct {
	serializer.GeneralBlueprintMask
	// ID is the unique name of the set over all components. This blueprint mask ID should be used to distinguish
	// from similar blueprint masks between humans in an easy way. Must not be empty.
	ID string `json:"blueprintMaskId"`
	// Dogus contains a set of dogus which alters the states of the dogus in the blueprint this mask is applied on.
	// The names and target states of all dogus must not be empty.
	Dogus []MaskTargetDogu `json:"dogus"`
}

// MaskTargetDogu defines a Dogu, its version, and the installation state in which it is supposed to be after a blueprint
// was applied for a blueprintMask.
type MaskTargetDogu struct {
	// Name defines the name of the dogu including its namespace, f. i. "official/nginx". Must not be empty. If you set another namespace than in the normal blueprint, a
	Name string `json:"name"`
	// Version defines the version of the dogu that is to be installed. This version is optional and overrides
	// the version of the dogu from the blueprint.
	Version string `json:"version"`
	// TargetState defines a state of installation of this dogu. Optional field, but defaults to "TargetStatePresent"
	TargetState string `json:"targetState"`
}

func ConvertToBlueprintMaskV1(spec domain.BlueprintMask) (BlueprintMaskV1, error) {
	var errorList []error
	convertedDogus := util.Map(spec.Dogus, func(dogu domain.MaskDogu) MaskTargetDogu {
		newState, err := serializer.ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)
		return MaskTargetDogu{
			Name:        dogu.Name.String(),
			Version:     dogu.Version.Raw,
			TargetState: newState,
		}
	})

	err := errors.Join(errorList...)
	if err != nil {
		return BlueprintMaskV1{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return BlueprintMaskV1{
		GeneralBlueprintMask: serializer.GeneralBlueprintMask{API: serializer.BlueprintMaskAPIV1},
		Dogus:                convertedDogus,
	}, nil
}

func convertToBlueprintMask(blueprintMask BlueprintMaskV1) (domain.BlueprintMask, error) {
	switch blueprintMask.API {
	case serializer.BlueprintMaskAPIV1:
	default:
		return domain.BlueprintMask{}, fmt.Errorf("unsupported Blueprint Mask API Version: %s", blueprintMask.API)
	}
	convertedDogus, err := convertMaskDogus(blueprintMask.Dogus)
	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("syntax of blueprintMaskV1 is not correct: %w", err)
	}
	return domain.BlueprintMask{Dogus: convertedDogus}, nil
}

func convertMaskDogus(dogus []MaskTargetDogu) ([]domain.MaskDogu, error) {
	var convertedDogus []domain.MaskDogu
	var errorList []error

	for _, dogu := range dogus {
		doguName, err := cescommons.QualifiedNameFromString(dogu.Name)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		state, err := serializer.ToDomainTargetState(dogu.TargetState)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		var version core.Version
		if dogu.Version != "" {
			version, err = core.ParseVersion(dogu.Version)
			if err != nil {
				errorList = append(errorList, fmt.Errorf("could not parse version of MaskTargetDogu: %w", err))
				continue
			}
		}
		convertedDogus = append(convertedDogus, domain.MaskDogu{
			Name:        doguName,
			Version:     version,
			TargetState: state,
		})
	}

	err := errors.Join(errorList...)
	return convertedDogus, err
}
