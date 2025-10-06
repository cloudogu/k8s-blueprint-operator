package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	libconfig "github.com/cloudogu/k8s-registry-lib/config"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

var version3211, _ = core.ParseVersion("3.2.1-1")
var version3212, _ = core.ParseVersion("3.2.1-2")
var version3213, _ = core.ParseVersion("3.2.1-3")
var subfolder = "subfolder"

var k8sNginxStatic = cescommons.QualifiedName{Namespace: "k8s", SimpleName: "nginx-static"}
var officialNexus = cescommons.QualifiedName{
	Namespace:  "official",
	SimpleName: "nexus",
}
var premiumNexus = cescommons.QualifiedName{
	Namespace:  "premium",
	SimpleName: "nexus",
}

var (
	debugValue = libconfig.Value("DEBUG")
	infoValue  = libconfig.Value("INFO")
	someValue1 = libconfig.Value("value1")
	someValue2 = libconfig.Value("value2")
)

func Test_BlueprintSpec_Validate_allOk(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023"}

	err := spec.ValidateStatically()

	require.Nil(t, err)
	require.Equal(t, 0, len(spec.Events))
}

func Test_BlueprintSpec_Validate_emptyID(t *testing.T) {
	spec := BlueprintSpec{}

	err := spec.ValidateStatically()

	assert.True(t, meta.IsStatusConditionFalse(spec.Conditions, ConditionValid))
	require.NotNil(t, err, "No ID definition should lead to an error")
	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecInvalidEvent{err}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_combineErrors(t *testing.T) {
	name, _ := cescommons.QualifiedNameFromString("/noNamespace")
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{{Version: &core.Version{}, Absent: false}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: name}}},
	}

	err := spec.ValidateStatically()

	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorContains(t, err, "blueprint spec is invalid")
	assert.ErrorContains(t, err, "blueprint spec doesn't have an ID")
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "blueprint mask is invalid")
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint(t *testing.T) {
	t.Run("mask for dogu which is not in blueprint", func(t *testing.T) {
		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: []Dogu{}},
			BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: officialNexus}}},
		}

		err := spec.validateMaskAgainstBlueprint()

		assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
		assert.ErrorContains(t, err, "dogu \"official/nexus\" is missing in the blueprint")
	})
	t.Run("namespace switch allowed", func(t *testing.T) {
		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: []Dogu{{Name: officialNexus}}},
			BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: premiumNexus}}},
			Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
		}

		err := spec.validateMaskAgainstBlueprint()

		require.Nil(t, err)
	})
	t.Run("namespace switch not allowed", func(t *testing.T) {
		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: []Dogu{{Name: officialNexus}}},
			BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: premiumNexus}}},
			Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: false},
		}

		err := spec.validateMaskAgainstBlueprint()

		assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
		assert.ErrorContains(t, err, "namespace switch is not allowed by default for dogu \"premium/nexus\": activate the feature flag for that")
	})
	t.Run("absent dogus cannot be present in blueprint mask", func(t *testing.T) {
		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: []Dogu{{Name: officialNexus, Absent: true}}},
			BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: officialNexus, Absent: false}}},
		}

		err := spec.validateMaskAgainstBlueprint()

		assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
		assert.ErrorContains(t, err, "absent dogu \"nexus\" cannot be present in blueprint mask")
	})
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint(t *testing.T) {
	t.Run("no mask", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialDogu1, Version: &version3211, Absent: false},
			{Name: officialDogu2, Version: &version3212, Absent: false},
			{Name: officialDogu3, Version: &version3213, Absent: true},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus},
			BlueprintMask: BlueprintMask{Dogus: []MaskDogu{}},
		}

		err := spec.CalculateEffectiveBlueprint()

		require.Nil(t, err)
	})

	t.Run("change version", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialDogu1, Version: &version3211, Absent: false},
			{Name: officialDogu2, Version: &version3212, Absent: false},
		}

		maskedDogus := []MaskDogu{
			{Name: officialDogu1, Version: version3212, Absent: false},
			{Name: officialDogu2, Version: version3211, Absent: false},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.Nil(t, err)
		require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
		assert.Equal(t, Dogu{Name: officialDogu1, Version: &version3212, Absent: false}, spec.EffectiveBlueprint.Dogus[0])
		assert.Equal(t, Dogu{Name: officialDogu2, Version: &version3211, Absent: false}, spec.EffectiveBlueprint.Dogus[1])
	})

	t.Run("make dogu absent", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialDogu1, Version: &version3211, Absent: false},
			{Name: officialDogu2, Version: &version3212, Absent: false},
		}

		maskedDogus := []MaskDogu{
			{Name: officialDogu1, Version: version3211, Absent: true},
			{Name: officialDogu2, Absent: true},
		}

		config := Config{
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				officialDogu1.SimpleName: {ConfigEntry{Key: "test", Value: &val1}},
			},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus, Config: config},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.Nil(t, err)
		require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
		assert.Equal(t, Dogu{Name: officialDogu1, Version: &version3211, Absent: true}, spec.EffectiveBlueprint.Dogus[0])
		assert.Equal(t, Dogu{Name: officialDogu2, Version: &version3212, Absent: true}, spec.EffectiveBlueprint.Dogus[1])
		assert.NotContains(t, spec.EffectiveBlueprint.Config.Dogus, officialDogu1.SimpleName)
	})

	t.Run("change dogu namespace", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialNexus, Version: &version3211, Absent: false},
		}

		maskedDogus := []MaskDogu{
			{Name: premiumNexus, Version: version3211, Absent: false},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
			Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: false},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.Error(t, err, "without the feature flag, namespace changes are not allowed")
		require.ErrorContains(t, err, "changing the dogu namespace is forbidden by default and can be allowed by a flag: \"official/nexus\" -> \"premium/nexus\"")
	})

	t.Run("change dogu namespace with flag", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialNexus, Version: &version3211, Absent: false},
		}

		maskedDogus := []MaskDogu{
			{Name: premiumNexus, Version: version3211, Absent: false},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
			Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.NoError(t, err, "with the feature flag namespace changes should be allowed")
		require.Equal(t, 1, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
		assert.Equal(t, Dogu{Name: premiumNexus, Version: &version3211, Absent: false}, spec.EffectiveBlueprint.Dogus[0])
	})

	t.Run("validate only config for dogus in blueprint", func(t *testing.T) {
		config := Config{
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"my-dogu": {},
			},
		}

		spec := BlueprintSpec{
			Blueprint: Blueprint{Config: config},
		}

		err := spec.CalculateEffectiveBlueprint()

		assert.ErrorContains(t, err, "setting config for dogu \"my-dogu\" is not allowed as it will not be installed with the blueprint")
		assert.Equal(t, 1, len(spec.Events))
		assert.Equal(t, "BlueprintSpecInvalid", spec.Events[0].Name())
	})

	t.Run("add additionalMounts", func(t *testing.T) {
		dogus := []Dogu{
			{
				Name:    k8sNginxStatic,
				Version: &version3211,
				Absent:  false,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{SourceType: ecosystem.DataSourceConfigMap, Name: "html-config", Volume: "customhtml", Subfolder: subfolder},
				},
			},
		}

		spec := BlueprintSpec{
			Blueprint: Blueprint{Dogus: dogus},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.Nil(t, err)
		assert.Equal(t, dogus[0], spec.EffectiveBlueprint.Dogus[0], "effective blueprint should contain dogu with all field from the original blueprint")
	})
}

