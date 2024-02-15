package config

import (
	"github.com/cloudogu/cesapp-lib/registry"
)

type etcdStore interface {
	registry.Registry
}
type globalConfigStore interface {
	registry.ConfigurationContext
}
