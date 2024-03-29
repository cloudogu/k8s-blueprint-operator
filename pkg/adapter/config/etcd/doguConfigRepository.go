package etcd

import (
	"context"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type DoguConfigRepository struct {
	etcdStore etcdStore
}

func NewDoguConfigRepository(etcdStore etcdStore) *DoguConfigRepository {
	return &DoguConfigRepository{etcdStore: etcdStore}
}

func (e DoguConfigRepository) Get(_ context.Context, key common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error) {
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

func (e DoguConfigRepository) Save(_ context.Context, entry *ecosystem.DoguConfigEntry) error {
	strDoguName := string(entry.Key.DoguName)
	strValue := string(entry.Value)
	err := setEtcdKey(entry.Key.Key, strValue, e.etcdStore.DoguConfig(strDoguName))
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set %s with value %q in etcd", entry.Key, strValue)
	}

	return nil
}

func (e DoguConfigRepository) Delete(_ context.Context, key common.DoguConfigKey) error {
	strDoguName := string(key.DoguName)
	err := deleteEtcdKey(key.Key, e.etcdStore.DoguConfig(strDoguName))
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete %s from etcd", key)
	}

	return nil
}

func (e DoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error) {
	return getAllByKeyOrEntry(ctx, keys, e.Get)
}

func (e DoguConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.DoguConfigEntry) error {
	return mapKeyOrEntry(ctx, entries, e.Save, "failed to set given dogu config entries in etcd")
}

func (e DoguConfigRepository) DeleteAllByKeys(ctx context.Context, keys []common.DoguConfigKey) error {
	return mapKeyOrEntry(ctx, keys, e.Delete, "failed to delete given dogu config keys in etcd")
}
