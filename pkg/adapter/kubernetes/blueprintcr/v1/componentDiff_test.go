package v1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

const testNamespace = "k8s"

var testDeployConfig = map[string]interface{}{"key": "value", "key1": map[string]interface{}{"key": "value2"}}

func Test_convertToComponentDiffDTO(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Version: nil, InstallationState: bpv2.TargetStateAbsent},
			Expected:      domain.ComponentDiffState{Version: testVersionHigh, InstallationState: bpv2.TargetStatePresent, Namespace: testNamespace, DeployConfig: testDeployConfig},
			NeededActions: []domain.Action{domain.ActionInstall}}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := ComponentDiff{
			Actual:        ComponentDiffState{Version: "", InstallationState: "absent"},
			Expected:      ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present", Namespace: testNamespace, DeployConfig: testDeployConfig},
			NeededActions: []ComponentAction{domain.ActionInstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Version: testVersionHigh, InstallationState: bpv2.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Version: nil, InstallationState: bpv2.TargetStateAbsent},
			NeededActions: []domain.Action{domain.ActionUninstall}}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := ComponentDiff{
			Actual:        ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
			Expected:      ComponentDiffState{Version: "", InstallationState: "absent"},
			NeededActions: []ComponentAction{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}

func Test_convertToComponentDiffDomain(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:        ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			Expected:      ComponentDiffState{Namespace: "k8s", Version: testVersionHighRaw, InstallationState: "present", DeployConfig: testDeployConfig},
			NeededActions: []ComponentAction{domain.ActionInstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: bpv2.TargetStateAbsent},
			Expected:      domain.ComponentDiffState{Namespace: "k8s", Version: testVersionHigh, InstallationState: bpv2.TargetStatePresent, DeployConfig: testDeployConfig},
			NeededActions: []domain.Action{domain.ActionInstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:        ComponentDiffState{Namespace: "k8s", Version: testVersionHighRaw, InstallationState: "present"},
			Expected:      ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			NeededActions: []ComponentAction{domain.ActionUninstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "k8s", Version: testVersionHigh, InstallationState: bpv2.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: bpv2.TargetStateAbsent},
			NeededActions: []domain.Action{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should fail in all ways", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:        ComponentDiffState{Namespace: "", Version: "a-b-c", InstallationState: "☹"},
			Expected:      ComponentDiffState{Namespace: "", Version: "a-b-c", InstallationState: "☹"},
			NeededActions: []ComponentAction{domain.ActionUninstall},
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
		diff := ComponentDiff{
			Actual:        ComponentDiffState{Namespace: "k8s", Version: compVersion080dev.String(), InstallationState: "present"},
			Expected:      ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			NeededActions: []ComponentAction{domain.ActionUninstall},
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:          testComponentName,
			Actual:        domain.ComponentDiffState{Namespace: "k8s", Version: compVersion080dev, InstallationState: bpv2.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: bpv2.TargetStateAbsent},
			NeededActions: []domain.Action{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}
