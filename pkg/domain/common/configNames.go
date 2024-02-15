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
		errs = append(errs, fmt.Errorf("dogu name for present dogu config should not be empty"))
	}
	if string(k.Key) == "" {
		errs = append(errs, fmt.Errorf("key for present dogu config should not be empty"))
	}

	return errors.Join(errs...)
}

func (k DoguConfigKey) String() string {
	return fmt.Sprintf("key %q of dogu %q", k.Key, k.DoguName)
}

type SensitiveDoguConfigKey DoguConfigKey
type GlobalConfigValue string
type DoguConfigValue string
type SensitiveDoguConfigValue string
type EncryptedDoguConfigValue string
