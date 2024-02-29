package config

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/etcd"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCombinedSecretEtcdSensitiveDoguConfigRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		secretRepo := &kubernetes.SecretSensitiveDoguConfigRepository{}
		etcdRepo := &etcd.SensitiveDoguConfigRepository{}

		// when
		combinedRepo := NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo, secretRepo)

		// then
		require.NotNil(t, combinedRepo)
		assert.Equal(t, secretRepo, combinedRepo.SecretSensitiveDoguConfigRepository)
		assert.Equal(t, etcdRepo, combinedRepo.SensitiveDoguConfigRepository)
	})
}
