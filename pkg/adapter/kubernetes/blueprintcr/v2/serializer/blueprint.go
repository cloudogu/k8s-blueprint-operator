package serializer

import (
	"errors"
	"fmt"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToBlueprintDTO(blueprint domain.EffectiveBlueprint) crd.BlueprintManifest {
	return crd.BlueprintManifest{
		Dogus:      ConvertToDoguDTOs(blueprint.Dogus),
		Components: ConvertToComponentDTOs(blueprint.Components),
		Config:     ConvertToConfigDTO(blueprint.Config),
	}
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
	configDomain := ConvertToConfigDomain(blueprint.Config)
	return domain.Blueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     configDomain,
	}, nil
}

func ConvertToEffectiveBlueprintDomain(blueprint *crd.BlueprintManifest) (domain.EffectiveBlueprint, error) {
	if blueprint == nil {
		return domain.EffectiveBlueprint{}, nil
	}
	convertedDogus, doguErr := ConvertDogus(blueprint.Dogus)
	convertedComponents, compErr := ConvertComponents(blueprint.Components)

	err := errors.Join(doguErr, compErr)
	if err != nil {
		return domain.EffectiveBlueprint{}, fmt.Errorf("cannot deserialize effective blueprint: %w", err)
	}
	return domain.EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDomain(blueprint.Config),
	}, nil
}

func ConvertToBlueprintMaskDomain(mask *crd.BlueprintMask) (domain.BlueprintMask, error) {
	if mask == nil {
		return domain.BlueprintMask{}, nil
	}

	convertedDogus, err := ConvertMaskDogus(mask.Dogus)

	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("cannot deserialize blueprint mask: %w", err)
	}
	return domain.BlueprintMask{
		Dogus: convertedDogus,
	}, nil
}
