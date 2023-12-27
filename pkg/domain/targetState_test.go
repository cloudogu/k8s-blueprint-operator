package domain

import (
	"testing"
)

func TestTargetState_String(t *testing.T) {
	tests := []struct {
		name  string
		state TargetState
		want  string
	}{
		{
			"String() map enum to string",
			TargetStatePresent,
			"present",
		},
		{
			"String() map enum to string",
			TargetStateAbsent,
			"absent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("TargetState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
