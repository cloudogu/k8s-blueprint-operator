package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
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

const (
	testDistributionNamespace       = "k8s"
	testChangeDistributionNamespace = "k8s-testing"
)

var k8sNginxStatic = cescommons.QualifiedName{Namespace: "k8s", SimpleName: "nginx-static"}
var officialNexus = cescommons.QualifiedName{
	Namespace:  "official",
	SimpleName: "nexus",
}
var premiumNexus = cescommons.QualifiedName{
	Namespace:  "premium",
	SimpleName: "nexus",
}

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
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				officialDogu1.SimpleName: {},
			},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus, Config: &config},
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
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				"my-dogu": {},
			},
		}

		spec := BlueprintSpec{
			Blueprint: Blueprint{Config: &config},
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
					{SourceType: ecosystem.DataSourceConfigMap, Name: "html-config", Volume: "customhtml", Subfolder: &subfolder},
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
				Dogus:      []Dogu{},
				Components: []Component{},
			},
		}

		clusterState := ecosystem.EcosystemState{
			InstalledDogus:      map[cescommons.SimpleName]*ecosystem.DoguInstallation{},
			InstalledComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{})

		// then
		stateDiff := StateDiff{
			DoguDiffs:                DoguDiffs{},
			ComponentDiffs:           ComponentDiffs{},
			DoguConfigDiffs:          map[cescommons.SimpleName]DoguConfigDiffs{},
			SensitiveDoguConfigDiffs: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{},
		}

		assert.True(t, meta.IsStatusConditionTrue(spec.Conditions, ConditionExecutable))
		require.NoError(t, err)
		require.Equal(t, 5, len(spec.Events))
		assert.Equal(t, newStateDiffDoguEvent(stateDiff.DoguDiffs), spec.Events[0])
		assert.Equal(t, newStateDiffComponentEvent(stateDiff.ComponentDiffs), spec.Events[1])
		assert.Equal(t, GlobalConfigDiffDeterminedEvent{GlobalConfigDiffs: GlobalConfigDiffs(nil)}, spec.Events[2])
		assert.Equal(t, DoguConfigDiffDeterminedEvent{
			DoguConfigDiffs: map[cescommons.SimpleName]DoguConfigDiffs{},
		}, spec.Events[3])
		assert.Equal(t, SensitiveDoguConfigDiffDeterminedEvent{
			SensitiveDoguConfigDiffs: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{},
		}, spec.Events[4])
		assert.Equal(t, stateDiff, spec.StateDiff)
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
			InstalledComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{})

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
			InstalledComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{})

		// then
		assert.True(t, meta.IsStatusConditionFalse(spec.Conditions, ConditionExecutable))
		require.Error(t, err)
		assert.ErrorContains(t, err, "action \"dogu namespace switch\" is not allowed")
	})

	t.Run("should return error with not allowed component namespace switch action", func(t *testing.T) {
		// given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Components: []Component{
					{
						Name: common.QualifiedComponentName{
							Namespace:  testChangeDistributionNamespace,
							SimpleName: testComponentName.SimpleName,
						},
						Version: compVersion3211,
					},
				},
			},
		}
		clusterState := ecosystem.EcosystemState{
			InstalledDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{},
			InstalledComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
				testComponentName.SimpleName: {
					Name:            testComponentName,
					ExpectedVersion: compVersion3211,
				},
			},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{})

		// then
		assert.True(t, meta.IsStatusConditionFalse(spec.Conditions, ConditionExecutable))
		require.Error(t, err)
		assert.ErrorContains(t, err, "action \"component namespace switch\" is not allowed")
	})

	t.Run("should return error with not allowed component downgrade action", func(t *testing.T) {
		// given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Components: []Component{
					{
						Name: common.QualifiedComponentName{
							Namespace:  testComponentName.Namespace,
							SimpleName: testComponentName.SimpleName,
						},
						Version: compVersion3210,
					},
				},
			},
		}
		clusterState := ecosystem.EcosystemState{
			InstalledDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{},
			InstalledComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
				testComponentName.SimpleName: {
					Name:            testComponentName,
					ExpectedVersion: compVersion3211,
				},
			},
		}

		// when
		err := spec.DetermineStateDiff(clusterState, map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue{})

		// then
		assert.True(t, meta.IsStatusConditionFalse(spec.Conditions, ConditionExecutable))
		require.Error(t, err)
		assert.ErrorContains(t, err, "action \"downgrade\" is not allowed")
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
		assert.Equal(t, "", condition.Message)
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
		assert.Equal(t, "", condition.Message)
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
	})
}