func TestBlueprintSpec_MissingConfigReferences(t *testing.T) {
	blueprint := BlueprintSpec{}
	blueprint.MissingConfigReferences(assert.AnError)
	require.Equal(t, 1, len(blueprint.Events))
	assert.Equal(t, "MissingConfigReferences", blueprint.Events[0].Name())
	assert.Equal(t, assert.AnError.Error(), blueprint.Events[0].Message())
}

func TestBlueprintSpec_DetermineStateDiff(t *testing.T) {
	// not every single case is tested here as this is a rather coarse-grained function
	// have a look at the tests for the more specialized functions used in the command, to see all possible combinations of diffs.
	t.Run("all ok with empty blueprint", func(t *testing.T) {
		// given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Dogus: []Dogu{},
			},
		}

		clusterState := ecosystem.EcosystemState{
			InstalledDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{}, false)

		// then
		stateDiff := StateDiff{
			DoguDiffs: DoguDiffs{},
		}

		assert.True(t, meta.IsStatusConditionTrue(spec.Conditions, ConditionExecutable))
		require.NoError(t, err)
		require.Empty(t, spec.Events)
		assert.Equal(t, stateDiff, spec.StateDiff)
	})

	t.Run("all ok with filled blueprint", func(t *testing.T) {
		// given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Dogus:  []Dogu{{Name: officialNexus, Version: &version3211}},
				Config: Config{Global: GlobalConfigEntries{{Key: "test", Value: &val1}}},
			},
		}

		clusterState := ecosystem.EcosystemState{
			InstalledDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{}, false)

		// then
		stateDiff := StateDiff{
			DoguDiffs: DoguDiffs{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						Absent: true,
					},
					Expected: DoguDiffState{
						Namespace: "official",
						Version:   &version3211,
					},
					NeededActions: []Action{ActionInstall},
				},
			},
			GlobalConfigDiffs: GlobalConfigDiffs{
				{
					Key:    "test",
					Actual: GlobalConfigValueState{},
					Expected: GlobalConfigValueState{
						Value:  (*string)(&val1),
						Exists: true,
					},
					NeededAction: ConfigActionSet,
				},
			},
		}

		assert.True(t, meta.IsStatusConditionTrue(spec.Conditions, ConditionExecutable))
		require.NoError(t, err)
		require.Equal(t, 1, len(spec.Events))
		assert.Equal(t, newStateDiffEvent(stateDiff), spec.Events[0])
		assert.Empty(t, cmp.Diff(stateDiff, spec.StateDiff))
	})

	t.Run("ok with allowed dogu namespace switch", func(t *testing.T) {
		// given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Dogus: []Dogu{
					{
						Name: cescommons.QualifiedName{
							Namespace:  "namespace-change",
							SimpleName: "name",
						},
					},
				},
			},
			Config: BlueprintConfiguration{
				AllowDoguNamespaceSwitch: true,
			},
		}

		clusterState := ecosystem.EcosystemState{
			InstalledDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{
				"name": {Name: cescommons.QualifiedName{
					Namespace:  "namespace",
					SimpleName: "name",
				}},
			},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{}, false)

		// then
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(spec.Conditions, ConditionExecutable))
	})

	t.Run("invalid blueprint state with not allowed dogu namespace switch", func(t *testing.T) {
		// given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Dogus: []Dogu{
					{
						Name: cescommons.QualifiedName{
							Namespace:  "namespace-change",
							SimpleName: "name",
						},
					},
				},
			},
			Config: BlueprintConfiguration{
				AllowDoguNamespaceSwitch: false,
			},
		}

		clusterState := ecosystem.EcosystemState{
			InstalledDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{
				"name": {Name: cescommons.QualifiedName{
					Namespace:  "namespace",
					SimpleName: "name",
				}},
			},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.DoguConfigKey]common.SensitiveDoguConfigValue{}, false)

		// then
		assert.True(t, meta.IsStatusConditionFalse(spec.Conditions, ConditionExecutable))
		require.Error(t, err)
		assert.ErrorContains(t, err, "action \"dogu namespace switch\" is not allowed")
	})
}

