package blueprintV2

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/blueprint-lib/bpcore"
	bpv2 "github.com/cloudogu/blueprint-lib/v2/blueprintV2"
	"github.com/cloudogu/blueprint-lib/v2/entities"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
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
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu4"}, Version: version1233},
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

	convertedDogus := []entities.TargetDogu{
		{Name: "official/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "official/dogu2", TargetState: "absent"},
		{Name: "official/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "official/dogu4", Version: version1233.Raw, TargetState: "present"},
	}

	convertedComponents := []entities.TargetComponent{
		{Name: "k8s/component1", Version: "", TargetState: "absent"},
		{Name: "k8s/component2", Version: "", TargetState: "absent"},
		{Name: "k8s/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "k8s/component4", Version: compVersion3213.String(), TargetState: "present"},
	}

	require.NoError(t, err)
	assert.Equal(t, bpv2.BlueprintV2{
		GeneralBlueprint: bpcore.GeneralBlueprint{API: bpcore.V2},
		Dogus:            convertedDogus,
		Components:       convertedComponents,
		Config: entities.TargetConfig{
			Dogus: entities.DoguConfigMap{
				"my-dogu": {
					Config: entities.DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: entities.SensitiveDoguConfig{
						Present: map[string]string{
							"config-encrypted": "42",
						},
					},
				},
			},
			Global: entities.GlobalConfig{
				Absent: []string{"test/key"},
			},
		},
	}, blueprintV2)
}

func Test_ConvertToBlueprint(t *testing.T) {
	dogus := []entities.TargetDogu{
		{Name: "official/dogu1", Version: version3211.Raw, TargetState: "absent"},
		{Name: "official/dogu2", TargetState: "absent"},
		{Name: "official/dogu3", Version: version3212.Raw, TargetState: "present"},
		{Name: "official/dogu4", Version: version1233.Raw},
	}

	components := []entities.TargetComponent{
		{Name: "k8s/component1", Version: compVersion3211.String(), TargetState: "absent"},
		{Name: "k8s/component2", TargetState: "absent"},
		{Name: "k8s/component3", Version: compVersion3212.String(), TargetState: "present"},
		{Name: "k8s/component4", Version: compVersion1233.String()},
	}

	blueprintV2 := bpv2.BlueprintV2{
		GeneralBlueprint: bpcore.GeneralBlueprint{API: bpcore.V2},
		Dogus:            dogus,
		Components:       components,
		Config: entities.TargetConfig{
			Dogus: entities.DoguConfigMap{
				"my-dogu": {
					Config: entities.DoguConfig{
						Present: map[string]string{
							"config": "42",
						},
					},
					SensitiveConfig: entities.SensitiveDoguConfig{
						Present: map[string]string{
							"config-encrypted": "42",
						},
					},
				},
			},
			Global: entities.GlobalConfig{Absent: []string{"test/key"}},
		},
	}
	blueprint, err := convertToBlueprintDomain(blueprintV2)

	require.NoError(t, err)

	convertedDogus := []domain.Dogu{
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu1"}, Version: version3211, TargetState: domain.TargetStateAbsent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu2"}, TargetState: domain.TargetStateAbsent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu3"}, Version: version3212, TargetState: domain.TargetStatePresent},
		{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "dogu4"}, Version: version1233},
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
	blueprintV2 := bpv2.BlueprintV2{
		GeneralBlueprint: bpcore.GeneralBlueprint{API: bpcore.V2},
		Dogus: []entities.TargetDogu{
			{Name: "dogu1", Version: version3211.Raw},
			{Name: "official/dogu1", Version: version3211.Raw, TargetState: "unknown"},
			{Name: "name/space/dogu2", Version: version3212.Raw},
			{Name: "official/dogu3", Version: "abc"},
		},
		Components: []entities.TargetComponent{
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
