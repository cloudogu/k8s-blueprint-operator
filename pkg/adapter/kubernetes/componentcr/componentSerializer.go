package componentcr

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	deployConfigKeyDeployNamespace = "deployNamespace"
	deployConfigKeyOverwriteConfig = "overwriteConfig"
	ComponentNameLabelKey          = "k8s.cloudogu.com/component.name"
	ComponentVersionLabelKey       = "k8s.cloudogu.com/component.version"
)

func parseComponentCR(cr *compV1.Component) (*ecosystem.ComponentInstallation, error) {
	if cr == nil {
		return nil, domainservice.NewInternalError(nil, "cannot parse component CR as it is nil")
	}

	expectedVersion, err := semver.NewVersion(cr.Spec.Version)
	if err != nil {
		return nil, domainservice.NewInternalError(err, "cannot load component CR as it cannot be parsed correctly")
	}

	// ignore error as the actual version could be not set and a nil pointer as the version is exactly what we want then
	actualVersion, _ := semver.NewVersion(cr.Status.InstalledVersion)

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
		ExpectedVersion:    expectedVersion,
		ActualVersion:      actualVersion,
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

func toComponentCR(component *ecosystem.ComponentInstallation) (*compV1.Component, error) {
	deployNamespace, err := toDeployNamespace(component.DeployConfig)
	if err != nil {
		return nil, err
	}

	valuesYamlOverwrite, err := toValuesYamlOverwrite(component.DeployConfig)
	if err != nil {
		return nil, err
	}

	spec := compV1.ComponentSpec{
		Namespace: string(component.Name.Namespace),
		Name:      string(component.Name.SimpleName),
		Version:   component.ExpectedVersion.String(),
	}
	if deployNamespace != "" {
		spec.DeployNamespace = deployNamespace
	}
	if valuesYamlOverwrite != "" {
		spec.ValuesYamlOverwrite = valuesYamlOverwrite
	}

	return &compV1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(component.Name.SimpleName),
			Labels: map[string]string{
				ComponentNameLabelKey:          string(component.Name.SimpleName),
				ComponentVersionLabelKey:       component.ExpectedVersion.String(),
				"app":                          "ces",
				"k8s.cloudogu.com/app":         "ces",
				"component.name":               string(component.Name.SimpleName),
				"app.kubernetes.io/name":       string(component.Name.SimpleName),
				"app.kubernetes.io/version":    component.ExpectedVersion.String(),
				"app.kubernetes.io/part-of":    "ces",
				"app.kubernetes.io/managed-by": "k8s-blueprint-operator",
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
		Version:   component.ExpectedVersion.String(),
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