func TestBlueprintSpec_CompletePostProcessing(t *testing.T) {
	t.Run("ok with event", func(t *testing.T) {
		// given
		blueprint := &BlueprintSpec{}
		// when
		changed := blueprint.Complete()
		// then
		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionCompleted)
		assert.Equal(t, metav1.ConditionTrue, condition.Status)
		assert.Equal(t, "Completed", condition.Reason)
		assert.Empty(t, condition.Message)
		assert.Equal(t, []Event{CompletedEvent{}}, blueprint.Events)
	})
	t.Run("no change if executed twice", func(t *testing.T) {
		// given
		blueprint := &BlueprintSpec{}
		// when
		changed := blueprint.Complete()
		assert.True(t, changed)
		blueprint.Events = nil
		changed = blueprint.Complete()
		// then
		assert.False(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionCompleted)
		assert.Equal(t, metav1.ConditionTrue, condition.Status)
		assert.Equal(t, "Completed", condition.Reason)
		assert.Empty(t, condition.Message)
		assert.Equal(t, 0, len(blueprint.Events))
	})
}

func TestBlueprintSpec_ValidateDynamically(t *testing.T) {
	t.Run("all ok, no errors", func(t *testing.T) {
		blueprint := BlueprintSpec{}

		blueprint.ValidateDynamically(nil)

		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, ConditionValid))
		require.Equal(t, 0, len(blueprint.Events))

	})

	t.Run("given dependency error", func(t *testing.T) {
		blueprint := BlueprintSpec{}
		givenErr := assert.AnError

		blueprint.ValidateDynamically(givenErr)

		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, "BlueprintSpecInvalid", blueprint.Events[0].Name())
	})
}

