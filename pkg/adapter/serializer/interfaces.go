package serializer

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain"

type BlueprintSerializer interface {
	Serialize(blueprint domain.Blueprint) (string, error)
	Deserialize(rawBlueprint string) (domain.Blueprint, error)
}

type BlueprintMaskSerializer interface {
	Serialize(mask domain.BlueprintMask) (string, error)
	Deserialize(rawBlueprint string) (domain.BlueprintMask, error)
}

type BlueprintApi string

const (
	V1 BlueprintApi = "v1"
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

type BlueprintMaskApi string

const (
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
