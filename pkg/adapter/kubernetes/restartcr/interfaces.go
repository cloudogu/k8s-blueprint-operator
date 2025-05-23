package restartcr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-lib/v2/client"
)

type DoguRestartInterface interface {
	ecosystemclient.DoguRestartInterface
}
