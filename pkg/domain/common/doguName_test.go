package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQualifiedDoguNameFromString(t *testing.T) {
	tests := []struct {
		test     string
		given    string
		expected QualifiedDoguName
		wantErr  assert.ErrorAssertionFunc
	}{
		{test: "ok", given: "official/postgres", expected: QualifiedDoguName{DoguNamespace("official"), SimpleDoguName("postgres")}, wantErr: assert.NoError},
		{test: "no ns", given: "postgres", expected: QualifiedDoguName{}, wantErr: assert.Error},
		{test: "no name", given: "official/", expected: QualifiedDoguName{}, wantErr: assert.Error},
		{test: "double namespace", given: "official/test/postgres", expected: QualifiedDoguName{}, wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			got, err := QualifiedDoguNameFromString(tt.given)
			if !tt.wantErr(t, err, fmt.Sprintf("TestQualifiedDoguNameFromString(%v)", tt.given)) {
				return
			}
			assert.Equalf(t, tt.expected, got, "TestQualifiedDoguNameFromString(%v)", tt.given)
		})
	}
}
