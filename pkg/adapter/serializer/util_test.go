package serializer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitComponentName(t *testing.T) {
	tests := []struct {
		test          string
		given         string
		ns            string
		componentName string
		wantErr       assert.ErrorAssertionFunc
	}{
		{test: "ok", given: "k8s/my-component", ns: "k8s", componentName: "my-component", wantErr: assert.NoError},
		{test: "no ns", given: "my-component", ns: "", componentName: "", wantErr: assert.Error},
		{test: "no name", given: "k8s/", ns: "k8s", componentName: "", wantErr: assert.NoError},
		{test: "double namespace", given: "official/k8s/my-component", ns: "", componentName: "", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.componentName, func(t *testing.T) {
			got, got1, err := SplitComponentName(tt.given)
			if !tt.wantErr(t, err, fmt.Sprintf("SplitDoguName(%v)", tt.given)) {
				return
			}
			assert.Equalf(t, tt.ns, got, "SplitDoguName(%v)", tt.given)
			assert.Equalf(t, tt.componentName, got1, "SplitDoguName(%v)", tt.given)
		})
	}
}

func TestJoinComponentName(t *testing.T) {
	t.Run("should return error on missing dist namespace", func(t *testing.T) {
		_, err := JoinComponentName("my-component", "")

		require.Error(t, err)
		assert.ErrorContains(t, err, "distribution namespace of component my-component must not be empty")
	})
	t.Run("should return joined component name", func(t *testing.T) {
		actual, err := JoinComponentName("my-component", "k8s")

		// then
		require.NoError(t, err)
		assert.Equal(t, "k8s/my-component", actual)
	})
}
