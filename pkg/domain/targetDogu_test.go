package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_validateDogu_errorOnMissingDoguName(t *testing.T) {
	dogus := []TargetDogu{
		{Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu field Name must not be empty")
}

func Test_validateDogu_errorOnEmptyDoguName(t *testing.T) {
	dogu := TargetDogu{Name: "", Version: "3.2.1-2", TargetState: TargetStatePresent}

	err := dogu.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu field Name must not be empty")
}

func Test_validateDogu_errorOnMissingVersionForPresentDogu(t *testing.T) {
	dogu := TargetDogu{Name: "present/dogu", TargetState: TargetStatePresent}

	err := dogu.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "dogu field Version must not be empty")
}

func Test_validateDogu_missingVersionOkayForAbsentDogu(t *testing.T) {
	dogu := TargetDogu{Name: "present/dogu", TargetState: TargetStateAbsent}

	err := dogu.Validate()

	require.Nil(t, err)
}

func Test_validateDogu_missingStateOkayForPresentDogu(t *testing.T) {
	dogu := TargetDogu{Name: "present/dogu", Version: "2018-1"}

	err := dogu.Validate()

	require.Nil(t, err)
}
