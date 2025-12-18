package serializer

import (
	"testing"

	crd "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_convertToDoguDiffStateDTO(t *testing.T) {
	t.Run("should convert empty reverse proxy config to nil", func(t *testing.T) {
		// given
		domainDiffState := domain.DoguDiffState{
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{},
		}
		// when
		result := convertToDoguDiffStateDTO(domainDiffState)
		// then
		assert.Nil(t, result.ReverseProxyConfig)
	})

	t.Run("should convert reverse proxy config", func(t *testing.T) {
		// given
		domainDiffState := domain.DoguDiffState{
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySize,
				RewriteTarget:    ecosystem.RewriteTarget(rewriteTarget),
				AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig),
			},
		}
		// when
		result := convertToDoguDiffStateDTO(domainDiffState)
		// then
		want := crd.DoguDiffState{
			ReverseProxyConfig: &crd.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySizeString,
				RewriteTarget:    &rewriteTarget,
				AdditionalConfig: &additionalConfig,
			},
		}

		assert.NotNil(t, result.ReverseProxyConfig)
		assert.Empty(t, cmp.Diff(want, result))
	})

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
