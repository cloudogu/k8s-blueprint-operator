package componentcr

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/cesapp-lib/core"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	compCli "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

const (
	ComponentNameLabelKey               = "k8s.cloudogu.com/component.name"
	ComponentVersionLabelKey            = "k8s.cloudogu.com/component.version"
	componentInstallationRepoContextKey = "componentInstallationRepoContext"
)

type componentInstallationRepoContext struct {
	resourceVersion string
}

type componentInstallationRepo struct {
	componentClient compCli.ComponentInterface
}

// NewComponentInstallationRepo creates a new component repo adapter.
func NewComponentInstallationRepo(componentClient compCli.ComponentInterface) domainservice.ComponentInstallationRepository {
	return &componentInstallationRepo{componentClient: componentClient}
}

// GetByName fetches a named component resource and returns it as ecosystem.ComponentInstallation.
func (repo *componentInstallationRepo) GetByName(ctx context.Context, componentName string) (*ecosystem.ComponentInstallation, error) {
	cr, err := repo.componentClient.Get(ctx, componentName, metav1.GetOptions{})
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

// GetAll fetches all installed component resources and returns them as list of ecosystem.ComponentInstallation.
func (repo *componentInstallationRepo) GetAll(ctx context.Context) (map[string]*ecosystem.ComponentInstallation, error) {
	list, err := repo.componentClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	componentInstallations := make(map[string]*ecosystem.ComponentInstallation, len(list.Items))
	for _, componentCr := range list.Items {
		cr, err := parseComponentCR(&componentCr)
		if err != nil {
			return nil, domainservice.NewInternalError(err, "failed to parse component CR %#v", componentCr)
		}
		componentInstallations[componentCr.Name] = cr
	}

	return componentInstallations, nil
}

func (repo *componentInstallationRepo) Update(ctx context.Context, component *ecosystem.ComponentInstallation) error {
	logger := log.FromContext(ctx).WithName("doguInstallationRepo.Update")
	patch, err := toComponentCRPatchBytes(component)
	if err != nil {
		return domainservice.NewInternalError(err, "failed to get patch bytes from component %q", component.Name)
	}

	logger.Info("patch component CR", "doguName", component.Name, "doguPatch", string(patch))

	_, err = repo.componentClient.Patch(ctx, component.Name, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return domainservice.NewInternalError(err, "failed to patch component %q", component.Name)
	}

	return nil
}

func (repo *componentInstallationRepo) Delete(ctx context.Context, componentName string) error {
	err := repo.componentClient.Delete(ctx, componentName, metav1.DeleteOptions{})
	if err != nil {
		return domainservice.NewInternalError(err, "failed to delete component CR %q", componentName)
	}

	return nil
}

func (repo *componentInstallationRepo) Create(ctx context.Context, component *ecosystem.ComponentInstallation) error {
	_, err := repo.componentClient.Create(ctx, toComponentCR(component), metav1.CreateOptions{})
	if err != nil {
		return domainservice.NewInternalError(err, "failed to create component CR %q", component.Name)
	}

	return nil
}

func parseComponentCR(cr *compV1.Component) (*ecosystem.ComponentInstallation, error) {
	if cr == nil {
		return nil, domainservice.NewInternalError(nil, "cannot parse component CR as it is nil")
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
		Name:                  cr.Name,
		DistributionNamespace: cr.Spec.Namespace,
		Version:               version,
		Status:                cr.Status.Status,
		Health:                ecosystem.HealthStatus(cr.Status.Health),
		PersistenceContext:    persistenceContext,
	}, nil
}
