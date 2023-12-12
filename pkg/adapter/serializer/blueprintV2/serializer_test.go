package blueprintV2

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	version1_2_0_1, _ = core.ParseVersion("1.2.0-1")
	version3_0_2_2, _ = core.ParseVersion("3.0.2-2")
	version0_2_1_1, _ = core.ParseVersion("0.2.1-1")
)

func TestSerializeBlueprint_ok(t *testing.T) {
	serializer := Serializer{}
	type args struct {
		spec domain.Blueprint
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"empty blueprint",
			args{spec: domain.Blueprint{}},
			`{"blueprintApi":"v2"}`,
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: domain.Blueprint{
				Dogus: []domain.Dogu{
					{Namespace: "official", Name: "nginx", Version: version1_2_0_1, TargetState: domain.TargetStatePresent},
					{Namespace: "premium", Name: "jira", Version: version3_0_2_2, TargetState: domain.TargetStateAbsent},
				},
			}},
			`{"blueprintApi":"v2","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present"},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent"}]}`,
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: domain.Blueprint{
				Components: []domain.Component{
					{Name: "blueprint-operator", Version: version0_2_1_1, TargetState: domain.TargetStatePresent},
					{Name: "dogu-operator", Version: version3_2_1_1, TargetState: domain.TargetStateAbsent},
				},
			}},
			`{"blueprintApi":"v2","components":[{"name":"blueprint-operator","version":"0.2.1-1","targetState":"present"},{"name":"dogu-operator","version":"3.2.1-1","targetState":"absent"}]}`,
			assert.NoError,
		},
		{
			"RegistryConfig in blueprint",
			args{spec: domain.Blueprint{
				RegistryConfig: domain.RegistryConfig{
					"dogu": map[string]interface{}{
						"test": "42",
					},
				},
			}},
			`{"blueprintApi":"v2","registryConfig":{"dogu":{"test":"42"}}}`,
			assert.NoError,
		},
		{
			"RegistryConfigAbsent in blueprint",
			args{spec: domain.Blueprint{
				RegistryConfigAbsent: []string{
					"dogu/jenkins/java_mem",
					"second/key",
				},
			}},
			`{"blueprintApi":"v2","registryConfigAbsent":["dogu/jenkins/java_mem","second/key"]}`,
			assert.NoError,
		},
		{
			"RegistryConfigEncrypted in blueprint",
			args{spec: domain.Blueprint{
				RegistryConfigEncrypted: domain.RegistryConfig{
					"dogu": map[string]interface{}{
						"jenkins": map[string]string{
							"privateKey": "==key to encrypt later==",
						},
					},
				},
			}},
			`{"blueprintApi":"v2","registryConfigEncrypted":{"dogu":{"jenkins":{"privateKey":"==key to encrypt later=="}}}}`,
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serializer.Serialize(tt.args.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("SerializeBlueprint(%v)", tt.args.spec)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SerializeBlueprint(%v)", tt.args.spec)
		})
	}
}

func TestSerializeBlueprint_error(t *testing.T) {
	serializer := Serializer{}
	blueprint := domain.Blueprint{
		Dogus: []domain.Dogu{
			{Namespace: "official", Name: "nginx", Version: version1_2_0_1, TargetState: -1},
		},
	}

	_, err := serializer.Serialize(blueprint)

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot serialize blueprint: ")
	assert.ErrorContains(t, err, "unknown target state ID: '-1'")
}

func TestDeserializeBlueprint_ok(t *testing.T) {
	serializer := Serializer{}
	type args struct {
		spec string
	}
	tests := []struct {
		name    string
		args    args
		want    domain.Blueprint
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"empty blueprint",
			args{spec: `{"blueprintApi":"v2"}`},
			domain.Blueprint{},
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: `{"blueprintApi":"v2","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present"},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent"}]}`},
			domain.Blueprint{
				Dogus: []domain.Dogu{
					{Namespace: "official", Name: "nginx", Version: version1_2_0_1, TargetState: domain.TargetStatePresent},
					{Namespace: "premium", Name: "jira", Version: version3_0_2_2, TargetState: domain.TargetStateAbsent},
				}},
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: `{"blueprintApi":"v2","components":[{"name":"blueprint-operator","version":"0.2.1-1","targetState":"present"},{"name":"dogu-operator","version":"3.2.1-1","targetState":"absent"}]}`},
			domain.Blueprint{
				Components: []domain.Component{
					{Name: "blueprint-operator", Version: version0_2_1_1, TargetState: domain.TargetStatePresent},
					{Name: "dogu-operator", Version: version3_2_1_1, TargetState: domain.TargetStateAbsent},
				},
			},
			assert.NoError,
		},
		{
			"RegistryConfig in blueprint",
			args{spec: `{"blueprintApi":"v2","registryConfig":{"dogu":{"test":"42"}}}`},
			domain.Blueprint{
				RegistryConfig: domain.RegistryConfig{
					"dogu": map[string]interface{}{
						"test": "42",
					},
				},
			},
			assert.NoError,
		},
		{
			"RegistryConfigAbsent in blueprint",
			args{spec: `{"blueprintApi":"v2","registryConfigAbsent":["dogu/jenkins/java_mem","second/key"]}`},
			domain.Blueprint{
				RegistryConfigAbsent: []string{
					"dogu/jenkins/java_mem",
					"second/key",
				},
			},
			assert.NoError,
		},
		{
			"RegistryConfigEncrypted in blueprint",
			args{spec: `{"blueprintApi":"v2","registryConfigEncrypted":{"dogu":{"privateKey":"==key to encrypt later=="}}}`},
			domain.Blueprint{
				RegistryConfigEncrypted: domain.RegistryConfig{
					"dogu": map[string]interface{}{
						"privateKey": "==key to encrypt later==",
					},
				},
			},
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

func TestDeserializeBlueprint_errors(t *testing.T) {
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
			want:    "cannot deserialize blueprint: invalid character 'a' looking for beginning of object key string",
			wantErr: assert.Error,
		},
		{
			name:    "deserialize API V1",
			args:    args{`{"blueprintApi":"v1"}`},
			want:    "blueprint API V1 is deprecated and got removed: packages and cesapp version got removed in favour of components",
			wantErr: assert.Error,
		},
		{
			name:    "deserialize unknown API version",
			args:    args{`{"blueprintApi":"v0"}`},
			want:    "unsupported Blueprint API Version: v0",
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

func TestDeserializeBlueprint_testErrorType(t *testing.T) {
	serializer := Serializer{}

	_, err := serializer.Deserialize(`{}`)
	require.Error(t, err)
	assert.ErrorContains(t, err, "cannot deserialize blueprint")
	var expectedErrorType *domain.InvalidBlueprintError
	assert.ErrorAs(t, err, &expectedErrorType)
}
