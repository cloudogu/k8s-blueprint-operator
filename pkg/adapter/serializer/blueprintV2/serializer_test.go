package blueprintV2

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

var (
	version1201, _ = core.ParseVersion("1.2.0-1")
	version3022, _ = core.ParseVersion("3.0.2-2")
)

var (
	compVersion0211 = semver.MustParse("0.2.1-1")
)

func TestSerializeBlueprint_ok(t *testing.T) {
	serializer := Serializer{}
	quantity := resource.MustParse("2Gi")
	type args struct {
		spec v2.Blueprint
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"empty blueprint",
			args{spec: v2.Blueprint{}},
			`{"blueprintApi":"v2","config":{"global":{}}}`,
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: v2.Blueprint{
				Dogus: []v2.Dogu{
					{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "nginx"}, Version: version1201, TargetState: v2.TargetStatePresent, MinVolumeSize: &quantity, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &quantity, AdditionalConfig: "additional", RewriteTarget: "/"}},
					{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "jira"}, Version: version3022, TargetState: v2.TargetStateAbsent},
				},
			}},
			`{"blueprintApi":"v2","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present","platformConfig":{"resource":{"minVolumeSize":"2Gi"},"reverseProxy":{"maxBodySize":"2Gi","rewriteTarget":"/","additionalConfig":"additional"}}},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent","platformConfig":{"resource":{},"reverseProxy":{}}}],"config":{"global":{}}}`,
			assert.NoError,
		},
		{
			"components in blueprint",
			args{spec: v2.Blueprint{
				Components: []v2.Component{
					{Name: v2.QualifiedComponentName{Namespace: "k8s", SimpleName: "blueprint-operator"}, Version: compVersion0211, TargetState: v2.TargetStatePresent},
					{Name: v2.QualifiedComponentName{Namespace: "k8s", SimpleName: "dogu-operator"}, Version: compVersion3211, TargetState: v2.TargetStateAbsent, DeployConfig: map[string]interface{}{"deployNamespace": "ecosystem", "overwriteConfig": map[string]string{"key": "value"}}},
				},
			}},
			`{"blueprintApi":"v2","components":[{"name":"k8s/blueprint-operator","version":"0.2.1-1","targetState":"present"},{"name":"k8s/dogu-operator","version":"","targetState":"absent","deployConfig":{"deployNamespace":"ecosystem","overwriteConfig":{"key":"value"}}}],"config":{"global":{}}}`,
			assert.NoError,
		},
		{
			"regular dogu config in blueprint",
			args{spec: v2.Blueprint{
				Config: v2.Config{
					Dogus: map[cescommons.SimpleName]v2.CombinedDoguConfig{
						"ldap": {
							Config: v2.DoguConfig{
								Present: map[v2.DoguConfigKey]v2.DoguConfigValue{
									{
										DoguName: "ldap",
										Key:      "container_config/memory_limit",
									}: "500m",
									{
										DoguName: "ldap",
										Key:      "container_config/swap_limit",
									}: "500m",
									{
										DoguName: "ldap",
										Key:      "password_change/notification_enabled",
									}: "true",
								},
								Absent: []v2.DoguConfigKey{
									{
										DoguName: "ldap",
										Key:      "password_change/mail_subject",
									},
									{
										DoguName: "ldap",
										Key:      "password_change/mail_text",
									},
									{
										DoguName: "ldap",
										Key:      "user_search_size_limit",
									},
								},
							},
						},
					},
				},
			}},
			`{"blueprintApi":"v2","config":{"dogus":{"ldap":{"config":{"present":{"container_config/memory_limit":"500m","container_config/swap_limit":"500m","password_change/notification_enabled":"true"},"absent":["password_change/mail_subject","password_change/mail_text","user_search_size_limit"]},"sensitiveConfig":{}}},"global":{}}}`,
			assert.NoError,
		},
		{
			"sensitive dogu config in blueprint",
			args{spec: v2.Blueprint{
				Config: v2.Config{
					Dogus: map[cescommons.SimpleName]v2.CombinedDoguConfig{
						"redmine": {
							SensitiveConfig: v2.SensitiveDoguConfig{
								Present: map[v2.SensitiveDoguConfigKey]v2.SensitiveDoguConfigValue{
									{
										DoguName: "redmine",
										Key:      "my-secret-password",
									}: "password-value",
									{
										DoguName: "redmine",
										Key:      "my-secret-password-2",
									}: "password-value-2",
								},
								Absent: []v2.SensitiveDoguConfigKey{{
									DoguName: "redmine",
									Key:      "my-secret-password-3",
								}},
							},
						},
					},
				},
			}},
			`{"blueprintApi":"v2","config":{"dogus":{"redmine":{"config":{},"sensitiveConfig":{"present":{"my-secret-password":"password-value","my-secret-password-2":"password-value-2"},"absent":["my-secret-password-3"]}}},"global":{}}}`,
			assert.NoError,
		},
		{
			"global config in blueprint",
			args{spec: v2.Blueprint{
				Config: v2.Config{
					Global: v2.GlobalConfig{
						Present: map[v2.GlobalConfigKey]v2.GlobalConfigValue{
							"key_provider": "pkcs1v15",
							"fqdn":         "ces.example.com",
							"admin_group":  "ces-admin",
						},
						Absent: []v2.GlobalConfigKey{
							"default_dogu",
							"some_other_key",
						},
					},
				},
			}},
			`{"blueprintApi":"v2","config":{"global":{"present":{"admin_group":"ces-admin","fqdn":"ces.example.com","key_provider":"pkcs1v15"},"absent":["default_dogu","some_other_key"]}}}`,
			assert.NoError,
		},
		{
			name: "component config",
			args: args{
				spec: v2.Blueprint{
					Components: []v2.Component{
						{Name: v2.QualifiedComponentName{SimpleName: "name", Namespace: "k8s"}, Version: compVersion3211, DeployConfig: map[string]interface{}{"key": "value"}},
					},
				},
			},
			want:    "{\"blueprintApi\":\"v2\",\"components\":[{\"name\":\"k8s/name\",\"version\":\"3.2.1-1\",\"targetState\":\"present\",\"deployConfig\":{\"key\":\"value\"}}],\"config\":{\"global\":{}}}",
			wantErr: assert.NoError,
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
	blueprint := v2.Blueprint{
		Dogus: []v2.Dogu{
			{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "nginx"}, Version: version1201, TargetState: -1},
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
		want    v2.Blueprint
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"empty blueprint",
			args{spec: `{"blueprintApi":"v2"}`},
			v2.Blueprint{},
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: `{"blueprintApi":"v2","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present"},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent"}]}`},
			v2.Blueprint{
				Dogus: []v2.Dogu{
					{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "nginx"}, Version: version1201, TargetState: v2.TargetStatePresent},
					{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "jira"}, Version: version3022, TargetState: v2.TargetStateAbsent},
				}},
			assert.NoError,
		},
		{
			"components in blueprint",
			args{spec: `{"blueprintApi":"v2","components":[{"name":"k8s/blueprint-operator","version":"0.2.1-1","targetState":"present"},{"name":"k8s/dogu-operator","version":"3.2.1-1","targetState":"absent"}]}`},
			v2.Blueprint{
				Components: []v2.Component{
					{Name: v2.QualifiedComponentName{Namespace: "k8s", SimpleName: "blueprint-operator"}, Version: compVersion0211, TargetState: v2.TargetStatePresent},
					{Name: v2.QualifiedComponentName{Namespace: "k8s", SimpleName: "dogu-operator"}, Version: compVersion3211, TargetState: v2.TargetStateAbsent},
				},
			},
			assert.NoError,
		},
		{
			"regular dogu config in blueprint",
			args{spec: `{"blueprintApi":"v2","config":{"dogus":{"ldap":{"config":{"present":{"container_config/memory_limit":"500m","container_config/swap_limit":"500m","password_change/notification_enabled":"true"},"absent":["password_change/mail_subject","password_change/mail_text","user_search_size_limit"]}}}}}`},
			v2.Blueprint{
				Config: v2.Config{
					Dogus: map[cescommons.SimpleName]v2.CombinedDoguConfig{
						"ldap": {
							DoguName: "ldap",
							Config: v2.DoguConfig{
								Present: map[v2.DoguConfigKey]v2.DoguConfigValue{
									{
										DoguName: "ldap",
										Key:      "container_config/memory_limit",
									}: "500m",
									{
										DoguName: "ldap",
										Key:      "container_config/swap_limit",
									}: "500m",
									{
										DoguName: "ldap",
										Key:      "password_change/notification_enabled",
									}: "true",
								},
								Absent: []v2.DoguConfigKey{
									{
										DoguName: "ldap",
										Key:      "password_change/mail_subject",
									},
									{
										DoguName: "ldap",
										Key:      "password_change/mail_text",
									},
									{
										DoguName: "ldap",
										Key:      "user_search_size_limit",
									},
								},
							},
						},
					},
				},
			},
			assert.NoError,
		},
		{
			"sensitive dogu config in blueprint",
			args{spec: `{"blueprintApi":"v2","config":{"dogus":{"redmine":{"sensitiveConfig":{"present":{"my-secret-password":"password-value","my-secret-password-2":"password-value-2"},"absent":["my-secret-password-3"]}}}}}`},
			v2.Blueprint{
				Config: v2.Config{
					Dogus: map[cescommons.SimpleName]v2.CombinedDoguConfig{
						"redmine": {
							DoguName: "redmine",
							SensitiveConfig: v2.SensitiveDoguConfig{
								Present: map[v2.SensitiveDoguConfigKey]v2.SensitiveDoguConfigValue{
									{
										DoguName: "redmine",
										Key:      "my-secret-password",
									}: "password-value",
									{
										DoguName: "redmine",
										Key:      "my-secret-password-2",
									}: "password-value-2",
								},
								Absent: []v2.SensitiveDoguConfigKey{{
									DoguName: "redmine",
									Key:      "my-secret-password-3",
								}},
							},
						},
					},
				},
			},
			assert.NoError,
		},
		{
			"global config in blueprint",
			args{spec: `{"blueprintApi":"v2","config":{"global":{"present":{"admin_group":"ces-admin","fqdn":"ces.example.com","key_provider":"pkcs1v15"},"absent":["default_dogu","some_other_key"]}}}`},
			v2.Blueprint{
				Config: v2.Config{
					Global: v2.GlobalConfig{
						Present: map[v2.GlobalConfigKey]v2.GlobalConfigValue{
							"key_provider": "pkcs1v15",
							"fqdn":         "ces.example.com",
							"admin_group":  "ces-admin",
						},
						Absent: []v2.GlobalConfigKey{
							"default_dogu",
							"some_other_key",
						},
					},
				},
			},
			assert.NoError,
		},
		{
			"component package config",
			args{spec: "{\"blueprintApi\":\"v2\",\"components\":[{\"name\":\"k8s/name\",\"version\":\"3.2.1-1\",\"targetState\":\"present\",\"deployConfig\":{\"key\":\"value\"}}],\"config\":{\"global\":{}}}"},
			v2.Blueprint{
				Components: []v2.Component{
					{Name: v2.QualifiedComponentName{SimpleName: "name", Namespace: "k8s"}, Version: compVersion3211, DeployConfig: map[string]interface{}{"key": "value"}},
				},
			},
			assert.NoError,
		},
		{
			"component package config error",
			args{spec: "{\"blueprintApi\":\"v2\",\"components\":[{\"name\":\"k8s/name\",\"version\":\"3.2.1-1\",\"targetState\":\"present\",\"deployConfig\":{\"key\"\":\"\"::\"value:{\"}}],\"config\":{\"global\":{}}}"},
			v2.Blueprint{},
			assert.Error,
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
