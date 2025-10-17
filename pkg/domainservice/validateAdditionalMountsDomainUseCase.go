package domainservice

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ValidateAdditionalMountsDomainUseCase struct {
	remoteDoguRegistry RemoteDoguRegistry
}

func NewValidateAdditionalMountsDomainUseCase(remoteDoguRegistry RemoteDoguRegistry) *ValidateAdditionalMountsDomainUseCase {
	return &ValidateAdditionalMountsDomainUseCase{
		remoteDoguRegistry,
	}
}

// ValidateAdditionalMounts checks if for all additional mounts of dogus fit to their volumes described in the dogu specifications.
// The dependencies are validated against dogu specifications in a remote dogu registry.
// This functions returns no error if everything is ok or
// a domain.InvalidBlueprintError if there are invalid additional mounts
// an InternalError if there is any other error, e.g. with the connection to the remote dogu registry
func (useCase *ValidateAdditionalMountsDomainUseCase) ValidateAdditionalMounts(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error {
	logger := log.FromContext(ctx).WithName("ValidateAdditionalMountsDomainUseCase.ValidateAdditionalMounts")
	dogusWithMounts := filterDogusWithAdditionalMounts(effectiveBlueprint.GetWantedDogus())
	if len(dogusWithMounts) == 0 {
		logger.V(2).Info("skip additional mounts validation as no dogus have additional mounts")
		return nil
	}
	logger.V(2).Info("load dogu specifications...", "dogusWithMounts", dogusWithMounts)
	doguSpecs, err := loadDoguSpecifications(ctx, useCase.remoteDoguRegistry, dogusWithMounts)
	if err != nil {
		return err
	}
	logger.V(2).Info("dogu specifications loaded", "specs", doguSpecs)

	var errorList []error
	for _, wantedDogu := range dogusWithMounts {
		doguSpec := doguSpecs[wantedDogu.Name]
		errorList = append(errorList, validateAdditionalMountsForDogu(wantedDogu, doguSpec))
	}
	err = errors.Join(errorList...)
	if err != nil {
		err = &domain.InvalidBlueprintError{
			WrappedError: err,
			Message:      "additionalMounts are invalid in effective blueprint",
		}
	}
	return err
}

func filterDogusWithAdditionalMounts(dogus []domain.Dogu) []domain.Dogu {
	var result []domain.Dogu
	for _, dogu := range dogus {
		if len(dogu.AdditionalMounts) > 0 {
			result = append(result, dogu)
		}
	}
	return result
}

func validateAdditionalMountsForDogu(dogu domain.Dogu, doguSpec *core.Dogu) error {
	possibleVolumes := filterClientVolumes(doguSpec.Volumes)
	possibleVolumeNames := util.Map(possibleVolumes, func(volume core.Volume) string {
		return volume.Name
	})
	var errorList []error
	for _, mount := range dogu.AdditionalMounts {
		if !slices.Contains(possibleVolumeNames, mount.Volume) {
			errorList = append(errorList, fmt.Errorf("volume %q in additional mount for dogu %q is invalid "+
				"because either the volume does not exist or it has volume clients. "+
				"It needs to be one of %+q", mount.Volume, doguSpec.Name, possibleVolumeNames))
		}
	}
	return errors.Join(errorList...)
}

func filterClientVolumes(volumes []core.Volume) []core.Volume {
	// At the moment, clients can be just configMaps, and they get mounted directly without an empty dir volume etc.
	// Therefore, additional mounts cannot work in this case.
	var result []core.Volume
	for _, volume := range volumes {
		if len(volume.Clients) == 0 {
			result = append(result, volume)
		}
	}
	return result
}
