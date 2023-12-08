package serializer

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

var ToID = map[string]domain.TargetState{
	"":        domain.TargetStatePresent,
	"present": domain.TargetStatePresent,
	"absent":  domain.TargetStateAbsent,
}

func ToDomainTargetState(stateString string) (domain.TargetState, error) {
	// Note that if the string is not found then it will be set to the zero value, which is 'Created'.
	id := ToID[stateString]
	var err error
	if id == domain.TargetStatePresent && stateString != "present" && stateString != "" {
		err = fmt.Errorf("unknown targetState '%s'", stateString)
	}
	return id, err
}

func ToSerializerTargetState(domainState domain.TargetState) (string, error) {
	convertedString := domainState.String()
	if convertedString != "present" && ToID[convertedString] == 0 {
		return "", fmt.Errorf("unknown target state ID: '%d'", domainState)
	}
	return domainState.String(), nil
}
