package serializer

import (
	"testing"

	crd "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_convertToDoguDiffStateDTO(t *testing.T) {
	t.Run("should convert resource config", func(t *testing.T) {
		// given
		storageClassName := "example-storage-class"
		domainDiffState := domain.DoguDiffState{
			MinVolumeSize:    &volumeSize,
			StorageClassName: &storageClassName,
		}
		// when
		result := convertToDoguDiffStateDTO(domainDiffState)
		// then
		want := crd.DoguDiffState{
			ResourceConfig: &crd.ResourceConfig{
				MinVolumeSize:    &volumeSizeString,
				StorageClassName: &storageClassName,
			},
		}

		assert.NotNil(t, result.ResourceConfig)
		assert.Empty(t, cmp.Diff(want, result))
	})
}
