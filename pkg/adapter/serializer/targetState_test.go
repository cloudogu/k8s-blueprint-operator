package serializer

import (
	"fmt"
	"github.com/cloudogu/blueprint-lib/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_toDomainTargetState(t *testing.T) {
	type args struct {
		state string
	}
	tests := []struct {
		name    string
		args    args
		want    v2.TargetState
		wantErr assert.ErrorAssertionFunc
	}{
		{"convert present state", args{"present"}, v2.TargetState(v2.TargetStatePresent), assert.NoError},
		{"convert absent state", args{"absent"}, v2.TargetState(v2.TargetStateAbsent), assert.NoError},
		{"convert empty state", args{""}, v2.TargetState(v2.TargetStatePresent), assert.NoError},
		{"error on unknown state", args{"unknown"}, v2.TargetState(v2.TargetStatePresent), assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDomainTargetState(tt.args.state)
			if !tt.wantErr(t, err, fmt.Sprintf("toDomainTargetState(%v)", tt.args.state)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toDomainTargetState(%v)", tt.args.state)
		})
	}
}

func Test_toSerializerTargetState(t *testing.T) {
	type args struct {
		domainState v2.TargetState
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"convert present state", args{v2.TargetState(v2.TargetStatePresent)}, "present", assert.NoError},
		{"convert absent state", args{v2.TargetState(v2.TargetStateAbsent)}, "absent", assert.NoError},
		{"error on unknown state", args{v2.TargetState(-1)}, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToSerializerTargetState(tt.args.domainState)
			if !tt.wantErr(t, err, fmt.Sprintf("toSerializerTargetState(%v)", tt.args.domainState)) {
				return
			}
			assert.Equalf(t, tt.want, got, "toSerializerTargetState(%v)", tt.args.domainState)
		})
	}
}
