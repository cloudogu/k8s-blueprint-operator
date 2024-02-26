package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getNeededSensitiveConfigAction(t *testing.T) {
	tests := []struct {
		name                 string
		expected             ConfigValueState
		actual               ConfigValueState
		doguAlreadyInstalled bool
		want                 ConfigAction
	}{
		{
			name:                 "none, does not exist",
			expected:             ConfigValueState{Value: "", Exists: false},
			actual:               ConfigValueState{Value: "", Exists: false},
			doguAlreadyInstalled: false,
			want:                 ConfigActionNone,
		},
		{
			name:                 "none, exists",
			expected:             ConfigValueState{Value: "", Exists: true},
			actual:               ConfigValueState{Value: "", Exists: true},
			doguAlreadyInstalled: false,
			want:                 ConfigActionNone,
		},
		{
			name:                 "set to encrypt",
			expected:             ConfigValueState{Value: "", Exists: true},
			actual:               ConfigValueState{Value: "", Exists: false},
			doguAlreadyInstalled: false,
			want:                 ConfigActionSetToEncrypt,
		},
		{
			name:                 "set encrypted",
			expected:             ConfigValueState{Value: "", Exists: true},
			actual:               ConfigValueState{Value: "", Exists: false},
			doguAlreadyInstalled: true,
			want:                 ConfigActionSetEncrypted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t,
				tt.want, getNeededSensitiveConfigAction(tt.expected, tt.actual, tt.doguAlreadyInstalled),
				"getNeededSensitiveConfigAction(%v, %v, %v)", tt.expected, tt.actual, tt.doguAlreadyInstalled)
		})
	}
}
