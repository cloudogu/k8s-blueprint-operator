package blueprintV2

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

// BlueprintV2 describes an abstraction of CES components that should be absent or present within one or more CES
// instances. When the same Blueprint is applied to two different CES instances it is required to leave two equal
// instances in terms of the components.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type BlueprintV2 struct {
	serializer.GeneralBlueprint
	// Dogus contains a set of exact dogu versions, which should be present or absent
	// in the CES instance after this blueprint was applied.
	// Optional.
	Dogus []serializer.TargetDogu `json:"dogus,omitempty"`
	// Components are a set of exact package versions,
	// which should be present or absent in the CES instance after which this blueprint was applied.
	// The packages must correspond to the used package manager.
	// Optional.
	Components []serializer.TargetComponent `json:"components,omitempty"`
	// Config is used for ecosystem configuration to be applied.
	// Optional.
	Config Config `json:"config,omitempty"`
}

type RegistryConfig map[string]map[string]interface{}

func ConvertToBlueprintDTO(blueprint domain.Blueprint) (BlueprintV2, error) {
	convertedDogus, doguError := serializer.ConvertToDoguDTOs(blueprint.Dogus)
	convertedComponents, compError := serializer.ConvertToComponentDTOs(blueprint.Components)

	err := errors.Join(doguError, compError)
	if err != nil {
		return BlueprintV2{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return BlueprintV2{
		GeneralBlueprint: serializer.GeneralBlueprint{API: serializer.V2},
		Dogus:            convertedDogus,
		Components:       convertedComponents,
		Config:           ConvertToConfigDTO(blueprint.Config),
	}, nil
}

func convertToBlueprintDomain(blueprint BlueprintV2) (domain.Blueprint, error) {
	switch blueprint.API {
	case serializer.V1:
		return domain.Blueprint{}, fmt.Errorf("blueprint API V1 is deprecated and got removed: " +
			"packages and cesapp version got removed in favour of components")
	case serializer.V2:
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
		Config:     ConvertToConfigDomain(blueprint.Config),
	}, nil
}
