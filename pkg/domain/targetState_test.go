package domain

import (
	"encoding/json"
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
