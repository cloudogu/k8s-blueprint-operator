package componentcr

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	return &ecosystem.ComponentInstallation{
		Name:               name,
		DeployNamespace:    cr.Spec.DeployNamespace,
		Version:            version,
		Status:             cr.Status.Status,
		Health:             ecosystem.HealthStatus(cr.Status.Health),
		PersistenceContext: persistenceContext,
	}, nil
}

func toComponentCR(componentInstallation *ecosystem.ComponentInstallation) *compV1.Component {
	return &compV1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(componentInstallation.Name.SimpleName),
			Labels: map[string]string{
				ComponentNameLabelKey:    string(componentInstallation.Name.SimpleName),
				ComponentVersionLabelKey: componentInstallation.Version.String(),
			},
		},
		Spec: compV1.ComponentSpec{
			Namespace: string(componentInstallation.Name.Namespace),
			Name:      string(componentInstallation.Name.SimpleName),
			Version:   componentInstallation.Version.String(),
			// TODO
			DeployNamespace: "",
			// TODO
			ValuesYamlOverwrite: "",
		},
	}
}

type componentCRPatch struct {
	Spec componentSpecPatch `json:"spec"`
}

type componentSpecPatch struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
}

func toComponentCRPatch(component *ecosystem.ComponentInstallation) *componentCRPatch {
	return &componentCRPatch{
		Spec: componentSpecPatch{
			Namespace: string(component.Name.Namespace),
			Name:      string(component.Name.SimpleName),
			Version:   component.Version.String(),
		},
	}
}

func toComponentCRPatchBytes(component *ecosystem.ComponentInstallation) ([]byte, error) {
	crPatch := toComponentCRPatch(component)
	patch, err := json.Marshal(crPatch)
	if err != nil {
		return []byte{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot patch component CR for component %q", component.Name),
		}
	}
	return patch, nil
}
