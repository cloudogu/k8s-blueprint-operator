package serializer

import (
	"testing"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_convertToDoguDiffStateDTO(t *testing.T) {
	t.Run("should convert empty reverse proxy config", func(t *testing.T) {
		// given
		domainDiffState := domain.DoguDiffState{
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{},
		}
		// when
		result := convertToDoguDiffStateDTO(domainDiffState)
		// then
		assert.NotNil(t, result.ReverseProxyConfig)
		assert.Nil(t, result.ReverseProxyConfig.RewriteTarget)
		assert.Nil(t, result.ReverseProxyConfig.AdditionalConfig)
		assert.Nil(t, result.ReverseProxyConfig.MaxBodySize)
	})

	t.Run("should convert reverse proxy config", func(t *testing.T) {
		// given
		domainDiffState := domain.DoguDiffState{
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySize,
				RewriteTarget:    &rewriteTarget,
				AdditionalConfig: &additionalConfig,
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
		domainDiffState := domain.DoguDiffState{
			MinVolumeSize: &volumeSize,
		}
		// when
		result := convertToDoguDiffStateDTO(domainDiffState)
		// then
		want := crd.DoguDiffState{
			ResourceConfig: &crd.ResourceConfig{
				MinVolumeSize: &volumeSizeString,
			},
		}

		assert.NotNil(t, result.ResourceConfig)
		assert.Empty(t, cmp.Diff(want, result))
	})
}
