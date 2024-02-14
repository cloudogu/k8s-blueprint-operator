package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQualifiedComponentNameFromString(t *testing.T) {
	tests := []struct {
		test     string
		given    string
		expected QualifiedComponentName
		wantErr  assert.ErrorAssertionFunc
	}{
		{test: "ok", given: "k8s/k8s-dogu-operator", expected: QualifiedComponentName{ComponentNamespace("k8s"), SimpleComponentName("k8s-dogu-operator")}, wantErr: assert.NoError},
		{test: "no ns", given: "k8s-dogu-operator", expected: QualifiedComponentName{}, wantErr: assert.Error},
		{test: "no name", given: "k8s/", expected: QualifiedComponentName{}, wantErr: assert.Error},
		{test: "double namespace", given: "k8s/test/k8s-dogu-operator", expected: QualifiedComponentName{}, wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			got, err := QualifiedComponentNameFromString(tt.given)
			if !tt.wantErr(t, err, fmt.Sprintf("TestQualifiedComponentNameFromString(%v)", tt.given)) {
				return
			}
			assert.Equalf(t, tt.expected, got, "TestQualifiedComponentNameFromString(%v)", tt.given)
		})
	}
}
