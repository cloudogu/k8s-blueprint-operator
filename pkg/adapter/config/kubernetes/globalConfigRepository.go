package kubernetes

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/cloudogu/k8s-registry-lib/config"
	liberrors "github.com/cloudogu/k8s-registry-lib/errors"
)

type GlobalConfigRepository struct {
	repo k8sGlobalConfigRepo
}

func NewGlobalConfigRepository(repo k8sGlobalConfigRepo) *GlobalConfigRepository {
	return &GlobalConfigRepository{repo: repo}
}

func (e GlobalConfigRepository) Get(ctx context.Context) (config.GlobalConfig, error) {
	loadedConfig, err := e.repo.Get(ctx)
	if err != nil {
		if liberrors.IsNotFoundError(err) {
			return loadedConfig, domainservice.NewNotFoundError(err, "could not find global config. Check if your ecosystem is ready for operation")
		} else if liberrors.IsConnectionError(err) {
			return loadedConfig, domainservice.NewInternalError(err, "could not load global config due to connection problems")
		} else {
			// GenericError and fallback if even that would not match the error
			return loadedConfig, domainservice.NewInternalError(err, "could not load global config due to an unknown problem")
		}
	}
	return loadedConfig, nil
}

func (e GlobalConfigRepository) Update(ctx context.Context, config config.GlobalConfig) (config.GlobalConfig, error) {
	updatedConfig, err := e.repo.Update(ctx, config)
	if err != nil {
		if liberrors.IsNotFoundError(err) {
			return updatedConfig, domainservice.NewNotFoundError(err, "could not update global config. Check if your ecosystem is ready for operation")
		} else if liberrors.IsConnectionError(err) {
			return updatedConfig, domainservice.NewInternalError(err, "could not update global config due to connection problems")
		} else if liberrors.IsConflictError(err) {
			return updatedConfig, domainservice.NewInternalError(err, "could not update global config due to conflicting changes")
		} else {
			// GenericError and fallback if even that would not match the error
			return updatedConfig, domainservice.NewInternalError(err, "could not update global config due to an unknown problem")
		}
	}
	return updatedConfig, nil
}
