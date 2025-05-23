package dogucr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-lib/v2/client"
)

// interface replication for generating mocks

//nolint:unused
type DoguInterface interface {
	ecosystemclient.DoguInterface
}
