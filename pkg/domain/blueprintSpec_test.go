package domain

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/api/meta"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

var version3211, _ = core.ParseVersion("3.2.1-1")
var version3212, _ = core.ParseVersion("3.2.1-2")
var version3213, _ = core.ParseVersion("3.2.1-3")

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
	require.Equal(t, StatusPhaseNew, spec.Status, "Status new should be the default")

	err := spec.ValidateStatically()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseNew, spec.Status)
	require.Equal(t, 0, len(spec.Events))
}

func Test_BlueprintSpec_Validate_emptyID(t *testing.T) {
	spec := BlueprintSpec{
		Conditions: &[]Condition{},
	}

	err := spec.ValidateStatically()

	assert.True(t, meta.IsStatusConditionFalse(*spec.Conditions, ConditionValid))
	require.NotNil(t, err, "No ID definition should lead to an error")
	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecInvalidEvent{err}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_combineErrors(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{{Version: core.Version{}, TargetState: TargetStatePresent}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{TargetState: 666}}},
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
			Blueprint:     Blueprint{Dogus: []Dogu{{Name: officialNexus, TargetState: TargetStateAbsent}}},
			BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: officialNexus, TargetState: TargetStatePresent}}},
		}

		err := spec.validateMaskAgainstBlueprint()

		assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
		assert.ErrorContains(t, err, "absent dogu \"nexus\" cannot be present in blueprint mask")
	})
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint(t *testing.T) {
	t.Run("no mask", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialDogu1, Version: version3211, TargetState: TargetStatePresent},
			{Name: officialDogu2, Version: version3212, TargetState: TargetStatePresent},
			{Name: officialDogu3, Version: version3213, TargetState: TargetStateAbsent},
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
			{Name: officialDogu1, Version: version3211, TargetState: TargetStatePresent},
			{Name: officialDogu2, Version: version3212, TargetState: TargetStatePresent},
		}

		maskedDogus := []MaskDogu{
			{Name: officialDogu1, Version: version3212, TargetState: TargetStatePresent},
			{Name: officialDogu2, Version: version3211, TargetState: TargetStatePresent},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.Nil(t, err)
		require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
		assert.Equal(t, Dogu{Name: officialDogu1, Version: version3212, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
		assert.Equal(t, Dogu{Name: officialDogu2, Version: version3211, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[1])
	})

	t.Run("make dogu absent", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialDogu1, Version: version3211, TargetState: TargetStatePresent},
			{Name: officialDogu2, Version: version3212, TargetState: TargetStatePresent},
		}

		maskedDogus := []MaskDogu{
			{Name: officialDogu1, Version: version3211, TargetState: TargetStateAbsent},
			{Name: officialDogu2, TargetState: TargetStateAbsent},
		}

		config := Config{
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				officialDogu1.SimpleName: {},
			},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus, Config: config},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.Nil(t, err)
		require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
		assert.Equal(t, Dogu{Name: officialDogu1, Version: version3211, TargetState: TargetStateAbsent}, spec.EffectiveBlueprint.Dogus[0])
		assert.Equal(t, Dogu{Name: officialDogu2, Version: version3212, TargetState: TargetStateAbsent}, spec.EffectiveBlueprint.Dogus[1])
		assert.NotContains(t, spec.EffectiveBlueprint.Config.Dogus, officialDogu1.SimpleName)
	})

	t.Run("change dogu namespace", func(t *testing.T) {
		dogus := []Dogu{
			{Name: officialNexus, Version: version3211, TargetState: TargetStatePresent},
		}

		maskedDogus := []MaskDogu{
			{Name: premiumNexus, Version: version3211, TargetState: TargetStatePresent},
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
			{Name: officialNexus, Version: version3211, TargetState: TargetStatePresent},
		}

		maskedDogus := []MaskDogu{
			{Name: premiumNexus, Version: version3211, TargetState: TargetStatePresent},
		}

		spec := BlueprintSpec{
			Blueprint:     Blueprint{Dogus: dogus},
			BlueprintMask: BlueprintMask{Dogus: maskedDogus},
			Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
		}
		err := spec.CalculateEffectiveBlueprint()

		require.NoError(t, err, "with the feature flag namespace changes should be allowed")
		require.Equal(t, 1, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
		assert.Equal(t, Dogu{Name: premiumNexus, Version: version3211, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
	})

	t.Run("validate only config for dogus in blueprint", func(t *testing.T) {
		config := Config{
			Dogus: map[cescommons.SimpleName]CombinedDoguConfig{
				"my-dogu": {},
			},
		}

		spec := BlueprintSpec{
			Blueprint: Blueprint{Config: config},
		}

		err := spec.CalculateEffectiveBlueprint()

		assert.ErrorContains(t, err, "setting config for dogu \"my-dogu\" is not allowed as it will not be installed with the blueprint")
		assert.Equal(t, 0, len(spec.Events))
	})
	t.Run("add additionalMounts", func(t *testing.T) {
		dogus := []Dogu{
			{
				Name:        k8sNginxStatic,
				Version:     version3211,
				TargetState: TargetStatePresent,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{SourceType: ecosystem.DataSourceConfigMap, Name: "html-config", Volume: "customhtml", Subfolder: "test"},
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
			Conditions: &[]Condition{},
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

		assert.True(t, meta.IsStatusConditionTrue(*spec.Conditions, ConditionExecutable))
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
			Conditions: &[]Condition{},
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
		assert.True(t, meta.IsStatusConditionTrue(*spec.Conditions, ConditionExecutable))
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
			Conditions: &[]Condition{},
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
		assert.True(t, meta.IsStatusConditionFalse(*spec.Conditions, ConditionExecutable))
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
			Conditions: &[]Condition{},
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
		assert.True(t, meta.IsStatusConditionFalse(*spec.Conditions, ConditionExecutable))
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
			Conditions: &[]Condition{},
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
		assert.True(t, meta.IsStatusConditionFalse(*spec.Conditions, ConditionExecutable))
		require.Error(t, err)
		assert.ErrorContains(t, err, "action \"downgrade\" is not allowed")
	})
}

func TestBlueprintSpec_CheckEcosystemHealthUpfront(t *testing.T) {
	t.Run("all healthy", func(t *testing.T) {
		blueprint := &BlueprintSpec{
			Conditions: &[]Condition{},
		}
		health := ecosystem.HealthResult{}
		err := blueprint.CheckEcosystemHealthUpfront(health)
		require.NoError(t, err)
		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, ConditionEcosystemHealthy))
		condition := meta.FindStatusCondition(*blueprint.Conditions, ConditionEcosystemHealthy)
		require.NotNil(t, condition)
		assert.Equal(t, "dogu health ignored: false; component health ignored: false", condition.Message)

		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, EcosystemHealthyUpfrontEvent{
			doguHealthIgnored:      false,
			componentHealthIgnored: false,
		}, blueprint.Events[0])
	})

	t.Run("no healthy event if condition status stays the same", func(t *testing.T) {
		blueprint := &BlueprintSpec{
			Conditions: &[]Condition{},
		}
		health := ecosystem.HealthResult{}

		err := blueprint.CheckEcosystemHealthUpfront(health)
		require.NoError(t, err)

		//reset events
		blueprint.Events = []Event{}
		require.Equal(t, 0, len(blueprint.Events))
		// call again and check if another event was created
		err = blueprint.CheckEcosystemHealthUpfront(health)
		require.NoError(t, err)
		require.Equal(t, 0, len(blueprint.Events))
	})

	t.Run("some unhealthy", func(t *testing.T) {
		blueprint := &BlueprintSpec{
			Conditions: &[]Condition{},
		}
		health := ecosystem.HealthResult{
			DoguHealth: ecosystem.DoguHealthResult{
				DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
					ecosystem.AvailableHealthStatus:   {"postfix"},
					ecosystem.UnavailableHealthStatus: {"ldap"},
					ecosystem.PendingHealthStatus:     {"postgresql"},
				},
			},
			ComponentHealth: ecosystem.ComponentHealthResult{
				ComponentsByStatus: map[ecosystem.HealthStatus][]common.SimpleComponentName{
					ecosystem.AvailableHealthStatus:   {"dogu-operator"},
					ecosystem.UnavailableHealthStatus: {"component-operator"},
					ecosystem.PendingHealthStatus:     {"service-discovery"},
				},
			},
		}
		err := blueprint.CheckEcosystemHealthUpfront(health)
		require.Error(t, err)
		assert.ErrorContains(t, err, "ecosystem is unhealthy before applying the blueprint")
		assert.True(t, meta.IsStatusConditionFalse(*blueprint.Conditions, ConditionEcosystemHealthy))
		condition := meta.FindStatusCondition(*blueprint.Conditions, ConditionEcosystemHealthy)
		require.NotNil(t, condition)
		assert.Contains(t, condition.Message, "2 dogu(s) are unhealthy: ldap, postgresql")
		assert.Contains(t, condition.Message, "2 component(s) are unhealthy: component-operator, service-discovery")

		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, EcosystemUnhealthyUpfrontEvent{HealthResult: health}, blueprint.Events[0])
	})

	t.Run("no unhealthy event if condition status stays the same", func(t *testing.T) {
		blueprint := &BlueprintSpec{
			Conditions: &[]Condition{},
		}
		health := ecosystem.HealthResult{
			DoguHealth: ecosystem.DoguHealthResult{
				DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
					ecosystem.UnavailableHealthStatus: {"ldap"},
				},
			},
		}

		err := blueprint.CheckEcosystemHealthUpfront(health)
		require.Error(t, err)

		//reset events
		blueprint.Events = []Event{}
		require.Equal(t, 0, len(blueprint.Events))
		// call again and check if another event was created
		err = blueprint.CheckEcosystemHealthUpfront(health)
		require.Error(t, err)
		require.Equal(t, 0, len(blueprint.Events))
	})
}

