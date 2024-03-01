package dogucr

import (
	ecosystemclient "github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// interface replication for generating mocks

//nolint:unused
type DoguInterface interface {
	ecosystemclient.DoguInterface
}

type PvcInterface interface {
	v1.PersistentVolumeClaimInterface
}
