package domain

// EffectiveBlueprint declaratively describes what is the wanted state after evaluating the blueprint and the blueprintMask.
// This is still a static description, so no actual state of the ecosystem is taken into consideration here.
type EffectiveBlueprint struct {
	// Dogus contains a set of exact dogu versions which should be present or absent in the CES instance after which this
	// blueprint was applied. Optional.
	Dogus []TargetDogu
	// Components contains a set of exact components versions which should be present or absent in the CES instance after which
	// this blueprint was applied. Optional.
	Components []Component
	// Used to configure registry globalRegistryEntries on blueprint upgrades
	RegistryConfig RegistryConfig
	// Used to remove registry globalRegistryEntries on blueprint upgrades
	RegistryConfigAbsent []string
	// Used to configure encrypted registry globalRegistryEntries on blueprint upgrades
	RegistryConfigEncrypted RegistryConfig
}

// GetWantedDogus returns a list of all dogus which should be installed
func (effectiveBlueprint *EffectiveBlueprint) GetWantedDogus() []TargetDogu {
	var wantedDogus []TargetDogu
	for _, dogu := range effectiveBlueprint.Dogus {
		if dogu.TargetState == TargetStatePresent {
			wantedDogus = append(wantedDogus, dogu)
		}
	}
	return wantedDogus
}
