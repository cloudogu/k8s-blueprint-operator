package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCombinedSecretEtcdSensitiveDoguConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		secretRepo := &SecretSensitiveDoguConfigRepository{}
		etcdRepo := &EtcdSensitiveDoguConfigRepository{}

		// when
		combinedRepo := NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo, secretRepo)

		// then
		require.NotNil(t, combinedRepo)
		assert.Equal(t, secretRepo, combinedRepo.SecretSensitiveDoguConfigRepository)
		assert.Equal(t, etcdRepo, combinedRepo.EtcdSensitiveDoguConfigRepository)
	})
}
