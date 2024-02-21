package config

import (
	"context"
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

func (e EtcdSensitiveDoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error) {
	// TODO implement me
	panic("implement me")
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

func (e EtcdSensitiveDoguConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.SensitiveDoguConfigEntry) error {
	// TODO implement me
	panic("implement me")
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
