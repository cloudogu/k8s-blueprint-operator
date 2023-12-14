package domainservice

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
)

var testDataDoguRegistry = stubRemoteDoguRegistry{
	dogus: map[string]map[string]*core.Dogu{
		"official/postgres": {
			"1.0.0-1": &core.Dogu{
				Name:         "official/postgres",
				Version:      "1.0.0-1",
				Dependencies: []core.Dependency{},
			},
		},
		"premium/postgres": { // to test namespace changes
			"1.0.0-1": &core.Dogu{
				Name:         "official/postgres",
				Version:      "1.0.0-1",
				Dependencies: []core.Dependency{},
			},
		},
		"official/redmine": {
			"1.0.0-1": &core.Dogu{
				Name:    "official/redmine",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypeDogu, Name: "postgres", Version: "1.0.0-1"},
				},
			},
		},
		"helloworld/bluespice": {
			"1.0.0-1": &core.Dogu{
				Name:    "helloworld/bluespice",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypeDogu, Name: "official/mysql", Version: "1.0.0-1"},
				},
			},
		},
		"official/k8s-ces-control": {
			"1.0.0-1": &core.Dogu{
				Name:    "official/k8s-ces-control",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypePackage, Name: "jq", Version: "1.0.0-1"},
				},
			},
		},
	},
}

type stubRemoteDoguRegistry struct {
	dogus map[string]map[string]*core.Dogu
}

func (registry stubRemoteDoguRegistry) GetDogu(qualifiedDoguName string, version string) (*core.Dogu, error) {
	dogu := registry.dogus[qualifiedDoguName][version]
	if dogu == nil {
		return nil, &NotFoundError{Message: fmt.Sprintf("dogu %s in version %s not found", qualifiedDoguName, version)}
	}
	return dogu, nil
}

func (registry stubRemoteDoguRegistry) GetDogus(dogusToLoad []DoguToLoad) (map[string]*core.Dogu, error) {
	dogus := map[string]*core.Dogu{}
	for _, doguToLoad := range dogusToLoad {
		doguSpec, err := registry.GetDogu(doguToLoad.QualifiedDoguName, doguToLoad.Version)
		if err != nil {
			return nil, err
		}
		dogus[doguToLoad.QualifiedDoguName] = doguSpec
	}
	return dogus, nil
}
