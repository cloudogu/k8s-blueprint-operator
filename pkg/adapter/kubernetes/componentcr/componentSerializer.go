package componentcr

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	deployConfigKeyDeployNamespace = "deployNamespace"
	deployConfigKeyOverwriteConfig = "overwriteConfig"
)

func parseComponentCR(cr *compV1.Component) (*ecosystem.ComponentInstallation, error) {
	if cr == nil {
		return nil, domainservice.NewInternalError(nil, "cannot parse component CR as it is nil")
	}

	version, err := semver.NewVersion(cr.Spec.Version)
	if err != nil {
		return nil, domainservice.NewInternalError(err, "cannot load component CR as it cannot be parsed correctly")
	}

	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[componentInstallationRepoContextKey] = componentInstallationRepoContext{
		resourceVersion: cr.GetResourceVersion(),
	}

	name, err := common.NewQualifiedComponentName(common.ComponentNamespace(cr.Spec.Namespace), common.SimpleComponentName(cr.Name))
	if err != nil {
		return nil, err
	}

	componentConfig, err := parseDeployConfig(cr)
	if err != nil {
		return nil, err
	}

	return &ecosystem.ComponentInstallation{
		Name:               name,
		Version:            version,
		Status:             cr.Status.Status,
		Health:             ecosystem.HealthStatus(cr.Status.Health),
		PersistenceContext: persistenceContext,
		DeployConfig:       componentConfig,
	}, nil
}

func parseDeployConfig(cr *compV1.Component) (ecosystem.DeployConfig, error) {
	componentConfig := ecosystem.DeployConfig{}
	if cr.Spec.DeployNamespace != "" {
		componentConfig[deployConfigKeyDeployNamespace] = cr.Spec.DeployNamespace
	}

	if cr.Spec.ValuesYamlOverwrite != "" {
		valuesYamlOverwrite := map[string]interface{}{}
		// We need to use k8syaml here because goyaml unmarshals to map[interface{}]interface {} which is not supported setting in a k8s resource.
		err := k8syaml.Unmarshal([]byte(cr.Spec.ValuesYamlOverwrite), &valuesYamlOverwrite)
		if err != nil {
			return nil, domainservice.NewInternalError(err, "failed to unmarshal values yaml overwrite %q", cr.Spec.ValuesYamlOverwrite)
		}
		componentConfig[deployConfigKeyOverwriteConfig] = valuesYamlOverwrite
	}

	return componentConfig, nil
}

func toComponentCR(componentInstallation *ecosystem.ComponentInstallation) (*compV1.Component, error) {
	deployNamespace, err := toDeployNamespace(componentInstallation.DeployConfig)
	if err != nil {
		return nil, err
	}

	valuesYamlOverwrite, err := toValuesYamlOverwrite(componentInstallation.DeployConfig)
	if err != nil {
		return nil, err
	}

	spec := compV1.ComponentSpec{
		Namespace: string(componentInstallation.Name.Namespace),
		Name:      string(componentInstallation.Name.SimpleName),
		Version:   componentInstallation.Version.String(),
	}
	if deployNamespace != "" {
		spec.DeployNamespace = deployNamespace
	}
	if valuesYamlOverwrite != "" {
		spec.ValuesYamlOverwrite = valuesYamlOverwrite
	}

	return &compV1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(componentInstallation.Name.SimpleName),
			Labels: map[string]string{
				ComponentNameLabelKey:    string(componentInstallation.Name.SimpleName),
				ComponentVersionLabelKey: componentInstallation.Version.String(),
			},
		},
		Spec: spec,
	}, nil
}

func toDeployNamespace(deployConfig ecosystem.DeployConfig) (string, error) {
	deployNamespace, found := deployConfig[deployConfigKeyDeployNamespace]
	if !found {
		return "", nil
	}
	deployNamespaceStr, ok := deployNamespace.(string)
	if !ok {
		return "", fmt.Errorf("deployNamespace is not type of string")
	}

	return deployNamespaceStr, nil
}

func toValuesYamlOverwrite(deployConfig ecosystem.DeployConfig) (string, error) {
	in, found := deployConfig[deployConfigKeyOverwriteConfig]
	if !found {
		return "", nil
	}
	valuesYamlOverwriteBytes, err := yaml.Marshal(in)
	if err != nil {
		return "", fmt.Errorf("failed to marshal overwrite config %q", in)
	}

	return string(valuesYamlOverwriteBytes), nil
}

type componentCRPatch struct {
	Spec componentSpecPatch `json:"spec"`
}

type componentSpecPatch struct {
	Namespace           string  `json:"namespace"`
	Name                string  `json:"name"`
	Version             string  `json:"version"`
	DeployNamespace     *string `json:"deployNamespace"`
	ValuesYamlOverwrite *string `json:"valuesYamlOverwrite"`
}

func toComponentCRPatch(component *ecosystem.ComponentInstallation) (*componentCRPatch, error) {
	deployNamespace, err := toDeployNamespace(component.DeployConfig)
	if err != nil {
		return nil, err
	}

	valuesYamlOverwrite, err := toValuesYamlOverwrite(component.DeployConfig)
	if err != nil {
		return nil, err
	}

	spec := componentSpecPatch{
		Namespace: string(component.Name.Namespace),
		Name:      string(component.Name.SimpleName),
		Version:   component.Version.String(),
	}
	if deployNamespace != "" {
		spec.DeployNamespace = &deployNamespace
	}
	if valuesYamlOverwrite != "" {
		spec.ValuesYamlOverwrite = &valuesYamlOverwrite
	}

	return &componentCRPatch{
		Spec: spec,
	}, nil
}

func toComponentCRPatchBytes(component *ecosystem.ComponentInstallation) ([]byte, error) {
	crPatch, err := toComponentCRPatch(component)
	if err != nil {
		return nil, domainservice.NewInternalError(err, "failed to create component CR patch for component %q", component.Name)
	}
	patch, err := json.Marshal(crPatch)

	if err != nil {
		return []byte{}, domainservice.NewInternalError(err, "cannot patch component CR for component %q", component.Name)
	}
	return patch, nil
}
