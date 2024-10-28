package ecosystem

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

type EcosystemConfigEntry interface {
	*GlobalConfigEntry | *DoguConfigEntry | *SensitiveDoguConfigEntry
}

type GlobalConfigEntry struct {
	Key   common.GlobalConfigKey
	Value common.GlobalConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}

type DoguConfigEntry struct {
	Key   common.DoguConfigKey
	Value common.DoguConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}

type SensitiveDoguConfigEntry struct {
	Key   common.SensitiveDoguConfigKey
	Value common.SensitiveDoguConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}
