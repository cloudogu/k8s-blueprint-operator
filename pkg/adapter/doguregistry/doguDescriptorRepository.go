package doguregistry

import (
	"context"
	"errors"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	cloudoguerrors "github.com/cloudogu/ces-commons-lib/errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DoguDescriptorRepository struct {
	remoteRepository remoteDoguDescriptorRepository
	localRepository  localDoguDescriptorRepository
}

func NewDoguDescriptorRepository(remoteRepository remoteDoguDescriptorRepository, localRepository localDoguDescriptorRepository) *DoguDescriptorRepository {
	return &DoguDescriptorRepository{remoteRepository: remoteRepository, localRepository: localRepository}
}

func (r *DoguDescriptorRepository) GetDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedVersion) (*core.Dogu, error) {
	logger := log.FromContext(ctx).
		WithName("DoguDescriptorRepository.GetDogu").
		WithValues("dogu", qualifiedDoguVersion.Name.SimpleName)

	// Try to get the dogu from the local repository first.
	dogu := r.getLocalDogu(ctx, qualifiedDoguVersion, logger)
	if dogu != nil {
		return dogu, nil
	}

	dogu, err := r.getRemoteDogu(ctx, qualifiedDoguVersion)
	if err != nil {
		return nil, err
	}

	err = r.localRepository.Add(ctx, qualifiedDoguVersion.Name.SimpleName, dogu)
	if err != nil {
		// just log the error, no need to fail the reconcilation
		logger.Info("failed to add dogu to local repository",
			"error", err,
			"dogu", qualifiedDoguVersion.Name.SimpleName,
			"version", qualifiedDoguVersion.Version.Raw)
	}
	return dogu, nil
}

func (r *DoguDescriptorRepository) getRemoteDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedVersion) (*core.Dogu, error) {
	// do not retry here. If any error happens, just reconcile later. We only do retries in application level.
	// This makes the code way easier and non-blocking.
	dogu, err := r.remoteRepository.Get(ctx, qualifiedDoguVersion)
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

func (r *DoguDescriptorRepository) getLocalDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedVersion, logger logr.Logger) *core.Dogu {
	dogu, err := r.localRepository.Get(ctx, cescommons.NewSimpleNameVersion(qualifiedDoguVersion.Name.SimpleName, qualifiedDoguVersion.Version))
	if err == nil {
		logger.V(2).Info("local dogu descriptor hit", "dogu", qualifiedDoguVersion.Name.SimpleName)
		return dogu
	} else {
		logger.V(2).Info("local dogu descriptor miss", "error", err)
		return nil
	}
}

func (r *DoguDescriptorRepository) GetDogus(ctx context.Context, dogusToLoad []cescommons.QualifiedVersion) (map[cescommons.QualifiedName]*core.Dogu, error) {
	dogus := make(map[cescommons.QualifiedName]*core.Dogu)

	var errs []error
	for _, doguRef := range dogusToLoad {
		dogu, err := r.GetDogu(ctx, doguRef)
		errs = append(errs, err)

		dogus[doguRef.Name] = dogu
	}

	return dogus, errors.Join(errs...)
}
