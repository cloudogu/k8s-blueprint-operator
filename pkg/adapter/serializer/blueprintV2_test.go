package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ConvertToBlueprintV2(t *testing.T) {
	dogus := []domain.TargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: "1.2.3-3"},
	}

	components := []domain.Component{
		{Name: "component1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/component2", TargetState: TargetStateAbsent},
		{Name: "present-component3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/component4", Version: "1.2.3-3"},
	}
	blueprint := domain.Blueprint{
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

	blueprintV2 := ConvertToBlueprintV2(blueprint)

	convertedDogus := []TargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/dogu4", Version: "1.2.3-3"},
	}

	convertedComponents := []TargetComponent{
		{Name: "component1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/component2", TargetState: TargetStateAbsent},
		{Name: "present-component3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/component4", Version: "1.2.3-3"},
	}

	assert.Equal(t, BlueprintV2{
		GeneralBlueprint: GeneralBlueprint{V2},
		Dogus:            convertedDogus,
		Components:       convertedComponents,
		RegistryConfig: RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
		RegistryConfigAbsent: []string{"_global/test/key"},
		RegistryConfigEncrypted: RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
	}, blueprintV2)
}

func Test_ConvertToBlueprint(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/dogu4", Version: "1.2.3-3"},
	}

	components := []TargetComponent{
		{Name: "component1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/component2", TargetState: TargetStateAbsent},
		{Name: "present-component3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/component4", Version: "1.2.3-3"},
	}

	blueprintV2 := BlueprintV2{
		GeneralBlueprint: GeneralBlueprint{V2},
		Dogus:            dogus,
		Components:       components,
		RegistryConfig: RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
		RegistryConfigAbsent: []string{"_global/test/key"},
		RegistryConfigEncrypted: RegistryConfig{
			"dogu": map[string]interface{}{
				"config": "42",
			},
		},
	}
	blueprint, err := convertToBlueprint(blueprintV2)

	require.Nil(t, err)

	convertedDogus := []domain.TargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: "1.2.3-3"},
	}

	convertedComponents := []domain.Component{
		{Name: "component1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
		{Name: "absent/component2", TargetState: TargetStateAbsent},
		{Name: "present-component3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/component4", Version: "1.2.3-3"},
	}

	assert.Equal(t, domain.Blueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
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
	}, blueprint)
}

func Test_ConvertToBlueprint_invalidDoguName(t *testing.T) {
	blueprintV2 := BlueprintV2{
		GeneralBlueprint: GeneralBlueprint{V2},
		Dogus: []TargetDogu{
			{Name: "dogu1", Version: "3.2.1-1"},
			{Name: "name/space/dogu2", Version: "3.2.1-2"},
		},
	}

	_, err := convertToBlueprint(blueprintV2)

	require.ErrorContains(t, err, "syntax of blueprintV2 is not correct")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'dogu1'")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'name/space/dogu2'")
}
