package blueprintV2

import (
	"encoding/json"
	"fmt"
	"github.com/cloudogu/blueprint-lib/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

type Serializer struct{}

func (b Serializer) Serialize(blueprint v2.Blueprint) (string, error) {
	blueprintDTO, err := ConvertToBlueprintDTO(blueprint)

	if err != nil {
		return "", fmt.Errorf("cannot serialize blueprint: %w", err)
	}

	serializedBlueprint, err := json.Marshal(blueprintDTO)

	if err != nil {
		return "", fmt.Errorf("cannot serialize blueprint: %w", err)
	}
	return string(serializedBlueprint), nil
}

func (b Serializer) Deserialize(rawBlueprint string) (v2.Blueprint, error) {
	blueprintDTO := BlueprintV2{}

	err := json.Unmarshal([]byte(rawBlueprint), &blueprintDTO)

	if err != nil {
		return v2.Blueprint{}, &domain.InvalidBlueprintError{WrappedError: err, Message: "cannot deserialize blueprint"}
	}
	blueprint, err := convertToBlueprintDomain(blueprintDTO)

	if err != nil {
		return v2.Blueprint{}, &domain.InvalidBlueprintError{WrappedError: err, Message: "cannot deserialize blueprint"}
	}

	return blueprint, nil
}
