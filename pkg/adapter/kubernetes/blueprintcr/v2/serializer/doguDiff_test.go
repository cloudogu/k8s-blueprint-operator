package serializer

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

func Test_convertMinimumVolumeSizeToDTO(t *testing.T) {
	volumeSize1g := resource.MustParse("1Gi")
	val1Gi := "1Gi"
	tests := []struct {
		name       string
		minVolSize *ecosystem.VolumeSize
		want       *string
	}{
		{
			name:       "nil",
			minVolSize: nil,
			want:       nil,
		},
		{
			name:       "empty",
			minVolSize: &ecosystem.VolumeSize{},
			want:       nil,
		},
		{
			name:       "1Gi",
			minVolSize: &volumeSize1g,
			want:       &val1Gi,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertMinimumVolumeSizeToDTO(tt.minVolSize), "convertMinimumVolumeSizeToDTO(%v)", tt.minVolSize)
		})
	}
}
