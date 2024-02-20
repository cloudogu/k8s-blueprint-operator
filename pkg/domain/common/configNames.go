package common

import (
	"errors"
	"fmt"
)

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

type GlobalConfigValue string
type DoguConfigValue string
type SensitiveDoguConfigValue string
type EncryptedDoguConfigValue string
