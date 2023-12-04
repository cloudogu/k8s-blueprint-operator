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
