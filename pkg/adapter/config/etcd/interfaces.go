package etcd

import (
	"github.com/cloudogu/cesapp-lib/registry"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type etcdStore interface {
	registry.Registry
}
type globalConfigStore interface {
	registry.ConfigurationContext
}

type secret interface {
	v1.SecretInterface
}

type etcdRegistry interface {
	registry.Registry
}

type configurationContext interface {
	registry.ConfigurationContext
}
