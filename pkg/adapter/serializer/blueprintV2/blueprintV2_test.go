package blueprintV2

import (
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

func Test_ConvertToBlueprintV2(t *testing.T) {
	dogus := []domain.Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: "component1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: "absent/component2", TargetState: domain.TargetStateAbsent},
		{Name: "present-component3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: "present/component4", Version: version1_2_3_3},
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

	blueprintV2, err := ConvertToBlueprintV2(blueprint)

	convertedDogus := []serializer.TargetDogu{
		{Name: "absent/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	convertedComponents := []serializer.TargetComponent{
		{Name: "component1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present-component3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/component4", Version: version1_2_3_3.Raw, TargetState: "present"},
	}

	require.Nil(t, err)
	assert.Equal(t, BlueprintV2{
		GeneralBlueprint: serializer.GeneralBlueprint{API: serializer.V2},
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
	dogus := []serializer.TargetDogu{
		{Name: "absent/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/dogu2", TargetState: "absent"},
		{Name: "present/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/dogu4", Version: version1_2_3_3.Raw},
	}

	components := []serializer.TargetComponent{
		{Name: "component1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present-component3", Version: version3212.Raw, TargetState: "present"},
		{Name: "present/component4", Version: version1_2_3_3.Raw},
	}

	blueprintV2 := BlueprintV2{
		GeneralBlueprint: serializer.GeneralBlueprint{API: serializer.V2},
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

	convertedDogus := []domain.Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: domain.TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version1_2_3_3},
	}

	convertedComponents := []domain.Component{
		{Name: "component1", Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: "absent/component2", TargetState: domain.TargetStateAbsent},
		{Name: "present-component3", Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: "present/component4", Version: version1_2_3_3},
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

func Test_ConvertToBlueprint_errors(t *testing.T) {
	blueprintV2 := BlueprintV2{
		GeneralBlueprint: serializer.GeneralBlueprint{API: serializer.V2},
		Dogus: []serializer.TargetDogu{
			{Name: "dogu1", Version: version3211.Raw},
			{Name: "official/dogu1", Version: version3211.Raw, TargetState: "unknown"},
			{Name: "name/space/dogu2", Version: version3212.Raw},
			{Name: "official/dogu3", Version: "abc"},
		},
		Components: []serializer.TargetComponent{
			{Name: "component1", Version: version3211.Raw, TargetState: "not known state"},
			{Name: "official/dogu3", Version: "abc"},
		},
	}

	_, err := convertToBlueprint(blueprintV2)

	require.ErrorContains(t, err, "syntax of blueprintV2 is not correct: ")

	require.ErrorContains(t, err, "cannot convert blueprint dogus: ")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'dogu1'")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'name/space/dogu2'")
	require.ErrorContains(t, err, "unknown target state \"unknown\"")

	require.ErrorContains(t, err, "cannot convert blueprint components: ")
	require.ErrorContains(t, err, "unknown target state \"not known state\"")
	require.ErrorContains(t, err, "could not parse version of target dogu \"official/dogu3\": failed to parse major version abc")
	require.ErrorContains(t, err, "could not parse version of target component \"official/dogu3\": failed to parse major version abc")
}
