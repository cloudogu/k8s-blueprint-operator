package serializer

import (
	"fmt"
	"github.com/cloudogu/blueprint-lib/v2"
)

// ToID provides common mappings from strings to domain.TargetState, e.g. for dogus.
var ToID = map[string]v2.TargetState{
	"":        v2.TargetStatePresent,
	"present": v2.TargetStatePresent,
	"absent":  v2.TargetStateAbsent,
}

// ToDomainTargetState maps a string to a domain.TargetState or returns an error if this is not possible.
func ToDomainTargetState(stateString string) (v2.TargetState, error) {
	// Note that if the string is not found then it will be set to the zero value, which is 'Created'.
	id := ToID[stateString]
	var err error
	if id == v2.TargetStatePresent && stateString != "present" && stateString != "" {
		err = fmt.Errorf("unknown target state %q", stateString)
	}
	return id, err
}

// ToSerializerTargetState maps a domain.TargetState to a string or returns an error if this is not possible.
func ToSerializerTargetState(domainState v2.TargetState) (string, error) {
	convertedString := domainState.String()
	if convertedString != "present" && ToID[convertedString] == 0 {
		return "", fmt.Errorf("unknown target state ID: '%d'", domainState)
	}
	return domainState.String(), nil
}
