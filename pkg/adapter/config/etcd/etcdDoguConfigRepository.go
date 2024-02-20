package etcd

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type EtcdDoguConfigRepository struct {
	etcdStore etcdStore
}

func NewEtcdDoguConfigRepository(etcdStore etcdStore) *EtcdDoguConfigRepository {
	return &EtcdDoguConfigRepository{etcdStore: etcdStore}
}

func (e EtcdDoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error) {
	// TODO implement me
	panic("implement me")
}

func (e EtcdDoguConfigRepository) Save(_ context.Context, entry *ecosystem.DoguConfigEntry) error {
	strDoguName := string(entry.Key.DoguName)
	strKey := entry.Key.Key
	strValue := string(entry.Value)
	err := setEtcdKey(strKey, strValue, e.etcdStore.DoguConfig(strDoguName))
	if err != nil {
		return fmt.Errorf("failed to set config key %q with value %q for dogu %q: %w", strKey, strValue, strDoguName, err)
	}

	return nil
}

func (e EtcdDoguConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.DoguConfigEntry) error {
	// TODO implement me
	panic("implement me")
}

func (e EtcdDoguConfigRepository) Delete(_ context.Context, key common.DoguConfigKey) error {
	strDoguName := string(key.DoguName)
	strKey := key.Key
	err := deleteEtcdKey(strKey, e.etcdStore.DoguConfig(strDoguName))
	if err != nil && !registry.IsKeyNotFoundError(err) {
		return fmt.Errorf("failed to delete config key %q for dogu %q: %w", strKey, strDoguName, err)
	}

	return nil
}
