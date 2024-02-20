package config

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/etcd"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/k8s"
)

type SecretEtcdSensitiveDoguConfigRepository struct {
	*etcd.EtcdSensitiveDoguConfigRepository
	*k8s.SecretSensitiveDoguConfigRepository
}

func NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo *etcd.EtcdSensitiveDoguConfigRepository, secretRepo *k8s.SecretSensitiveDoguConfigRepository) *SecretEtcdSensitiveDoguConfigRepository {
	return &SecretEtcdSensitiveDoguConfigRepository{
		EtcdSensitiveDoguConfigRepository:   etcdRepo,
		SecretSensitiveDoguConfigRepository: secretRepo,
	}
}
