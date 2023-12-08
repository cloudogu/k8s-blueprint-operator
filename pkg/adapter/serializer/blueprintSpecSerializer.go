package serializer

import (
	"encoding/json"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

func SerializeBlueprint(blueprint domain.Blueprint) (string, error) {
	blueprintDTO, err := ConvertToBlueprintV2(blueprint)

	if err != nil {
		return "", fmt.Errorf("cannot serialize blueprint: %w", err)
	}

	serializedBlueprint, err := json.Marshal(blueprintDTO)

	if err != nil {
		return "", fmt.Errorf("cannot serialize blueprint: %w", err)
	}
	return string(serializedBlueprint), nil
}

func DeserializeBlueprint(rawBlueprint []byte) (domain.Blueprint, error) {
	blueprintDTO := BlueprintV2{}

	err := json.Unmarshal(rawBlueprint, &blueprintDTO)

	if err != nil {
		return domain.Blueprint{}, fmt.Errorf("cannot deserialize blueprint: %w", err)
	}
	blueprint, err := convertToBlueprint(blueprintDTO)

	if err != nil {
		return domain.Blueprint{}, fmt.Errorf("cannot deserialize blueprint: %w", err)
	}

	return blueprint, nil
}

func SerializeBlueprintMask(mask domain.BlueprintMask) (string, error) {
	blueprintDTO, err := ConvertToBlueprintMaskV1(mask)

	if err != nil {
		return "", fmt.Errorf("cannot serialize blueprint mask: %w", err)
	}

	serializedMask, err := json.Marshal(blueprintDTO)

	if err != nil {
		return "", fmt.Errorf("cannot serialize blueprint mask: %w", err)
	}
	return string(serializedMask), nil
}

func DeserializeBlueprintMask(rawBlueprint []byte) (domain.BlueprintMask, error) {
	blueprintMaskDTO := BlueprintMaskV1{}

	err := json.Unmarshal(rawBlueprint, &blueprintMaskDTO)

	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("cannot deserialize blueprint mask: %w", err)
	}
	mask, err := convertToBlueprintMask(blueprintMaskDTO)

	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("cannot deserialize blueprint mask: %w", err)
	}

	return mask, nil
}