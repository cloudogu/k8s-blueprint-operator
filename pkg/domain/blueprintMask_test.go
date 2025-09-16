package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/require"
)

var version3_2_1_0, _ = core.ParseVersion("3.2.1-0")

var (
	officialNamespace = cescommons.Namespace("official")
	k8sNamespace      = cescommons.Namespace("k8s")
	officialDogu1     = cescommons.QualifiedName{Namespace: officialNamespace, SimpleName: cescommons.SimpleName("dogu1")}
	officialDogu2     = cescommons.QualifiedName{Namespace: officialNamespace, SimpleName: cescommons.SimpleName("dogu2")}
	officialDogu3     = cescommons.QualifiedName{Namespace: officialNamespace, SimpleName: cescommons.SimpleName("dogu3")}
)

func Test_Validate(t *testing.T) {
	dogus := []MaskDogu{
		{Name: officialDogu1, Version: version3_2_1_0, Absent: true},
		{Name: officialDogu2, Absent: true},
		{Name: officialDogu3, Version: version3_2_1_0, Absent: false},
	}
	blueprintMask := BlueprintMask{
		Dogus: dogus,
	}

	err := blueprintMask.Validate()

	require.Nil(t, err)
}

func Test_ValidateWithDuplicatedDoguNames(t *testing.T) {
	dogus := []MaskDogu{
		{Name: officialDogu1, Absent: false},
		{Name: officialDogu1, Absent: true},
	}
	blueprintMask := BlueprintMask{Dogus: dogus}

	err := blueprintMask.Validate()

	require.NotNil(t, err, "Multiple definitions for the same dogu should lead to an error")
}
