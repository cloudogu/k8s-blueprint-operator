package domain

import (
	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
)

var (
	compVersion123 = semver.MustParse("1.2.3")
	version123, _  = core.ParseVersion("1.2.3")
)

func TestComponent_Validate(t *testing.T) {
	t.Run("errorOnMissingComponentVersion", func(t *testing.T) {
		component := bpv2.Component{Name: testComponentName, TargetState: bpv2.TargetStatePresent}

		err := NewComponentValidator(component).validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), `version of component "k8s/my-component" must not be empty`)
	})

	t.Run("errorOnEmptyComponentVersion", func(t *testing.T) {
		component := bpv2.Component{Name: testComponentName, Version: nil, TargetState: bpv2.TargetStatePresent}

		err := NewComponentValidator(component).validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "version of component \"k8s/my-component\" must not be empty")
	})

	t.Run("errorOnMissingComponentName", func(t *testing.T) {
		component := bpv2.Component{Version: compVersion123, TargetState: bpv2.TargetStatePresent}

		err := NewComponentValidator(component).validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "component name must not be empty")
	})

	t.Run("errorOnEmptyComponentNamespace", func(t *testing.T) {
		component := bpv2.Component{Name: bpv2.QualifiedComponentName{Namespace: "", SimpleName: "test"}, Version: compVersion123, TargetState: bpv2.TargetStatePresent}
		err := NewComponentValidator(component).validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "namespace of component \"test\" must not be empty")
	})

	t.Run("errorOnEmptyComponentName", func(t *testing.T) {
		component := bpv2.Component{Name: bpv2.QualifiedComponentName{Namespace: "k8s"}, Version: compVersion123, TargetState: bpv2.TargetStatePresent}

		err := NewComponentValidator(component).validate()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "component name must not be empty")
	})

	t.Run("emptyComponentStateDefaultsToPresent", func(t *testing.T) {
		component := bpv2.Component{Name: testComponentName, Version: compVersion123}

		err := NewComponentValidator(component).validate()

		require.NoError(t, err)
		assert.Equal(t, bpv2.TargetState(bpv2.TargetStatePresent), component.TargetState)
	})

	t.Run("missingComponentVersionOkayForAbsent", func(t *testing.T) {
		component := bpv2.Component{Name: testComponentName, TargetState: bpv2.TargetStateAbsent}

		err := NewComponentValidator(component).validate()

		require.NoError(t, err)
	})
}
