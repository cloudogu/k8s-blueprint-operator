package v1

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	version3211, _    = core.ParseVersion("3.2.1-1")
	version3212, _    = core.ParseVersion("3.2.1-2")
	version1_2_3_3, _ = core.ParseVersion("1.2.3-3")
)

func TestConvertToEffectiveBlueprint(t *testing.T) {
	//given
	dogus := []domain.Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: "absent/component1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: "absent/component2", TargetState: domain.TargetStateAbsent},
		{Name: "present/component3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: "present/component4", Version: version1_2_3_3},
	}
	blueprint := domain.EffectiveBlueprint{
		Dogus:      dogus,
		Components: components,
		RegistryConfig: domain.RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
		RegistryConfigAbsent: []string{"_global/test/key"},
		RegistryConfigEncrypted: domain.RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
	}

	//when
	blueprintV2, err := ConvertToEffectiveBlueprintDTO(blueprint)

	//then
	convertedDogus := []serializer.TargetDogu{
		{Name: "absent/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []serializer.TargetComponent{
		{Name: "absent/component1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present/component3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/component4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	require.NoError(t, err)
	assert.Equal(t, EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		RegistryConfig: map[string]string{
			"dogu/config": "42",
		},
		RegistryConfigAbsent: []string{"_global/test/key"},
		RegistryConfigEncrypted: map[string]string{
			"dogu/config": "42",
		},
	}, blueprintV2)
}

func TestConvertToEffectiveBlueprintV1(t *testing.T) {
	//given
	convertedDogus := []serializer.TargetDogu{
		{Name: "absent/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []serializer.TargetComponent{
		{Name: "absent/component1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present/component3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/component4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	dto := EffectiveBlueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
		RegistryConfig: map[string]string{
			"dogu/config": "42",
		},
		RegistryConfigAbsent: []string{"_global/test/key"},
		RegistryConfigEncrypted: map[string]string{
			"dogu/config": "42",
		},
	}
	//when
	blueprint, err := ConvertToEffectiveBlueprintDomain(dto)

	//then
	dogus := []domain.Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: "component1", DistributionNamespace: "absent", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: "component2", DistributionNamespace: "absent", TargetState: domain.TargetStateAbsent},
		{Name: "component3", DistributionNamespace: "present", Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: "component4", DistributionNamespace: "present", Version: version1_2_3_3},
	}
	expected := domain.EffectiveBlueprint{
		Dogus:      dogus,
		Components: components,
		RegistryConfig: domain.RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
		RegistryConfigAbsent: []string{"_global/test/key"},
		RegistryConfigEncrypted: domain.RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
	}

	require.NoError(t, err)
	assert.Equal(t, expected, blueprint)
}

func Test_widenMap(t *testing.T) {
}

func Test_convertToRegistryConfig(t *testing.T) {
	tests := []struct {
		name    string
		dto     map[string]string
		want    domain.RegistryConfig
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			dto:     nil,
			want:    domain.RegistryConfig{},
			wantErr: assert.NoError,
		},
		{
			name:    "no keys",
			dto:     map[string]string{},
			want:    domain.RegistryConfig{},
			wantErr: assert.NoError,
		},
		{
			name: "depth 1",
			dto: map[string]string{
				"key1": "val1",
			},
			want:    domain.RegistryConfig{},
			wantErr: assert.Error,
		},
		{
			name: "multiple deep keys",
			dto: map[string]string{
				"key1/key2/key3.1": "val1",
				"key1/key2/key3.2": "val2",
			},
			want: domain.RegistryConfig{
				"key1": {
					"key2": map[string]interface{}{
						"key3.1": "val1",
						"key3.2": "val2",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToRegistryConfig(tt.dto)
			if !tt.wantErr(t, err, fmt.Sprintf("convertToRegistryConfig(%v)", tt.dto)) {
				return
			}
			assert.Equalf(t, tt.want, got, "convertToRegistryConfig(%v)", tt.dto)
		})
	}
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
