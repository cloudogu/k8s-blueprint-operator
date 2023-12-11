package doguregistry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/remote"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

type Remote struct {
	registry cesappLibRemoteRegistry
}

func NewRemote(registry remote.Registry) *Remote {
	return &Remote{registry: registry}
}

func (r *Remote) GetDogu(qualifiedDoguName string, version string) (*core.Dogu, error) {
	dogu, err := r.registry.GetVersion(qualifiedDoguName, version)
	if err != nil {
		// this is ugly, maybe do it better in cesapp-lib?
		if strings.Contains(err.Error(), "404 not found") {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("dogu %q with version %q could not be found", qualifiedDoguName, version),
			}
		}

		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("failed to get dogu %q with version %q", qualifiedDoguName, version),
		}
	}

	return dogu, nil
}

func (r *Remote) GetDogus(dogusToLoad []domainservice.DoguToLoad) (map[string]*core.Dogu, error) {
	dogus := make(map[string]*core.Dogu)

	var errs []error
	for _, doguRef := range dogusToLoad {
		dogu, err := r.GetDogu(doguRef.QualifiedDoguName, doguRef.Version)
		errs = append(errs, err)

		dogus[doguRef.QualifiedDoguName] = dogu
	}

	return dogus, errors.Join(errs...)
}
