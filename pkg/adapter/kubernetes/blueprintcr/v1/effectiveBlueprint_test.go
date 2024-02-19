package v1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

func TestConvertToEffectiveBlueprint(t *testing.T) {
	//given
	dogus := []domain.Dogu{
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "premium", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "premium", SimpleName: "dogu4"}, Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: common.QualifiedComponentName{SimpleName: "component1", Namespace: "k8s"}, Version: nil, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{SimpleName: "component2", Namespace: "k8s"}, Version: nil, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{SimpleName: "component3", Namespace: "k8s-testing"}, Version: compVersion3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedComponentName{SimpleName: "component4", Namespace: "k8s-testing"}, Version: compVersion1233},
	}
	blueprint := domain.EffectiveBlueprint{
		Dogus:      dogus,
		Components: components,
		Config: domain.Config{
			Dogus: map[common.SimpleDoguName]domain.CombinedDoguConfig{
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
						Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
							{DoguConfigKey: common.DoguConfigKey{
								DoguName: "my-dogu",
								Key:      "config-encrypted",
							}}: "42",
						},
					},
				},
			},
			Global: domain.GlobalConfig{Absent: []common.GlobalConfigKey{"test/key"}},
		},
	}

	//when
	blueprintV2, err := ConvertToEffectiveBlueprintDTO(blueprint)

	//then
	convertedDogus := []serializer.TargetDogu{
		{Name: "official/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "official/dogu2", TargetState: "absent"},
		{Name: "premium/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "premium/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []serializer.TargetComponent{
		{Name: "k8s/component1", TargetState: "absent"},
		{Name: "k8s/component2", TargetState: "absent"},
		{Name: "k8s-testing/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "k8s-testing/component4", Version: compVersion1233.String(), TargetState: "present"},
	}

	require.NoError(t, err)
	assert.Equal(t, EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config: Config{
			Dogus: map[string]CombinedDoguConfig{
				"my-dogu": {
					Config: DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: SensitiveDoguConfig{
						Present: map[string]string{
							"config-encrypted": "42",
						},
					},
				},
			},
			Global: GlobalConfig{
				Absent: []string{"test/key"},
			},
		},
	}, blueprintV2)
}

func TestConvertToEffectiveBlueprintV1(t *testing.T) {
	//given
	convertedDogus := []serializer.TargetDogu{
		{Name: "official/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "official/dogu2", TargetState: "absent"},
		{Name: "premium/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "premium/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []serializer.TargetComponent{
		{Name: "k8s/component1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "k8s/component2", TargetState: "absent"},
		{Name: "k8s-testing/component3", Version: version3212.Raw, TargetState: "present"},
		{Name: "k8s-testing/component4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	dto := EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		Config: Config{
			Dogus: map[string]CombinedDoguConfig{
				"my-dogu": {
					Config: DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: SensitiveDoguConfig{
						Present: map[string]string{
							"config-encrypted": "42",
						},
					},
				},
			},
			Global: GlobalConfig{Absent: []string{"test/key"}},
		},
	}
	//when
	blueprint, err := ConvertToEffectiveBlueprintDomain(dto)

	//then
	dogus := []domain.Dogu{
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "premium", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "premium", SimpleName: "dogu4"}, Version: version1_2_3_3},
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
			Dogus: map[common.SimpleDoguName]domain.CombinedDoguConfig{
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
						Present: map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{
							{DoguConfigKey: common.DoguConfigKey{
								DoguName: "my-dogu",
								Key:      "config-encrypted",
							}}: "42",
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

func Test_setKey(t *testing.T) {
	t.Run("test simple key", func(t *testing.T) {
		theMap := map[string]interface{}{}

		setKey([]string{"key1"}, "val", theMap)

		assert.Equal(t, map[string]interface{}{
			"key1": "val",
		}, theMap)
	})

	t.Run("with filled map", func(t *testing.T) {
		theMap := map[string]interface{}{
			"key1": "val",
		}

		setKey([]string{"key2"}, "val", theMap)

		assert.Equal(t, map[string]interface{}{
			"key1": "val",
			"key2": "val",
		}, theMap)
	})

	t.Run("depth 2", func(t *testing.T) {
		theMap := map[string]interface{}{}

		setKey([]string{"key1", "key2"}, "val", theMap)
		setKey([]string{"key1", "key2"}, "val", theMap)

		assert.Equal(t, map[string]interface{}{
			"key1": map[string]interface{}{
				"key2": "val",
			},
		}, theMap)
	})
}

func Test_widenMap1(t *testing.T) {
	type args struct {
		currentMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "nil",
			args: args{
				currentMap: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "depth 3",
			args: args{
				currentMap: map[string]string{
					"1/2/3.1":     "v1",
					"1/2/3.2":     "v2",
					"1/2/3.3/4.1": "v3",
					"1/2/3.3/4.2": "v4",
				},
			},
			want: map[string]interface{}{
				"1": map[string]interface{}{
					"2": map[string]interface{}{
						"3.1": "v1",
						"3.2": "v2",
						"3.3": map[string]interface{}{
							"4.1": "v3",
							"4.2": "v4"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, widenMap(tt.args.currentMap), "widenMap(%v)", tt.args.currentMap)
		})
	}
}
