package common

import (
	"errors"
	"fmt"
)

type RegistryConfigKey interface {
	GlobalConfigKey | DoguConfigKey | SensitiveDoguConfigKey
}

type GlobalConfigKey string

type DoguConfigKey struct {
	DoguName SimpleDoguName
	Key      string
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

type SensitiveDoguConfigKey struct {
	DoguConfigKey
}

// GlobalConfigValue is a single global config value
type GlobalConfigValue string

// DoguConfigValue  is a single dogu config value, which is no sensitive configuration
type DoguConfigValue string

// SensitiveDoguConfigValue is a single unencrypted sensitive dogu config value
type SensitiveDoguConfigValue string

// EncryptedDoguConfigValue is a single encrypted sensitive dogu config value
type EncryptedDoguConfigValue string
