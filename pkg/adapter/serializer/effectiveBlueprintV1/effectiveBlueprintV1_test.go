package effectiveBlueprintV1

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	version3_2_1_1, _ = core.ParseVersion("3.2.1-1")
	version3_2_1_2, _ = core.ParseVersion("3.2.1-2")
	version1_2_3_3, _ = core.ParseVersion("1.2.3-3")
)

func TestConvertToEffectiveBlueprint(t *testing.T) {
	//given
	dogus := []domain.Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3_2_1_1, TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3_2_1_2, TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: "component1", Version: version3_2_1_1, TargetState: domain.TargetStateAbsent},
		{Name: "absent/component2", TargetState: domain.TargetStateAbsent},
		{Name: "present-component3", Version: version3_2_1_2, TargetState: domain.TargetStatePresent},
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
	blueprintV2, err := ConvertToEffectiveBlueprintV1(blueprint)

	//then
	convertedDogus := []TargetDogu{
		{Name: "absent/dogu1", Version: version3_2_1_1.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3_2_1_2.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []TargetComponent{
		{Name: "component1", Version: version3_2_1_1.Raw, TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present-component3", Version: version3_2_1_2.Raw, TargetState: "present"},
		{Name: "present/component4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	require.Nil(t, err)
	assert.Equal(t, EffectiveBlueprintV1{
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
	convertedDogus := []TargetDogu{
		{Name: "absent/dogu1", Version: version3_2_1_1.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3_2_1_2.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []TargetComponent{
		{Name: "component1", Version: version3_2_1_1.Raw, TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present-component3", Version: version3_2_1_2.Raw, TargetState: "present"},
		{Name: "present/component4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	dto := EffectiveBlueprintV1{
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
	blueprint, err := ConvertToEffectiveBlueprint(dto)

	//then
	dogus := []domain.Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3_2_1_1, TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3_2_1_2, TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: "component1", Version: version3_2_1_1, TargetState: domain.TargetStateAbsent},
		{Name: "absent/component2", TargetState: domain.TargetStateAbsent},
		{Name: "present-component3", Version: version3_2_1_2, TargetState: domain.TargetStatePresent},
		{Name: "present/component4", Version: version1_2_3_3},
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

	require.Nil(t, err)
	assert.Equal(t, expected, blueprint)
}

func Test_widenMap(t *testing.T) {
	tests := []struct {
		name  string
		given map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name:  "nil",
			given: nil,
			want:  map[string]interface{}{},
		},
		{
			name:  "no keys",
			given: map[string]interface{}{},
			want:  map[string]interface{}{},
		},
		{
			name: "depth 1",
			given: map[string]interface{}{
				"key": "val",
			},
			want: map[string]interface{}{
				"key": "val",
			},
		},
		{
			name: "depth 2",
			given: map[string]interface{}{
				"key/key2": "val",
			},
			want: map[string]interface{}{
				"key": map[string]interface{}{
					"key2": "val",
				},
			},
		},
		{
			name: "depth 3",
			given: map[string]interface{}{
				"key/key2/key3": "val",
			},
			want: map[string]interface{}{
				"key": map[string]interface{}{
					"key2": map[string]interface{}{
						"key3": "val",
					},
				},
			},
		},
		{
			name: "multiple sub keys",
			given: map[string]interface{}{
				"key1/key2/key3.1": "val",
				"key1/key2/key3.2": "val",
			},
			want: map[string]interface{}{
				"key1": map[string]interface{}{
					"key2": map[string]interface{}{
						"key3.1": "val",
						"key3.2": "val",
					},
				},
			},
		},
		{
			name: "there is a value and a sub map for a key",
			given: map[string]interface{}{
				"key1/key2":      "val",
				"key1/key2/key3": "val",
			},
			want: map[string]interface{}{
				"key1": map[string]interface{}{
					"key2": "val",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, widenMap(tt.given), "widenMap(%v)", tt.given)
		})
	}
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
