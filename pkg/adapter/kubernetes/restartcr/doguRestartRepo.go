package restartcr

import (
	"context"
	"errors"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type doguRestartRepository struct {
	// TODO: Import restartClient from dogu operator
	restartInterface DoguRestartInterface
}

func NewDoguRestartRepository(restartInterface DoguRestartInterface) *doguRestartRepository {
	return &doguRestartRepository{restartInterface: restartInterface}
}

func (d doguRestartRepository) RestartAll(ctx context.Context, names []common.SimpleDoguName) error {
	var createErrors []error
	for _, doguName := range names {
		_, err := d.restartInterface.Create(ctx, &v1.DoguRestart{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(doguName), // TODO: name: doguname+blueprintID+randomID
			},
			Spec:   v1.DoguRestartSpec{},
			Status: v1.DoguRestartStatus{},
		}, metav1.CreateOptions{})
		if err != nil {
			createErrors = append(createErrors, err)
		}
	}
	return errors.Join(createErrors...)
}
