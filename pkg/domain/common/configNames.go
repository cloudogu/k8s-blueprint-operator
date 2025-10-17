package common

import (
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type GlobalConfigKey = config.Key

type DoguConfigKey struct {
	DoguName cescommons.SimpleName
	Key      config.Key
}

func (k DoguConfigKey) String() string {
	return fmt.Sprintf("key %q of dogu %q", k.Key, k.DoguName)
}

// GlobalConfigValue is a single global config value
type GlobalConfigValue = config.Value

// DoguConfigValue  is a single dogu config value, which is no sensitive configuration
type DoguConfigValue = config.Value

// SensitiveDoguConfigValue is a single unencrypted sensitive dogu config value
type SensitiveDoguConfigValue = config.Value
