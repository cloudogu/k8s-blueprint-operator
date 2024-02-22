package config

import (
	"context"
	"errors"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type EtcdDoguConfigRepository struct {
	etcdStore etcdStore
}

func NewEtcdDoguConfigRepository(etcdStore etcdStore) *EtcdDoguConfigRepository {
	return &EtcdDoguConfigRepository{etcdStore: etcdStore}
}

func (e EtcdDoguConfigRepository) Get(_ context.Context, key common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error) {
	entry, err := e.etcdStore.DoguConfig(string(key.DoguName)).Get(key.Key)
	if registry.IsKeyNotFoundError(err) {
		return nil, domainservice.NewNotFoundError(err, "could not find %s in etcd", key)
	} else if err != nil {
		return nil, domainservice.NewInternalError(err, "failed to get %s from etcd", key)
	}

	return &ecosystem.DoguConfigEntry{
		Key:   key,
		Value: common.DoguConfigValue(entry),
	}, nil
}

func (e EtcdDoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error) {
	var errs []error
	entries := make(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry)
	for _, key := range keys {
		entry, err := e.Get(ctx, key)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		entries[key] = entry
	}

	return entries, errors.Join(errs...)
}

func (e EtcdDoguConfigRepository) Save(_ context.Context, entry *ecosystem.DoguConfigEntry) error {
	strDoguName := string(entry.Key.DoguName)
	strValue := string(entry.Value)
	err := setEtcdKey(entry.Key.Key, strValue, e.etcdStore.DoguConfig(strDoguName))
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set %s with value %q in etcd", entry.Key, strValue)
	}

	return nil
}

func (e EtcdDoguConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.DoguConfigEntry) error {
	var errs []error
	for _, entry := range entries {
		err := e.Save(ctx, entry)
		errs = append(errs, err)
	}

	err := errors.Join(errs...)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set given dogu config entries in etcd")
	}

	return nil
}

func (e EtcdDoguConfigRepository) Delete(_ context.Context, key common.DoguConfigKey) error {
	strDoguName := string(key.DoguName)
	err := deleteEtcdKey(key.Key, e.etcdStore.DoguConfig(strDoguName))
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete %s from etcd", key)
	}

	return nil
}
