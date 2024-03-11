package restartcr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
)

type DoguRestartInterface interface {
	ecosystemclient.DoguRestartInterface
}
