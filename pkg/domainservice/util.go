package domainservice

import (
	"context"
	"errors"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

func loadDoguSpecifications(ctx context.Context, remoteDoguRegistry RemoteDoguRegistry, wantedDogus []domain.Dogu) (map[cescommons.QualifiedName]*core.Dogu, error) {
	dogusToLoad := util.Map(wantedDogus, func(dogu domain.Dogu) cescommons.QualifiedVersion {
		doguVersion := core.Version{}
		if dogu.Version != nil {
			doguVersion = *dogu.Version
		}
		return cescommons.QualifiedVersion{
			Name:    dogu.Name,
			Version: doguVersion,
		}
	})
	doguSpecsOfWantedDogus, err := remoteDoguRegistry.GetDogus(ctx, dogusToLoad)
	if err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			return nil, &domain.InvalidBlueprintError{WrappedError: err, Message: "remote dogu registry has no dogu specification for at least one wanted dogu"}
		} else { // should be InternalError
			return nil, &InternalError{WrappedError: err, Message: "cannot load dogu specifications from remote registry for dogu dependency validation"}
		}
	}
	return doguSpecsOfWantedDogus, nil
}
