package config

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type EtcdGlobalConfigRepository struct {
	configStore globalConfigStore
}

func NewEtcdGlobalConfigRepository(configStore globalConfigStore) *EtcdGlobalConfigRepository {
	return &EtcdGlobalConfigRepository{configStore: configStore}
}

func (e EtcdGlobalConfigRepository) GetAll(ctx context.Context) ([]*ecosystem.GlobalConfigEntry, error) {
	// TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) Save(ctx context.Context, entry *ecosystem.GlobalConfigEntry) error {
	// TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.GlobalConfigEntry) error {
	// TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) Delete(ctx context.Context, key common.GlobalConfigKey) error {
	// TODO implement me
	panic("implement me")
}
