package config

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type EtcdGlobalConfigRepository struct {
}

func (e EtcdGlobalConfigRepository) GetAll(ctx context.Context) ([]*ecosystem.GlobalConfigEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) Save(ctx context.Context, entry *ecosystem.GlobalConfigEntry) error {
	//TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) SaveAll(ctx context.Context, keys []*ecosystem.GlobalConfigEntry) error {
	//TODO implement me
	panic("implement me")
}

func (e EtcdGlobalConfigRepository) Delete(ctx context.Context, key common.SimpleDoguName) error {
	//TODO implement me
	panic("implement me")
}
