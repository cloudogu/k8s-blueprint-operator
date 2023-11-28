package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_blueprintMaskV1Validator_Validate(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-0", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-0", TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMaskV1{
		GeneralBlueprintMask: GeneralBlueprintMask{API: "v1"},
		ID:                   "testMask1",
		Dogus:                dogus,
	}

	err := blueprintMask.Validate()

	require.Nil(t, err)
}

func Test_blueprintMaskV1Validator_ValidateWithMissingDoguName(t *testing.T) {
	dogus := []MaskTargetDogu{
		{TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMaskV1{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "A missing dogu name definition in a dogu should lead to an error")
}

func Test_blueprintMaskV1Validator_ValidateWithDuplicatedDoguNames(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "present/dogu4", TargetState: TargetStatePresent},
		{Name: "present/dogu4", TargetState: TargetStateAbsent},
	}
	blueprintMask := BlueprintMaskV1{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "Multiple definitions for the same dogu should lead to an error")
}

func Test_blueprintMaskV1Validator_ValidateWithEmptyMaskAPI(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "present/dogu4", TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMaskV1{ID: "testMask1", Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "No mask API definition should lead to an error")
}

func Test_blueprintMaskV1Validator_ValidateWithEmptyMaskID(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "present/dogu4", TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMaskV1{GeneralBlueprintMask: GeneralBlueprintMask{API: "v1"}, Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "No mask ID definition should lead to an error")
}
