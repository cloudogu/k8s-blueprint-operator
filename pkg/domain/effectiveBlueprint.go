package domain

// EffectiveBlueprintV2 declaratively describes what is the wanted state after evaluating the blueprint and the blueprintMask.
// This is still a static description, so no actual state of the ecosystem is taken into consideration here.
type EffectiveBlueprint struct {
	// ID is the unique name of the set over all components. This blueprint ID should be used to distinguish from similar
	// blueprints between humans in an easy way. Must not be empty.
	ID string `json:"blueprintId"`
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu `json:"dogus,omitempty"`
	// Components contains a set of exact components versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []Component `json:"components,omitempty"`
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig RegistryConfig `json:"registryConfig,omitempty"`
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string `json:"registryConfigAbsent,omitempty"`
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted RegistryConfig `json:"registryConfigEncrypted,omitempty"`
}
