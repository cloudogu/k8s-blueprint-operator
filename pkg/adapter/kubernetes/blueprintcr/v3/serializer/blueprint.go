package serializer

import (
	"errors"
	"fmt"

	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func SerializeBlueprintAndMask(blueprintSpec *domain.BlueprintSpec, manifest bpv3.BlueprintManifest, maskManifest *bpv3.BlueprintMaskManifest) error {
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

func ConvertBlueprintStatus(blueprintCR *bpv3.Blueprint) (domain.EffectiveBlueprint, error) {
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

func ConvertToBlueprintDTO(blueprint domain.EffectiveBlueprint) bpv3.BlueprintManifest {
	return bpv3.BlueprintManifest{
		Dogus:  ConvertToDoguDTOs(blueprint.Dogus),
		Config: ConvertToConfigDTO(blueprint.Config),
	}
}

func ConvertToBlueprintDomain(blueprint bpv3.BlueprintManifest) (domain.Blueprint, error) {
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

func ConvertToEffectiveBlueprintDomain(blueprint *bpv3.BlueprintManifest) (domain.EffectiveBlueprint, error) {
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

func ConvertToBlueprintMaskDomain(mask *bpv3.BlueprintMaskManifest) (domain.BlueprintMask, error) {
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