func TestBlueprintSpec_ShouldBeApplied(t *testing.T) {
	t.Run("should be applied", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				DryRun: false,
			},
			StateDiff: StateDiff{
				GlobalConfigDiffs: []GlobalConfigEntryDiff{
					{
						Key:          "test",
						NeededAction: ConfigActionSet,
					},
				},
			},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should not be applied without any changes", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				DryRun: false,
			},
		}
		assert.Falsef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
	})
	t.Run("should not be applied due to dry run", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				DryRun: true,
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
		assert.Equal(t, "dogu health ignored: false; component health ignored: false", condition.Message)
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
		assert.Contains(t, condition.Message, "0 component(s) are unhealthy:")
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

func TestBlueprintSpec_SetComponentAppliedCondition(t *testing.T) {
	diff := StateDiff{
		ComponentDiffs: ComponentDiffs{
			ComponentDiff{
				Name: "k8s-dogu-operator",
				NeededActions: []Action{
					ActionUpgrade, ActionSwitchComponentNamespace,
				},
			},
		},
	}

	t.Run("applied", func(t *testing.T) {
		blueprint := BlueprintSpec{
			StateDiff: diff,
		}

		changed := blueprint.SetComponentsAppliedCondition(nil)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionComponentsApplied)
		assert.Equal(t, metav1.ConditionTrue, condition.Status)
		assert.Equal(t, "Applied", condition.Reason)
		assert.Equal(t, "components applied: \"k8s-dogu-operator\": [upgrade, component namespace switch]", condition.Message)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, ComponentsAppliedEvent{Diffs: diff.ComponentDiffs}, blueprint.Events[0])
	})

	t.Run("error", func(t *testing.T) {
		blueprint := BlueprintSpec{
			StateDiff: diff,
		}

		changed := blueprint.SetComponentsAppliedCondition(assert.AnError)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionComponentsApplied)
		assert.Equal(t, metav1.ConditionFalse, condition.Status)
		assert.Equal(t, "CannotApply", condition.Reason)
		assert.Equal(t, assert.AnError.Error(), condition.Message)
	})

	t.Run("no condition change", func(t *testing.T) {
		blueprint := BlueprintSpec{
			StateDiff: diff,
		}

		changed := blueprint.SetComponentsAppliedCondition(assert.AnError)
		assert.True(t, changed)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, ExecutionFailedEvent{err: assert.AnError}, blueprint.Events[0])
		blueprint.Events = nil

		changed = blueprint.SetComponentsAppliedCondition(assert.AnError)
		assert.False(t, changed)
		require.Equal(t, 0, len(blueprint.Events))
	})
}

func TestBlueprintSpec_SetDogusAppliedCondition(t *testing.T) {
	diff := StateDiff{
		DoguDiffs: DoguDiffs{
			DoguDiff{
				DoguName: "cas",
				NeededActions: []Action{
					ActionUpgrade, ActionSwitchDoguNamespace,
				},
			},
		},
	}

	t.Run("applied", func(t *testing.T) {
		blueprint := BlueprintSpec{
			StateDiff: diff,
		}

		changed := blueprint.SetDogusAppliedCondition(nil)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionDogusApplied)
		assert.Equal(t, metav1.ConditionTrue, condition.Status)
		assert.Equal(t, "Applied", condition.Reason)
		assert.Equal(t, "dogus applied: \"cas\": [upgrade, dogu namespace switch]", condition.Message)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, DogusAppliedEvent{Diffs: diff.DoguDiffs}, blueprint.Events[0])
	})

	t.Run("error", func(t *testing.T) {
		blueprint := BlueprintSpec{
			StateDiff: diff,
		}

		changed := blueprint.SetDogusAppliedCondition(assert.AnError)

		assert.True(t, changed)
		condition := meta.FindStatusCondition(blueprint.Conditions, ConditionDogusApplied)
		assert.Equal(t, metav1.ConditionFalse, condition.Status)
		assert.Equal(t, "CannotApply", condition.Reason)
		assert.Equal(t, assert.AnError.Error(), condition.Message)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, ExecutionFailedEvent{err: assert.AnError}, blueprint.Events[0])
	})

	t.Run("no condition change", func(t *testing.T) {
		blueprint := BlueprintSpec{
			StateDiff: diff,
		}

		changed := blueprint.SetDogusAppliedCondition(assert.AnError)
		assert.True(t, changed)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, ExecutionFailedEvent{err: assert.AnError}, blueprint.Events[0])
		blueprint.Events = nil

		changed = blueprint.SetDogusAppliedCondition(assert.AnError)
		assert.False(t, changed)
		require.Equal(t, 0, len(blueprint.Events))
	})
}

