package serializer

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"strings"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

var configKeySeparator = "/"

func ConvertToBlueprintDTO(blueprint domain.EffectiveBlueprint) (crd.Blueprint, error) {
	var errorList []error
	convertedDogus, doguError := ConvertToDoguDTOs(blueprint.Dogus)
	convertedComponents, componentError := ConvertToComponentDTOs(blueprint.Components)
	errorList = append(errorList, doguError, componentError)

	err := errors.Join(errorList...)
	if err != nil {
		return crd.Blueprint{}, fmt.Errorf("cannot convert blueprintMask to BlueprintMaskV1 DTO: %w", err)
	}

	return crd.Blueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config:     ConvertToConfigDTO(blueprint.Config),
	}, nil
}

func convertToTargetState(absent bool) domain.TargetState {
	if absent {
		return domain.TargetStateAbsent
	} else {
		return domain.TargetStatePresent
	}
}

func ConvertToEffectiveBlueprintDomain(blueprint crd.Blueprint) (domain.EffectiveBlueprint, error) {
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
