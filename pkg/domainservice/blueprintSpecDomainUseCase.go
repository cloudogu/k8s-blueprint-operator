package domainservice

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
)

type BlueprintSpecDomainUseCase struct {
	remoteDoguRegistry RemoteDoguRegistry
}

func (useCase *BlueprintSpecDomainUseCase) ValidateDoguDependencies(spec domain.BlueprintSpec) error {
	spec.Blueprint.GetWantedDogus()
	//TODO
	//_, _ = useCase.findAllDependencies(spec.Blueprint)
}

func (useCase *BlueprintSpecDomainUseCase) findDirectDependencies(blueprint domain.Blueprint) (map[string]core.Dogu, error) {
	neededDependencies := make(map[string]core.Dogu)

	for _, wantedDogu := range blueprint.GetWantedDogus() {
		doguSpec, err := useCase.remoteDoguRegistry.GetDogu(wantedDogu.GetQualifiedName(), wantedDogu.Version)
		if err != nil {
			return nil, fmt.Errorf("cannot load dogu specifications from remote registry for dogu dependency validation: %w", err)
		}
		neededDependencies[wantedDogu.GetQualifiedName()] = doguSpec
	}
	//TODO
	return nil, nil
}
