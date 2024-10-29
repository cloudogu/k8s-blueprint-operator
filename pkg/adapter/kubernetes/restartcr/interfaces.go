package restartcr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-operator/v2/api/ecoSystem"
)

type DoguRestartInterface interface {
	ecosystemclient.DoguRestartInterface
}
