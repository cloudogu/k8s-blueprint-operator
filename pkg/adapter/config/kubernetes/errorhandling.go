package kubernetes

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	liberrors "github.com/cloudogu/k8s-registry-lib/errors"
)

func mapToBlueprintError(err error) error {
	if err != nil {
		if liberrors.IsNotFoundError(err) {
			return domainservice.NewNotFoundError(err, "could not find config. Check if your ecosystem is ready for operation")
		} else if liberrors.IsConflictError(err) {
			return domainservice.NewConflictError(err, "could not update config due to conflicting changes")
		} else if liberrors.IsConnectionError(err) {
			return domainservice.NewInternalError(err, "could not load/update config due to connection problems")
		} else if liberrors.IsAlreadyExistsError(err) {
			return domainservice.NewConflictError(err, "could not create config as it already exists")
		} else {
			// GenericError and fallback if even that would not match the error
			return domainservice.NewInternalError(err, "could not load/update config due to an unknown problem")
		}
	}
	return nil
}
