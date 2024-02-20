package etcd

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
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
		return fmt.Errorf("failed to set global config key %q with value %q: %w", strKey, strValue, err)
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
		return fmt.Errorf("failed to delete global config key %q: %w", strKey, err)
	}

	return nil
}
