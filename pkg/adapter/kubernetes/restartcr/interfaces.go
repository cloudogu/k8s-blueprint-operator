package restartcr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-operator/v3/api/ecoSystem"
)

type DoguRestartInterface interface {
	ecosystemclient.DoguRestartInterface
}
