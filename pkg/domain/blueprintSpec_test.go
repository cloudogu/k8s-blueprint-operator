package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BlueprintSpec_Validate_allOk(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023"}
	require.Equal(t, StatusPhaseNew, spec.Status, "Status new should be the default")

	err := spec.Validate()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseValidated, spec.Status)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecValidatedEvent{}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_inStatusValidated(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023", Status: StatusPhaseValidated}

	err := spec.Validate()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseValidated, spec.Status)
	require.Equal(t, 0, len(spec.Events), "there should be no additional Events generated")
}

func Test_BlueprintSpec_Validate_inStatusInProgress(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023", Status: StatusPhaseInProgress}

	err := spec.Validate()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseInProgress, spec.Status, "should stay in the old status")
	require.Equal(t, 0, len(spec.Events), "there should be no additional Events generated")
}

func Test_BlueprintSpec_Validate_inStatusInvalid(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023", Status: StatusPhaseInvalid}

	err := spec.Validate()

	require.NotNil(t, err, "should not evaluate again and should stop with an error")
	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorContains(t, err, "blueprint spec was marked invalid before. Do not revalidate")
}

func Test_BlueprintSpec_Validate_emptyID(t *testing.T) {
	spec := BlueprintSpec{}

	err := spec.Validate()

	require.NotNil(t, err, "No ID definition should lead to an error")
	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecInvalidEvent{err}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_combineErrors(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{{Name: "no namespace"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{{Name: "no namespace"}}},
	}

	err := spec.Validate()

	var invalidError *InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorContains(t, err, "blueprint spec is invalid")
	assert.ErrorContains(t, err, "blueprint spec don't have an ID")
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "blueprint mask is invalid")
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint_maskForDoguWhichIsNotInBlueprint(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{{Namespace: "official", Name: "nexus"}}},
	}

	err := spec.validateMaskAgainstBlueprint()

	assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
	assert.ErrorContains(t, err, "dogu nexus is missing in the blueprint")
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint_namespaceSwitchAllowed(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{{Namespace: "official", Name: "nexus"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{{Namespace: "premium", Name: "nexus"}}},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
	}

	err := spec.validateMaskAgainstBlueprint()

	require.Nil(t, err)
}

func Test_BlueprintSpec_validateMaskAgainstBlueprint_namespaceSwitchNotAllowed(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{{Namespace: "official", Name: "nexus"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{{Namespace: "premium", Name: "nexus"}}},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: false},
	}

	err := spec.validateMaskAgainstBlueprint()

	assert.ErrorContains(t, err, "blueprint mask does not match the blueprint")
	assert.ErrorContains(t, err, "namespace switch is not allowed by default for dogu nexus. Activate feature flag for that")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_noMask(t *testing.T) {
	dogus := []TargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Namespace: "absent", Name: "dogu3", Version: "3.2.1-3", TargetState: TargetStateAbsent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{}},
		Status:        StatusPhaseValidated,
	}

	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_statusNew(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{}},
		Status:        StatusPhaseNew,
	}

	err := spec.CalculateEffectiveBlueprint()

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot calculate effective blueprint before the blueprint spec is validated")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_statusInvalid(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{}},
		Status:        StatusPhaseInvalid,
	}

	err := spec.CalculateEffectiveBlueprint()

	require.NotNil(t, err)
	assert.ErrorContains(t, err, "cannot calculate effective blueprint on invalid blueprint spec")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_changeVersion(t *testing.T) {
	dogus := []TargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: "3.2.1-2", TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskTargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: "3.2.1-1", TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)
	require.Equal(t, 2, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	assert.Equal(t, TargetDogu{Namespace: "official", Name: "dogu1", Version: "3.2.1-2", TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
	assert.Equal(t, TargetDogu{Namespace: "official", Name: "dogu2", Version: "3.2.1-1", TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[1])
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_makeDoguAbsent(t *testing.T) {
	dogus := []TargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
		{Namespace: "official", Name: "dogu2", Version: "3.2.1-2", TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskTargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent},
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
	assert.Equal(t, TargetDogu{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStateAbsent}, spec.EffectiveBlueprint.Dogus[0])
	assert.Equal(t, TargetDogu{Namespace: "official", Name: "dogu2", Version: "3.2.1-2", TargetState: TargetStateAbsent}, spec.EffectiveBlueprint.Dogus[1])
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_makeAbsentDoguPresent(t *testing.T) {
	dogus := []TargetDogu{
		{Namespace: "official", Name: "dogu1", TargetState: TargetStateAbsent},
	}

	maskedDogus := []MaskTargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
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
	assert.Equal(t, TargetDogu{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_changeDoguNamespace(t *testing.T) {
	dogus := []TargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskTargetDogu{
		{Namespace: "premium", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: false},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.NotNil(t, err, "without the feature flag, namespace changes are not allowed")
	require.ErrorContains(t, err, "changing the dogu namespace is only allowed with the changeDoguNamespace flag")
}

func Test_BlueprintSpec_CalculateEffectiveBlueprint_changeDoguNamespaceWithFlag(t *testing.T) {
	dogus := []TargetDogu{
		{Namespace: "official", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
	}

	maskedDogus := []MaskTargetDogu{
		{Namespace: "premium", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent},
	}

	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: dogus},
		BlueprintMask: BlueprintMask{Dogus: maskedDogus},
		Config:        BlueprintConfiguration{AllowDoguNamespaceSwitch: true},
		Status:        StatusPhaseValidated,
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err, "with the feature flag namespace changes should be allowed")
	require.Equal(t, 1, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	assert.Equal(t, TargetDogu{Namespace: "premium", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
}