func TestBlueprintSpec_ShouldBeApplied(t *testing.T) {
	conditionCompleted := Condition{
		Type:   ConditionCompleted,
		Status: metav1.ConditionTrue,
	}
	t.Run("should be applied on global config change", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				Stopped: false,
			},
			StateDiff: StateDiff{
				GlobalConfigDiffs: []GlobalConfigEntryDiff{
					{
						Key:          "test",
						NeededAction: ConfigActionSet,
					},
				},
			},
			Conditions: []Condition{conditionCompleted},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should be applied on dogu config change", func(t *testing.T) {
		doguKey := common.DoguConfigKey{
			DoguName: "testDogu",
			Key:      "testKey",
		}
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				Stopped: false,
			},
			StateDiff: StateDiff{
				DoguConfigDiffs: map[cescommons.SimpleName]DoguConfigDiffs{
					cescommons.SimpleName("testDogu"): {
						{
							Key:          doguKey,
							NeededAction: ConfigActionSet,
						},
					},
				},
			},
			Conditions: []Condition{conditionCompleted},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should be applied on sensitive dogu config change", func(t *testing.T) {
		doguKey := common.DoguConfigKey{
			DoguName: "testDogu",
			Key:      "testKey",
		}
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				Stopped: false,
			},
			StateDiff: StateDiff{
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
					cescommons.SimpleName("testDogu"): {
						{
							Key:          doguKey,
							NeededAction: ConfigActionSet,
						},
					},
				},
			},
			Conditions: []Condition{conditionCompleted},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should be applied on condition completed false", func(t *testing.T) {
		spec := &BlueprintSpec{
			Conditions: []Condition{{
				Type:   ConditionCompleted,
				Status: metav1.ConditionFalse,
			}},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should be applied on condition completed unknown", func(t *testing.T) {
		spec := &BlueprintSpec{
			Conditions: []Condition{{
				Type:   ConditionCompleted,
				Status: metav1.ConditionUnknown,
			}},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should be applied on condition completed not set", func(t *testing.T) {
		spec := &BlueprintSpec{
			Conditions: []Condition{},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should not be applied without any changes", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				Stopped: false,
			},
			Conditions: []Condition{conditionCompleted},
		}
		assert.Falsef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should not be applied due to being stopped", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				Stopped: true,
			},
		}
		assert.Falsef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})

}

func TestBlueprintSpec_MarkWaitingForSelfUpgrade(t *testing.T) {
	t.Run("first call -> new event", func(t *testing.T) {
		blueprint := BlueprintSpec{}
		blueprint.MarkWaitingForSelfUpgrade()

		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, ConditionSelfUpgradeCompleted))
		assert.Equal(t, []Event{AwaitSelfUpgradeEvent{}}, blueprint.Events)
	})

	t.Run("repeated call -> no event", func(t *testing.T) {
		blueprint := BlueprintSpec{}

		blueprint.MarkWaitingForSelfUpgrade()
		blueprint.Events = []Event(nil)
		blueprint.MarkWaitingForSelfUpgrade()

		assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, ConditionSelfUpgradeCompleted))
		assert.Equal(t, []Event(nil), blueprint.Events, "no additional event if condition did not change")
	})
}

func TestBlueprintSpec_MarkSelfUpgradeCompleted(t *testing.T) {
	t.Run("first call -> new event", func(t *testing.T) {
		blueprint := BlueprintSpec{}
		blueprint.MarkSelfUpgradeCompleted()

		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, ConditionSelfUpgradeCompleted))
		assert.Equal(t, []Event{SelfUpgradeCompletedEvent{}}, blueprint.Events)
	})

	t.Run("repeated call -> no event", func(t *testing.T) {
		blueprint := BlueprintSpec{}

		blueprint.MarkSelfUpgradeCompleted()
		blueprint.Events = []Event(nil)
		blueprint.MarkSelfUpgradeCompleted()

		assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, ConditionSelfUpgradeCompleted))
		assert.Equal(t, []Event(nil), blueprint.Events, "no additional event if status already was AwaitSelfUpgrade")
	})
}

