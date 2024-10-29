package doguregistry

import (
	"context"
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/retry"
	remotedogudescriptor "github.com/cloudogu/remote-dogu-descriptor-lib/repository"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

var maxTries = 20

type Remote struct {
	repository cescommons.RemoteDoguDescriptorRepository
}

func NewRemote(repository cescommons.RemoteDoguDescriptorRepository) *Remote {
	return &Remote{repository: repository}
}

func (r *Remote) GetDogu(qualifiedDoguVersion cescommons.QualifiedDoguVersion) (*core.Dogu, error) {
	dogu := &core.Dogu{}
	err := retry.OnError(maxTries, isConnectionError, func() error {
		var err error
		dogu, err = r.repository.Get(context.TODO(), qualifiedDoguVersion) // TODO context
		return err
	})
	if err != nil {
		// this is ugly, maybe do it better in cesapp-lib?
		if strings.Contains(err.Error(), "404 not found") {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("dogu %q with version %q could not be found", qualifiedDoguVersion.Name.SimpleName, qualifiedDoguVersion.Version),
			}
		}

		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("failed to get dogu %q with version %q", qualifiedDoguVersion.Name.SimpleName, qualifiedDoguVersion.Version),
		}
	}

	return dogu, nil
}

func (r *Remote) GetDogus(dogusToLoad []cescommons.QualifiedDoguVersion) (map[cescommons.QualifiedDoguName]*core.Dogu, error) {
	dogus := make(map[cescommons.QualifiedDoguName]*core.Dogu)

	var errs []error
	for _, doguRef := range dogusToLoad {
		dogu, err := r.GetDogu(doguRef)
		errs = append(errs, err)

		dogus[doguRef.Name] = dogu
	}

	return dogus, errors.Join(errs...)
}

func isConnectionError(err error) bool {
	return !strings.Contains(err.Error(), remotedogudescriptor.ConnectionError.Error())
}