// TODO: write tests
//func TestBlueprintSpec_setComponentsAppliedConditionAfterStateDiff(t *testing.T) {
//	//	decision table:
//	//	current condition   withDiff     withoutDiff  DiffError
//	//	=======================================================
//	//	none                NeedToApply	 Applied      Unknown
//	//	Unknown             NeedToApply	 Applied      Unknown
//	//	NeedToApply         NeedToApply	 Applied      Unknown
//	//	Applied             NeedToApply	 no change    Unknown
//	//	CannotApply         no change    no change    Unknown
//	t.Run("no condition before but changes", func(t *testing.T) {
//		diff := StateDiff{
//			DoguDiffs: DoguDiffs{
//				DoguDiff{
//					DoguName: "cas",
//					NeededActions: []Action{
//						ActionUpgrade, ActionSwitchDoguNamespace,
//					},
//				},
//			},
//		}
//		blueprint := BlueprintSpec{
//			Conditions: &[]Condition{},
//			StateDiff:  diff,
//		}
//
//		blueprint.setDogusAppliedConditionAfterStateDiff()
//
//		event := newStateDiffDoguEvent(diff.DoguDiffs)
//		require.Equal(t, 1, len(blueprint.Events))
//		assert.Equal(t, event, blueprint.Events[0])
//		condition := meta.FindStatusCondition(*blueprint.Conditions, ConditionDogusApplied)
//		assert.Equal(t, metav1.ConditionFalse, condition.Status)
//		assert.Equal(t, "NeedToApply", condition.Reason)
//		assert.Equal(t, event.Message(), condition.Message)
//	})
//	t.Run("no condition before and no changes", func(t *testing.T) {
//		diff := StateDiff{
//			DoguDiffs: DoguDiffs{
//				DoguDiff{
//					DoguName:      "cas",
//					NeededActions: nil,
//				},
//			},
//		}
//		blueprint := BlueprintSpec{
//			Conditions: &[]Condition{},
//			StateDiff:  diff,
//		}
//
//		blueprint.setDogusAppliedConditionAfterStateDiff()
//
//		event := newStateDiffDoguEvent(diff.DoguDiffs)
//		require.Equal(t, 1, len(blueprint.Events))
//		assert.Equal(t, event, blueprint.Events[0])
//		condition := meta.FindStatusCondition(*blueprint.Conditions, ConditionDogusApplied)
//		assert.Equal(t, metav1.ConditionFalse, condition.Status)
//		assert.Equal(t, "Applied", condition.Reason)
//		assert.Equal(t, event.Message(), condition.Message)
//	})
//}

