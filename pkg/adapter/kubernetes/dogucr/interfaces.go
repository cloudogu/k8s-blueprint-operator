package dogucr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-lib/v2/client"
	v2 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// interface replication for generating mocks

//nolint:unused
type DoguInterface interface {
	ecosystemclient.DoguInterface
}

type PvcInterface interface {
	v2.PersistentVolumeClaimInterface
}
