package config

import (
	"context"
	"errors"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type EtcdSensitiveDoguConfigRepository struct {
	etcdStore etcdStore
}

func NewEtcdSensitiveDoguConfigRepository(etcdStore etcdStore) *EtcdSensitiveDoguConfigRepository {
	return &EtcdSensitiveDoguConfigRepository{etcdStore: etcdStore}
}

func (e EtcdSensitiveDoguConfigRepository) Get(_ context.Context, key common.SensitiveDoguConfigKey) (*ecosystem.SensitiveDoguConfigEntry, error) {
	entry, err := e.etcdStore.DoguConfig(string(key.DoguName)).Get(key.Key)
	if registry.IsKeyNotFoundError(err) {
		return nil, domainservice.NewNotFoundError(err, "could not find sensitive %s in etcd", key)
	} else if err != nil {
		return nil, domainservice.NewInternalError(err, "failed to get sensitive %s from etcd", key)
	}

	return &ecosystem.SensitiveDoguConfigEntry{
		Key:   key,
		Value: common.EncryptedDoguConfigValue(entry),
	}, nil
}

func (e EtcdSensitiveDoguConfigRepository) Save(_ context.Context, entry *ecosystem.SensitiveDoguConfigEntry) error {
	strDoguName := string(entry.Key.DoguName)
	strValue := string(entry.Value)
	err := setEtcdKey(entry.Key.Key, strValue, e.etcdStore.DoguConfig(strDoguName))
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set encrypted %s with value %q in etcd", entry.Key, strValue)
	}

	return nil
}

func (e EtcdSensitiveDoguConfigRepository) Delete(_ context.Context, key common.SensitiveDoguConfigKey) error {
	strDoguName := string(key.DoguName)
	err := deleteEtcdKey(key.Key, e.etcdStore.DoguConfig(strDoguName))
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete encrypted %s from etcd", key)
	}

	return nil
}

func (e EtcdSensitiveDoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.SensitiveDoguConfigKey) (map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, error) {
	return getAllByKeyOrEntry(ctx, keys, e.Get)
}

func (e EtcdSensitiveDoguConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.SensitiveDoguConfigEntry) error {
	return mapKeyOrEntry(ctx, entries, e.Save, "failed to set given sensitive dogu config entries in etcd")
}

func (e EtcdSensitiveDoguConfigRepository) DeleteAllByKeys(ctx context.Context, keys []common.SensitiveDoguConfigKey) error {
	return mapKeyOrEntry(ctx, keys, e.Delete, "failed to delete given sensitive dogu config keys in etcd")
}

type registryKey interface {
	common.SensitiveDoguConfigKey | common.GlobalConfigKey | common.DoguConfigKey
}

type registryEntry interface {
	*ecosystem.DoguConfigEntry | *ecosystem.GlobalConfigEntry | *ecosystem.SensitiveDoguConfigEntry
}

func getAllByKeyOrEntry[T registryKey, K registryEntry](ctx context.Context, collection []T, fn func(context.Context, T) (K, error)) (map[T]K, error) {
	var errs []error
	entries := make(map[T]K)
	for _, key := range collection {
		entry, err := fn(ctx, key)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		entries[key] = entry
	}

	return entries, errors.Join(errs...)
}

func mapKeyOrEntry[T registryKey | registryEntry](ctx context.Context, collection []T, fn func(context.Context, T) error, errorMsg string) error {
	var errs []error
	for _, key := range collection {
		err := fn(ctx, key)
		errs = append(errs, err)
	}

	err := errors.Join(errs...)
	if err != nil {
		return domainservice.NewInternalError(err, errorMsg)
	}

	return nil
}

func setEtcdKey(key, value string, config configurationContext) error {
	return config.Set(key, value)
}

func deleteEtcdKey(key string, config configurationContext) error {
	return config.Delete(key)
}
