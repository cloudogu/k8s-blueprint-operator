package domainservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	nginxDependencyName        = "nginx"
	nginxStaticDependencyName  = "nginx-static"
	nginxIngressDependencyName = "nginx-ingress"
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
	dogusToLoad := util.Map(wantedDogus, func(dogu domain.Dogu) DoguToLoad {
		return DoguToLoad{
			DoguName: dogu.Name,
			Version:  dogu.Version.Raw,
		}
	})
	logger.Info("load dogu specifications...", "wantedDogus", wantedDogus)
	doguSpecsOfWantedDogus, err := useCase.remoteDoguRegistry.GetDogus(dogusToLoad)
	if err != nil {
		var notFoundError *NotFoundError
		if errors.As(err, &notFoundError) {
			return &domain.InvalidBlueprintError{WrappedError: err, Message: "remote dogu registry has no dogu specification for at least one wanted dogu"}
		} else { // should be InternalError
			return &InternalError{WrappedError: err, Message: "cannot load dogu specifications from remote registry for dogu dependency validation"}
		}
	}
	logger.Info("dogu specifications loaded", "specs", doguSpecsOfWantedDogus)

	var errorList []error
	for _, wantedDogu := range wantedDogus {
		dependencyDoguSpec := doguSpecsOfWantedDogus[wantedDogu.Name]
		logger.Info(fmt.Sprintf("check dependencies of %q in version %q", wantedDogu.Name, wantedDogu.Version.Raw))
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
	knownDoguSpecs map[common.QualifiedDoguName]*core.Dogu,
	dependenciesOfWantedDogu []core.Dependency,
) error {
	logger := log.FromContext(ctx).WithName("ValidateDependenciesDomainUseCase.checkDoguDependencies")
	var problems []error

	for _, dependencyOfWantedDogu := range dependenciesOfWantedDogu {
		logger.Info(fmt.Sprintf("check dependency %q in version %q...", dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Version))
		if dependencyOfWantedDogu.Type != core.DependencyTypeDogu {
			logger.Info(fmt.Sprintf("dogu has a dependency %q of type %q. At the moment only dogu dependencies are validated.", dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Type))
			continue
		}

		// Exception for the old nginx dependency from the single node Cloudogu EcoSystem.
		// We only have to check if nginx-static and nginx-ingress are present.
		if dependencyOfWantedDogu.Name == nginxDependencyName {
			if !checkNginxIngressAndStatic(wantedDogus) {
				problems = append(problems, fmt.Errorf("dogu has %q dependency but %q and %q are missing in the effective blueprint", nginxDependencyName, nginxIngressDependencyName, nginxStaticDependencyName))
			}
			logger.Info(fmt.Sprintf("dogu has dependency %q. %q and %q are available.", nginxDependencyName, nginxIngressDependencyName, nginxStaticDependencyName))
			continue
		}

		// check if dogu exists in blueprint and version is ok
		err := checkDoguDependency(dependencyOfWantedDogu, wantedDogus, knownDoguSpecs)
		problems = append(problems, err)
	}
	err := errors.Join(problems...)
	return err
}

func checkNginxIngressAndStatic(wantedDogus []domain.Dogu) bool {
	foundNginxIngress := isDoguInSlice(wantedDogus, nginxIngressDependencyName)
	foundNginxStatic := isDoguInSlice(wantedDogus, nginxStaticDependencyName)

	return foundNginxIngress && foundNginxStatic
}

func isDoguInSlice(dogus []domain.Dogu, name common.SimpleDoguName) bool {
	for _, dogu := range dogus {
		if dogu.Name.Name == name {
			return true
		}
	}

	return false
}

func checkDoguDependency(
	dependencyOfWantedDogu core.Dependency,
	wantedDogus []domain.Dogu,
	knownDoguSpecs map[common.QualifiedDoguName]*core.Dogu,
) error {
	// this also works with namespace changes as only the simple dogu name get searched
	dependencyInBlueprint, err := domain.FindDoguByName(wantedDogus, common.SimpleDoguName(dependencyOfWantedDogu.Name))
	if err != nil {
		return fmt.Errorf("dependency '%s' in version '%s' is not a present dogu in the effective blueprint", dependencyOfWantedDogu.Name, dependencyOfWantedDogu.Version)
	}
	// dependencyDoguSpec := useCase.remoteDoguRegistry.GetDogu(dependencyInBlueprint.GetQualifiedName(), dependencyInBlueprint.Version)
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
	allows, err := comparator.Allows(doguInBlueprint.Version)
	if err != nil {
		return fmt.Errorf("an error occurred when comparing the versions: %w", err)
	}
	if !allows {
		return fmt.Errorf("parsed version %s does not fulfill version requirement of %s dogu %s", doguInBlueprint.Version.Raw, expectedVersion, doguInBlueprint.Name)
	}
	return nil // no error, dependency is ok
}
