package domain

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var version3_2_1_1, _ = core.ParseVersion("3.2.1-1")
var version3_2_1_2, _ = core.ParseVersion("3.2.1-2")
var version3_2_1_3, _ = core.ParseVersion("3.2.1-3")

var testContext = context.Background()

func Test_BlueprintSpec_Validate_allOk(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023"}
	require.Equal(t, StatusPhaseNew, spec.Status, "Status new should be the default")

	err := spec.ValidateStatically()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseStaticallyValidated, spec.Status)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecStaticallyValidatedEvent{}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_inStatusValidated(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023", Status: StatusPhaseStaticallyValidated}

	err := spec.ValidateStatically()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseStaticallyValidated, spec.Status)
	require.Equal(t, 0, len(spec.Events), "there should be no additional Events generated")
}

func Test_BlueprintSpec_Validate_inStatusInProgress(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023", Status: StatusPhaseInProgress}

	err := spec.ValidateStatically()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseInProgress, spec.Status, "should stay in the old status")
	require.Equal(t, 0, len(spec.Events), "there should be no additional Events generated")
}

func Test_BlueprintSpec_Validate_inStatusInvalid(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023", Status: StatusPhaseInvalid}

	err := spec.ValidateStatically()

	require.NotNil(t, err, "should not evaluate again and should stop with an error")
	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorContains(t, err, "blueprint spec was marked invalid before: do not revalidate")
}

