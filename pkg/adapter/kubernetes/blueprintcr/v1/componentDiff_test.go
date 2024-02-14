package v1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_convertToComponentDiffDTO(t *testing.T) {
	t.Run("should copy model diff to DTO diff - absent", func(t *testing.T) {
		// given
		domainDiff := domain.ComponentDiff{
			Name:         testComponentName,
			Actual:       domain.ComponentDiffState{Version: nil, InstallationState: domain.TargetStateAbsent},
			Expected:     domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			NeededAction: domain.ActionInstall}

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
			Name:         testComponentName,
			Actual:       domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			Expected:     domain.ComponentDiffState{Version: nil, InstallationState: domain.TargetStateAbsent},
			NeededAction: domain.ActionUninstall}

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
			Actual:       ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			Expected:     ComponentDiffState{Namespace: "k8s", Version: testVersionHighRaw, InstallationState: "present"},
			NeededAction: domain.ActionInstall,
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:         testComponentName,
			Actual:       domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: domain.TargetStateAbsent},
			Expected:     domain.ComponentDiffState{Namespace: "k8s", Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			NeededAction: domain.ActionInstall,
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should copy model diff to DTO diff - present", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:       ComponentDiffState{Namespace: "k8s", Version: testVersionHighRaw, InstallationState: "present"},
			Expected:     ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			NeededAction: domain.ActionUninstall,
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:         testComponentName,
			Actual:       domain.ComponentDiffState{Namespace: "k8s", Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
			Expected:     domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: domain.TargetStateAbsent},
			NeededAction: domain.ActionUninstall,
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should fail in all ways", func(t *testing.T) {
		// given
		diff := ComponentDiff{
			Actual:       ComponentDiffState{Namespace: "", Version: "a-b-c", InstallationState: "☹"},
			Expected:     ComponentDiffState{Namespace: "", Version: "a-b-c", InstallationState: "☹"},
			NeededAction: domain.ActionUninstall,
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
			Actual:       ComponentDiffState{Namespace: "k8s", Version: compVersion080dev.String(), InstallationState: "present"},
			Expected:     ComponentDiffState{Namespace: "", Version: "", InstallationState: "absent"},
			NeededAction: domain.ActionUninstall,
		}

		// when
		actual, err := convertToComponentDiffDomain(testComponentName, diff)

		// then
		require.NoError(t, err)
		expected := domain.ComponentDiff{
			Name:         testComponentName,
			Actual:       domain.ComponentDiffState{Namespace: "k8s", Version: compVersion080dev, InstallationState: domain.TargetStatePresent},
			Expected:     domain.ComponentDiffState{Namespace: "", Version: nil, InstallationState: domain.TargetStateAbsent},
			NeededAction: domain.ActionUninstall,
		}
		assert.Equal(t, expected, actual)
	})
}
