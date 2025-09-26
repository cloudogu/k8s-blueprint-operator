package serializer

import (
	"testing"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
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
			Expected:      crd.ComponentDiffState{Version: nil, Absent: false},
			NeededActions: []crd.ComponentAction{domain.ActionUninstall},
		}
		assert.Equal(t, expected, actual)
	})
}
