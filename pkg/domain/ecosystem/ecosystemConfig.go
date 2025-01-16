package ecosystem

import (
	"github.com/cloudogu/blueprint-lib/v2"
)

type EcosystemConfigEntry interface {
	*GlobalConfigEntry | *DoguConfigEntry | *SensitiveDoguConfigEntry
}

type GlobalConfigEntry struct {
	Key   v2.GlobalConfigKey
	Value v2.GlobalConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}

type DoguConfigEntry struct {
	Key   v2.DoguConfigKey
	Value v2.DoguConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}

type SensitiveDoguConfigEntry struct {
	Key   v2.SensitiveDoguConfigKey
	Value v2.SensitiveDoguConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}
