package debugmodecr

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const debugModeSingletonCRName = "debug-mode"

type debugModeRepo struct {
	debugModeClient DebugModeInterface
}

// NewDebugModeRepo returns a new debugModeRepo to interact with the debug mode CR.
func NewDebugModeRepo(debugModeClient DebugModeInterface) domainservice.DebugModeRepository {
	return &debugModeRepo{debugModeClient: debugModeClient}
}

func (repo *debugModeRepo) GetSingleton(ctx context.Context) (*ecosystem.DebugMode, error) {
	cr, err := repo.debugModeClient.Get(ctx, debugModeSingletonCRName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, &domainservice.NotFoundError{
				WrappedError: err,
				Message:      fmt.Sprintf("cannot load debug mode CR %q as it does not exist", debugModeSingletonCRName),
			}
		}
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("error while loading debug mode CR %q", debugModeSingletonCRName),
		}
	}
	return parseDebugModeCR(cr)
}
