package serializer

import (
	"errors"
	"fmt"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func SerializeBlueprintAndMask(blueprintSpec *domain.BlueprintSpec, manifest crd.BlueprintManifest, maskManifest *crd.BlueprintMaskManifest) error {
	blueprint, blueprintErr := ConvertToBlueprintDomain(manifest)
	blueprintMask, maskErr := ConvertToBlueprintMaskDomain(maskManifest)
	serializationErr := errors.Join(blueprintErr, maskErr)
	if serializationErr != nil {
		return serializationErr
	}

	blueprintSpec.Blueprint = blueprint
	blueprintSpec.BlueprintMask = blueprintMask
	return nil
}

func ConvertBlueprintStatus(blueprintCR *crd.Blueprint) (domain.EffectiveBlueprint, error) {
	var effectiveBlueprint domain.EffectiveBlueprint
	var err error
	if blueprintCR.Status != nil {
		effectiveBlueprint, err = ConvertToEffectiveBlueprintDomain(blueprintCR.Status.EffectiveBlueprint)
		if err != nil {
			return domain.EffectiveBlueprint{}, err
		}
	}
	return effectiveBlueprint, nil
}

func ConvertToBlueprintDTO(blueprint domain.EffectiveBlueprint) crd.BlueprintManifest {
	return crd.BlueprintManifest{
		Dogus:  ConvertToDoguDTOs(blueprint.Dogus),
		Config: ConvertToConfigDTO(blueprint.Config),
	}
}

func ConvertToBlueprintDomain(blueprint crd.BlueprintManifest) (domain.Blueprint, error) {
	convertedDogus, err := ConvertDogus(blueprint.Dogus)
	if err != nil {
		return domain.Blueprint{}, &domain.InvalidBlueprintError{
			WrappedError: err,
			Message:      "cannot deserialize blueprint",
		}
	}
	configDomain := ConvertToConfigDomain(blueprint.Config)
	return domain.Blueprint{
		Dogus:  convertedDogus,
		Config: configDomain,
	}, nil
}

func ConvertToEffectiveBlueprintDomain(blueprint *crd.BlueprintManifest) (domain.EffectiveBlueprint, error) {
	if blueprint == nil {
		return domain.EffectiveBlueprint{}, nil
	}
	convertedDogus, err := ConvertDogus(blueprint.Dogus)
	if err != nil {
		return domain.EffectiveBlueprint{}, fmt.Errorf("cannot deserialize effective blueprint: %w", err)
	}
	return domain.EffectiveBlueprint{
		Dogus:  convertedDogus,
		Config: ConvertToConfigDomain(blueprint.Config),
	}, nil
}

func ConvertToBlueprintMaskDomain(mask *crd.BlueprintMaskManifest) (domain.BlueprintMask, error) {
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
