package serializer

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestToDomainProxyBodySize(t *testing.T) {
	t.Run("should return nil on empty string", func(t *testing.T) {
		// when
		size, err := ToDomainProxyBodySize("")

		// then
		require.NoError(t, err)
		assert.Nil(t, size)
	})

	t.Run("should return zero size quantity on '0'", func(t *testing.T) {
		// when
		size, err := ToDomainProxyBodySize("0")

		// then
		require.NoError(t, err)
		assert.Equal(t, *size, resource.MustParse("0"))
	})

	t.Run("should return normal quantity", func(t *testing.T) {
		// when
		size, err := ToDomainProxyBodySize("100M")

		// then
		require.NoError(t, err)
		assert.Equal(t, *size, resource.MustParse("100M"))
	})
}
