package kubernetes

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

type GlobalConfigRepository struct {
	repo repository.GlobalConfigRepository
}

func NewGlobalConfigRepository(repo repository.GlobalConfigRepository) *GlobalConfigRepository {
	return &GlobalConfigRepository{repo: repo}
}

func (e GlobalConfigRepository) Get(ctx context.Context) (config.GlobalConfig, error) {
	return e.repo.Get(ctx)
	//TODO: add own error types again
	//if registry.IsKeyNotFoundError(err) {
	//	return nil, domainservice.NewNotFoundError(err, "could not find key %q from global config in etcd", key)
	//} else if err != nil {
	//	return nil, domainservice.NewInternalError(err, "failed to get value for key %q from global config in etcd", key)
	//}
}

func (e GlobalConfigRepository) Update(ctx context.Context, config config.GlobalConfig) (config.GlobalConfig, error) {
	updatedConfig, err := e.repo.Update(ctx, config)
	// TODO: we cannot see here, if there is a real conflict or there was a connection error.
	//  With a conflict, we can immediately restart the business process
	//  With an connection error we need a longer backoff (internalError)
	if err != nil {
		return config, domainservice.NewInternalError(err, "failed to update global config")
	}
	return updatedConfig, nil
}