func TestBlueprintSpec_CheckEcosystemHealthAfterwards(t *testing.T) {
	tests := []struct {
		name               string
		inputSpec          *BlueprintSpec
		healthResult       ecosystem.HealthResult
		expectedEventNames []string
		expectedEventMsgs  []string
	}{
		{
			name:      "should write unhealthy dogus in event",
			inputSpec: &BlueprintSpec{},
			healthResult: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{
					DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
						ecosystem.AvailableHealthStatus:   {"postfix"},
						ecosystem.UnavailableHealthStatus: {"ldap"},
						ecosystem.PendingHealthStatus:     {"postgresql"},
					},
				},
			},
			expectedEventNames: []string{"EcosystemUnhealthyAfterwards"},
			expectedEventMsgs:  []string{"ecosystem health:\n  2 dogu(s) are unhealthy: ldap, postgresql\n  0 component(s) are unhealthy: "},
		},
		{
			name:      "ecosystem healthy",
			inputSpec: &BlueprintSpec{},
			healthResult: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{
					DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
						ecosystem.AvailableHealthStatus: {"postfix", "ldap", "postgresql"},
					},
				},
			},
			expectedEventNames: []string{"EcosystemHealthyAfterwards"},
			expectedEventMsgs:  []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.inputSpec.CheckEcosystemHealthAfterwards(tt.healthResult)
			eventNames := util.Map(tt.inputSpec.Events, Event.Name)
			eventMsgs := util.Map(tt.inputSpec.Events, Event.Message)
			assert.ElementsMatch(t, tt.expectedEventNames, eventNames)
			assert.ElementsMatch(t, tt.expectedEventMsgs, eventMsgs)
		})
	}
}