func TestBlueprintSpec_setDogusAppliedConditionAfterStateDiff(t *testing.T) {
	//	current condition   withDiff     withoutDiff  DiffError
	//	=======================================================
	//	none                NeedToApply	 Applied      Unknown
	//	Unknown             NeedToApply	 Applied      Unknown
	//	NeedToApply         NeedToApply	 Applied      Unknown
	//	Applied             NeedToApply	 no change    Unknown
	//	CannotApply         no change    no change    Unknown

	diffEventMsg := "dogu state diff determined: 2 actions (\"dogu namespace switch\": 1, \"upgrade\": 1)"
	noDiffEventMsg := "dogu state diff determined: 0 actions ()"
	type given struct {
		hasDiff   bool
		condition Condition
		diffErr   error
	}
	type expected struct {
		changed   bool
		condition *Condition
		hasEvent  bool
	}
	tests := []struct {
		name     string
		given    given
		expected expected
	}{
		{
			name: "none, diff -> NeedToApply, event",
			given: given{
				hasDiff: true,
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: diffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "none, no diff -> Applied, event",
			given: given{
				hasDiff: false,
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "Applied",
					Message: noDiffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "Unknown, diff -> NeedToApply, event",
			given: given{
				hasDiff: true,
				condition: Condition{
					Type:   ConditionDogusApplied,
					Status: metav1.ConditionUnknown,
					Reason: "CannotDetermineStateDiff",
				},
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: diffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "Unknown, no diff -> Applied, event",
			given: given{
				hasDiff: false,
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "Applied",
					Message: noDiffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "NeedToApply, diff -> NeedToApply, event",
			given: given{
				hasDiff: true,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionTrue,
					Reason:  "NeedToApply",
					Message: "other msg",
				},
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: diffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "NeedToApply, diff -> NeedToApply, no event",
			given: given{
				hasDiff: true,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: diffEventMsg,
				},
			},
			expected: expected{
				changed: false,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: diffEventMsg,
				},
				hasEvent: false,
			},
		},
		{
			name: "NeedToApply, no diff -> Applied, event",
			given: given{
				hasDiff: false,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: "error msg",
				},
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "Applied",
					Message: noDiffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "Applied, diff -> NeedToApply, event",
			given: given{
				hasDiff: true,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionTrue,
					Reason:  "Applied",
					Message: "other msg",
				},
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "NeedToApply",
					Message: diffEventMsg,
				},
				hasEvent: true,
			},
		},
		{
			name: "Applied, no diff -> no change",
			given: given{
				hasDiff: false,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  "Applied",
					Message: "any msg",
				},
			},
			expected: expected{
				changed:  false,
				hasEvent: false,
			},
		},
		{
			name: "CannotApply, no diff -> no change",
			given: given{
				hasDiff: false,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  ReasonCannotApply,
					Message: "error msg",
				},
			},
			expected: expected{
				changed:  false,
				hasEvent: false,
			},
		},
		{
			name: "CannotApply, diff -> no change",
			given: given{
				hasDiff: true,
				condition: Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionFalse,
					Reason:  ReasonCannotApply,
					Message: "error msg",
				},
			},
			expected: expected{
				changed:  false,
				hasEvent: false,
			},
		},
		{
			name: "error given",
			given: given{
				hasDiff: false,
				diffErr: assert.AnError,
			},
			expected: expected{
				changed: true,
				condition: &Condition{
					Type:    ConditionDogusApplied,
					Status:  metav1.ConditionUnknown,
					Reason:  "CannotDetermineStateDiff",
					Message: assert.AnError.Error(),
				},
				hasEvent: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			spec := &BlueprintSpec{
				Conditions: []Condition{tt.given.condition},
			}
			if tt.given.hasDiff {
				spec.StateDiff = StateDiff{
					DoguDiffs: DoguDiffs{
						DoguDiff{
							DoguName: "redmine",
							NeededActions: []Action{
								ActionUpgrade, ActionSwitchDoguNamespace,
							},
						},
					},
				}
			}
			//when
			hasChanged := spec.setDogusAppliedConditionAfterStateDiff(tt.given.diffErr)
			//then
			assert.Equal(t, tt.expected.changed, hasChanged, "function should return %v", hasChanged)
			//condition
			condition := meta.FindStatusCondition(spec.Conditions, ConditionDogusApplied)
			if tt.expected.changed && tt.expected.condition != nil {
				require.NotNil(t, condition)
				assert.Equal(t, tt.expected.condition.Status, condition.Status)
				assert.Equal(t, tt.expected.condition.Message, condition.Message)
				assert.Equal(t, tt.expected.condition.Reason, condition.Reason)
			} else {
				conditions := spec.Conditions
				condition := conditions[0]
				assert.Equal(t, tt.given.condition, condition)
			}
			//events
			if tt.expected.hasEvent {
				require.Equal(t, 1, len(spec.Events), "should have an event")
				assert.Equal(t, newStateDiffDoguEvent(spec.StateDiff.DoguDiffs), spec.Events[0])
			}
		})
	}
}
