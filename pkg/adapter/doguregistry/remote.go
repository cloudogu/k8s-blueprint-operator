package doguregistry

import (
	"context"
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

var maxTries = 20

type Remote struct {
	repository remoteDoguDescriptorRepository
}

func NewRemote(repository remoteDoguDescriptorRepository) *Remote {
	return &Remote{repository: repository}
}

func (r *Remote) GetDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedDoguVersion) (*core.Dogu, error) {
	dogu := &core.Dogu{}
	err := retry.OnError(maxTries, isConnectionError, func() error {
		var err error
		dogu, err = r.repository.Get(ctx, qualifiedDoguVersion)
		return err
	})
	if err != nil {
		// this is ugly, maybe do it better in cesapp-lib?
		if strings.Contains(err.Error(), cescommons.DoguDescriptorNotFoundError.Error()) {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("dogu %q with version %q could not be found", qualifiedDoguVersion.Name, qualifiedDoguVersion.Version.Raw),
			}
		}

		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("failed to get dogu %q with version %q", qualifiedDoguVersion.Name, qualifiedDoguVersion.Version.Raw),
		}
	}

	return dogu, nil
}

func (r *Remote) GetDogus(ctx context.Context, dogusToLoad []cescommons.QualifiedDoguVersion) (map[cescommons.QualifiedDoguName]*core.Dogu, error) {
	dogus := make(map[cescommons.QualifiedDoguName]*core.Dogu)

	var errs []error
	for _, doguRef := range dogusToLoad {
		dogu, err := r.GetDogu(ctx, doguRef)
		errs = append(errs, err)

		dogus[doguRef.Name] = dogu
	}

	return dogus, errors.Join(errs...)
}

func isConnectionError(err error) bool {
	return strings.Contains(err.Error(), cescommons.ConnectionError.Error())
}
