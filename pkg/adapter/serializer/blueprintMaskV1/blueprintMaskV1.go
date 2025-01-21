package blueprintMaskV1

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	bpmask "github.com/cloudogu/k8s-blueprint-lib/json/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-lib/json/bpcore"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

func ConvertToBlueprintMaskV1(spec domain.BlueprintMask) (bpmask.BlueprintMaskV1, error) {
	var errorList []error
	convertedDogus := util.Map(spec.Dogus, func(dogu domain.MaskDogu) bpmask.MaskTargetDogu {
		newState, err := serializer.ToSerializerTargetState(dogu.TargetState)
		errorList = append(errorList, err)
		return bpmask.MaskTargetDogu{
			Name:        dogu.Name.String(),
			Version:     dogu.Version.Raw,
			TargetState: newState,
		}
	})

	err := errors.Join(errorList...)
	if err != nil {
		return bpmask.BlueprintMaskV1{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return bpmask.BlueprintMaskV1{
		GeneralBlueprintMask: bpcore.GeneralBlueprintMask{API: bpcore.MaskV1},
		Dogus:                convertedDogus,
	}, nil
}

func convertToBlueprintMask(blueprintMask bpmask.BlueprintMaskV1) (domain.BlueprintMask, error) {
	switch blueprintMask.API {
	case bpcore.MaskV1:
	default:
		return domain.BlueprintMask{}, fmt.Errorf("unsupported Blueprint Mask API Version: %s", blueprintMask.API)
	}

	convertedDogus, err := convertMaskDogus(blueprintMask.Dogus)
	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("syntax of blueprintMaskV1 is not correct: %w", err)
	}

	return domain.BlueprintMask{Dogus: convertedDogus}, nil
}

func convertMaskDogus(dogus []bpmask.MaskTargetDogu) ([]domain.MaskDogu, error) {
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