func TestBlueprintSpec_HandleHealthResult(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		blueprint := BlueprintSpec{}

		changed := blueprint.HandleHealthResult(ecosystem.HealthResult{}, nil)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionEcosystemHealthy)
		assert.Equal(t, metav1.ConditionTrue, condition.Status)
		assert.Equal(t, "Healthy", condition.Reason)
		assert.Equal(t, "dogu health ignored: false", condition.Message)
	})

	t.Run("unhealthy", func(t *testing.T) {
		blueprint := BlueprintSpec{}
		health := ecosystem.HealthResult{
			DoguHealth: ecosystem.DoguHealthResult{
				DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
					ecosystem.UnavailableHealthStatus: {"ldap"},
				},
			},
		}

		changed := blueprint.HandleHealthResult(health, nil)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionEcosystemHealthy)
		assert.Equal(t, metav1.ConditionFalse, condition.Status)
		assert.Equal(t, "Unhealthy", condition.Reason)
		assert.Contains(t, condition.Message, "ecosystem health:")
		assert.Contains(t, condition.Message, "1 dogu(s) are unhealthy: ldap")
	})

	t.Run("error given, condition Unknown", func(t *testing.T) {
		blueprint := BlueprintSpec{}

		changed := blueprint.HandleHealthResult(ecosystem.HealthResult{}, assert.AnError)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionEcosystemHealthy)
		assert.Equal(t, metav1.ConditionUnknown, condition.Status)
		assert.Equal(t, "CannotCheckHealth", condition.Reason)
		assert.Equal(t, assert.AnError.Error(), condition.Message)
	})

	t.Run("no condition change", func(t *testing.T) {
		blueprint := BlueprintSpec{}

		changed := blueprint.HandleHealthResult(ecosystem.HealthResult{}, assert.AnError)
		assert.True(t, changed, "condition should change after the first call")
		changed = blueprint.HandleHealthResult(ecosystem.HealthResult{}, assert.AnError)
		assert.False(t, changed, "condition should not change here")
	})
}

func Test_removeLogLevelChangesFromConfig(t *testing.T) {
	t.Run("should remove logging config entries and keep others", func(t *testing.T) {
		// given
		config := Config{
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"dogu1": {
					{Key: "logging/root", Value: &debugValue},
					{Key: "other/setting", Value: &someValue1},
					{Key: "another/config", Value: &someValue2},
				},
				"dogu2": {
					{Key: "logging/root", Value: &infoValue},
					{Key: "database/host", Value: &someValue1},
				},
			},
			Global: GlobalConfigEntries{
				{Key: "global/setting", Value: &someValue1},
			},
		}

		// when
		result := removeLogLevelChangesFromConfig(config)

		// then
		assert.Equal(t, config.Global, result.Global, "Global config should be preserved")

		// Check dogu1 config
		dogu1Config := result.Dogus["dogu1"]
		assert.Len(t, dogu1Config, 2, "dogu1 should have 2 config entries after removing logging")
		assert.Contains(t, dogu1Config, ConfigEntry{Key: "other/setting", Value: &someValue1})
		assert.Contains(t, dogu1Config, ConfigEntry{Key: "another/config", Value: &someValue2})
		assert.NotContains(t, dogu1Config, ConfigEntry{Key: "logging/root", Value: &debugValue})

		// Check dogu2 config
		dogu2Config := result.Dogus["dogu2"]
		assert.Len(t, dogu2Config, 1, "dogu2 should have 1 config entry after removing logging")
		assert.Contains(t, dogu2Config, ConfigEntry{Key: "database/host", Value: &someValue1})
	})

	t.Run("should return empty config entries for dogu with only logging config", func(t *testing.T) {
		// given
		config := Config{
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"dogu-with-only-logging": {
					{Key: "logging/root", Value: &debugValue},
				},
			},
		}

		// when
		result := removeLogLevelChangesFromConfig(config)

		// then
		assert.Len(t, result.Dogus["dogu-with-only-logging"], 0, "dogu with only logging should have empty config")
	})

	t.Run("should handle empty dogu config", func(t *testing.T) {
		// given
		config := Config{
			Dogus: map[cescommons.SimpleName]DoguConfigEntries{
				"empty-dogu": {},
			},
		}

		// when
		result := removeLogLevelChangesFromConfig(config)

		// then
		assert.Len(t, result.Dogus["empty-dogu"], 0, "empty dogu config should remain empty")
	})

	t.Run("should handle empty config", func(t *testing.T) {
		// given
		config := Config{
			Dogus:  map[cescommons.SimpleName]DoguConfigEntries{},
			Global: GlobalConfigEntries{},
		}

		// when
		result := removeLogLevelChangesFromConfig(config)

		// then
		assert.Empty(t, result.Dogus, "empty dogus map should remain empty")
		assert.Empty(t, result.Global, "empty global config should remain empty")
	})
}
