package config

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

func (e EtcdGlobalConfigRepository) GetAll(ctx context.Context) ([]*ecosystem.GlobalConfigEntry, error) {
	// TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) Save(_ context.Context, entry *ecosystem.GlobalConfigEntry) error {
	strKey := string(entry.Key)
	strValue := string(entry.Value)
	err := e.configStore.Set(strKey, strValue)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to set global config key %q with value %q", strKey, strValue)
	}

	return nil
}

func (e EtcdGlobalConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.GlobalConfigEntry) error {
	// TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) Delete(ctx context.Context, key common.GlobalConfigKey) error {
	strKey := string(key)
	err := e.configStore.Delete(strKey)
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return domainservice.NewInternalError(err, "failed to delete global config key %q", strKey)
	}

	return nil
}
