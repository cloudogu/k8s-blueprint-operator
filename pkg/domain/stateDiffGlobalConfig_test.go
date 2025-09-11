package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalConfigDiffs_HasChanges(t *testing.T) {
	valChanged := "changed"
	valInitial := "initial"
	tests := []struct {
		name  string
		diffs GlobalConfigDiffs
		want  bool
	}{
		{
			name:  "false on empty input",
			diffs: GlobalConfigDiffs{},
			want:  false,
		},
		{
			name: "true on GlobalConfigDiff with action",
			diffs: []GlobalConfigEntryDiff{
				{
					Key:          "testkey",
					Actual:       GlobalConfigValueState{Value: &valChanged, Exists: true},
					Expected:     GlobalConfigValueState{Value: &valInitial, Exists: true},
					NeededAction: ConfigActionSet,
				},
			},
			want: true,
		},
		{
			name: "false on GlobalConfigDiff without action",
			diffs: []GlobalConfigEntryDiff{
				{
					Key:          "testkey",
					Actual:       GlobalConfigValueState{Value: &valChanged, Exists: true},
					Expected:     GlobalConfigValueState{Value: &valInitial, Exists: true},
					NeededAction: ConfigActionNone,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.diffs.HasChanges(), "HasChanges()")
		})
	}
}
