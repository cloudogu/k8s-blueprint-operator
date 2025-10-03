package debugmodecr

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
)

func parseDebugModeCR(cr *v1.DebugMode) (*ecosystem.DebugMode, error) {
	if cr == nil {
		return nil, &domainservice.InternalError{
			WrappedError: nil,
			Message:      "cannot parse debug mode CR as it is nil",
		}
	}
	// parse debug mode fields
	return &ecosystem.DebugMode{Phase: string(cr.Status.Phase)}, nil
}
