package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
)

func Test_MaskTargetDogu_validate_noErrorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := bpv2.MaskDogu{Name: officialDogu1, TargetState: TargetStatePresent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_MaskTargetDogu_validate_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := bpv2.MaskDogu{Name: officialDogu1, TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_MaskTargetDogu_validate_defaultToPresentState(t *testing.T) {
	version, _ := core.ParseVersion("2018-1")
	dogu := bpv2.MaskDogu{Name: officialDogu1, Version: version}

	err := dogu.validate()

	require.Nil(t, err)
	assert.Equal(t, TargetState(TargetStatePresent), dogu.TargetState)
}

func Test_MaskTargetDogu_validate_errorOnMissingNameForDogu(t *testing.T) {
	dogu := bpv2.MaskDogu{Name: cescommons.QualifiedName{Namespace: "official"}}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu name must not be empty")
}

func Test_MaskTargetDogu_validate_errorOnUnknownTargetState(t *testing.T) {
	dogu := bpv2.MaskDogu{Name: officialDogu1, TargetState: -1}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu target state is invalid: official/dogu1")
}
