package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_blueprintMaskV1Validator_Validate(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Namespace: "absent", Name: "dogu1", Version: "3.2.1-0", TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: "3.2.1-0", TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMask{
		Dogus: dogus,
	}

	err := blueprintMask.Validate()

	require.Nil(t, err)
}

func Test_blueprintMaskV1Validator_ValidateWithMissingDoguName(t *testing.T) {
	dogus := []MaskTargetDogu{
		{TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMask{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "A missing dogu name definition in a dogu should lead to an error")
}

func Test_blueprintMaskV1Validator_ValidateWithDuplicatedDoguNames(t *testing.T) {
	dogus := []MaskTargetDogu{
		{Name: "present/dogu4", TargetState: TargetStatePresent},
		{Name: "present/dogu4", TargetState: TargetStateAbsent},
	}
	blueprintMask := BlueprintMask{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "Multiple definitions for the same dogu should lead to an error")
}
