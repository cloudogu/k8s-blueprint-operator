package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_validateMaskTargetDogu_noErrorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := MaskTargetDogu{Name: "present/dogu", TargetState: TargetStatePresent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateMaskTargetDogu_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := MaskTargetDogu{Name: "present/dogu", TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateMaskTargetDogu_missingStateOkayForPresentDogu(t *testing.T) {
	dogu := MaskTargetDogu{Name: "present/dogu", Version: "2018-1"}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateMaskTargetDogu_errorOnMissingNameForDogu(t *testing.T) {
	dogu := MaskTargetDogu{}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "could not Validate blueprint mask, dogu field Name must not be empty")
}
