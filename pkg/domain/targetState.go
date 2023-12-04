package domain

// TargetState defines an enum of values that determines a state of installation.
type TargetState int

const (
	// TargetStatePresent is the default state. If selected the chosen item must be present after the blueprint was
	// applied.
	TargetStatePresent = iota
	// TargetStateAbsent sets the state of the item to absent. If selected the chosen item must be absent after the
	// blueprint was applied.
	TargetStateAbsent
)

var PossbileTargetStates = []TargetState{
	TargetStatePresent, TargetStateAbsent,
}

// String returns a string representation of the given TargetState enum value.
func (state TargetState) String() string {
	return toString[state]
}

var toString = map[TargetState]string{
	TargetStatePresent: "present",
	TargetStateAbsent:  "absent",
}
