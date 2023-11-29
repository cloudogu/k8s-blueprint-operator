package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_validateMaskTargetDogu_noErrorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := MaskTargetDogu{Namespace: "present", Name: "dogu", TargetState: TargetStatePresent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateMaskTargetDogu_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := MaskTargetDogu{Namespace: "present", Name: "dogu", TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateMaskTargetDogu_missingStateOkayForPresentDogu(t *testing.T) {
	dogu := MaskTargetDogu{Namespace: "present", Name: "dogu", Version: "2018-1"}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateMaskTargetDogu_errorOnMissingNameForDogu(t *testing.T) {
	dogu := MaskTargetDogu{}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu is invalid: dogu field Namespace must not be empty")
}
