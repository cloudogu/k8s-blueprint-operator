package serializer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitDoguName(t *testing.T) {
	tests := []struct {
		test     string
		given    string
		ns       string
		doguName string
		wantErr  assert.ErrorAssertionFunc
	}{
		{test: "ok", given: "official/postgres", ns: "official", doguName: "postgres", wantErr: assert.NoError},
		{test: "no ns", given: "postgres", ns: "", doguName: "", wantErr: assert.Error},
		{test: "no name", given: "official/", ns: "official", doguName: "", wantErr: assert.NoError},
		{test: "double namespace", given: "official/test/postgres", ns: "", doguName: "", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.doguName, func(t *testing.T) {
			got, got1, err := SplitDoguName(tt.given)
			if !tt.wantErr(t, err, fmt.Sprintf("SplitDoguName(%v)", tt.given)) {
				return
			}
			assert.Equalf(t, tt.ns, got, "SplitDoguName(%v)", tt.given)
			assert.Equalf(t, tt.doguName, got1, "SplitDoguName(%v)", tt.given)
		})
	}
}
