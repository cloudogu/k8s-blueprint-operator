package dogucr

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
	ecosystemclient "github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

const doguInstallationRepoContextKey = "doguInstallationRepoContext"

type doguInstallationRepoContext struct {
	resourceVersion string
}

type doguInstallationRepo struct {
	doguClient ecosystemclient.DoguInterface
}

// NewDoguInstallationRepo returns a new doguInstallationRepo to interact on BlueprintSpecs.
func NewDoguInstallationRepo(doguClient ecosystemclient.DoguInterface) domainservice.DoguInstallationRepository {
	return &doguInstallationRepo{doguClient: doguClient}
}
func (repo *doguInstallationRepo) GetByName(ctx context.Context, doguName string) (*ecosystem.DoguInstallation, error) {
	cr, err := repo.doguClient.Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot load dogu CR %q as it does not exist", doguName),
			}
		}
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading dogu CR %q", doguName),
		}
	}

	return parseDoguCR(cr)
}

func parseDoguCR(cr *v1.Dogu) (*ecosystem.DoguInstallation, error) {
	if cr == nil {
		return nil, &domainservice.InternalError{
			WrappedError: nil,
			Message:      "Cannot parse dogu CR as it is nil",
		}
	}
	// parse dogu fields
	version, versionErr := core.ParseVersion(cr.Spec.Version)
	namespace, _, nameErr := serializer.SplitDoguName(cr.Spec.Name)
	err := errors.Join(versionErr, nameErr)
	if err != nil {
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      "Cannot load dogu CR as it cannot be parsed correctly",
		}
	}
	// parse persistence context
	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[doguInstallationRepoContextKey] = doguInstallationRepoContext{
		resourceVersion: cr.GetResourceVersion(),
	}
	return &ecosystem.DoguInstallation{
		Namespace:          namespace,
		Name:               cr.Name,
		Version:            version,
		Status:             cr.Status.Status,
		Health:             ecosystem.HealthStatus(cr.Status.Health),
		UpgradeConfig:      ecosystem.UpgradeConfig{AllowNamespaceSwitch: cr.Spec.UpgradeConfig.AllowNamespaceSwitch},
		PersistenceContext: persistenceContext,
	}, nil
}

func (repo *doguInstallationRepo) GetAll(ctx context.Context) (map[string]*ecosystem.DoguInstallation, error) {
	crList, err := repo.doguClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      "error while listing dogu CRs",
		}
	}

	var errs []error
	doguInstallations := make(map[string]*ecosystem.DoguInstallation, len(crList.Items))
	for _, cr := range crList.Items {
		doguInstallation, err := parseDoguCR(&cr)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		doguInstallations[doguInstallation.Name] = doguInstallation
	}

	err = errors.Join(errs...)
	if err != nil {
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to parse some dogu CRs",
		}
	}

	return doguInstallations, nil
}