func TestBlueprintSpec_StartApplying(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// given
		spec := &BlueprintSpec{}
		// when
		spec.StartApplying()
		// then
		assert.Equal(t, spec, &BlueprintSpec{
			Status: StatusPhaseInProgress,
			Events: []Event{InProgressEvent{}},
		})
	})
}

func TestBlueprintSpec_MarkBlueprintApplicationFailed(t *testing.T) {
	// given
	spec := &BlueprintSpec{}
	err := fmt.Errorf("test-error")
	// when
	spec.MarkBlueprintApplicationFailed(err)
	// then
	assert.Equal(t, spec, &BlueprintSpec{
		Status: StatusPhaseBlueprintApplicationFailed,
		Events: []Event{ExecutionFailedEvent{err: err}},
	})
}

func TestBlueprintSpec_MarkBlueprintApplied(t *testing.T) {
	// given
	spec := &BlueprintSpec{}
	// when
	spec.MarkBlueprintApplied()
	// then
	assert.Equal(t, spec, &BlueprintSpec{
		Status: StatusPhaseBlueprintApplied,
		Events: []Event{BlueprintAppliedEvent{}},
	})
}

func TestBlueprintSpec_CensorSensitiveData(t *testing.T) {
	// given
	spec := &BlueprintSpec{
		StateDiff: StateDiff{
			SensitiveDoguConfigDiffs: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
				"ldapDiff": []SensitiveDoguConfigEntryDiff{{
					Actual:   DoguConfigValueState{Value: "Test1"},
					Expected: DoguConfigValueState{Value: "Test2"},
				}},
			},
		},
	}
	// when
	spec.CensorSensitiveData()

	// then
	require.Len(t, spec.StateDiff.SensitiveDoguConfigDiffs, 1)
	assert.Contains(t, maps.Keys(spec.StateDiff.SensitiveDoguConfigDiffs), cescommons.SimpleName("ldapDiff"))
	require.Len(t, spec.StateDiff.SensitiveDoguConfigDiffs["ldapDiff"], 1)
	assert.Equal(t, censorValue, spec.StateDiff.SensitiveDoguConfigDiffs["ldapDiff"][0].Actual.Value)
	assert.Equal(t, censorValue, spec.StateDiff.SensitiveDoguConfigDiffs["ldapDiff"][0].Expected.Value)
}

