package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

// ToDomainTargetState maps a string to a domain.TargetState or returns an error if this is not possible.
func ToDomainTargetState(absent bool) domain.TargetState {
	if absent {
		return domain.TargetStateAbsent
	} else {
		return domain.TargetStatePresent
	}
}

// ToSerializerAbsentState maps a domain.TargetState to the absent state of dogus in the CR.
// If the state is not present, it will be interpreted as absent
func ToSerializerAbsentState(domainState domain.TargetState) bool {
	return domainState == domain.TargetStatePresent
}
