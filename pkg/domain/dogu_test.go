package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_TargetDogu_validate_errorOnMissingDoguName(t *testing.T) {
	dogus := []Dogu{
		{Version: version3212, TargetState: TargetStatePresent},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu name must not be empty")
}

func Test_TargetDogu_validate_errorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, TargetState: TargetStatePresent}

	err := dogu.validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_TargetDogu_validate_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, TargetState: TargetStateAbsent}

	err := dogu.validate()

	require.Nil(t, err)
}

func Test_TargetDogu_validate_defaultToPresentState(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, Version: version1_2_3}

	err := dogu.validate()

	require.Nil(t, err)
	assert.Equal(t, TargetState(TargetStatePresent), dogu.TargetState)
}

func Test_TargetDogu_validate_errorOnUnknownTargetState(t *testing.T) {
	dogu := Dogu{Name: officialDogu1, TargetState: -1}

	err := dogu.validate()

	require.Error(t, err)
	require.ErrorContains(t, err, "dogu target state is invalid: official/dogu1")
}
