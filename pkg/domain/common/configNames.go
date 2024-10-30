package common

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
)

type RegistryConfigKey interface {
	GlobalConfigKey | DoguConfigKey
}

type GlobalConfigKey = config.Key

type DoguConfigKey struct {
	DoguName cescommons.SimpleDoguName
	Key      config.Key
}

func (k DoguConfigKey) Validate() error {
	var errs []error
	if k.DoguName == "" {
		errs = append(errs, fmt.Errorf("dogu name for dogu config key %q should not be empty", k.Key))
	}
	if string(k.Key) == "" {
		errs = append(errs, fmt.Errorf("key for dogu config of dogu %q should not be empty", k.DoguName))
	}

	return errors.Join(errs...)
}

func (k DoguConfigKey) String() string {
	return fmt.Sprintf("key %q of dogu %q", k.Key, k.DoguName)
}

type SensitiveDoguConfigKey = DoguConfigKey

// GlobalConfigValue is a single global config value
type GlobalConfigValue = config.Value

// DoguConfigValue  is a single dogu config value, which is no sensitive configuration
type DoguConfigValue = config.Value

// SensitiveDoguConfigValue is a single unencrypted sensitive dogu config value
type SensitiveDoguConfigValue = config.Value
