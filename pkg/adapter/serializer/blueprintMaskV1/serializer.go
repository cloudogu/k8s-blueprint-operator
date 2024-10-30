package blueprintMaskV1

import (
	"encoding/json"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
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

func (b Serializer) Deserialize(rawBlueprintMask string) (domain.BlueprintMask, error) {
	blueprintMaskDTO := BlueprintMaskV1{}

	err := json.Unmarshal([]byte(rawBlueprintMask), &blueprintMaskDTO)

	if err != nil {
		return domain.BlueprintMask{}, &domain.InvalidBlueprintError{WrappedError: err, Message: "cannot deserialize blueprint mask"}
	}
	mask, err := convertToBlueprintMask(blueprintMaskDTO)

	if err != nil {
		return domain.BlueprintMask{}, &domain.InvalidBlueprintError{WrappedError: err, Message: "cannot deserialize blueprint mask"}
	}

	return mask, nil
}
