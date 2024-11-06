package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoguConfigDiffs_HasChanges(t *testing.T) {
	tests := []struct {
		name  string
		diffs DoguConfigDiffs
		want  bool
	}{
		{
			name: "no changes",
			diffs: DoguConfigDiffs{
				SensitiveDoguConfigEntryDiff{
					Key: dogu1Key1,
					Actual: DoguConfigValueState{
						Exists: false,
					},
					Expected: DoguConfigValueState{
						Exists: false,
					},
					NeededAction: ConfigActionNone,
				},
			},
			want: false,
		},
		{
			name: "set config",
			diffs: DoguConfigDiffs{
				SensitiveDoguConfigEntryDiff{
					Key: dogu1Key1,
					Actual: DoguConfigValueState{
						Exists: false,
					},
					Expected: DoguConfigValueState{
						Value:  "test",
						Exists: true,
					},
					NeededAction: ConfigActionSet,
				},
			},
			want: true,
		},
		{
			name: "remove config",
			diffs: DoguConfigDiffs{
				SensitiveDoguConfigEntryDiff{
					Key: dogu1Key1,
					Actual: DoguConfigValueState{
						Value:  "test",
						Exists: true,
					},
					Expected: DoguConfigValueState{
						Exists: true,
					},
					NeededAction: ConfigActionRemove,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.diffs.HasChanges(), "HasChanges()")
		})
	}
}
