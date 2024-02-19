package domain

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/require"
	"testing"
)

var version3_2_1_0, _ = core.ParseVersion("3.2.1-0")

var (
	officialNamespace = common.DoguNamespace("official")
	officialDogu1     = common.QualifiedDoguName{Namespace: officialNamespace, Name: common.SimpleDoguName("dogu1")}
	officialDogu2     = common.QualifiedDoguName{Namespace: officialNamespace, Name: common.SimpleDoguName("dogu2")}
	officialDogu3     = common.QualifiedDoguName{Namespace: officialNamespace, Name: common.SimpleDoguName("dogu3")}
)

func Test_Validate(t *testing.T) {
	dogus := []MaskDogu{
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: TargetStateAbsent},
		{Name: officialDogu2, TargetState: TargetStateAbsent},
		{Name: officialDogu3, Version: version3_2_1_0, TargetState: TargetStatePresent},
	}
	blueprintMask := BlueprintMask{
		Dogus: dogus,
	}

	err := blueprintMask.Validate()

	require.Nil(t, err)
}

func Test_ValidateWithDuplicatedDoguNames(t *testing.T) {
	dogus := []MaskDogu{
		{Name: officialDogu1, TargetState: TargetStatePresent},
		{Name: officialDogu1, TargetState: TargetStateAbsent},
	}
	blueprintMask := BlueprintMask{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "Multiple definitions for the same dogu should lead to an error")
}
