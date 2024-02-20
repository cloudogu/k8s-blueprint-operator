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
	var entriesByDogu map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry
	for _, key := range keys {
		entryRaw, err := e.etcdStore.DoguConfig(string(key.DoguName)).Get(key.Key)
		if registry.IsKeyNotFoundError(err) {
			errs = append(errs, domainservice.NewNotFoundError(err, "could not find %q in etcd", key))
			continue
		} else if err != nil {
			errs = append(errs, domainservice.NewInternalError(err, "failed to get %q from etcd", key))
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
	strKey := entry.Key.Key
	strValue := string(entry.Value)
	err := setEtcdKey(strKey, strValue, e.etcdStore.DoguConfig(strDoguName))
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set encrypted config key %q with value %q for dogu %q", strKey, strValue, strDoguName)
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
		return domainservice.NewInternalError(err, "failed to set given dogu config entries in etcd")
	}

	return nil
}

func (e EtcdSensitiveDoguConfigRepository) Delete(_ context.Context, key common.SensitiveDoguConfigKey) error {
	strDoguName := string(key.DoguName)
	strKey := key.Key
	err := deleteEtcdKey(strKey, e.etcdStore.DoguConfig(strDoguName))
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete encrypted config key %q for dogu %q", strKey, strDoguName)
	}

	return nil
}

func setEtcdKey(key, value string, config configurationContext) error {
	return config.Set(key, value)
}

func deleteEtcdKey(key string, config configurationContext) error {
	return config.Delete(key)
}
