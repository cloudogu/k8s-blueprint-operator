package config

import (
	"context"
	"errors"

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
		return nil, domainservice.NewInternalError(err, "failed to get %q from global config in etcd", key)
	}

	return &ecosystem.GlobalConfigEntry{
		Key:   key,
		Value: common.GlobalConfigValue(entry),
	}, nil
}

func (e EtcdGlobalConfigRepository) GetAll(_ context.Context) ([]*ecosystem.GlobalConfigEntry, error) {
	globalConfigEntriesRaw, err := e.configStore.GetAll()
	if registry.IsKeyNotFoundError(err) {
		return nil, domainservice.NewNotFoundError(err, "could not find global config in etcd")
	} else if err != nil {
		return nil, domainservice.NewInternalError(err, "failed to get global config from etcd")
	}

	var globalConfigEntries []*ecosystem.GlobalConfigEntry
	for key, value := range globalConfigEntriesRaw {
		globalConfigEntries = append(globalConfigEntries, &ecosystem.GlobalConfigEntry{
			Key:   common.GlobalConfigKey(key),
			Value: common.GlobalConfigValue(value),
		})
	}

	return globalConfigEntries, nil
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

func (e EtcdGlobalConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.GlobalConfigEntry) error {
	var errs []error
	for _, entry := range entries {
		err := e.Save(ctx, entry)
		errs = append(errs, err)
	}

	err := errors.Join(errs...)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set given global config entries in etcd")
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
