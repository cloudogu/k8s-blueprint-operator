package adapter

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/etcd"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/config"
)

type SecretEtcdSensitiveDoguConfigRepository struct {
	*etcd.EtcdSensitiveDoguConfigRepository
	*config.SecretSensitiveDoguConfigRepository
}

func NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo *etcd.EtcdSensitiveDoguConfigRepository, secretRepo *config.SecretSensitiveDoguConfigRepository) *SecretEtcdSensitiveDoguConfigRepository {
	return &SecretEtcdSensitiveDoguConfigRepository{
		EtcdSensitiveDoguConfigRepository:   etcdRepo,
		SecretSensitiveDoguConfigRepository: secretRepo,
	}
}
