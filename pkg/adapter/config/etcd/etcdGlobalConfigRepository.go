package etcd

import (
	"context"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type EtcdGlobalConfigRepository struct {
	configStore globalConfigStore
}

func NewEtcdGlobalConfigRepository(configStore globalConfigStore) *EtcdGlobalConfigRepository {
	return &EtcdGlobalConfigRepository{configStore: configStore}
}

func (e EtcdGlobalConfigRepository) Get(_ context.Context, key common.GlobalConfigKey) (*ecosystem.GlobalConfigEntry, error) {
	entry, err := e.configStore.Get(string(key))
	if registry.IsKeyNotFoundError(err) {
		return nil, domainservice.NewNotFoundError(err, "could not find key %q from global config in etcd", key)
	} else if err != nil {
		return nil, domainservice.NewInternalError(err, "failed to get value for key %q from global config in etcd", key)
	}

	return &ecosystem.GlobalConfigEntry{
		Key:   key,
		Value: common.GlobalConfigValue(entry),
	}, nil
}

func (e EtcdGlobalConfigRepository) Save(_ context.Context, entry *ecosystem.GlobalConfigEntry) error {
	strKey := string(entry.Key)
	strValue := string(entry.Value)
	err := e.configStore.Set(strKey, strValue)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set global config key %q with value %q in etcd", strKey, strValue)
	}

	return nil
}

func (e EtcdGlobalConfigRepository) Delete(_ context.Context, key common.GlobalConfigKey) error {
	strKey := string(key)
	err := e.configStore.Delete(strKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete global config key %q from etcd", strKey)
	}

	return nil
}

func (e EtcdGlobalConfigRepository) GetAllByKey(ctx context.Context, keys []common.GlobalConfigKey) (map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry, error) {
	return getAllByKey(ctx, keys, e.Get)
}

func (e EtcdGlobalConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.GlobalConfigEntry) error {
	return saveOrDeleteAllByRegistryKeys(ctx, entries, e.Save, "failed to set given global config entries in etcd")
}

func (e EtcdGlobalConfigRepository) DeleteAllByKeys(ctx context.Context, keys []common.GlobalConfigKey) error {
	return saveOrDeleteAllByRegistryKeys(ctx, keys, e.Delete, "failed to delete given global config keys in etcd")
}
