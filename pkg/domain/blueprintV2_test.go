package domain

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetState_String(t *testing.T) {
	tests := []struct {
		name  string
		state TargetState
		want  string
	}{
		{
			"String() map enum to string",
			TargetStatePresent,
			"present",
		},
		{
			"String() map enum to string",
			TargetStateAbsent,
			"absent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("TargetState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargetState_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		state   TargetState
		want    []byte
		wantErr bool
	}{
		{
			"MarshalJSON to bytes",
			TargetStatePresent,
			[]byte("\"present\""),
			false,
		},
		{
			"MarshalJSON to bytes",
			TargetStateAbsent,
			[]byte("\"absent\""),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.state.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TargetState.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TargetState.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTargetState_UnmarshalJSON_secondValue(t *testing.T) {
	jsonBlob := []byte("\"absent\"")
	var sut TargetState
	err := json.Unmarshal(jsonBlob, &sut)

	assert.Nil(t, err)
	assert.Equal(t, TargetState(TargetStateAbsent), sut)
}

func TestTargetState_UnmarshalJSON_firstValue(t *testing.T) {
	jsonBlob := []byte("\"present\"")
	var sut TargetState
	err := json.Unmarshal(jsonBlob, &sut)

	assert.Nil(t, err)
	assert.Equal(t, TargetState(TargetStatePresent), sut)
}

func TestTargetState_UnmarshalJSON_unknownValueParsesToFirstState(t *testing.T) {
	jsonBlob := []byte("\"test\"")
	var sut TargetState
	err := json.Unmarshal(jsonBlob, &sut)

	assert.Nil(t, err)
	assert.Equal(t, TargetState(TargetStatePresent), sut)
}

func TestTargetState_UnmarshalJSON_error(t *testing.T) {
	jsonBlob := []byte("test")
	var sut TargetState
	err := json.Unmarshal(jsonBlob, &sut)

	assert.NotNil(t, err)
}

func Test_validate_ok(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "absent/dogu1", Version: "3.2.1-0", TargetState: TargetStateAbsent},
		{Name: "absent/dogu2", TargetState: TargetStateAbsent},
		{Name: "present/dogu3", Version: "3.2.1-0", TargetState: TargetStatePresent},
		{Name: "present/dogu4", Version: "1.2.3-4"},
	}

	components := []Component{
		{Name: "absent/pkg1", Version: "3.2.1-0", TargetState: TargetStateAbsent},
		{Name: "absent/pkg2", TargetState: TargetStateAbsent},
		{Name: "present-pkg3", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Name: "present/pkg4", Version: "1.2.3-4"},
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
	assert.Contains(t, err.Error(), "could not validate blueprint, dogu field Name must not be empty")
}

func Test_validate_errorOnPackageProblem(t *testing.T) {
	components := []Component{
		{Name: "absent-component", TargetState: TargetStateAbsent},
		{Name: "present-component", Version: "3.2.1-2", TargetState: TargetStatePresent},
		{Version: "3.2.1-2"},
	}
	blueprint := BlueprintV2{Components: components}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not validate blueprint, component field Name must not be empty")
}
func Test_validateDogu_errorOnMissingDoguName(t *testing.T) {
	dogus := []TargetDogu{
		{Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not validate blueprint, dogu field Name must not be empty")
}

func Test_validateDogu_errorOnEmptyDoguName(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "", Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not validate blueprint, dogu field Name must not be empty")
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

func Test_validateComponents_errorOnMissingDoguName(t *testing.T) {
	components := []Component{
		{Name: "absent-component", TargetState: TargetStateAbsent},
		{Name: "present-component", Version: "3.2.1-2", TargetState: TargetStatePresent},
	}
	blueprint := BlueprintV2{Components: components}
	err := blueprint.validateComponents()

	require.Nil(t, err)
}

func Test_validateComponents_errorOnCesappInComponentsSection(t *testing.T) {
	components := []Component{
		{Name: "cesapp"},
	}
	blueprint := BlueprintV2{Components: components}
	err := blueprint.validateComponents()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not validate blueprint, "+
		"package cesapp does not belong into the Components section")
}

func Test_bluePrintValidator_validateDoguUniqueness(t *testing.T) {
	dogus := []TargetDogu{
		{Name: "present/dogu1", Version: "3.2.1-0", TargetState: TargetStatePresent},
		{Name: "present/dogu1", Version: "1.2.3-4"},
	}

	blueprint := BlueprintV2{Dogus: dogus}

	err := blueprint.validateDoguUniqueness()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not validate blueprint, there is at least one duplicate")
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

	assert.Contains(t, err.Error(), "could not validate blueprint, a config key is empty")
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
