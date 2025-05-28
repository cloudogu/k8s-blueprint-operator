package v1

import (
	"errors"
	"fmt"
	"strings"

	crd "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

var configKeySeparator = "/"

func ConvertToEffectiveBlueprintDTO(blueprint domain.EffectiveBlueprint) (crd.EffectiveBlueprint, error) {
	var errorList []error
	convertedDogus, doguError := serializer.ConvertToDoguDTOs(blueprint.Dogus)
	convertedComponents, componentError := serializer.ConvertToComponentDTOs(blueprint.Components)
	errorList = append(errorList, doguError, componentError)

	err := errors.Join(errorList...)
	if err != nil {
		return crd.EffectiveBlueprint{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return crd.EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDTO(blueprint.Config),
	}, nil
}

func ConvertToEffectiveBlueprintDomain(blueprint crd.EffectiveBlueprint) (domain.EffectiveBlueprint, error) {
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
