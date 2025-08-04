package serializer

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func Test_convertMinimumVolumeSizeToDTO(t *testing.T) {
	tests := []struct {
		name       string
		minVolSize ecosystem.VolumeSize
		want       string
	}{
		{
			name:       "empty",
			minVolSize: ecosystem.VolumeSize{},
			want:       "",
		},
		{
			name:       "empty",
			minVolSize: resource.MustParse("1Gi"),
			want:       "1Gi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, convertMinimumVolumeSizeToDTO(tt.minVolSize), "convertMinimumVolumeSizeToDTO(%v)", tt.minVolSize)
		})
	}
}
