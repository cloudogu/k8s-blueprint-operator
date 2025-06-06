package dogucr

import (
	"context"
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	ecosystemclient "github.com/cloudogu/k8s-dogu-lib/v2/client"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	doguInstallationRepoContextKey = "doguInstallationRepoContext"
)

type doguInstallationRepoContext struct {
	resourceVersion string
}

type doguInstallationRepo struct {
	doguClient DoguInterface
}

// NewDoguInstallationRepo returns a new doguInstallationRepo to interact on BlueprintSpecs.
func NewDoguInstallationRepo(doguClient ecosystemclient.DoguInterface) domainservice.DoguInstallationRepository {
	return &doguInstallationRepo{doguClient: doguClient}
}
func (repo *doguInstallationRepo) GetByName(ctx context.Context, doguName cescommons.SimpleName) (*ecosystem.DoguInstallation, error) {
	cr, err := repo.doguClient.Get(ctx, string(doguName), metav1.GetOptions{})
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

func (repo *doguInstallationRepo) GetAll(ctx context.Context) (map[cescommons.SimpleName]*ecosystem.DoguInstallation, error) {
	crList, err := repo.doguClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      "error while listing dogu CRs",
		}
	}

	var errs []error
	doguInstallations := make(map[cescommons.SimpleName]*ecosystem.DoguInstallation, len(crList.Items))
	for _, cr := range crList.Items {
		doguInstallation, err := parseDoguCR(&cr)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		doguInstallations[doguInstallation.Name.SimpleName] = doguInstallation
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

func (repo *doguInstallationRepo) Create(ctx context.Context, dogu *ecosystem.DoguInstallation) error {
	cr := toDoguCR(dogu)
	_, err := repo.doguClient.Create(ctx, cr, metav1.CreateOptions{})
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot create dogu CR for dogu %q", dogu.Name),
		}
	}
	return nil
}

func (repo *doguInstallationRepo) Update(ctx context.Context, dogu *ecosystem.DoguInstallation) error {
	logger := log.FromContext(ctx).WithName("doguInstallationRepo.Update")
	patch, err := toDoguCRPatchBytes(dogu)
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot create patch for dogu CR for dogu %q", dogu.Name),
		}
	}
	logger.Info("patch dogu CR", "doguName", dogu.Name, "doguPatch", string(patch))
	_, err = repo.doguClient.Patch(ctx, string(dogu.Name.SimpleName), types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot patch dogu CR for dogu %q", dogu.Name),
		}
	}
	return nil
}

func (repo *doguInstallationRepo) Delete(ctx context.Context, doguName cescommons.SimpleName) error {
	err := repo.doguClient.Delete(ctx, string(doguName), metav1.DeleteOptions{})
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot delete dogu CR for dogu %q", doguName),
		}
	}
	return nil
}
