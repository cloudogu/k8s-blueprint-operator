package kubernetes

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

type DoguConfigRepository struct {
	repo repository.DoguConfigRepository
}

func NewDoguConfigRepository(repo repository.DoguConfigRepository) *DoguConfigRepository {
	return &DoguConfigRepository{repo: repo}
}

func (e DoguConfigRepository) Get(ctx context.Context, doguName common.SimpleDoguName) (config.DoguConfig, error) {
	// TODO: There seems to be no way to know, if we have a NotFoundError or a connection error.
	return e.repo.Get(ctx, doguName)
}

func (e DoguConfigRepository) Update(ctx context.Context, config config.DoguConfig) (config.DoguConfig, error) {
	updatedConfig, err := e.repo.Update(ctx, config)
	// TODO: we cannot see here, if there is a real conflict or there was a connection error.
	//  With a conflict, we can immediately restart the business process
	//  With an connection error we need a longer backoff (internalError)
	if err != nil {
		return config, domainservice.NewInternalError(err, "failed to update config for %s", config.DoguName)
	}
	return updatedConfig, nil
}
