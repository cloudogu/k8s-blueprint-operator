package domainservice

import (
	"context"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
)

var testDataDoguRegistry = stubRemoteDoguRegistry{
	dogus: map[cescommons.QualifiedDoguName]map[string]*core.Dogu{
		officialPostgres: {
			"1.0.0-1": &core.Dogu{
				Name:         "official/postgres",
				Version:      "1.0.0-1",
				Dependencies: []core.Dependency{},
			},
		},
		premiumPostgres: { // to test namespace changes
			"1.0.0-1": &core.Dogu{
				Name:         "official/postgres",
				Version:      "1.0.0-1",
				Dependencies: []core.Dependency{},
			},
		},
		officialRedmine: {
			"1.0.0-1": &core.Dogu{
				Name:    "official/redmine",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypeDogu, Name: "postgres", Version: "1.0.0-1"},
				},
			},
		},
		helloworldBluespice: {
			"1.0.0-1": &core.Dogu{
				Name:    "helloworld/bluespice",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypeDogu, Name: "official/mysql", Version: "1.0.0-1"},
				},
			},
		},
		officialK8sCesControl: {
			"1.0.0-1": &core.Dogu{
				Name:    "official/k8s-ces-control",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypePackage, Name: "jq", Version: "1.0.0-1"},
				},
			},
		},
		officialPlantuml: {
			"1.0.0-1": &core.Dogu{
				Name:    "official/plantuml",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypeDogu, Name: "nginx", Version: "1.0.0-1"},
				},
			},
		},
		k8sNginxStatic: {
			"1.0.0-1": &core.Dogu{
				Name:    "k8s/nginx-static",
				Version: "1.0.0-1",
			},
		},
		k8sNginxIngress: {
			"1.0.0-1": &core.Dogu{
				Name:    "k8s/nginx-ingress",
				Version: "1.0.0-1",
			},
		},
		ldapMapper: {
			"1.0.0-1": &core.Dogu{
				Name:    "official/ldap-mapper",
				Version: "1.0.0-1",
				Dependencies: []core.Dependency{
					{Type: core.DependencyTypeDogu, Name: "registrator"},
				},
			},
		},
	},
}

type stubRemoteDoguRegistry struct {
	dogus map[cescommons.QualifiedDoguName]map[string]*core.Dogu
}

func (registry stubRemoteDoguRegistry) GetDogu(ctx context.Context, qualifiedDoguVersion cescommons.QualifiedDoguVersion) (*core.Dogu, error) {
	dogu := registry.dogus[qualifiedDoguVersion.Name][qualifiedDoguVersion.Version.Raw]
	if dogu == nil {
		return nil, &NotFoundError{Message: fmt.Sprintf("dogu %s in version %s not found", qualifiedDoguVersion.Name, qualifiedDoguVersion.Version.Raw)}
	}
	return dogu, nil
}

func (registry stubRemoteDoguRegistry) GetDogus(ctx context.Context, dogusToLoad []cescommons.QualifiedDoguVersion) (map[cescommons.QualifiedDoguName]*core.Dogu, error) {
	dogus := map[cescommons.QualifiedDoguName]*core.Dogu{}
	for _, doguToLoad := range dogusToLoad {
		doguSpec, err := registry.GetDogu(ctx, doguToLoad)
		if err != nil {
			return nil, err
		}
		dogus[doguToLoad.Name] = doguSpec
	}
	return dogus, nil
}
