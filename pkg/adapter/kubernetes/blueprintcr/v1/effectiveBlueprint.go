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
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig map[string]string `json:"registryConfig,omitempty"`
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string `json:"registryConfigAbsent,omitempty"`
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted map[string]string `json:"registryConfigEncrypted,omitempty"`
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

	registryConfigAbsent := blueprint.RegistryConfigAbsent
	if registryConfigAbsent == nil {
		registryConfigAbsent = []string{}
	}

	return EffectiveBlueprint{
		Dogus:                   convertedDogus,
		Components:              convertedComponents,
		RegistryConfig:          flattenRegistryConfig(blueprint.RegistryConfig),
		RegistryConfigAbsent:    registryConfigAbsent,
		RegistryConfigEncrypted: flattenRegistryConfig(blueprint.RegistryConfigEncrypted),
	}, nil
}

func ConvertToEffectiveBlueprintDomain(blueprint EffectiveBlueprint) (domain.EffectiveBlueprint, error) {
	convertedDogus, doguErr := serializer.ConvertDogus(blueprint.Dogus)
	convertedComponents, compErr := serializer.ConvertComponents(blueprint.Components)
	convertedConfig, configError := convertToRegistryConfig(blueprint.RegistryConfig)
	convertedEncryptedConfig, encryptedConfigError := convertToRegistryConfig(blueprint.RegistryConfig)
	err := errors.Join(doguErr, compErr, configError, encryptedConfigError)
	if err != nil {
		return domain.EffectiveBlueprint{}, fmt.Errorf("syntax of blueprintV2 is not correct: %w", err)
	}
	return domain.EffectiveBlueprint{
		Dogus:                   convertedDogus,
		Components:              convertedComponents,
		RegistryConfig:          convertedConfig,
		RegistryConfigAbsent:    blueprint.RegistryConfigAbsent,
		RegistryConfigEncrypted: convertedEncryptedConfig,
	}, nil
}

func convertToRegistryConfig(flattenedConfig map[string]string) (domain.RegistryConfig, error) {
	// expand key structure
	widenedMap := widenMap(flattenedConfig)
	//convert it to domain.RegistryConfig (which has at least depth 2)
	config := domain.RegistryConfig{}
	for key1, val1 := range widenedMap {
		switch subMap := val1.(type) {
		case map[string]interface{}:
			for key2, val2 := range subMap {
				if config[key1] == nil {
					config[key1] = make(map[string]interface{})
				}
				config[key1][key2] = val2
			}
		default:
			return domain.RegistryConfig{}, fmt.Errorf("registry config is invalid: values need to be at least at depth 2: key %q is invalid", key1)
		}
	}
	return config, nil
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

func flattenRegistryConfig(config domain.RegistryConfig) map[string]string {
	intermediateResult := make(map[string]interface{})
	for key, value := range config {
		for key2, value2 := range value {
			intermediateResult[key+configKeySeparator+key2] = value2
		}
	}
	keyToValueConfig := make(map[string]string)
	flattenMap("", intermediateResult, keyToValueConfig)

	return keyToValueConfig
}

func flattenMap(prefix string, src map[string]interface{}, dest map[string]string) {
	if len(prefix) > 0 {
		prefix += configKeySeparator
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flattenMap(prefix+k, child, dest)
		//case []interface{}: there should be no arrays in the config
		default:
			dest[prefix+k] = fmt.Sprintf("%s", v)
		}
	}
}
