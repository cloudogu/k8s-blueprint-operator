package serializer

import (
	"fmt"
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
	return domainState != domain.TargetStatePresent
}

//FIXME: remove old TargetState types, we need to change the domain.StateDiff for that
// we do this in #54968 if these changes on the blueprint CRD got merged, so we do not have to revert everything

// ToID provides common mappings from strings to domain.TargetState, e.g. for dogus.
var ToID = map[string]domain.TargetState{
	"":        domain.TargetStatePresent,
	"present": domain.TargetStatePresent,
	"absent":  domain.TargetStateAbsent,
}

// ToOldDomainTargetState maps a string to a domain.TargetState or returns an error if this is not possible.
func ToOldDomainTargetState(stateString string) (domain.TargetState, error) {
	// Note that if the string is not found then it will be set to the zero value, which is 'Created'.
	id := ToID[stateString]
	var err error
	if id == domain.TargetStatePresent && stateString != "present" && stateString != "" {
		err = fmt.Errorf("unknown target state %q", stateString)
	}
	return id, err
}
