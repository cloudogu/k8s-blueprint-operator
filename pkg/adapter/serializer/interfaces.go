package serializer

import "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"

// BlueprintSerializer can serialize a domain.Blueprint to a string. The format is implementation specific.
// Add a new implementation if you want either another format, e.g. json, xml, or if you change your specific structure in that format.
type BlueprintSerializer interface {
	// Serialize translates a domain.Blueprint into a string representation.
	// Returns an error if the blueprint cannot be deserialized for any reason.
	Serialize(blueprint domain.Blueprint) (string, error)

	// Deserialize translates a string into a domain.Blueprint.
	// Returns a domain.InvalidBlueprintError if the given string has syntax or simple semantic errors.
	Deserialize(rawBlueprint string) (domain.Blueprint, error)
}

type BlueprintMaskSerializer interface {
	// Serialize translates a domain.BlueprintMask into a string representation.
	// Returns an error if the blueprint cannot be deserialized for any reason.
	Serialize(mask domain.BlueprintMask) (string, error)

	// Deserialize translates a string into a domain.BlueprintMask.
	// Returns a domain.InvalidBlueprintError if the given string has syntax or simple semantic errors.
	Deserialize(rawBlueprintMask string) (domain.BlueprintMask, error)
}
