package serializer

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

var (
	version3211, _    = core.ParseVersion("3.2.1-1")
	version3212, _    = core.ParseVersion("3.2.1-2")
	version1_2_3_3, _ = core.ParseVersion("1.2.3-3")
)

var (
	compVersion1233 = semver.MustParse("1.2.3-3")
	compVersion3211 = semver.MustParse("3.2.1-1")
	compVersion3212 = semver.MustParse("3.2.1-2")
)

var emptyPlatformConfig = crd.PlatformConfig{
	ResourceConfig: crd.ResourceConfig{
		MinVolumeSize: "0",
	},
}

func TestConvertToBlueprintDTO(t *testing.T) {
	t.Run("convert dogus", func(t *testing.T) {
		dogus := []domain.Dogu{
			{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
			{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		}
		blueprint := domain.EffectiveBlueprint{
			Dogus: dogus,
		}

		//when
		blueprintV2 := ConvertToBlueprintDTO(blueprint)

		//then
		convertedDogus := []crd.Dogu{
			{Name: "official/dogu1", PlatformConfig: emptyPlatformConfig, Version: version3211.Raw, Absent: true},
			{Name: "premium/dogu3", PlatformConfig: emptyPlatformConfig, Version: version3212.Raw, Absent: false},
		}

		assert.Equal(t, crd.BlueprintManifest{
			Dogus:      convertedDogus,
			Components: []crd.Component{},
		}, blueprintV2)
	})

	t.Run("convert components", func(t *testing.T) {
		components := []domain.Component{
			{Name: common.QualifiedComponentName{SimpleName: "component1", Namespace: "k8s"}, Version: nil, TargetState: domain.TargetStateAbsent},
			{Name: common.QualifiedComponentName{SimpleName: "component3", Namespace: "k8s-testing"}, Version: compVersion3212, TargetState: domain.TargetStatePresent},
		}
		blueprint := domain.EffectiveBlueprint{
			Components: components,
		}

		//when
		blueprintV2 := ConvertToBlueprintDTO(blueprint)

		//then
		convertedComponents := []crd.Component{
			{Name: "k8s/component1", Absent: true},
			{Name: "k8s-testing/component3", Version: compVersion3212.String(), Absent: false},
		}

		assert.Equal(t, crd.BlueprintManifest{
			Dogus:      []crd.Dogu{},
			Components: convertedComponents,
		}, blueprintV2)
	})

	t.Run("convert config", func(t *testing.T) {
		blueprint := domain.EffectiveBlueprint{
			Config: domain.Config{
				Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{
					"my-dogu": {
						Config: domain.DoguConfig{
							Present: map[common.DoguConfigKey]common.DoguConfigValue{
								{
									DoguName: "my-dogu",
									Key:      "config",
								}: "42",
							},
						},
						SensitiveConfig: domain.SensitiveDoguConfig{
							Present: map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef{
								{
									DoguName: "my-dogu",
									Key:      "sensitive-config",
								}: {
									SecretName: "mySecret",
									SecretKey:  "myKey",
								},
							},
						},
					},
				},
				Global: domain.GlobalConfig{Absent: []common.GlobalConfigKey{"test/key"}},
			},
		}
		blueprintV2 := ConvertToBlueprintDTO(blueprint)

		assert.Equal(t, crd.BlueprintManifest{
			Dogus:      []crd.Dogu{},
			Components: []crd.Component{},
			Config: crd.Config{
				Dogus: map[string]crd.CombinedDoguConfig{
					"my-dogu": {
						Config: &crd.DoguConfig{
							Present: map[string]string{
								"config": "42",
							},
						},
						SensitiveConfig: &crd.SensitiveDoguConfig{
							Present: []crd.SensitiveConfigEntry{
								{
									Key:        "sensitive-config",
									SecretName: "mySecret",
									SecretKey:  "myKey",
								},
							},
						},
					},
				},
				Global: crd.GlobalConfig{
					Absent: []string{"test/key"},
				},
			},
		}, blueprintV2)
	})
}

func TestConvertToEffectiveBlueprintDomain(t *testing.T) {
	//given
	convertedDogus := []crd.Dogu{
		{Name: "official/dogu1", Version: version3211.Raw, Absent: true},
		{Name: "official/dogu2", Absent: true},
		{Name: "premium/dogu3", Version: version3212.Raw, Absent: false},
		{Name: "premium/dogu4", Version: version1_2_3_3.Raw, Absent: false},
		{
			Name:    "premium/dogu5",
			Version: version1_2_3_3.Raw,
			Absent:  false,
			PlatformConfig: crd.PlatformConfig{
				ResourceConfig:     crd.ResourceConfig{},
				ReverseProxyConfig: crd.ReverseProxyConfig{},
				AdditionalMountsConfig: []crd.AdditionalMount{
					{
						SourceType: crd.DataSourceConfigMap,
						Name:       "config",
						Volume:     "volume",
						Subfolder:  "subfolder",
					},
				},
			},
		},
	}

	convertedComponents := []crd.Component{
		{Name: "k8s/component1", Version: version3211.Raw, Absent: true},
		{Name: "k8s/component2", Absent: true},
		{Name: "k8s-testing/component3", Version: version3212.Raw, Absent: false},
		{Name: "k8s-testing/component4", Version: version1_2_3_3.Raw, Absent: false},
	}

	dto := crd.BlueprintManifest{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config: crd.Config{
			Dogus: map[string]crd.CombinedDoguConfig{
				"my-dogu": {
					Config: &crd.DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: &crd.SensitiveDoguConfig{
						Present: []crd.SensitiveConfigEntry{
							{
								Key:        "sensitive-config",
								SecretName: "mySecret",
								SecretKey:  "myKey",
							},
						},
					},
				},
			},
			Global: crd.GlobalConfig{Absent: []string{"test/key"}},
		},
	}
	//when
	blueprint, err := ConvertToEffectiveBlueprintDomain(dto)

	//then
	dogus := []domain.Dogu{
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu4"}, Version: version1_2_3_3},
		{
			Name:    cescommons.QualifiedName{Namespace: "premium", SimpleName: "dogu5"},
			Version: version1_2_3_3,
			AdditionalMounts: []ecosystem.AdditionalMount{
				{
					SourceType: ecosystem.DataSourceConfigMap,
					Name:       "config",
					Volume:     "volume",
					Subfolder:  "subfolder",
				},
			},
		},
	}

	components := []domain.Component{
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component1"}, Version: compVersion3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{Namespace: "k8s-testing", SimpleName: "component3"}, Version: compVersion3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedComponentName{Namespace: "k8s-testing", SimpleName: "component4"}, Version: compVersion1233},
	}
	expected := domain.EffectiveBlueprint{
		Dogus:      dogus,
		Components: components,
		Config: domain.Config{
			Dogus: map[cescommons.SimpleName]domain.CombinedDoguConfig{
				"my-dogu": {
					DoguName: "my-dogu",
					Config: domain.DoguConfig{
						Present: map[common.DoguConfigKey]common.DoguConfigValue{
							{
								DoguName: "my-dogu",
								Key:      "config",
							}: "42",
						},
					},
					SensitiveConfig: domain.SensitiveDoguConfig{
						Present: map[common.SensitiveDoguConfigKey]domain.SensitiveValueRef{
							{
								DoguName: "my-dogu",
								Key:      "sensitive-config",
							}: {
								SecretName: "mySecret",
								SecretKey:  "myKey",
							},
						},
					},
				},
			},
			Global: domain.GlobalConfig{Absent: []common.GlobalConfigKey{"test/key"}},
		},
	}

	require.NoError(t, err)
	assert.Equal(t, expected, blueprint)
}

func TestConvertToBlueprintDomain(t *testing.T) {
	t.Run("convert empty", func(t *testing.T) {
		given := crd.BlueprintManifest{}
		converted, err := ConvertToBlueprintDomain(given)

		require.NoError(t, err)
		assert.Equal(t, domain.Blueprint{}, converted)
	})

	t.Run("error converting", func(t *testing.T) {
		given := crd.BlueprintManifest{
			Dogus: []crd.Dogu{
				{
					Name:    "official/redmine",
					Version: "unparsable",
				},
			},
		}
		_, err := ConvertToBlueprintDomain(given)

		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot deserialize blueprint")
	})

	t.Run("convert dogu", func(t *testing.T) {
		given := crd.BlueprintManifest{
			Dogus: []crd.Dogu{
				{
					Name:    "official/redmine",
					Version: version1_2_3_3.Raw,
				},
			},
		}
		converted, err := ConvertToBlueprintDomain(given)

		require.NoError(t, err)
		assert.Equal(t, domain.Blueprint{
			Dogus: []domain.Dogu{
				{
					Name:        cescommons.QualifiedName{Namespace: "official", SimpleName: "redmine"},
					Version:     version1_2_3_3,
					TargetState: domain.TargetStatePresent,
				},
			},
		}, converted)
	})

	t.Run("convert component", func(t *testing.T) {
		given := crd.BlueprintManifest{
			Components: []crd.Component{
				{
					Name:    "k8s/loki",
					Version: compVersion1233.Original(),
				},
			},
		}
		converted, err := ConvertToBlueprintDomain(given)

		require.NoError(t, err)
		assert.Equal(t, domain.Blueprint{
			Components: []domain.Component{
				{
					Name: common.QualifiedComponentName{
						Namespace:  "k8s",
						SimpleName: "loki",
					},
					Version:     compVersion1233,
					TargetState: domain.TargetStatePresent,
				},
			},
		}, converted)
	})

}

func TestConvertToBlueprintMaskDomain(t *testing.T) {
	type args struct {
		mask crd.BlueprintMask
	}
	tests := []struct {
		name    string
		args    args
		want    domain.BlueprintMask
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			args:    args{mask: crd.BlueprintMask{}},
			want:    domain.BlueprintMask{},
			wantErr: assert.NoError,
		},
		{
			name: "will convert a MaskDogu",
			args: args{mask: crd.BlueprintMask{
				Dogus: []crd.MaskDogu{
					{
						Name:    "official/ldap",
						Version: version1_2_3_3.Raw,
					},
				},
			}},
			want: domain.BlueprintMask{
				Dogus: []domain.MaskDogu{
					{
						Name: cescommons.QualifiedName{
							Namespace:  "official",
							SimpleName: "ldap",
						},
						Version:     version1_2_3_3,
						TargetState: domain.TargetStatePresent,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error if invalid mask",
			args: args{mask: crd.BlueprintMask{
				Dogus: []crd.MaskDogu{
					{
						Name:    "invalid name",
						Version: version1_2_3_3.Raw,
					},
				},
			}},
			want: domain.BlueprintMask{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "cannot deserialize blueprint mask", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToBlueprintMaskDomain(tt.args.mask)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertToBlueprintMaskDomain(%v)", tt.args.mask)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertToBlueprintMaskDomain(%v)", tt.args.mask)
		})
	}
}
