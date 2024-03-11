package etcd

import (
	"context"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type GlobalConfigRepository struct {
	configStore globalConfigStore
}

func NewGlobalConfigRepository(configStore globalConfigStore) *GlobalConfigRepository {
	return &GlobalConfigRepository{configStore: configStore}
}

func (e GlobalConfigRepository) Get(_ context.Context, key common.GlobalConfigKey) (*ecosystem.GlobalConfigEntry, error) {
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

func (e GlobalConfigRepository) Save(_ context.Context, entry *ecosystem.GlobalConfigEntry) error {
	strKey := string(entry.Key)
	strValue := string(entry.Value)
	err := e.configStore.Set(strKey, strValue)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set global config key %q with value %q in etcd", strKey, strValue)
	}

	return nil
}

func (e GlobalConfigRepository) Delete(_ context.Context, key common.GlobalConfigKey) error {
	strKey := string(key)
	err := e.configStore.Delete(strKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete global config key %q from etcd", strKey)
	}

	return nil
}

func (e GlobalConfigRepository) GetAllByKey(ctx context.Context, keys []common.GlobalConfigKey) (map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry, error) {
	return getAllByKeyOrEntry(ctx, keys, e.Get)
}

func (e GlobalConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.GlobalConfigEntry) error {
	return mapKeyOrEntry(ctx, entries, e.Save, "failed to set given global config entries in etcd")
}

func (e GlobalConfigRepository) DeleteAllByKeys(ctx context.Context, keys []common.GlobalConfigKey) error {
	return mapKeyOrEntry(ctx, keys, e.Delete, "failed to delete given global config keys in etcd")
}
