package domain

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var version1_2_3, _ = core.ParseVersion("1.2.3")

func Test_validateComponents_errorOnMissingComponentVersion(t *testing.T) {
	component := Component{Name: "present-component", DistributionNamespace: "k8s", TargetState: TargetStatePresent}

	err := component.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), `version of component "present-component" must not be empty`)
}

func Test_validateComponents_errorOnEmptyComponentVersion(t *testing.T) {
	component := Component{Name: "present-component", Version: core.Version{}, TargetState: TargetStatePresent}

	err := component.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "version of component \"present-component\" must not be empty")
}

func Test_validateComponents_errorOnMissingComponentName(t *testing.T) {
	component := Component{Version: version1_2_3, TargetState: TargetStatePresent}

	err := component.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
}

func Test_validateComponents_errorOnEmptyComponentName(t *testing.T) {
	component := Component{Name: "", Version: version1_2_3, TargetState: TargetStatePresent}

	err := component.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
}

func Test_validateComponents_emptyComponentStateDefaultsToPresent(t *testing.T) {
	component := Component{Name: "present-component", Version: version1_2_3, DistributionNamespace: "k8s"}

	err := component.Validate()

	require.NoError(t, err)
}

func Test_validateComponents_missingComponentVersionOkayForAbsent(t *testing.T) {
	component := Component{Name: "present-component", TargetState: TargetStateAbsent}

	err := component.Validate()

	require.NoError(t, err)
}
