package adapter

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/etcd"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/config"
)

type SecretEtcdSensitiveDoguConfigRepository struct {
	*etcd.SensitiveDoguConfigRepository
	*config.SecretSensitiveDoguConfigRepository
}

func NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo *etcd.SensitiveDoguConfigRepository, secretRepo *config.SecretSensitiveDoguConfigRepository) *SecretEtcdSensitiveDoguConfigRepository {
	return &SecretEtcdSensitiveDoguConfigRepository{
		SensitiveDoguConfigRepository:       etcdRepo,
		SecretSensitiveDoguConfigRepository: secretRepo,
	}
}
