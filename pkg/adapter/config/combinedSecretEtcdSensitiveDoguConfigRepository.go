package config

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/etcd"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/config/kubernetes"
)

type SecretEtcdSensitiveDoguConfigRepository struct {
	*etcd.SensitiveDoguConfigRepository
	*kubernetes.SecretSensitiveDoguConfigRepository
}

func NewCombinedSecretEtcdSensitiveDoguConfigRepository(etcdRepo *etcd.SensitiveDoguConfigRepository, secretRepo *kubernetes.SecretSensitiveDoguConfigRepository) *SecretEtcdSensitiveDoguConfigRepository {
	return &SecretEtcdSensitiveDoguConfigRepository{
		SensitiveDoguConfigRepository:       etcdRepo,
		SecretSensitiveDoguConfigRepository: secretRepo,
	}
}
