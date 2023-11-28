package domain

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validate_ok(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-0", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-0", TargetState: TargetStatePresent},
		{Name: "present/dogu4", Version: "1.2.3-4"},
	}

	components := []Component{
		{Name: "absent/component1", Version: "3.2.1-0", TargetState: TargetStateAbsent},
		{Name: "absent/component2", TargetState: TargetStateAbsent},
		{Name: "present-component3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/component4", Version: "1.2.3-4"},
	}
	blueprint := BlueprintV2{Dogus: dogus, Components: components}

	err := blueprint.Validate()

	require.Nil(t, err)
}
func Test_validate_errorOnDoguProblem(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "absent/dogu", TargetState: TargetStateAbsent},
		{Name: "present/dogu", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Version: "3.2.1-2"},
	}
	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not Validate blueprint: dogu field Name must not be empty")
}

func Test_validateDogus_ok(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "absent/dogu", Version: "1.2.3-9", TargetState: TargetStateAbsent},
		{Name: "absent/versionIsOptionalForStateAbsent", TargetState: TargetStateAbsent},
		{Name: "present/dogu", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/StateDefaultsToPresent", Version: "3.2.1-2"},
	}
	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.validateDogus()

	require.Nil(t, err)
}

func Test_validateDogus_errorOnMissingDoguName(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "absent/dogu", TargetState: TargetStateAbsent},
		{Name: "present/dogu", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Version: "3.2.1-2"},
	}
	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.validateDogus()

	require.NotNil(t, err)
}

func Test_validateComponents_ok(t *testing.T) {
	components := []Component{
		{Name: "absent-component", TargetState: TargetStateAbsent},
		{Name: "present-component", Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := BlueprintV2{Components: components}

	err := blueprint.validateComponents()

	require.Nil(t, err)
}

func Test_validateComponents_errorOnMissingName(t *testing.T) {
	components := []Component{
		{Name: "absent-component", TargetState: TargetStateAbsent},
		{Name: "present-component", Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := BlueprintV2{Components: components}
	err := blueprint.validateComponents()

	require.Nil(t, err)
}

func Test_validateDoguUniqueness(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "present/dogu1", Version: "3.2.1-0", TargetState: TargetStatePresent},
		{Name: "present/dogu1", Version: "1.2.3-4"},
	}

	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.validateDoguUniqueness()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not Validate blueprint, there are duplicate dogus: [present/dogu1]")
}

func Test_validateComponentUniqueness(t *testing.T) {
	components := []Component{
		{Name: "present/component1", Version: "3.2.1-0", TargetState: TargetStatePresent},
		{Name: "present/component1", Version: "1.2.3-4"},
	}

	blueprint := BlueprintV2{Components: components}

	err := blueprint.validateComponentUniqueness()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not Validate blueprint, there are duplicate components: [present/component1]")
}

func TestValidationGlobalValid(t *testing.T) {
	registryConfig := map[string]map[string]interface{}{
		"_global": {
			"key1": "value1",
		},
	}
	blueprint := BlueprintV2{
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
	blueprint := BlueprintV2{
		RegistryConfig: globalRegistryKeys,
	}

	err := blueprint.Validate()
	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "could not Validate blueprint, a config key is empty")
}

func TestSetDoguRegistryKeysSuccessful(t *testing.T) {
	dogu1 := map[string]interface{}{"keyDogu1": "valDogu1"}
	dogu2 := map[string]interface{}{"keyDogu2": "valDogu2"}
	dogu3 := map[string]interface{}{"keyDogu3": "valDogu3"}
	dogus := []TargetDogu{
		{Name: "present/dogu1", Version: "3.2.1-0", TargetState: TargetStatePresent},
		{Name: "present/dogu2", Version: "1.2.3-4", TargetState: TargetStatePresent},
		{Name: "present/dogu3", Version: "1.2.3-4", TargetState: TargetStatePresent},
	}

	blueprint := BlueprintV2{
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
