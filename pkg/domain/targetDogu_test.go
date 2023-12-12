package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_TargetDogu_validate_errorOnMissingDoguName(t *testing.T) {
	dogus := []Dogu{
		{Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu field Name must not be empty")
}

func Test_TargetDogu_validate_errorOnEmptyDoguName(t *testing.T) {
	dogu := Dogu{Name: "", Version: "3.2.1-2", TargetState: TargetStatePresent}

	err := dogu.validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu field Name must not be empty")
}

func Test_TargetDogu_validate_errorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := Dogu{Name: "present/dogu", TargetState: TargetStatePresent}

	err := dogu.validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu field Version must not be empty")
}

func Test_TargetDogu_validate_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := Dogu{Namespace: "present", Name: "dogu", TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_TargetDogu_validate_defaultToPresentState(t *testing.T) {
	dogu := Dogu{Namespace: "present", Name: "dogu", Version: "2018-1"}

	err := dogu.validate()

	require.Nil(t, err)
	assert.Equal(t, TargetState(TargetStatePresent), dogu.TargetState)
}

func Test_TargetDogu_validate_errorOnUnknownTargetState(t *testing.T) {
	dogu := Dogu{Namespace: "official", Name: "dogu1", TargetState: -1}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu target state is invalid: official/dogu1")
}
