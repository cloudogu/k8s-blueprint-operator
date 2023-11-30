package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BlueprintSpec_Validate_allOk(t *testing.T) {
	spec := BlueprintSpec{Id: "29.11.2023"}

	err := spec.Validate()

	require.Nil(t, err)
	assert.Equal(t, StatusPhaseValidated, spec.Status)
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecValidatedEvent{}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_emptyID(t *testing.T) {
	spec := BlueprintSpec{}

	err := spec.Validate()

	require.NotNil(t, err, "No ID definition should lead to an error")
	require.Equal(t, 1, len(spec.Events))
	assert.Equal(t, BlueprintSpecInvalidEvent{err}, spec.Events[0])
}

func Test_BlueprintSpec_Validate_combineErrors(t *testing.T) {
	spec := BlueprintSpec{
		Blueprint:     Blueprint{Dogus: []TargetDogu{{Name: "no namespace"}}},
		BlueprintMask: BlueprintMask{Dogus: []MaskTargetDogu{{Name: "no namespace"}}},
	}

	err := spec.Validate()

	assert.ErrorContains(t, err, "blueprint spec is invalid")
	assert.ErrorContains(t, err, "blueprint spec don't have an ID")
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "blueprint mask is invalid")
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
	}

	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err)

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
		config:        BlueprintConfiguration{allowDoguNamespaceSwitch: false},
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
		config:        BlueprintConfiguration{allowDoguNamespaceSwitch: true},
	}
	err := spec.CalculateEffectiveBlueprint()

	require.Nil(t, err, "with the feature flag namespace changes should be allowed")
	require.Equal(t, 1, len(spec.EffectiveBlueprint.Dogus), "effective blueprint should contain the elements from the mask")
	assert.Equal(t, TargetDogu{Namespace: "premium", Name: "dogu1", Version: "3.2.1-1", TargetState: TargetStatePresent}, spec.EffectiveBlueprint.Dogus[0])
}