func Test_BlueprintSpec_Validate_emptyID(t *testing.T) {
	spec := BlueprintSpec{}

	err := spec.ValidateStatically()

	require.NotNil(t, err, "No ID definition should lead to an error")
	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecInvalidEvent{err}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_combineErrors(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{{Name: "no namespace"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Name: "no namespace"}}},
	}

	err := spec.ValidateStatically()

	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorContains(t, err, "blueprint spec is invalid")
	assert.ErrorContains(t, err, "blueprint spec doesn't have an ID")
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "blueprint mask is invalid")
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint_maskForDoguWhichIsNotInBlueprint(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Namespace: "official", Name: "nexus"}}},
	}

	err := spec.validateMaskAgainstBlueprint()

	assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
	assert.ErrorContains(t, err, "dogu \"nexus\" is missing in the blueprint")
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint_namespaceSwitchAllowed(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{{Namespace: "official", Name: "nexus"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Namespace: "premium", Name: "nexus"}}},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
	}

	err := spec.validateMaskAgainstBlueprint()

	require.Nil(t, err)
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint_namespaceSwitchNotAllowed(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{{Namespace: "official", Name: "nexus"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{{Namespace: "premium", Name: "nexus"}}},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: false},
	}

	err := spec.validateMaskAgainstBlueprint()

	assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
	assert.ErrorContains(t, err, "namespace switch is not allowed by default for dogu \"nexus\": activate the feature flag for that")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_noMask(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: version3_2_1_2, TargetState: TargetStatePresent},
		{Namespace: "absent", Name: "dogu3", Version: version3_2_1_3, TargetState: TargetStateAbsent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{}},
		Status:        StatusPhaseValidated,
	}

	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_statusNew(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{}},
		Status:        StatusPhaseNew,
	}

	err := spec.CalculateEffectiveBlueprint()

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot calculate effective blueprint before the blueprint spec is validated")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_statusInvalid(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []Dogu{}},
		BlueprintMask: BlueprintMask{Dogus: []MaskDogu{}},
		Status:        StatusPhaseInvalid,
	}

	err := spec.CalculateEffectiveBlueprint()

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot calculate effective blueprint on invalid blueprint spec")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_changeVersion(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: version3_2_1_2, TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskDogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_2, TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: version3_2_1_1, TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)
	require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	assert.Equal(t, Dogu{Namespace: "official", Name: "dogu1", Version: version3_2_1_2, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
	assert.Equal(t, Dogu{Namespace: "official", Name: "dogu2", Version: version3_2_1_1, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[1])
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_makeDoguAbsent(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: version3_2_1_2, TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskDogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStateAbsent},
		{Namespace: "official", Name: "dogu2", TargetState: TargetStateAbsent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)
	require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	assert.Equal(t, Dogu{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStateAbsent}, spec.EffectiveBlueprint.Dogus[0])
	assert.Equal(t, Dogu{Namespace: "official", Name: "dogu2", Version: version3_2_1_2, TargetState: TargetStateAbsent}, spec.EffectiveBlueprint.Dogus[1])
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_makeAbsentDoguPresent(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "official", Name: "dogu1", TargetState: TargetStateAbsent},
	}

	maskedDogus := []MaskDogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)
	require.Equal(t, 1, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	//TODO: Is that the correct behavior? (absent dogus can be made present?)
	assert.Equal(t, Dogu{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_changeDoguNamespace(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskDogu{
		{Namespace: "premium", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: false},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Error(t, err, "without the feature flag, namespace changes are not allowed")
	require.ErrorContains(t, err, "changing the dogu namespace is forbidden by default and can be allowed by a flag: \"official/dogu1\" -> \"premium/dogu1\"")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_changeDoguNamespaceWithFlag(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "official", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskDogu{
		{Namespace: "premium", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.NoError(t, err, "with the feature flag namespace changes should be allowed")
	require.Equal(t, 1, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	assert.Equal(t, Dogu{Namespace: "premium", Name: "dogu1", Version: version3_2_1_1, TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
}

func TestBlueprintSpec_MarkInvalid(t *testing.T) {
	spec := BlueprintSpec{
		Config: BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
		Status: StatusPhaseValidated,
	}
	expectedErr := &InvalidBlueprintError{
		WrappedError: nil,
		Message:      "test-error",
	}
	spec.MarkInvalid(expectedErr)

	assert.Equal(t, StatusPhaseInvalid, spec.Status)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecInvalidEvent{ValidationError: expectedErr}, spec.Events[0])
}

func TestBlueprintSpec_DetermineStateDiff(t *testing.T) {
	// not every single case is tested here as this is a rather coarse-grained function
	// have a look at the tests for the more specialized functions used in the command, to see all possible combinations of diffs.
	t.Run("all ok with empty blueprint", func(t *testing.T) {
		//given
		spec := BlueprintSpec{
			EffectiveBlueprint: EffectiveBlueprint{
				Dogus:                   []Dogu{},
				Components:              []Component{},
				RegistryConfig:          nil,
				RegistryConfigAbsent:    nil,
				RegistryConfigEncrypted: nil,
			},
			Status: StatusPhaseValidated,
		}

		installedDogus := map[string]*ecosystem.DoguInstallation{}
		installedComponents := map[string]*ecosystem.ComponentInstallation{}

		//when
		err := spec.DetermineStateDiff(testContext, installedDogus, installedComponents)

		//then
		stateDiff := StateDiff{DoguDiffs: DoguDiffs{}, ComponentDiffs: ComponentDiffs{}}
		require.NoError(t, err)
		assert.Equal(t, StatusPhaseStateDiffDetermined, spec.Status)
		require.Equal(t, 1, len(spec.Events))
		assert.Equal(t, StateDiffDeterminedEvent{stateDiff}, spec.Events[0])
		assert.Equal(t, stateDiff, spec.StateDiff)
	})

	notAllowedStatus := []StatusPhase{StatusPhaseNew, StatusPhaseStaticallyValidated, StatusPhaseEffectiveBlueprintGenerated}
	for _, initialStatus := range notAllowedStatus {
		t.Run(fmt.Sprintf("cannot determine state diff in status %q", initialStatus), func(t *testing.T) {
			//given
			spec := BlueprintSpec{
				Status: initialStatus,
			}
			installedDogus := map[string]*ecosystem.DoguInstallation{}
			//when
			err := spec.DetermineStateDiff(testContext, installedDogus, nil)

			//then
			assert.Error(t, err)
			assert.Equal(t, spec.Status, initialStatus)
			require.Equal(t, 0, len(spec.Events))
			assert.ErrorContains(t, err, fmt.Sprintf("cannot determine state diff in status phase %q", initialStatus))
		})
	}
	t.Run("do not re-determine state diff", func(t *testing.T) {
		initialStatus := StatusPhaseCompleted
		//given
		spec := BlueprintSpec{
			Status: initialStatus,
		}
		installedDogus := map[string]*ecosystem.DoguInstallation{}
		//when
		err := spec.DetermineStateDiff(testContext, installedDogus, nil)

		//then
		assert.NoError(t, err)
		assert.Equal(t, spec.Status, initialStatus)
		require.Equal(t, 0, len(spec.Events))
	})
}

func TestBlueprintSpec_CheckEcosystemHealthUpfront(t *testing.T) {
	tests := []struct {
		name               string
		inputSpec          *BlueprintSpec
		healthResult       ecosystem.HealthResult
		expectedStatus     StatusPhase
		expectedEventNames []string
		expectedEventMsgs  []string
	}{
		{
			name:               "should return early if ignore dogu health is configured",
			inputSpec:          &BlueprintSpec{Config: BlueprintConfiguration{IgnoreDoguHealth: true}},
			healthResult:       ecosystem.HealthResult{},
			expectedStatus:     StatusPhaseEcosystemHealthyUpfront,
			expectedEventNames: []string{"EcosystemHealthyUpfront"},
			expectedEventMsgs:  []string{"dogu health ignored: true"},
		},
		{
			name:      "should write unhealthy dogus in event",
			inputSpec: &BlueprintSpec{},
			healthResult: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{
					DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
						ecosystem.AvailableHealthStatus:   {"postfix"},
						ecosystem.UnavailableHealthStatus: {"ldap"},
						ecosystem.PendingHealthStatus:     {"postgresql"},
					},
				},
			},
			expectedStatus:     StatusPhaseEcosystemUnhealthyUpfront,
			expectedEventNames: []string{"EcosystemUnhealthyUpfront"},
			expectedEventMsgs:  []string{"ecosystem is unhealthy: 2 dogus are unhealthy: ldap, postgresql"},
		},
		{
			name:      "all dogus healthy",
			inputSpec: &BlueprintSpec{},
			healthResult: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{
					DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
						ecosystem.AvailableHealthStatus: {"postfix", "ldap", "postgresql"},
					},
				},
			},
			expectedStatus:     StatusPhaseEcosystemHealthyUpfront,
			expectedEventNames: []string{"EcosystemHealthyUpfront"},
			expectedEventMsgs:  []string{"dogu health ignored: false"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.inputSpec.CheckEcosystemHealthUpfront(tt.healthResult)
			eventNames := util.Map(tt.inputSpec.Events, Event.Name)
			eventMsgs := util.Map(tt.inputSpec.Events, Event.Message)
			assert.ElementsMatch(t, tt.expectedEventNames, eventNames)
			assert.ElementsMatch(t, tt.expectedEventMsgs, eventMsgs)
		})
	}
}

func TestBlueprintSpec_CheckEcosystemHealthAfterwards(t *testing.T) {
	tests := []struct {
		name               string
		inputSpec          *BlueprintSpec
		healthResult       ecosystem.HealthResult
		expectedStatus     StatusPhase
		expectedEventNames []string
		expectedEventMsgs  []string
	}{
		{
			name:      "should write unhealthy dogus in event",
			inputSpec: &BlueprintSpec{},
			healthResult: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{
					DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
						ecosystem.AvailableHealthStatus:   {"postfix"},
						ecosystem.UnavailableHealthStatus: {"ldap"},
						ecosystem.PendingHealthStatus:     {"postgresql"},
					},
				},
			},
			expectedStatus:     StatusPhaseEcosystemUnhealthyUpfront,
			expectedEventNames: []string{"EcosystemUnhealthyAfterwards"},
			expectedEventMsgs:  []string{"ecosystem is unhealthy: 2 dogus are unhealthy: ldap, postgresql"},
		},
		{
			name:      "ecosystem healthy",
			inputSpec: &BlueprintSpec{},
			healthResult: ecosystem.HealthResult{
				DoguHealth: ecosystem.DoguHealthResult{
					DogusByStatus: map[ecosystem.HealthStatus][]ecosystem.DoguName{
						ecosystem.AvailableHealthStatus: {"postfix", "ldap", "postgresql"},
					},
				},
			},
			expectedStatus:     StatusPhaseEcosystemHealthyAfterwards,
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

func TestBlueprintSpec_MarkInProgress(t *testing.T) {
	//given
	spec := &BlueprintSpec{}
	//when
	spec.MarkInProgress()
	//then
	assert.Equal(t, spec, &BlueprintSpec{
		Status: StatusPhaseInProgress,
		Events: []Event{InProgressEvent{}},
	})
}

func TestBlueprintSpec_MarkFailed(t *testing.T) {
	//given
	spec := &BlueprintSpec{}
	err := fmt.Errorf("test-error")
	//when
	spec.MarkFailed(err)
	//then
	assert.Equal(t, spec, &BlueprintSpec{
		Status: StatusPhaseFailed,
		Events: []Event{ExecutionFailedEvent{err: err}},
	})
}

func TestBlueprintSpec_MarkBlueprintApplied(t *testing.T) {
	//given
	spec := &BlueprintSpec{}
	//when
	spec.MarkBlueprintApplied()
	//then
	assert.Equal(t, spec, &BlueprintSpec{
		Status: StatusPhaseBlueprintApplied,
		Events: []Event{BlueprintAppliedEvent{}},
	})
}

func TestBlueprintSpec_MarkCompleted(t *testing.T) {
	//given
	spec := &BlueprintSpec{}
	//when
	spec.MarkCompleted()
	//then
	assert.Equal(t, spec, &BlueprintSpec{
		Status: StatusPhaseCompleted,
		Events: []Event{CompletedEvent{}},
	})
}
