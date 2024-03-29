package restartcr

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type doguRestartRepository struct {
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
				GenerateName: fmt.Sprintf("%s-", string(doguName)),
			},
			Spec: v1.DoguRestartSpec{
				DoguName: string(doguName),
			},
		}, metav1.CreateOptions{})
		if err != nil {
			createErrors = append(createErrors, err)
		}
	}
	return errors.Join(createErrors...)
}
