package serializer

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testNamespace = common.ComponentNamespace("k8s")
var testNamespaceString = "k8s"

var testDeployConfig = map[string]interface{}{"key": "value", "key1": map[string]interface{}{"key": "value2"}}

func Test_convertToComponentDiffDTO(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Version: nil, Absent: true},
			Expected:      domain.ComponentDiffState{Version: testSemverVersionHigh, Namespace: testNamespace, DeployConfig: testDeployConfig},
			NeededActions: []domain.Action{domain.ActionInstall}}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Version: nil, Absent: true},
			Expected:      crd.ComponentDiffState{Version: &testSemverVersionHighRaw, Namespace: testNamespaceString, DeployConfig: testDeployConfig},
			NeededActions: []crd.ComponentAction{domain.ActionInstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Version: testSemverVersionHigh},
			Expected:      domain.ComponentDiffState{Version: nil, Absent: false},
			NeededActions: []domain.Action{domain.ActionUninstall}}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Version: &testSemverVersionHighRaw},
			Expected:      crd.ComponentDiffState{Version: nil, Absent: true},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}

func Test_convertToComponentDiffDomain(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: "", Version: nil, Absent: true},
			Expected:      crd.ComponentDiffState{Namespace: testNamespaceString, Version: &testSemverVersionHighRaw, DeployConfig: testDeployConfig},
			NeededActions: []crd.ComponentAction{domain.ActionInstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "", Version: nil, Absent: true},
			Expected:      domain.ComponentDiffState{Namespace: testNamespace, Version: testSemverVersionHigh, DeployConfig: testDeployConfig},
			NeededActions: []domain.Action{domain.ActionInstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: testNamespaceString, Version: &testSemverVersionHighRaw},
			Expected:      crd.ComponentDiffState{Namespace: "", Version: nil, Absent: true},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: testNamespace, Version: testSemverVersionHigh},
			Expected:      domain.ComponentDiffState{Namespace: "", Version: nil, Absent: true},
			NeededActions: []domain.Action{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should fail to parse version", func(t *testing.T) {
		// given
		versionABC := "a-b-c"
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: "", Version: &versionABC},
			Expected:      crd.ComponentDiffState{Namespace: "", Version: &versionABC},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}

		// when
		_, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse actual version")
		assert.ErrorContains(t, err, "failed to parse expected version")
	})
	t.Run("should accept dev version", func(t *testing.T) {
		compVersion080dev := semver.MustParse("0.8.0-dev")
		// given
		Version080String := compVersion080dev.String()
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: testNamespaceString, Version: &Version080String},
			Expected:      crd.ComponentDiffState{Namespace: "", Version: nil, Absent: true},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: testNamespace, Version: compVersion080dev},
			Expected:      domain.ComponentDiffState{Namespace: "", Version: nil, Absent: true},
			NeededActions: []domain.Action{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}
