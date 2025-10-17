package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MaskTargetDogu_validate_noErrorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := MaskDogu{Name: officialDogu1, Absent: false}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_MaskTargetDogu_validate_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := MaskDogu{Name: officialDogu1, Absent: true}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_MaskTargetDogu_validate_defaultToPresentState(t *testing.T) {
	version, _ := core.ParseVersion("2018-1")
	dogu := MaskDogu{Name: officialDogu1, Version: version}

	err := dogu.validate()

	require.Nil(t, err)
	assert.False(t, dogu.Absent)
}

func Test_MaskTargetDogu_validate_errorOnMissingNameForDogu(t *testing.T) {
	dogu := MaskDogu{Name: cescommons.QualifiedName{Namespace: "official"}}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu name must not be empty")
}
