package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_validateDogu_errorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := TargetDogu{Name: "present/dogu", TargetState: TargetStatePresent}

	err := dogu.validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not validate blueprint, dogu field Version must not be empty")
}

func Test_validateDogu_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := TargetDogu{Name: "present/dogu", TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_validateDogu_missingStateOkayForPresentDogu(t *testing.T) {
	dogu := TargetDogu{Name: "present/dogu", Version: "2018-1"}

	err := dogu.validate()

	require.Nil(t, err)
}
