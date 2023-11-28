package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_validateComponents_errorOnMissingComponentVersion(t *testing.T) {
	component := Component{Name: "present-component", TargetState: TargetStatePresent}

	err := component.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "component version must not be empty")
}

func Test_validateComponents_errorOnEmptyComponentVersion(t *testing.T) {
	component := Component{Name: "present/component", Version: "", TargetState: TargetStatePresent}

	err := component.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "component version must not be empty")
}

func Test_validateComponents_errorOnMissingComponentName(t *testing.T) {
	component := Component{Version: "1.2.3", TargetState: TargetStatePresent}

	err := component.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
}

func Test_validateComponents_errorOnEmptyComponentName(t *testing.T) {
	component := Component{Name: "", Version: "1.2.3", TargetState: TargetStatePresent}

	err := component.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
}

func Test_validateComponents_emptyComponentStateDefaultsToPresent(t *testing.T) {
	component := Component{Name: "present-component", Version: "1.2.3"}

	err := component.Validate()

	require.Nil(t, err)
}

func Test_validateComponents_missingComponentVersionOkayForAbsent(t *testing.T) {
	component := Component{Name: "present-component", TargetState: TargetStateAbsent}

	err := component.Validate()

	require.Nil(t, err)
}
