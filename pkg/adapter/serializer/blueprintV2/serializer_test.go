package blueprintV2

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
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
			`{"blueprintApi":"v2","config":{"global":{}}}`,
			assert.NoError,
		},
		{
			"dogus in blueprint",
			args{spec: domain.Blueprint{
				Dogus: []domain.Dogu{
					{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "nginx"}, Version: version1201, TargetState: domain.TargetStatePresent, MinVolumeSize: &quantity, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &quantity, AdditionalConfig: "additional", RewriteTarget: "/"}},
					{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "jira"}, Version: version3022, TargetState: domain.TargetStateAbsent},
				},
			}},
			`{"blueprintApi":"v2","dogus":[{"name":"official/nginx","version":"1.2.0-1","targetState":"present","platformConfig":{"resource":{"minVolumeSize":"2Gi"},"reverseProxy":{"maxBodySize":"2Gi","rewriteTarget":"/","additionalConfig":"additional"}}},{"name":"premium/jira","version":"3.0.2-2","targetState":"absent","platformConfig":{"resource":{},"reverseProxy":{}}}],"config":{"global":{}}}`,
			assert.NoError,
		},
		{
			"components in blueprint",
			args{spec: domain.Blueprint{
				Components: []domain.Component{
					{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "blueprint-operator"}, Version: compVersion0211, TargetState: domain.TargetStatePresent},
					{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "dogu-operator"}, Version: compVersion3211, TargetState: domain.TargetStateAbsent, DeployConfig: map[string]interface{}{"deployNamespace": "ecosystem", "overwriteConfig": map[string]string{"key": "value"}}},
				},
			}},
			`{"blueprintApi":"v2","components":[{"name":"k8s/blueprint-operator","version":"0.2.1-1","targetState":"present"},{"name":"k8s/dogu-operator","version":"","targetState":"absent","deployConfig":{"deployNamespace":"ecosystem","overwriteConfig":{"key":"value"}}}],"config":{"global":{}}}`,
			assert.NoError,
		},
		{
			"regular dogu config in blueprint",
			args{spec: domain.Blueprint{
				Config: domain.Config{
					Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{
						"ldap": {
							Config: domain.DoguConfig{
								Present: map[common.DoguConfigKey]common.DoguConfigValue{
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
								Absent: []common.DoguConfigKey{
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
			args{spec: domain.Blueprint{
				Config: domain.Config{
					Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{
						"redmine": {
							SensitiveConfig: domain.SensitiveDoguConfig{
								Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
									{
										DoguName: "redmine",
										Key:      "my-secret-password",
									}: "password-value",
									{
										DoguName: "redmine",
										Key:      "my-secret-password-2",
									}: "password-value-2",
								},
								Absent: []common.SensitiveDoguConfigKey{{
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
			args{spec: domain.Blueprint{
				Config: domain.Config{
					Global: domain.GlobalConfig{
						Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
							"key_provider": "pkcs1v15",
							"fqdn":         "ces.example.com",
							"admin_group":  "ces-admin",
						},
						Absent: []common.GlobalConfigKey{
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
				spec: domain.Blueprint{
					Components: []domain.Component{
						{Name: common.QualifiedComponentName{SimpleName: "name", Namespace: "k8s"}, Version: compVersion3211, DeployConfig: map[string]interface{}{"key": "value"}},
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
	blueprint := domain.Blueprint{
		Dogus: []domain.Dogu{
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
					{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "nginx"}, Version: version1201, TargetState: domain.TargetStatePresent},
					{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "jira"}, Version: version3022, TargetState: domain.TargetStateAbsent},
				}},
			assert.NoError,
		},
		{
			"components in blueprint",
			args{spec: `{"blueprintApi":"v2","components":[{"name":"k8s/blueprint-operator","version":"0.2.1-1","targetState":"present"},{"name":"k8s/dogu-operator","version":"3.2.1-1","targetState":"absent"}]}`},
			domain.Blueprint{
				Components: []domain.Component{
					{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "blueprint-operator"}, Version: compVersion0211, TargetState: domain.TargetStatePresent},
					{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "dogu-operator"}, Version: compVersion3211, TargetState: domain.TargetStateAbsent},
				},
			},
			assert.NoError,
		},
		{
			"regular dogu config in blueprint",
			args{spec: `{"blueprintApi":"v2","config":{"dogus":{"ldap":{"config":{"present":{"container_config/memory_limit":"500m","container_config/swap_limit":"500m","password_change/notification_enabled":"true"},"absent":["password_change/mail_subject","password_change/mail_text","user_search_size_limit"]}}}}}`},
			domain.Blueprint{
				Config: domain.Config{
					Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{
						"ldap": {
							DoguName: "ldap",
							Config: domain.DoguConfig{
								Present: map[common.DoguConfigKey]common.DoguConfigValue{
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
								Absent: []common.DoguConfigKey{
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
			domain.Blueprint{
				Config: domain.Config{
					Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{
						"redmine": {
							DoguName: "redmine",
							SensitiveConfig: domain.SensitiveDoguConfig{
								Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
									{
										DoguName: "redmine",
										Key:      "my-secret-password",
									}: "password-value",
									{
										DoguName: "redmine",
										Key:      "my-secret-password-2",
									}: "password-value-2",
								},
								Absent: []common.SensitiveDoguConfigKey{{
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
			domain.Blueprint{
				Config: domain.Config{
					Global: domain.GlobalConfig{
						Present: map[common.GlobalConfigKey]common.GlobalConfigValue{
							"key_provider": "pkcs1v15",
							"fqdn":         "ces.example.com",
							"admin_group":  "ces-admin",
						},
						Absent: []common.GlobalConfigKey{
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
			domain.Blueprint{
				Components: []domain.Component{
					{Name: common.QualifiedComponentName{SimpleName: "name", Namespace: "k8s"}, Version: compVersion3211, DeployConfig: map[string]interface{}{"key": "value"}},
				},
			},
			assert.NoError,
		},
		{
			"component package config error",
			args{spec: "{\"blueprintApi\":\"v2\",\"components\":[{\"name\":\"k8s/name\",\"version\":\"3.2.1-1\",\"targetState\":\"present\",\"deployConfig\":{\"key\"\":\"\"::\"value:{\"}}],\"config\":{\"global\":{}}}"},
			domain.Blueprint{},
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
