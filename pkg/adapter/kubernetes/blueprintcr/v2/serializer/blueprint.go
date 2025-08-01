package serializer

import (
	"errors"
	"fmt"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToBlueprintDTO(blueprint domain.EffectiveBlueprint) (crd.BlueprintManifest, error) {
	var errorList []error
	convertedDogus := ConvertToDoguDTOs(blueprint.Dogus)
	convertedComponents, componentError := ConvertToComponentDTOs(blueprint.Components)
	errorList = append(errorList, componentError)

	err := errors.Join(errorList...)
	if err != nil {
		return crd.BlueprintManifest{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return crd.BlueprintManifest{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDTO(blueprint.Config),
	}, nil
}

func ConvertToBlueprintDomain(blueprint crd.BlueprintManifest) (domain.Blueprint, error) {
	convertedDogus, doguErr := ConvertDogus(blueprint.Dogus)
	convertedComponents, compErr := ConvertComponents(blueprint.Components)

	err := errors.Join(doguErr, compErr)
	if err != nil {
		return domain.Blueprint{}, &domain.InvalidBlueprintError{
			WrappedError: err,
			Message:      "cannot deserialize blueprint",
		}
	}
	return domain.Blueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDomain(blueprint.Config),
	}, nil
}

func ConvertToEffectiveBlueprintDomain(blueprint crd.BlueprintManifest) (domain.EffectiveBlueprint, error) {
	convertedDogus, doguErr := ConvertDogus(blueprint.Dogus)
	convertedComponents, compErr := ConvertComponents(blueprint.Components)

	err := errors.Join(doguErr, compErr)
	if err != nil {
		return domain.EffectiveBlueprint{}, fmt.Errorf("syntax of blueprintV2 is not correct: %w", err)
	}
	return domain.EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDomain(blueprint.Config),
	}, nil
}

func ConvertToBlueprintMaskDomain(mask crd.BlueprintMask) (domain.BlueprintMask, error) {
	convertedDogus, err := ConvertMaskDogus(mask.Dogus)

	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("cannot deserialize blueprint mask: %w", err)
	}
	return domain.BlueprintMask{
		Dogus: convertedDogus,
	}, nil
}
