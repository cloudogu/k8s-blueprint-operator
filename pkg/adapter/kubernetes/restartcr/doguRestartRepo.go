package restartcr

import (
	"context"
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type doguRestartRepository struct {
	restartInterface DoguRestartInterface
}

func NewDoguRestartRepository(restartInterface DoguRestartInterface) *doguRestartRepository {
	return &doguRestartRepository{restartInterface: restartInterface}
}

func (d doguRestartRepository) RestartAll(ctx context.Context, names []cescommons.SimpleDoguName) error {
	var createErrors []error
	for _, doguName := range names {
		_, err := d.restartInterface.Create(ctx, &v2.DoguRestart{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("%s-", string(doguName)),
			},
			Spec: v2.DoguRestartSpec{
				DoguName: string(doguName),
			},
		}, metav1.CreateOptions{})
		if err != nil {
			createErrors = append(createErrors, err)
		}
	}
	return errors.Join(createErrors...)
}
