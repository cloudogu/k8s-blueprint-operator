package domain

import (
	"testing"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/require"
)

var version3_2_1_0, _ = core.ParseVersion("3.2.1-0")

var (
	officialNamespace = cescommons.Namespace("official")
	officialDogu1     = cescommons.QualifiedName{Namespace: officialNamespace, SimpleName: cescommons.SimpleName("dogu1")}
	officialDogu2     = cescommons.QualifiedName{Namespace: officialNamespace, SimpleName: cescommons.SimpleName("dogu2")}
	officialDogu3     = cescommons.QualifiedName{Namespace: officialNamespace, SimpleName: cescommons.SimpleName("dogu3")}
)

func Test_Validate(t *testing.T) {
	dogus := []bpv2.MaskDogu{
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: bpv2.TargetStateAbsent},
		{Name: officialDogu2, TargetState: bpv2.TargetStateAbsent},
		{Name: officialDogu3, Version: version3_2_1_0, TargetState: bpv2.TargetStatePresent},
	}
	blueprintMask := bpv2.BlueprintMask{
		Dogus: dogus,
	}

	err := newBlueprintMaskValidator(blueprintMask).validate()

	require.Nil(t, err)
}

func Test_ValidateWithDuplicatedDoguNames(t *testing.T) {
	dogus := []bpv2.MaskDogu{
		{Name: officialDogu1, TargetState: bpv2.TargetStatePresent},
		{Name: officialDogu1, TargetState: bpv2.TargetStateAbsent},
	}
	blueprintMask := bpv2.BlueprintMask{Dogus: dogus}

	err := newBlueprintMaskValidator(blueprintMask).validate()

	require.NotNil(t, err, "Multiple definitions for the same dogu should lead to an error")
}
