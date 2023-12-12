package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Validate(t *testing.T) {
	dogus := []MaskDogu{
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

func Test_ValidateWithMissingDoguName(t *testing.T) {
	dogus := []MaskDogu{
		{TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMask{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "A missing dogu name definition in a dogu should lead to an error")
}

func Test_ValidateWithDuplicatedDoguNames(t *testing.T) {
	dogus := []MaskDogu{
		{Name: "present/dogu4", TargetState: TargetStatePresent},
		{Name: "present/dogu4", TargetState: TargetStateAbsent},
	}
	blueprintMask := BlueprintMask{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "Multiple definitions for the same dogu should lead to an error")
}
