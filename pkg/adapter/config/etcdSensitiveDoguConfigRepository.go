package config

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type EtcdSensitiveDoguConfigRepository struct {
	etcdStore etcdStore
}

func NewEtcdSensitiveDoguConfigRepository(etcdStore etcdStore) *EtcdSensitiveDoguConfigRepository {
	return &EtcdSensitiveDoguConfigRepository{etcdStore: etcdStore}
}

func (e EtcdSensitiveDoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (e EtcdSensitiveDoguConfigRepository) Save(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry) error {
	//TODO implement me
	panic("implement me")
}

func (e EtcdSensitiveDoguConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.SensitiveDoguConfigEntry) error {
	//TODO implement me
	panic("implement me")
}

func (e EtcdSensitiveDoguConfigRepository) Delete(ctx context.Context, key common.SensitiveDoguConfigKey) error {
	//TODO implement me
	panic("implement me")
}
