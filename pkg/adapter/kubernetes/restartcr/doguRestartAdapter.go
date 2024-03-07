package restartcr

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

type doguRestartAdapter struct {
	// TODO: Import restartClient from dogu operator
}

// TODO: Generate a constructor

func (d doguRestartAdapter) RestartAll(ctx context.Context, names []common.QualifiedDoguName) error {
	//TODO implement me
	// restart all dogus with the restart client of the dogu operator
	panic("implement me")
}
