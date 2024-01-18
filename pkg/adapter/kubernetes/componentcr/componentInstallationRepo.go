package componentcr

import (
	"context"
	"fmt"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/cesapp-lib/core"
	compCli "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

const componentInstallationRepoContextKey = "componentInstallationRepoContext"

type componentInstallationRepoContext struct {
	resourceVersion string
}

type componentInstallationRepo struct {
	namespace       string
	componentClient compCli.ComponentV1Alpha1Interface
}

// NewComponentInstallationRepo creates a new component repo adapter.
func NewComponentInstallationRepo(namespace string, componentClient compCli.ComponentV1Alpha1Interface) domainservice.ComponentInstallationRepository {
	return &componentInstallationRepo{namespace: namespace, componentClient: componentClient}
}

func (repo *componentInstallationRepo) GetByName(ctx context.Context, componentName string) (*ecosystem.ComponentInstallation, error) {
	cr, err := repo.componentClient.Components(repo.namespace).Get(ctx, componentName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot read component CR %q as it does not exist", componentName),
			}
		}
		return nil, domainservice.NewInternalError(err, "error while reading component CR %q", componentName)
	}

	return parseComponentCR(cr)
}

func (repo *componentInstallationRepo) GetAll(ctx context.Context) (map[string]*ecosystem.ComponentInstallation, error) {
	return nil, nil
}

func parseComponentCR(cr *compV1.Component) (*ecosystem.ComponentInstallation, error) {
	if cr == nil {
		return nil, &domainservice.InternalError{
			WrappedError: nil,
			Message:      "cannot parse component CR as it is nil",
		}
	}

	version, err := core.ParseVersion(cr.Spec.Version)
	if err != nil {
		return nil, domainservice.NewInternalError(err, "cannot load component CR as it cannot be parsed correctly")
	}

	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[componentInstallationRepoContextKey] = componentInstallationRepoContext{
		resourceVersion: cr.GetResourceVersion(),
	}
	return &ecosystem.ComponentInstallation{
		Namespace: cr.Namespace,
		Name:      cr.Name,
		Version:   version,
		Status:    cr.Status.Status,
		//Health:             ecosystem.HealthStatus(cr.Status.Health), # coming soon
		PersistenceContext: persistenceContext,
	}, nil
}
