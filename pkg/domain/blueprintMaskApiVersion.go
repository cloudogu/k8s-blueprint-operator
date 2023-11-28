package domain

type BlueprintMaskApi string

const (
	BlueprintMaskAPIV1 BlueprintApi = "v1"
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
