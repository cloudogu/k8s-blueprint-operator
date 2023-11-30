package domain

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
)

// TargetState defines an enum of values that determines a state of installation.
type TargetState int

const (
	// TargetStatePresent is the default state. If selected the chosen item must be present after the blueprint was
	// applied.
	TargetStatePresent = iota
	// TargetStateAbsent sets the state of the item to absent. If selected the chosen item must be absent after the
	// blueprint was applied.
	TargetStateAbsent
	// TargetStateIgnore is currently only internally used to mark items that are present in the CES instance at hand
	// but not mentioned in the blueprint.
	TargetStateIgnore
)

var PossbileTargetStates = []TargetState{
	TargetStatePresent, TargetStateAbsent, TargetStateIgnore,
}

// String returns a string representation of the given TargetState enum value.
func (state TargetState) String() string {
	return toString[state]
}

var toString = map[TargetState]string{
	TargetStatePresent: "present",
	TargetStateAbsent:  "absent",
}

var toID = map[string]TargetState{
	"present": TargetStatePresent,
	"absent":  TargetStateAbsent,
}

// MarshalJSON marshals the enum as a quoted json string
func (state TargetState) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[state])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value. Use it with usual json unmarshalling:
//
//	 jsonBlob := []byte("\"present\"")
//		var state TargetState
//		err := json.Unmarshal(jsonBlob, &state)
func (state *TargetState) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return errors.Wrapf(err, "cannot unmarshal value %s to a TargetState", string(b))
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*state = toID[j]
	return nil
}