func TestBlueprintSpec_CompletePostProcessing(t *testing.T) {
	t.Run("status change on success EcosystemHealthyAfterwards -> Completed", func(t *testing.T) {
		// given
		spec := &BlueprintSpec{
			Status: StatusPhaseEcosystemHealthyAfterwards,
		}
		// when
		spec.CompletePostProcessing()
		// then
		assert.Equal(t, spec, &BlueprintSpec{
			Status: StatusPhaseCompleted,
			Events: []Event{CompletedEvent{}},
		})
	})

	t.Run("status change on failure InProgress -> Failed", func(t *testing.T) {
		// given
		spec := &BlueprintSpec{
			Status: StatusPhaseInProgress,
		}
		// when
		spec.CompletePostProcessing()
		// then
		assert.Equal(t, spec, &BlueprintSpec{
			Status: StatusPhaseFailed,
			Events: []Event{ExecutionFailedEvent{errors.New(handleInProgressMsg)}},
		})
	})

	t.Run("status change on failure EcosystemUnhealthyAfterwards -> Failed", func(t *testing.T) {
		// given
		spec := &BlueprintSpec{
			Status: StatusPhaseEcosystemUnhealthyAfterwards,
		}
		// when
		spec.CompletePostProcessing()
		// then
		assert.Equal(t, spec, &BlueprintSpec{
			Status: StatusPhaseFailed,
			Events: []Event{ExecutionFailedEvent{errors.New("ecosystem is unhealthy")}},
		})
	})

	t.Run("status change on failure ApplicationFailed -> Failed", func(t *testing.T) {
		// given
		spec := &BlueprintSpec{
			Status: StatusPhaseBlueprintApplicationFailed,
		}
		// when
		spec.CompletePostProcessing()
		// then
		assert.Equal(t, spec, &BlueprintSpec{
			Status: StatusPhaseFailed,
			Events: []Event{ExecutionFailedEvent{errors.New("could not apply blueprint")}},
		})
	})
}

func TestBlueprintSpec_ValidateDynamically(t *testing.T) {
	t.Run("all ok, no errors", func(t *testing.T) {
		blueprint := BlueprintSpec{
			Conditions: &[]Condition{},
		}

		blueprint.ValidateDynamically(nil)

		assert.True(t, meta.IsStatusConditionTrue(*blueprint.Conditions, ConditionValid))
		require.Equal(t, 0, len(blueprint.Events))

	})

	t.Run("given dependency error", func(t *testing.T) {
		blueprint := BlueprintSpec{
			Conditions: &[]Condition{},
		}
		givenErr := assert.AnError

		blueprint.ValidateDynamically(givenErr)

		assert.True(t, meta.IsStatusConditionFalse(*blueprint.Conditions, ConditionValid))
		require.Equal(t, 1, len(blueprint.Events))
	})
}

