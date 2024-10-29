package serializer

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
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
		want    domain.TargetState
		wantErr assert.ErrorAssertionFunc
	}{
		{"convert present state", args{"present"}, domain.TargetState(domain.TargetStatePresent), assert.NoError},
		{"convert absent state", args{"absent"}, domain.TargetState(domain.TargetStateAbsent), assert.NoError},
		{"convert empty state", args{""}, domain.TargetState(domain.TargetStatePresent), assert.NoError},
		{"error on unknown state", args{"unknown"}, domain.TargetState(domain.TargetStatePresent), assert.Error},
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
		domainState domain.TargetState
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"convert present state", args{domain.TargetState(domain.TargetStatePresent)}, "present", assert.NoError},
		{"convert absent state", args{domain.TargetState(domain.TargetStateAbsent)}, "absent", assert.NoError},
		{"error on unknown state", args{domain.TargetState(-1)}, "", assert.Error},
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
