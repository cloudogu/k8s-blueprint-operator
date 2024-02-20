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

func (e EtcdSensitiveDoguConfigRepository) GetAllByKey(_ context.Context, keys []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error) {
	var errs []error
	entriesByDogu := make(map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry)
	for _, key := range keys {
		entryRaw, err := e.etcdStore.DoguConfig(string(key.DoguName)).Get(key.Key)
		if registry.IsKeyNotFoundError(err) {
			errs = append(errs, domainservice.NewNotFoundError(err, "could not find %s in etcd", key))
			continue
		} else if err != nil {
			errs = append(errs, domainservice.NewInternalError(err, "failed to get %s from etcd", key))
			continue
		}

		entriesByDogu[key.DoguName] = append(entriesByDogu[key.DoguName], &ecosystem.SensitiveDoguConfigEntry{
			Key:   common.SensitiveDoguConfigKey{DoguConfigKey: common.DoguConfigKey{DoguName: key.DoguName, Key: key.Key}},
			Value: common.EncryptedDoguConfigValue(entryRaw),
		})
	}

	return entriesByDogu, errors.Join(errs...)
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

func (e EtcdSensitiveDoguConfigRepository) SaveAll(ctx context.Context, entries []*ecosystem.SensitiveDoguConfigEntry) error {
	var errs []error
	for _, entry := range entries {
		err := e.Save(ctx, entry)
		errs = append(errs, err)
	}

	err := errors.Join(errs...)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set given sensitive dogu config entries in etcd")
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

func setEtcdKey(key, value string, config configurationContext) error {
	return config.Set(key, value)
}

func deleteEtcdKey(key string, config configurationContext) error {
	return config.Delete(key)
}
