package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

type Config struct {
	Dogus  map[common.SimpleDoguName]CombinedDoguConfig
	Global GlobalConfig
}

type CombinedDoguConfig struct {
	DoguName        common.SimpleDoguName
	Config          DoguConfig
	SensitiveConfig SensitiveDoguConfig
}

type DoguConfig struct {
	Present map[common.DoguConfigKey]common.DoguConfigValue
	Absent  []common.DoguConfigKey
}

type SensitiveDoguConfig struct {
	Present map[common.SensitiveDoguConfigKey]common.SensitiveDoguConfigValue
	Absent  []common.SensitiveDoguConfigKey
}

type GlobalConfig struct {
	Present map[common.GlobalConfigKey]common.GlobalConfigValue
	Absent  []common.GlobalConfigKey
}

func (config CombinedDoguConfig) validate() error {
	return nil //TODO
}

func (config DoguConfig) validate(doguName common.SimpleDoguName) error {
	var errs []error
	for configKey, _ := range config.Present {
		errs = append(errs, configKey.Validate())

		// validate that all keys are of the same dogu
		if doguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("dogu ")) // TODO: better error message
		}
	}

	for _, configKey := range config.Absent {
		errs = append(errs, configKey.Validate())

		// absent keys cannot be present
		_, isPresent := config.Present[configKey]
		if isPresent {
			errs = append(errs, fmt.Errorf("%q cannot be present and absent at the same time", configKey))
		}

		// validate that all keys are of the same dogu
		if doguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("dogu ")) // TODO: better error message
		}
	}

	return errors.Join(errs...)
}

func (config GlobalConfig) validate() error {
	var errs []error
	for configKey, _ := range config.Present {

		// empty value is allowed
		if string(configKey) == "" {
			errs = append(errs, fmt.Errorf("key for present global config should not be empty"))
		}
	}
	for _, configKey := range config.Absent {

		// empty value is allowed
		if string(configKey) == "" {
			errs = append(errs, fmt.Errorf("key for absent global config should not be empty"))
		}

		// absent keys cannot be present
		_, isPresent := config.Present[configKey]
		if isPresent {
			errs = append(errs, fmt.Errorf("global config key %q cannot be present and absent at the same time", configKey))
		}
	}

	return errors.Join(errs...)
}
