package serializer

import (
	"testing"

	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_convertDoguDiffStateDomain(t *testing.T) {
	t.Run("should convert empty reverse proxy config", func(t *testing.T) {
		// given
		crdDiffState := crd.DoguDiffState{
			ReverseProxyConfig: &crd.ReverseProxyConfig{},
		}
		// when
		result, err := convertDoguDiffStateDomain(crdDiffState)
		// then
		require.NoError(t, err)
		assert.NotNil(t, result.ReverseProxyConfig)
		assert.Nil(t, result.ReverseProxyConfig.RewriteTarget)
		assert.Nil(t, result.ReverseProxyConfig.AdditionalConfig)
		assert.Nil(t, result.ReverseProxyConfig.MaxBodySize)
	})

	t.Run("should convert reverse proxy config", func(t *testing.T) {
		// given
		crdDiffState := crd.DoguDiffState{
			ReverseProxyConfig: &crd.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySizeString,
				RewriteTarget:    &rewriteTarget,
				AdditionalConfig: &additionalConfig,
			},
		}
		// when
		result, err := convertDoguDiffStateDomain(crdDiffState)
		// then
		want := domain.DoguDiffState{
			ReverseProxyConfig: &ecosystem.ReverseProxyConfig{
				MaxBodySize:      &proxyBodySize,
				RewriteTarget:    &rewriteTarget,
				AdditionalConfig: &additionalConfig,
			},
		}

		require.NoError(t, err)
		assert.NotNil(t, result.ReverseProxyConfig)
		assert.Empty(t, cmp.Diff(want, result))
	})

	t.Run("should throw error on reverse proxy config convert error", func(t *testing.T) {
		// given
		wrongBodySize := "1Z"
		crdDiffState := crd.DoguDiffState{
			ReverseProxyConfig: &crd.ReverseProxyConfig{
				MaxBodySize: &wrongBodySize,
			},
		}
		// when
		result, err := convertDoguDiffStateDomain(crdDiffState)

		// then
		require.Error(t, err)
		assert.Equal(t, domain.DoguDiffState{}, result)
		assert.ErrorContains(t, err, "failed to parse maximum proxy body size")
	})

	t.Run("should convert resource config", func(t *testing.T) {
		// given
		crdDiffState := crd.DoguDiffState{
			ResourceConfig: &crd.ResourceConfig{
				MinVolumeSize: &volumeSizeString,
			},
		}
		// when
		result, err := convertDoguDiffStateDomain(crdDiffState)

		// then
		want := domain.DoguDiffState{
			MinVolumeSize: &volumeSize,
		}
		require.NoError(t, err)
		assert.NotNil(t, result.MinVolumeSize)
		assert.Empty(t, cmp.Diff(want, result))
	})

	t.Run("should return error on resource config convert error", func(t *testing.T) {
		// given
		wrongVolumeSize := "1Gu"
		crdDiffState := crd.DoguDiffState{
			ResourceConfig: &crd.ResourceConfig{
				MinVolumeSize: &wrongVolumeSize,
			},
		}
		// when
		result, err := convertDoguDiffStateDomain(crdDiffState)

		// then
		require.Error(t, err)
		assert.Equal(t, domain.DoguDiffState{}, result)
		assert.ErrorContains(t, err, "failed to parse minimum volume size")
	})
}
