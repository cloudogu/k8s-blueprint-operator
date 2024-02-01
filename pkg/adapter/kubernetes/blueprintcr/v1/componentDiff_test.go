package v1

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_convertToComponentDiffDTO(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			ComponentName: testComponentName,
			Actual:        domain.ComponentDiffState{Version: core.Version{}, InstallationState: domain.TargetStateAbsent},
			Expected:      domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			NeededAction:  domain.ActionInstall}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := ComponentDiff{
			Actual:       ComponentDiffState{Version: "", InstallationState: "absent"},
			Expected:     ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
			NeededAction: domain.ActionInstall,
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			ComponentName: testComponentName,
			Actual:        domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Version: core.Version{}, InstallationState: domain.TargetStateAbsent},
			NeededAction:  domain.ActionUninstall}

		// when
		actual := convertToComponentDiffDTO(domainDiff)

		// then
		expected := ComponentDiff{
			Actual:       ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
			Expected:     ComponentDiffState{Version: "", InstallationState: "absent"},
			NeededAction: domain.ActionUninstall,
		}
		assert.Equal(t, expected, actual)
	})
}

func Test_convertToComponentDiffDomain(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:       ComponentDiffState{Version: "", InstallationState: "absent"},
			Expected:     ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
			NeededAction: domain.ActionInstall,
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			ComponentName: testComponentName,
			Actual:        domain.ComponentDiffState{Version: core.Version{}, InstallationState: domain.TargetStateAbsent},
			Expected:      domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			NeededAction:  domain.ActionInstall,
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:       ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
			Expected:     ComponentDiffState{Version: "", InstallationState: "absent"},
			NeededAction: domain.ActionUninstall,
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			ComponentName: testComponentName,
			Actual:        domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			Expected:      domain.ComponentDiffState{Version: core.Version{}, InstallationState: domain.TargetStateAbsent},
			NeededAction:  domain.ActionUninstall,
		}
		assert.Equal(t, expected, actual)
	})
}
