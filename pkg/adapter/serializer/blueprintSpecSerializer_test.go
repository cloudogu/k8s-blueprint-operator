package serializer

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerializeBlueprint_ok(t *testing.T) {
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
				Dogus: []domain.TargetDogu{
					{Namespace: "official", Name: "nginx", Version: "1.2.0-1", TargetState: domain.TargetStatePresent},
					{Namespace: "premium", Name: "jira", Version: "3.0.2-2", TargetState: domain.TargetStateAbsent},
				},
			}},
			`{"blueprintApi":"v2","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present"},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent"}]}`,
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: domain.Blueprint{
				Components: []domain.Component{
					{Name: "blueprint-operator", Version: "0.2.1-1", TargetState: domain.TargetStatePresent},
					{Name: "dogu-operator", Version: "3.2.1-1", TargetState: domain.TargetStateAbsent},
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
			got, err := SerializeBlueprint(tt.args.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("SerializeBlueprint(%v)", tt.args.spec)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SerializeBlueprint(%v)", tt.args.spec)
		})
	}
}

func TestSerializeBlueprint_error(t *testing.T) {
	blueprint := domain.Blueprint{
		Dogus: []domain.TargetDogu{
			{Namespace: "official", Name: "nginx", Version: "1.2.0-1", TargetState: -1},
		},
	}

	_, err := SerializeBlueprint(blueprint)

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot serialize blueprint: ")
	assert.ErrorContains(t, err, "unknown target state ID: '-1'")
}
