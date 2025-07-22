package serializer

import (
	"github.com/Masterminds/semver/v3"
	crd "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const testNamespace = "k8s"

var testDeployConfig = map[string]interface{}{"key": "value", "key1": map[string]interface{}{"key": "value2"}}

func Test_convertToComponentDiffDTO(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Version: nil, InstallationState: domain.TargetStateAbsent},
			Expected:      domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent, Namespace: testNamespace, DeployConfig: testDeployConfig},
			NeededActions: []domain.Action{domain.ActionInstall}}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Version: "", InstallationState: "absent"},
			Expected:      crd.ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present", Namespace: testNamespace, DeployConfig: testDeployConfig},
			NeededActions: []crd.ComponentAction{domain.ActionInstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Version: nil, InstallationState: domain.TargetStateAbsent},
			NeededActions: []domain.Action{domain.ActionUninstall}}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
			Expected:      crd.ComponentDiffState{Version: "", InstallationState: "absent"},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}

func Test_convertToComponentDiffDomain(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			Expected:      crd.ComponentDiffState{Namespace: "k8s", Version: testVersionHighRaw, InstallationState: "present", DeployConfig: testDeployConfig},
			NeededActions: []crd.ComponentAction{domain.ActionInstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: domain.TargetStateAbsent},
			Expected:      domain.ComponentDiffState{Namespace: "k8s", Version: testVersionHigh, InstallationState: domain.TargetStatePresent, DeployConfig: testDeployConfig},
			NeededActions: []domain.Action{domain.ActionInstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: "k8s", Version: testVersionHighRaw, InstallationState: "present"},
			Expected:      crd.ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "k8s", Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: domain.TargetStateAbsent},
			NeededActions: []domain.Action{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should fail in all ways", func(t *testing.T) {
		// given
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: "", Version: "a-b-c", InstallationState: "☹"},
			Expected:      crd.ComponentDiffState{Namespace: "", Version: "a-b-c", InstallationState: "☹"},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}

		// when
		_, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse actual version")
		assert.ErrorContains(t, err, "failed to parse expected version")
		assert.ErrorContains(t, err, "failed to parse actual installation state")
		assert.ErrorContains(t, err, "failed to parse expected installation state")
	})
	t.Run("should accept dev version", func(t *testing.T) {
		compVersion080dev := semver.MustParse("0.8.0-dev")
		// given
		diff := crd.ComponentDiff{
			Actual:        crd.ComponentDiffState{Namespace: "k8s", Version: compVersion080dev.String(), InstallationState: "present"},
			Expected:      crd.ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "k8s", Version: compVersion080dev, InstallationState: domain.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: domain.TargetStateAbsent},
			NeededActions: []domain.Action{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}
