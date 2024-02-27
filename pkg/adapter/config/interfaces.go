package config

import (
	"github.com/cloudogu/cesapp-lib/registry"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type secret interface {
	v1.SecretInterface
}

type etcdRegistry interface {
	registry.Registry
}

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
