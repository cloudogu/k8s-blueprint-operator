package config

import (
	"github.com/cloudogu/cesapp-lib/registry"
)

//nolint:unused
//goland:noinspection GoUnusedType
type globalConfigStore interface {
	registry.ConfigurationContext
}

//nolint:unused
//goland:noinspection GoUnusedType
type configurationContext interface {
	registry.ConfigurationContext
}
