package domain

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	compVersion123 = semver.MustParse("1.2.3")
	version123, _  = core.ParseVersion("1.2.3")
)

func TestComponent_Validate(t *testing.T) {
	t.Run("errorOnMissingComponentVersion", func(t *testing.T) {
		component := Component{Name: testComponentName, TargetState: TargetStatePresent}

		err := component.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), `version of component "k8s/my-component" must not be empty`)
	})

	t.Run("errorOnEmptyComponentVersion", func(t *testing.T) {
		component := Component{Name: testComponentName, Version: nil, TargetState: TargetStatePresent}

		err := component.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "version of component \"k8s/my-component\" must not be empty")
	})

	t.Run("errorOnMissingComponentName", func(t *testing.T) {
		component := Component{Version: compVersion123, TargetState: TargetStatePresent}

		err := component.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "component name must not be empty")
	})

	t.Run("errorOnEmptyComponentNamespace", func(t *testing.T) {
		component := Component{Name: common.QualifiedComponentName{Namespace: "", SimpleName: "test"}, Version: compVersion123, TargetState: TargetStatePresent}
		err := component.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "namespace of component \"test\" must not be empty")
	})

	t.Run("errorOnEmptyComponentName", func(t *testing.T) {
		component := Component{Name: common.QualifiedComponentName{Namespace: "k8s"}, Version: compVersion123, TargetState: TargetStatePresent}

		err := component.Validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "component name must not be empty")
	})

	t.Run("emptyComponentStateDefaultsToPresent", func(t *testing.T) {
		component := Component{Name: testComponentName, Version: compVersion123}

		err := component.Validate()

		require.NoError(t, err)
		assert.Equal(t, TargetState(TargetStatePresent), component.TargetState)
	})

	t.Run("missingComponentVersionOkayForAbsent", func(t *testing.T) {
		component := Component{Name: testComponentName, TargetState: TargetStateAbsent}

		err := component.Validate()

		require.NoError(t, err)
	})
}
