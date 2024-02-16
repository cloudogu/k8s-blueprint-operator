package v1

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"strings"
)

var configKeySeparator = "/"

// EffectiveBlueprint describes an abstraction of CES components that should be absent or present within one or more CES
// instances after combining the blueprint with the blueprint mask.
//
// In general additions without changing the version are fine, as long as they don't change semantics. Removal or
// renaming are breaking changes and require a new blueprint API version.
type EffectiveBlueprint struct {
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []serializer.TargetDogu `json:"dogus,omitempty"`
	// Components contains a set of exact component versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []serializer.TargetComponent `json:"components,omitempty"`
	// Config is used for ecosystem configuration to be applied.
	// Optional.
	Config Config `json:"config,omitempty"`
}

func ConvertToEffectiveBlueprintDTO(blueprint domain.EffectiveBlueprint) (EffectiveBlueprint, error) {
	var errorList []error
	convertedDogus, doguError := serializer.ConvertToDoguDTOs(blueprint.Dogus)
	convertedComponents, componentError := serializer.ConvertToComponentDTOs(blueprint.Components)
	errorList = append(errorList, doguError, componentError)

	err := errors.Join(errorList...)
	if err != nil {
		return EffectiveBlueprint{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDTO(blueprint.Config),
	}, nil
}

func ConvertToEffectiveBlueprintDomain(blueprint EffectiveBlueprint) (domain.EffectiveBlueprint, error) {
	convertedDogus, doguErr := serializer.ConvertDogus(blueprint.Dogus)
	convertedComponents, compErr := serializer.ConvertComponents(blueprint.Components)

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

func widenMap(currentMap map[string]string) map[string]interface{} {
	newMap := map[string]interface{}{}
	for key, val := range currentMap {
		keys := strings.Split(key, configKeySeparator)
		setKey(keys, val, newMap)
	}
	return newMap
}

func setKey(keys []string, value string, initialMap map[string]interface{}) {
	currentMap := initialMap
	length := len(keys)
	for i, key := range keys {
		if i == length-1 {
			currentMap[key] = value
			break
		}
		if currentMap[key] == nil {
			currentMap[key] = map[string]interface{}{}
		}
		currentMap, _ = currentMap[key].(map[string]interface{})
	}
}