func TestBlueprintSpec_ShouldBeApplied(t *testing.T) {
	t.Run("should be applied", func(t *testing.T) {
		spec := &BlueprintSpec{
			Config: BlueprintConfiguration{
				DryRun: false,
			},
		}
		assert.Truef(t, spec.ShouldBeApplied(), "ShouldBeApplied()")
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

		assert.Equal(t, StatusPhaseAwaitSelfUpgrade, blueprint.Status)
		assert.Equal(t, []Event{AwaitSelfUpgradeEvent{}}, blueprint.Events)
	})

	t.Run("repeated call -> no event", func(t *testing.T) {
		blueprint := BlueprintSpec{
			Status: StatusPhaseAwaitSelfUpgrade,
		}

		blueprint.MarkWaitingForSelfUpgrade()

		assert.Equal(t, StatusPhaseAwaitSelfUpgrade, blueprint.Status)
		assert.Equal(t, []Event(nil), blueprint.Events, "no additional event if status already was AwaitSelfUpgrade")
	})
}

func TestBlueprintSpec_MarkSelfUpgradeCompleted(t *testing.T) {
	t.Run("first call -> new event", func(t *testing.T) {
		blueprint := BlueprintSpec{
			Status: StatusPhaseAwaitSelfUpgrade,
		}
		blueprint.MarkSelfUpgradeCompleted()

		assert.Equal(t, StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		assert.Equal(t, []Event{SelfUpgradeCompletedEvent{}}, blueprint.Events)
	})

	t.Run("repeated call -> no event", func(t *testing.T) {
		blueprint := BlueprintSpec{
			Status: StatusPhaseSelfUpgradeCompleted,
		}

		blueprint.MarkSelfUpgradeCompleted()

		assert.Equal(t, StatusPhaseSelfUpgradeCompleted, blueprint.Status)
		assert.Equal(t, []Event(nil), blueprint.Events, "no additional event if status already was AwaitSelfUpgrade")
	})
}

func TestBlueprintSpec_GetDogusThatNeedARestart(t *testing.T) {
	testDogu1 := Dogu{Name: cescommons.QualifiedName{Namespace: "testNamespace", SimpleName: "testDogu1"}}
	testBlueprint1 := Blueprint{Dogus: []Dogu{testDogu1}}
	testDoguConfigDiffsChanged := []DoguConfigEntryDiff{{
		Actual:       DoguConfigValueState{},
		Expected:     DoguConfigValueState{Value: "testValue", Exists: true},
		NeededAction: ConfigActionSet,
	}}
	testDoguConfigDiffsActionNone := []DoguConfigEntryDiff{{
		NeededAction: ConfigActionNone,
	}}

	testDoguConfigChangeDiffChanged := StateDiff{
		DoguConfigDiffs: map[cescommons.SimpleName]DoguConfigDiffs{testDogu1.Name.SimpleName: testDoguConfigDiffsChanged},
	}
	testDoguConfigChangeDiffActionNone := StateDiff{
		DoguConfigDiffs: map[cescommons.SimpleName]DoguConfigDiffs{testDogu1.Name.SimpleName: testDoguConfigDiffsActionNone},
	}

	type fields struct {
		Blueprint          Blueprint
		EffectiveBlueprint EffectiveBlueprint
		StateDiff          StateDiff
	}
	tests := []struct {
		name   string
		fields fields
		want   []cescommons.SimpleName
	}{
		{
			name:   "return nothing on empty blueprint",
			fields: fields{},
			want:   nil,
		},
		{
			name:   "return nothing on no config change",
			fields: fields{Blueprint: testBlueprint1},
			want:   nil,
		},
		{
			name: "return dogu on dogu config change",
			fields: fields{
				Blueprint:          testBlueprint1,
				EffectiveBlueprint: EffectiveBlueprint(testBlueprint1),
				StateDiff:          testDoguConfigChangeDiffChanged,
			},
			want: []cescommons.SimpleName{testDogu1.Name.SimpleName},
		},
		{
			name: "return nothing on dogu config unchanged",
			fields: fields{
				Blueprint:          testBlueprint1,
				EffectiveBlueprint: EffectiveBlueprint(testBlueprint1),
				StateDiff:          testDoguConfigChangeDiffActionNone,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &BlueprintSpec{
				Blueprint:          tt.fields.Blueprint,
				EffectiveBlueprint: tt.fields.EffectiveBlueprint,
				StateDiff:          tt.fields.StateDiff,
			}
			assert.Equalf(t, tt.want, spec.GetDogusThatNeedARestart(), "GetDogusThatNeedARestart()")
		})
	}
}
