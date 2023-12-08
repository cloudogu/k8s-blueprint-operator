package blueprintMaskV1

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerializeBlueprintMask_ok(t *testing.T) {
	serializer := Serializer{}
	type args struct {
		spec domain.BlueprintMask
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"empty blueprint mask",
			args{spec: domain.BlueprintMask{}},
			`{"blueprintMaskApi":"v1","blueprintMaskId":"","dogus":[]}`,
			assert.NoError,
		},
		{
			"dogus in blueprint mask",
			args{spec: domain.BlueprintMask{
				Dogus: []domain.MaskTargetDogu{
					{Namespace: "official", Name: "nginx", Version: "1.2.0-1", TargetState: domain.TargetStatePresent},
					{Namespace: "premium", Name: "jira", Version: "3.0.2-2", TargetState: domain.TargetStateAbsent},
				},
			}},
			`{"blueprintMaskApi":"v1","blueprintMaskId":"","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present"},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent"}]}`,
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializer.Serialize(tt.args.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("SerializeBlueprintMask(%v)", tt.args.spec)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SerializeBlueprintMask(%v)", tt.args.spec)
		})
	}
}

func TestSerializeBlueprintMask_error(t *testing.T) {
	serializer := Serializer{}
	mask := domain.BlueprintMask{
		Dogus: []domain.MaskTargetDogu{
			{Namespace: "official", Name: "nginx", Version: "1.2.0-1", TargetState: -1},
		},
	}

	_, err := serializer.Serialize(mask)

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot serialize blueprint mask: ")
	assert.ErrorContains(t, err, "unknown target state ID: '-1'")
}

func TestDeserializeBlueprintMask_ok(t *testing.T) {
	serializer := Serializer{}
	type args struct {
		spec string
	}
	tests := []struct {
		name    string
		args    args
		want    domain.BlueprintMask
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"empty blueprint mask",
			args{spec: `{"blueprintMaskApi":"v1"}`},
			domain.BlueprintMask{},
			assert.NoError,
		},
		{
			"dogus in blueprint mask",
			args{spec: `{"blueprintMaskApi":"v1","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present"},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent"}]}`},
			domain.BlueprintMask{
				Dogus: []domain.MaskTargetDogu{
					{Namespace: "official", Name: "nginx", Version: "1.2.0-1", TargetState: domain.TargetStatePresent},
					{Namespace: "premium", Name: "jira", Version: "3.0.2-2", TargetState: domain.TargetStateAbsent},
				}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializer.Deserialize(tt.args.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("SerializeBlueprint(%v)", tt.args.spec)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SerializeBlueprint(%v)", tt.args.spec)
		})
	}
}

func TestDeserializeBlueprintMask_errors(t *testing.T) {
	serializer := Serializer{}
	type args struct {
		rawBlueprint string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "json syntax error",
			args:    args{`{a}`},
			want:    "cannot deserialize blueprint mask: invalid character 'a' looking for beginning of object key string",
			wantErr: assert.Error,
		},
		{
			name:    "deserialize unknown API version",
			args:    args{`{"blueprintMaskApi":"v0"}`},
			want:    "cannot deserialize blueprint mask: unsupported Blueprint Mask API Version: v0",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := serializer.Deserialize(tt.args.rawBlueprint)
			if !tt.wantErr(t, err, fmt.Sprintf("DeserializeBlueprint(%v)", tt.args.rawBlueprint)) {
				return
			}
			assert.ErrorContains(t, err, tt.want, "DeserializeBlueprint(%v)", tt.args.rawBlueprint)
		})
	}
}
