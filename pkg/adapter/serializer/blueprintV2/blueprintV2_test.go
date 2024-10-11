package blueprintV2

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/cesapp-lib/core"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	version3211, _ = core.ParseVersion("3.2.1-1")
	version3212, _ = core.ParseVersion("3.2.1-2")
	version1233, _ = core.ParseVersion("1.2.3-3")
)

var (
	compVersion1233 = semver.MustParse("1.2.3-3")
	compVersion3211 = semver.MustParse("3.2.1-1")
	compVersion3212 = semver.MustParse("3.2.1-2")
	compVersion3213 = semver.MustParse("3.2.1-3")
)

func Test_ConvertToBlueprintV2(t *testing.T) {
	dogus := []domain.Dogu{
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu4"}, Version: version1233},
	}

	components := []domain.Component{
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component1"}, Version: compVersion3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component3"}, Version: compVersion3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component4"}, Version: compVersion3213},
	}
	blueprint := domain.Blueprint{
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
							{
								DoguName: "my-dogu",
								Key:      "config-encrypted",
							}: "42",
						},
					},
				},
			},
			Global: domain.GlobalConfig{Absent: []common.GlobalConfigKey{"test/key"}},
		},
	}

	blueprintV2, err := ConvertToBlueprintDTO(blueprint)

	convertedDogus := []serializer.TargetDogu{
		{Name: "official/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "official/dogu2", TargetState: "absent"},
		{Name: "official/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "official/dogu4", Version: version1233.Raw, TargetState: "present"},
	}

	convertedComponents := []serializer.TargetComponent{
		{Name: "k8s/component1", Version: "", TargetState: "absent"},
		{Name: "k8s/component2", Version: "", TargetState: "absent"},
		{Name: "k8s/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "k8s/component4", Version: compVersion3213.String(), TargetState: "present"},
	}

	require.NoError(t, err)
	assert.Equal(t, BlueprintV2{
		GeneralBlueprint: serializer.GeneralBlueprint{API: serializer.V2},
		Dogus:            convertedDogus,
		Components:       convertedComponents,
		Config: v1.Config{
			Dogus: map[string]v1.CombinedDoguConfig{
				"my-dogu": {
					Config: v1.DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: v1.SensitiveDoguConfig{
						Present: map[string]string{
							"config-encrypted": "42",
						},
					},
				},
			},
			Global: v1.GlobalConfig{
				Absent: []string{"test/key"},
			},
		},
	}, blueprintV2)
}

func Test_ConvertToBlueprint(t *testing.T) {
	dogus := []serializer.TargetDogu{
		{Name: "official/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "official/dogu2", TargetState: "absent"},
		{Name: "official/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "official/dogu4", Version: version1233.Raw},
	}

	components := []serializer.TargetComponent{
		{Name: "k8s/component1", Version: compVersion3211.String(), TargetState: "absent"},
		{Name: "k8s/component2", TargetState: "absent"},
		{Name: "k8s/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "k8s/component4", Version: compVersion1233.String()},
	}

	blueprintV2 := BlueprintV2{
		GeneralBlueprint: serializer.GeneralBlueprint{API: serializer.V2},
		Dogus:            dogus,
		Components:       components,
		Config: v1.Config{
			Dogus: map[string]v1.CombinedDoguConfig{
				"my-dogu": {
					Config: v1.DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: v1.SensitiveDoguConfig{
						Present: map[string]string{
							"config-encrypted": "42",
						},
					},
				},
			},
			Global: v1.GlobalConfig{Absent: []string{"test/key"}},
		},
	}
	blueprint, err := convertToBlueprintDomain(blueprintV2)

	require.NoError(t, err)

	convertedDogus := []domain.Dogu{
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "dogu4"}, Version: version1233},
	}

	convertedComponents := []domain.Component{
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component1"}, Version: compVersion3211, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component2"}, TargetState: domain.TargetStateAbsent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component3"}, Version: compVersion3212, TargetState: domain.TargetStatePresent},
		{Name: common.QualifiedComponentName{Namespace: "k8s", SimpleName: "component4"}, Version: compVersion1233},
	}

	assert.Equal(t, domain.Blueprint{
		Dogus:      convertedDogus,
		Components: convertedComponents,
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
							{
								DoguName: "my-dogu",
								Key:      "config-encrypted",
							}: "42",
						},
					},
				},
			},
			Global: domain.GlobalConfig{Absent: []common.GlobalConfigKey{"test/key"}},
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

	_, err := convertToBlueprintDomain(blueprintV2)

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
