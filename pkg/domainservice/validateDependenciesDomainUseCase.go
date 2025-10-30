package domainservice

import (
	"context"
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	nginxDependencyName       = "nginx"
	registratorDependencyName = "registrator"
)

type ValidateDependenciesDomainUseCase struct {
	remoteDoguRegistry RemoteDoguRegistry
}

func NewValidateDependenciesDomainUseCase(remoteDoguRegistry RemoteDoguRegistry) *ValidateDependenciesDomainUseCase {
	return &ValidateDependenciesDomainUseCase{
		remoteDoguRegistry,
	}
}

// ValidateDependenciesForAllDogus checks if for all dogus in the blueprint there are also all dependencies in the blueprint.
// The dependencies are validated against dogu specifications in a remote dogu registry.
// This functions returns no error if everything is ok or
// a domain.InvalidBlueprintError if there are dependencies missing or
// an InternalError if there is any other error, e.g. with the connection to the remote dogu registry
func (useCase *ValidateDependenciesDomainUseCase) ValidateDependenciesForAllDogus(ctx context.Context, effectiveBlueprint domain.EffectiveBlueprint) error {
	logger := log.FromContext(ctx).WithName("ValidateDependenciesDomainUseCase.ValidateDependenciesForAllDogus")
	wantedDogus := effectiveBlueprint.GetWantedDogus()
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
	logger.V(2).Info("load dogu specifications...", "wantedDogus", wantedDogus)
	doguSpecsOfWantedDogus, err := useCase.remoteDoguRegistry.GetDogus(ctx, dogusToLoad)
	if err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			return &domain.InvalidBlueprintError{WrappedError: err, Message: "remote dogu registry has no dogu specification for at least one wanted dogu"}
		} else { // should be InternalError
			return &InternalError{WrappedError: err, Message: "cannot load dogu specifications from remote registry for dogu dependency validation"}
		}
	}
	logger.V(2).Info("dogu specifications loaded", "specs", doguSpecsOfWantedDogus)

	var errorList []error
	for _, wantedDogu := range wantedDogus {
		dependencyDoguSpec := doguSpecsOfWantedDogus[wantedDogu.Name]
		logger.V(2).Info(fmt.Sprintf("check dependencies of %q in version %q", wantedDogu.Name, wantedDogu.Version.Raw))
		err = useCase.checkDoguDependencies(ctx, wantedDogus, doguSpecsOfWantedDogus, dependencyDoguSpec.Dependencies)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("dependencies for dogu '%s' are not satisfied in blueprint: %w", wantedDogu.Name, err))
		}
	}
	err = errors.Join(errorList...)
	if err != nil {
		err = &domain.InvalidBlueprintError{
			WrappedError: err,
			Message:      "dependencies are not satisfied in effective blueprint",
		}
	}
	return err
}

func (useCase *ValidateDependenciesDomainUseCase) checkDoguDependencies(
	ctx context.Context,
	wantedDogus []domain.Dogu,
	knownDoguSpecs map[cescommons.QualifiedName]*core.Dogu,
	dependenciesOfWantedDogu []core.Dependency,
) error {
	logger := log.FromContext(ctx).WithName("ValidateDependenciesDomainUseCase.checkDoguDependencies")
	var problems []error

	for _, dependencyOfWantedDogu := range dependenciesOfWantedDogu {
		logger.V(2).Info(fmt.Sprintf(
			"check dependency %q in version %q...",
			dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Version,
		))
		if dependencyOfWantedDogu.Type != core.DependencyTypeDogu {
			logger.V(1).Info(fmt.Sprintf(
				"dogu has a dependency %q of type %q. At the moment only dogu dependencies are validated.",
				dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Type,
			))
			continue
		}

		// Ignore registrator dogu because this is only needed in single node version of the Cloudogu EcoSystem.
		if dependencyOfWantedDogu.Name == registratorDependencyName {
			continue
		}

		// Exception for the old nginx dependency from the single node Cloudogu EcoSystem.
		// The nginx dependency was replaced by a the k8s-ces-gateway and k8s-ces-assets component
		if dependencyOfWantedDogu.Name == nginxDependencyName {
			continue
		}

		// check if dogu exists in blueprint and version is ok
		err := checkDoguDependency(dependencyOfWantedDogu, wantedDogus, knownDoguSpecs)
		problems = append(problems, err)
	}
	err := errors.Join(problems...)
	return err
}

func checkDoguDependency(
	dependencyOfWantedDogu core.Dependency,
	wantedDogus []domain.Dogu,
	knownDoguSpecs map[cescommons.QualifiedName]*core.Dogu,
) error {
	// this also works with namespace changes as only the simple dogu name get searched
	dependencyInBlueprint, found := domain.FindDoguByName(wantedDogus, cescommons.SimpleName(dependencyOfWantedDogu.Name))
	if !found {
		return fmt.Errorf("dependency '%s' in version '%s' is not a present dogu in the effective blueprint", dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Version)
	}
	dependencyDoguSpec := knownDoguSpecs[dependencyInBlueprint.Name]
	return checkDependencyVersion(dependencyInBlueprint, dependencyDoguSpec.Version)
}

func checkDependencyVersion(doguInBlueprint domain.Dogu, expectedVersion string) error {
	// it does not count as an error if no version is specified as the field is optional
	if expectedVersion == "" {
		return nil
	}
	comparator, err := core.ParseVersionComparator(expectedVersion)
	if err != nil {
		return fmt.Errorf("failed to parse version comparator of version %s for dogu dependency %s: %w", expectedVersion, doguInBlueprint.Name, err)
	}
	doguInBlueprintVersion := core.Version{}
	if doguInBlueprint.Version != nil {
		doguInBlueprintVersion = *doguInBlueprint.Version
	}
	allows, err := comparator.Allows(doguInBlueprintVersion)
	if err != nil {
		return fmt.Errorf("an error occurred when comparing the versions: %w", err)
	}
	if !allows {
		return fmt.Errorf("parsed version %s does not fulfill version requirement of %s dogu %s", doguInBlueprint.Version.Raw, expectedVersion, doguInBlueprint.Name)
	}
	return nil // no error, dependency is ok
}
