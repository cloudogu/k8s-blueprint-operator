package blueprintV2

import (
	"errors"
	"fmt"

	bpv2 "github.com/cloudogu/k8s-blueprint-lib/json/blueprintV2"
	"github.com/cloudogu/k8s-blueprint-lib/json/bpcore"

	v1 "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToBlueprintDTO(blueprint domain.Blueprint) (bpv2.BlueprintV2, error) {
	convertedDogus, doguError := serializer.ConvertToDoguDTOs(blueprint.Dogus)
	convertedComponents, compError := serializer.ConvertToComponentDTOs(blueprint.Components)

	err := errors.Join(doguError, compError)
	if err != nil {
		return bpv2.BlueprintV2{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return bpv2.BlueprintV2{
		GeneralBlueprint: bpcore.GeneralBlueprint{API: bpcore.V2},
		Dogus:            convertedDogus,
		Components:       convertedComponents,
		Config:           v1.ConvertToConfigDTO(blueprint.Config),
	}, nil
}

func convertToBlueprintDomain(blueprint bpv2.BlueprintV2) (domain.Blueprint, error) {
	switch blueprint.API {
	case "v1":
		return domain.Blueprint{}, fmt.Errorf("blueprint API V1 is deprecated and got removed: " +
			"packages and cesapp version got removed in favour of components")
	case bpcore.V2:
	default:
		return domain.Blueprint{}, fmt.Errorf("unsupported Blueprint API Version: %s", blueprint.API)
	}
	convertedDogus, doguErr := serializer.ConvertDogus(blueprint.Dogus)
	convertedComponents, compErr := serializer.ConvertComponents(blueprint.Components)
	err := errors.Join(doguErr, compErr)
	if err != nil {
		return domain.Blueprint{}, fmt.Errorf("syntax of blueprintV2 is not correct: %w", err)
	}

	return domain.Blueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     v1.ConvertToConfigDomain(blueprint.Config),
	}, nil
}
