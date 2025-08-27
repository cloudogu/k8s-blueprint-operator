package doguregistry

import (
	"context"
	"errors"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	cloudoguerrors "github.com/cloudogu/ces-commons-lib/errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

type Remote struct {
	repository remoteDoguDescriptorRepository
}

func NewRemote(repository remoteDoguDescriptorRepository) *Remote {
	return &Remote{repository: repository}
}

func (r *Remote) GetDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedVersion) (*core.Dogu, error) {
	// do not retry here. If any error happens, just reconcile later. We only do retries in application level.
	// This makes the code way easier and non-blocking.
	dogu, err := r.repository.Get(ctx, qualifiedDoguVersion)
	if err != nil {
		if cloudoguerrors.IsNotFoundError(err) {
			return nil, domainservice.NewNotFoundError(
				err,
				"dogu %q with version %q could not be found",
				qualifiedDoguVersion.Name, qualifiedDoguVersion.Version.Raw,
			)
		} else {
			return nil, domainservice.NewInternalError(
				err,
				"failed to get dogu %q with version %q",
				qualifiedDoguVersion.Name, qualifiedDoguVersion.Version.Raw,
			)
		}
	}
	return dogu, nil
}

func (r *Remote) GetDogus(ctx context.Context, dogusToLoad []cescommons.QualifiedVersion) (map[cescommons.QualifiedName]*core.Dogu, error) {
	dogus := make(map[cescommons.QualifiedName]*core.Dogu)

	var errs []error
	for _, doguRef := range dogusToLoad {
		dogu, err := r.GetDogu(ctx, doguRef)
		errs = append(errs, err)

		dogus[doguRef.Name] = dogu
	}

	return dogus, errors.Join(errs...)
}
