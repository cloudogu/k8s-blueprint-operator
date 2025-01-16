package serializer

import (
	"github.com/cloudogu/blueprint-lib/v2"
)

// BlueprintSerializer can serialize a domain.Blueprint to a string. The format is implementation specific.
// Add a new implementation if you want either another format, e.g. json, xml, or if you change your specific structure in that format.
type BlueprintSerializer interface {
	// Serialize translates a domain.Blueprint into a string representation.
	// Returns an error if the blueprint cannot be deserialized for any reason.
	Serialize(blueprint v2.Blueprint) (string, error)

	// Deserialize translates a string into a domain.Blueprint.
	// Returns a domain.InvalidBlueprintError if the given string has syntax or simple semantic errors.
	Deserialize(rawBlueprint string) (v2.Blueprint, error)
}

type BlueprintMaskSerializer interface {
	// Serialize translates a domain.BlueprintMask into a string representation.
	// Returns an error if the blueprint cannot be deserialized for any reason.
	Serialize(mask v2.BlueprintMask) (string, error)

	// Deserialize translates a string into a domain.BlueprintMask.
	// Returns a domain.InvalidBlueprintError if the given string has syntax or simple semantic errors.
	Deserialize(rawBlueprintMask string) (v2.BlueprintMask, error)
}

// BlueprintApi represents an API version for a specific serialized format of domain.Blueprint.
type BlueprintApi string

const (
	// V1 is the API version of the BlueprintV1 json format of the classic-CES and which is used with the cesapp.
	V1 BlueprintApi = "v1"
	// V2 is the API version of the BlueprintV2 json format used in the MultiNode-CES inside kubernetes, e.g. for Blueprint-CRs.
	V2 BlueprintApi = "v2"
)

// GeneralBlueprint defines the minimum set to parse the blueprint API version string in order to select the right
// blueprint handling strategy. This is necessary in order to accommodate maximal changes in different blueprint API
// versions.
type GeneralBlueprint struct {
	// API is used to distinguish between different versions of the used API and impacts directly the interpretation of
	// this blueprint. Must not be empty.
	//
	// This field MUST NOT be MODIFIED or REMOVED because the API is paramount for distinguishing between different
	// blueprint version implementations.
	API BlueprintApi `json:"blueprintApi"`
}

// BlueprintMaskApi represents an API version for a specific serialized format of domain.BlueprintMask.
type BlueprintMaskApi string

const (
	// BlueprintMaskAPIV1 is the API version used by the Classic-CES and the MultiNode-CES as it has no environment specific fields.
	BlueprintMaskAPIV1 BlueprintMaskApi = "v1"
)

// GeneralBlueprintMask defines the minimum set to parse the blueprint mask API version string in order to select the
// right blueprint mask handling strategy. This is necessary in order to accommodate maximal changes in different
// blueprint mask API versions.
type GeneralBlueprintMask struct {
	// API is used to distinguish between different versions of the used API and impacts directly the interpretation of
	// this blueprint mask. Must not be empty.
	//
	// This field MUST NOT be MODIFIED or REMOVED because the API is paramount for distinguishing between different
	// blueprint mask version implementations.
	API BlueprintMaskApi `json:"blueprintMaskApi"`
}
