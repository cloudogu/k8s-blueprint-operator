package domain

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

func Test_validate_ok(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "absent", Name: "dogu1", Version: version3_2_1_0, TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "dogu2", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu3", Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu4", Version: version3213},
	}

	components := []Component{
		{Name: "absent-component1", DistributionNamespace: "", Version: version3_2_1_0, TargetState: TargetStateAbsent},
		{Name: "absent-component2", TargetState: TargetStateAbsent},
		{Name: "present-component3", DistributionNamespace: "k8s", Version: version3212, TargetState: TargetStatePresent},
		{Name: "present-component4", DistributionNamespace: "k8s", Version: version3213},
	}
	blueprint := Blueprint{Dogus: dogus, Components: components}

	err := blueprint.Validate()

	require.NoError(t, err)
}
func Test_validate_multipleErrors(t *testing.T) {
	dogus := []Dogu{
		{Version: version3212},
	}
	components := []Component{
		{Version: version3212},
		{Name: testComponentName, Version: version3212},
	}
	blueprint := Blueprint{Dogus: dogus, Components: components}

	err := blueprint.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "blueprint is invalid")
	assert.Contains(t, err.Error(), "dogu name must not be empty")
	assert.Contains(t, err.Error(), "component name must not be empty")
	assert.Contains(t, err.Error(), `distribution namespace of component "my-component" must not be empty`)
}

func Test_validateDogus_ok(t *testing.T) {
	dogus := []Dogu{
		{Namespace: "absent", Name: "dogu", Version: version3_2_1_4, TargetState: TargetStateAbsent},
		{Namespace: "absent", Name: "versionIsOptionalForStateAbsent", TargetState: TargetStateAbsent},
		{Namespace: "present", Name: "dogu", Version: version3212, TargetState: TargetStatePresent},
		{Namespace: "present", Name: "StateDefaultsToPresent", Version: version3212},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.NoError(t, err)
}

func Test_validateDogus_multipleErrors(t *testing.T) {
	dogus := []Dogu{
		{Name: "test"},
		{Version: version3212},
	}
	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDogus()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dogu name must not be empty")
	assert.Contains(t, err.Error(), "dogu version must not be empty")
}

func Test_validateComponents_ok(t *testing.T) {
	components := []Component{
		{Name: "absent-component", TargetState: TargetStateAbsent},
		{Name: "present-component", DistributionNamespace: "k8s", Version: version3212, TargetState: TargetStatePresent},
	}
	blueprint := Blueprint{Components: components}

	err := blueprint.validateComponents()

	require.NoError(t, err)
}

func Test_validateComponents_multipleErrors(t *testing.T) {
	components := []Component{
		{Name: testComponentName},
		{Version: version3212},
	}
	blueprint := Blueprint{Components: components}
	err := blueprint.validateComponents()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "component name must not be empty")
	assert.Contains(t, err.Error(), `version of component "my-component" must not be empty`)
}

func Test_validateDoguUniqueness(t *testing.T) {
	dogus := []Dogu{
		{Name: "present/dogu1", Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Name: "present/dogu1", Version: version3213},
		{Name: "present/dogu2", Version: version3213},
		{Name: "present/dogu2", Version: version3213},
	}

	blueprint := Blueprint{Dogus: dogus}

	err := blueprint.validateDoguUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate dogus")
	assert.Contains(t, err.Error(), "present/dogu1")
	assert.Contains(t, err.Error(), "present/dogu2")
}

func Test_validateComponentUniqueness(t *testing.T) {
	components := []Component{
		{Name: "present/component1", Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Name: "present/component1", Version: version3213},
		{Name: "present/component2", Version: version3213},
		{Name: "present/component2", Version: version3213},
	}

	blueprint := Blueprint{Components: components}

	err := blueprint.validateComponentUniqueness()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "there are duplicate components")
	assert.Contains(t, err.Error(), "present/component1")
	assert.Contains(t, err.Error(), "present/component2")
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
		{Namespace: "present", Name: "dogu1", Version: version3_2_1_0, TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu2", Version: version3213, TargetState: TargetStatePresent},
		{Namespace: "present", Name: "dogu3", Version: version3213, TargetState: TargetStatePresent},
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
