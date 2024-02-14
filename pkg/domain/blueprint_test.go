package domain

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

var (
	compVersion3210 = semver.MustParse("3.2.1-0")
	compVersion3212 = semver.MustParse("3.2.1-2")
	compVersion3213 = semver.MustParse("3.2.1-3")

	testComponentName1 = common.QualifiedComponentName{Namespace: "k8s", Name: "my-component1"}
	testComponentName2 = common.QualifiedComponentName{Namespace: "official", Name: "my-component2"}
	testComponentName3 = common.QualifiedComponentName{Namespace: "testing", Name: "my-component3"}
	testComponentName4 = common.QualifiedComponentName{Namespace: "k8s", Name: "my-component4"}
)

func Test_validate_ok(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: TargetStateAbsent},
		{Name: officialDogu2, TargetState: TargetStateAbsent},
		{Name: officialDogu3, Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Name: officialNexus, Version: version3213},
	}

	components := []Component{
		{Name: testComponentName1, Version: compVersion3210, TargetState: TargetStateAbsent},
		{Name: testComponentName2, TargetState: TargetStateAbsent},
		{Name: testComponentName3, Version: compVersion3212, TargetState: TargetStatePresent},
		{Name: testComponentName4, Version: compVersion3213},
	}
	blueprint := Blueprint{Dogus: dogus, Components: components}

	err := blueprint.Validate()

	require.NoError(t, err)
}

func Test_validate_multipleErrors(t *testing.T) {
	dogus := []Dogu{
		{Version: version3212, TargetState: 666},
	}
	components := []Component{
		{Version: compVersion3212},
		{Name: testComponentName, Version: compVersion3212},
	}
	blueprint := Blueprint{Dogus: dogus, Components: components}

	err := blueprint.Validate()

	require.Error(t, err)
	assert.ErrorContains(t, err, "blueprint is invalid")
	assert.ErrorContains(t, err, "dogu is invalid")
	assert.ErrorContains(t, err, "namespace of dogu \"\" must not be empty")
	assert.ErrorContains(t, err, "dogu name must not be empty: '/'")
	assert.ErrorContains(t, err, "dogu target state is invalid")
	assert.ErrorContains(t, err, "component name must not be empty")
	assert.ErrorContains(t, err, `namespace of component "" must not be empty`)
}

func Test_validateDogus_ok(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1, Version: version3_2_1_4, TargetState: TargetStateAbsent},
		//versionIsOptionalForStateAbsent
		{Name: officialDogu2, TargetState: TargetStateAbsent},
		{Name: officialDogu3, Version: version3212, TargetState: TargetStatePresent},
		//StateDefaultsToPresent
		{Name: officialNexus, Version: version3212},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.NoError(t, err)
}

func Test_validateDogus_multipleErrors(t *testing.T) {
	dogus := []Dogu{
		{Name: officialDogu1},
		{Name: officialDogu2, TargetState: 666},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dogu target state is invalid")
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_validateComponents_ok(t *testing.T) {
	components := []Component{
		{
			Name: common.QualifiedComponentName{
				Namespace: "k8s",
				Name:      "absent-component",
			},
			TargetState: TargetStateAbsent,
		},
		{
			Name: common.QualifiedComponentName{
				Name:      "present-component",
				Namespace: "k8s",
			},
			Version:     compVersion3212,
			TargetState: TargetStatePresent,
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
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Name: officialDogu1, Version: version3213},
		{Name: officialDogu2, Version: version3213},
		{Name: officialDogu2, Version: version3213},
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
				Namespace: "present",
				Name:      "component1",
			},
			Version:     compVersion3210,
			TargetState: TargetStatePresent,
		},
		{
			Name: common.QualifiedComponentName{
				Namespace: "present",
				Name:      "component1",
			},
			Version: compVersion3213},
		{
			Name: common.QualifiedComponentName{
				Namespace: "present",
				Name:      "component2",
			},
			Version: compVersion3213},
		{
			Name: common.QualifiedComponentName{
				Namespace: "present",
				Name:      "component2",
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

func TestValidationGlobalValid(t *testing.T) {
	registryConfig := map[string]map[string]interface{}{
		"_global": {
			"key1": "value1",
		},
	}
	blueprint := Blueprint{
		RegistryConfig: registryConfig,
	}

	err := blueprint.Validate()
	assert.Nil(t, err)
}

func TestValidationGlobalEmptyKeyError(t *testing.T) {
	globalRegistryKeys := map[string]map[string]interface{}{
		"_global": {
			"": "myVal2",
		},
	}
	blueprint := Blueprint{
		RegistryConfig: globalRegistryKeys,
	}

	err := blueprint.Validate()
	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "blueprint is invalid: a Config key is empty")
}

func TestSetDoguRegistryKeysSuccessful(t *testing.T) {
	dogu1 := map[string]interface{}{"keyDogu1": "valDogu1"}
	dogu2 := map[string]interface{}{"keyDogu2": "valDogu2"}
	dogu3 := map[string]interface{}{"keyDogu3": "valDogu3"}
	dogus := []Dogu{
		{Name: officialDogu1, Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Name: officialDogu2, Version: version3213, TargetState: TargetStatePresent},
		{Name: officialDogu3, Version: version3213, TargetState: TargetStatePresent},
	}

	blueprint := Blueprint{
		Dogus: dogus,
		RegistryConfig: map[string]map[string]interface{}{
			"dogu1":         dogu1,
			"present/dogu2": dogu2,
			"dogu3":         dogu3,
		},
	}
	err := blueprint.Validate()
	assert.Nil(t, err)
}
