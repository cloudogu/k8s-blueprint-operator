package config

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type EtcdDoguConfigRepository struct {
}

func (e EtcdDoguConfigRepository) GetAllByKey(ctx context.Context, keys []common.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (e EtcdDoguConfigRepository) Save(ctx context.Context, entry *ecosystem.DoguConfigEntry) error {
	//TODO implement me
	panic("implement me")
}

func (e EtcdDoguConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.DoguConfigEntry) error {
	//TODO implement me
	panic("implement me")
}

func (e EtcdDoguConfigRepository) Delete(ctx context.Context, key common.DoguConfigKey) error {
	//TODO implement me
	panic("implement me")
}
