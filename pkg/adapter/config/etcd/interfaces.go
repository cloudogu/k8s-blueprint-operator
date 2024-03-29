package etcd

import (
	"github.com/cloudogu/cesapp-lib/registry"
)

type etcdStore interface {
	registry.Registry
}
type globalConfigStore interface {
	registry.ConfigurationContext
}

type configurationContext interface {
	registry.ConfigurationContext
}
