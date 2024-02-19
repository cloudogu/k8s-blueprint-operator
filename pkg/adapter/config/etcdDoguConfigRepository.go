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

func (e EtcdDoguConfigRepository) GetAllByKey(_ context.Context, keys []common.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error) {
	var errs []error
	var entriesByDogu map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry
	for _, key := range keys {
		entryRaw, err := e.etcdStore.DoguConfig(string(key.DoguName)).Get(key.Key)
		if registry.IsKeyNotFoundError(err) {
			errs = append(errs, domainservice.NewNotFoundError(err, "could not find %q in etcd", key))
			continue
		} else if err != nil {
			errs = append(errs, domainservice.NewInternalError(err, "failed to get %q from etcd", key))
			continue
		}

		entriesByDogu[key.DoguName] = append(entriesByDogu[key.DoguName], &ecosystem.DoguConfigEntry{
			Key:   common.DoguConfigKey{DoguName: key.DoguName, Key: key.Key},
			Value: common.DoguConfigValue(entryRaw),
		})
	}

	return entriesByDogu, errors.Join(errs...)
}

func (e EtcdDoguConfigRepository) Save(_ context.Context, entry *ecosystem.DoguConfigEntry) error {
	strDoguName := string(entry.Key.DoguName)
	strKey := entry.Key.Key
	strValue := string(entry.Value)
	err := setEtcdKey(strKey, strValue, e.etcdStore.DoguConfig(strDoguName))
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set config key %q with value %q for dogu %q", strKey, strValue, strDoguName)
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
	strKey := key.Key
	err := deleteEtcdKey(strKey, e.etcdStore.DoguConfig(strDoguName))
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete config key %q for dogu %q", strKey, strDoguName)
	}

	return nil
}
