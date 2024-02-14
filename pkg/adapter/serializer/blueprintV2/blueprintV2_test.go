package blueprintV2

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
	compVersion3213 = semver.MustParse("3.2.1-3")
)

func Test_ConvertToBlueprintV2(t *testing.T) {
	dogus := []domain.Dogu{
		{Name: common.QualifiedDoguName{Namespace: "absent", Name: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "absent", Name: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "present", Name: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "present", Name: "dogu4"}, Version: version1_2_3_3},
	}

	components := []domain.Component{
		{Name: "component1", DistributionNamespace: "absent", Version: compVersion3211, TargetState: domain.TargetStateAbsent},
		{Name: "component2", DistributionNamespace: "absent", TargetState: domain.TargetStateAbsent},
		{Name: "component3", DistributionNamespace: "present", Version: compVersion3212, TargetState: domain.TargetStatePresent},
		{Name: "component4", DistributionNamespace: "present", Version: compVersion3213},
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
		{Name: "absent/component1", Version: "", TargetState: "absent"},
		{Name: "absent/component2", Version: "", TargetState: "absent"},
		{Name: "present/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "present/component4", Version: compVersion3213.String(), TargetState: "present"},
	}

	require.NoError(t, err)
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
		{Name: "absent/component1", Version: compVersion3211.String(), TargetState: "absent"},
		{Name: "absent/component2", TargetState: "absent"},
		{Name: "present/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "present/component4", Version: compVersion1233.String()},
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

	require.NoError(t, err)

	convertedDogus := []domain.Dogu{
		{Name: common.QualifiedDoguName{Namespace: "absent", Name: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "absent", Name: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "present", Name: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "present", Name: "dogu4"}, Version: version1_2_3_3},
	}

	convertedComponents := []domain.Component{
		{Name: "component1", DistributionNamespace: "absent", Version: compVersion3211, TargetState: domain.TargetStateAbsent},
		{Name: "component2", DistributionNamespace: "absent", TargetState: domain.TargetStateAbsent},
		{Name: "component3", DistributionNamespace: "present", Version: compVersion3212, TargetState: domain.TargetStatePresent},
		{Name: "component4", DistributionNamespace: "present", Version: compVersion1233},
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
			{Name: "component1", Version: compVersion3211.String(), TargetState: "not known state"},
			{Name: "k8s/component2", Version: "abc"},
		},
	}

	_, err := convertToBlueprint(blueprintV2)

	require.ErrorContains(t, err, "syntax of blueprintV2 is not correct: ")

	require.ErrorContains(t, err, "cannot convert blueprint dogus: ")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'dogu1'")
	require.ErrorContains(t, err, "dogu name needs to be in the form 'namespace/dogu' but is 'name/space/dogu2'")
	require.ErrorContains(t, err, `unknown target state "unknown"`)

	require.ErrorContains(t, err, "cannot convert blueprint components: ")
	require.ErrorContains(t, err, `unknown target state "not known state"`)
	require.ErrorContains(t, err, `could not parse version of target dogu "official/dogu3": failed to parse major version abc`)
	require.ErrorContains(t, err, `could not parse version of target component "k8s/component2"`)
}
