package domain

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/stretchr/testify/assert"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

var (
	compVersion3210 = semver.MustParse("3.2.1-0")
	compVersion3212 = semver.MustParse("3.2.1-2")
	compVersion3213 = semver.MustParse("3.2.1-3")

	testComponentName1 = common.QualifiedComponentName{Namespace: "k8s", SimpleName: "my-component1"}
	testComponentName2 = common.QualifiedComponentName{Namespace: "official", SimpleName: "my-component2"}
	testComponentName3 = common.QualifiedComponentName{Namespace: "testing", SimpleName: "my-component3"}
	testComponentName4 = common.QualifiedComponentName{Namespace: "k8s", SimpleName: "my-component4"}
)

func Test_validate_ok(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: &version3_2_1_0, Absent: true},
		{Name: officialDogu2, Absent: true},
		{Name: officialDogu3, Version: &version3_2_1_0, Absent: false},
		{Name: officialNexus, Version: &version3213},
	}

	components := []Component{
		{Name: testComponentName1, Version: compVersion3210, Absent: true},
		{Name: testComponentName2, Absent: true},
		{Name: testComponentName3, Version: compVersion3212, Absent: false},
		{Name: testComponentName4, Version: compVersion3213},
	}
	blueprint := Blueprint{Dogus: dogus, Components: components}

	err := blueprint.Validate()

	require.NoError(t, err)
}

func Test_validate_multipleErrors(t *testing.T) {

	dogus := []Dogu{
		{Version: nil},
	}
	components := []Component{
		{Version: compVersion3212},
		{Name: testComponentName, Version: compVersion3212},
	}
	blueprint := Blueprint{
		Dogus:      dogus,
		Components: components,
		Config: &Config{
			Global: GlobalConfig{
				Present: nil,
				Absent: []common.GlobalConfigKey{
					"",
				},
			},
		},
	}

	err := blueprint.Validate()

	require.Error(t, err)
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "dogu is invalid")
	assert.ErrorContains(t, err, "dogu version must not be empty")
	assert.ErrorContains(t, err, "component name must not be empty")
	assert.ErrorContains(t, err, `namespace of component "" must not be empty`)
	assert.ErrorContains(t, err, `key for absent global config should not be empty`)
}

func Test_validateDogus_ok(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: &version3_2_1_4, Absent: true},
		//versionIsOptionalForStateAbsent
		{Name: officialDogu2, Absent: true},
		{Name: officialDogu3, Version: &version3212, Absent: false},
		//StateDefaultsToPresent
		{Name: officialNexus, Version: &version3212},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.NoError(t, err)
}

func Test_validateDogus_multipleErrors(t *testing.T) {
	wrongBodySize := resource.MustParse("1Ki")
	dogus := []Dogu{
		{Name: officialDogu1},
		{Name: officialDogu2, ReverseProxyConfig: &ecosystem.ReverseProxyConfig{MaxBodySize: &wrongBodySize}},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dogu proxy body size is not in Decimal SI (\"M\" or \"G\")")
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_validateComponents_ok(t *testing.T) {
	components := []Component{
		{
			Name: common.QualifiedComponentName{
				Namespace:  "k8s",
				SimpleName: "absent-component",
			},
			Absent: true,
		},
		{
			Name: common.QualifiedComponentName{
				SimpleName: "present-component",
				Namespace:  "k8s",
			},
			Version: compVersion3212,
			Absent:  false,
		},
	}
	blueprint := Blueprint{Components: components}

	err := blueprint.validateComponents()

	require.NoError(t, err)
}

func Test_validateComponents_multipleErrors(t *testing.T) {
	components := []Component{
		{Name: testComponentName},
		{Version: compVersion3212},
	}
	blueprint := Blueprint{Components: components}
	err := blueprint.validateComponents()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
	assert.Contains(t, err.Error(), `version of component "k8s/my-component" must not be empty`)
}

func Test_validateDoguUniqueness(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: &version3_2_1_0, Absent: false},
		{Name: officialDogu1, Version: &version3213},
		{Name: officialDogu2, Version: &version3213},
		{Name: officialDogu2, Version: &version3213},
	}

	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDoguUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate dogus")
	assert.Contains(t, err.Error(), "dogu1")
	assert.Contains(t, err.Error(), "dogu2")
}

func Test_validateComponentUniqueness(t *testing.T) {
	components := []Component{
		{
			Name: common.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component1",
			},
			Version: compVersion3210,
			Absent:  false,
		},
		{
			Name: common.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component1",
			},
			Version: compVersion3213},
		{
			Name: common.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component2",
			},
			Version: compVersion3213},
		{
			Name: common.QualifiedComponentName{
				Namespace:  "present",
				SimpleName: "component2",
			},
			Version: compVersion3213},
	}

	blueprint := Blueprint{Components: components}

	err := blueprint.validateComponentUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate components")
	assert.Contains(t, err.Error(), "component1")
	assert.Contains(t, err.Error(), "component2")
}
