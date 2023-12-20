package dogucr

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	ecosystemclient "github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
func (repo *doguInstallationRepo) GetByName(ctx context.Context, doguName string) (ecosystem.DoguInstallation, error) {
	cr, err := repo.doguClient.Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return ecosystem.DoguInstallation{}, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot load dogu CR %q as it does not exist", doguName),
			}
		}
		return ecosystem.DoguInstallation{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading dogu CR %q", doguName),
		}
	}

	return parseDoguCR(cr)
}

func parseDoguCR(cr *v1.Dogu) (ecosystem.DoguInstallation, error) {
	if cr == nil {
		return ecosystem.DoguInstallation{}, &domainservice.InternalError{
			WrappedError: nil,
			Message:      "Cannot parse dogu CR as it is nil",
		}
	}
	// parse dogu fields
	version, versionErr := core.ParseVersion(cr.Spec.Version)
	namespace, _, nameErr := serializer.SplitDoguName(cr.Spec.Name)
	err := errors.Join(versionErr, nameErr)
	if err != nil {
		return ecosystem.DoguInstallation{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      "Cannot load dogu CR as it cannot be parsed correctly",
		}
	}
	// parse persistence context
	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[doguInstallationRepoContextKey] = doguInstallationRepoContext{
		resourceVersion: cr.GetResourceVersion(),
	}
	return ecosystem.DoguInstallation{
		Namespace:          namespace,
		Name:               cr.Name,
		Version:            version,
		Status:             cr.Status.Status,
		Health:             ecosystem.HealthStatus(cr.Status.Health),
		UpgradeConfig:      ecosystem.UpgradeConfig{AllowNamespaceSwitch: cr.Spec.UpgradeConfig.AllowNamespaceSwitch},
		PersistenceContext: persistenceContext,
	}, nil
}

// getPersistenceContext reads the repo-specific resourceVersion from the domain.BlueprintSpec or returns an error.
func getPersistenceContext(ctx context.Context, spec ecosystem.DoguInstallation) (doguInstallationRepoContext, error) {
	logger := log.FromContext(ctx).WithName("doguInstallationRepo.Update")
	rawField, versionExists := spec.PersistenceContext[doguInstallationRepoContextKey]
	if versionExists {
		persistenceCtx, isContext := rawField.(doguInstallationRepoContext)
		if isContext {
			return persistenceCtx, nil
		} else {
			err := fmt.Errorf("persistence context in doguInstallation is not a 'doguInstallationRepoContext' but '%T'", rawField)
			logger.Error(err, "does this value come from a different repository?")
			return doguInstallationRepoContext{}, err
		}
	} else {
		err := errors.New("no doguInstallationRepoContext was provided over the persistenceContext in the given doguInstallation")
		logger.Error(err, "This is normally written while loading the doguInstallation over this repository. "+
			"Did you try to persist a new doguInstallation with repo.Update()?")
		return doguInstallationRepoContext{}, err
	}
}

func (repo *doguInstallationRepo) GetAllByName(ctx context.Context, doguNames []string) ([]ecosystem.DoguInstallation, error) {
	//TODO implement me
	panic("implement me")
}
