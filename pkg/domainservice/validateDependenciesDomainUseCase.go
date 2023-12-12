package domainservice

import (
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
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
func (useCase *ValidateDependenciesDomainUseCase) ValidateDependenciesForAllDogus(effectiveBlueprint domain.EffectiveBlueprint) error {
	wantedDogus := effectiveBlueprint.GetWantedDogus()
	dogusToLoad := util.Map(wantedDogus, func(dogu domain.Dogu) DoguToLoad {
		return DoguToLoad{
			QualifiedDoguName: dogu.GetQualifiedName(),
			Version:           dogu.Version.Raw,
		}
	})
	doguSpecsOfWantedDogus, err := useCase.remoteDoguRegistry.GetDogus(dogusToLoad)
	if err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			return &domain.InvalidBlueprintError{WrappedError: err, Message: "remote dogu registry has no dogu specification for at least one wanted dogu"}
		} else { //should be InternalError
			return &InternalError{WrappedError: err, Message: "cannot load dogu specifications from remote registry for dogu dependency validation"}
		}
	}

	for _, wantedDogu := range wantedDogus {
		dependencyDoguSpec := doguSpecsOfWantedDogus[wantedDogu.GetQualifiedName()]
		err = useCase.checkDoguDependencies(wantedDogus, doguSpecsOfWantedDogus, dependencyDoguSpec.Dependencies)
		if err != nil {
			err = fmt.Errorf("dependencies for dogu '%s' are not satisfied blueprint: %w", wantedDogu.Name, err)
		}
	}

	return err
}

func (useCase *ValidateDependenciesDomainUseCase) checkDoguDependencies(
	wantedDogus []domain.Dogu,
	knownDoguSpecs map[string]*core.Dogu,
	dependenciesOfWantedDogu []core.Dependency,
) error {
	var problems []error

	for _, dependencyOfWantedDogu := range dependenciesOfWantedDogu {
		if dependencyOfWantedDogu.Type != core.DependencyTypeDogu {
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
	knownDoguSpecs map[string]*core.Dogu,
) error {
	// this also works with namespace changes as only the simple dogu name get searched
	dependencyInBlueprint, err := domain.FindDoguByName(wantedDogus, dependencyOfWantedDogu.Name)
	if err != nil {
		return fmt.Errorf("dependency '%s' in version '%s' is not a present dogu in the effective blueprint", dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Version)
	}
	//dependencyDoguSpec := useCase.remoteDoguRegistry.GetDogu(dependencyInBlueprint.GetQualifiedName(), dependencyInBlueprint.Version)
	dependencyDoguSpec := knownDoguSpecs[dependencyInBlueprint.GetQualifiedName()]
	return checkDependencyVersion(dependencyInBlueprint, dependencyDoguSpec.Version)
}

func checkDependencyVersion(doguInBlueprint domain.Dogu, expectedVersion string) error {
	// it does not count as an error if no version is specified as the field is optional
	if expectedVersion == "" {
		return nil
	}
	comparator, err := core.ParseVersionComparator(expectedVersion)
	if err != nil {
		return fmt.Errorf("failed to parse ParseVersionComparator of version %s for doguDependency %s: %w", expectedVersion, doguInBlueprint.Name, err)
	}
	allows, err := comparator.Allows(doguInBlueprint.Version)
	if err != nil {
		return fmt.Errorf("an error occurred when comparing the versions: %w", err)
	}
	if !allows {
		return fmt.Errorf("%s parsed Version does not fulfill version requirement of %s dogu %s", doguInBlueprint.Version.Raw, expectedVersion, doguInBlueprint.Name)
	}
	return nil // no error, dependency is ok
}
