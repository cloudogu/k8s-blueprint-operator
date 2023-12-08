package blueprintMaskV1

import (
	"encoding/json"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

type Serializer struct{}

func (b Serializer) Serialize(mask domain.BlueprintMask) (string, error) {
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

func (b Serializer) Deserialize(rawBlueprint string) (domain.BlueprintMask, error) {
	blueprintMaskDTO := BlueprintMaskV1{}

	err := json.Unmarshal([]byte(rawBlueprint), &blueprintMaskDTO)

	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("cannot deserialize blueprint mask: %w", err)
	}
	mask, err := convertToBlueprintMask(blueprintMaskDTO)

	if err != nil {
		return domain.BlueprintMask{}, fmt.Errorf("cannot deserialize blueprint mask: %w", err)
	}

	return mask, nil
}
