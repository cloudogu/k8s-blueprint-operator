package ecosystem

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

type GlobalConfigKey string
type GlobalConfigValue string
type GlobalConfigEntry struct {
	Key   GlobalConfigKey
	Value GlobalConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}

type DoguConfigKey struct {
	DoguName common.SimpleDoguName
	Key      string
}
type DoguConfigValue string

type DoguConfigEntry struct {
	Key   DoguConfigKey
	Value DoguConfigValue
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext interface{}
}
