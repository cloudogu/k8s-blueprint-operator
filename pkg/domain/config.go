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

func (config Config) validate() error {
	var errs []error
	for doguName, doguConfig := range config.Dogus {
		if doguName != doguConfig.DoguName {
			errs = append(errs, fmt.Errorf("dogu name %q in map and dogu name %q in value are not equal", doguName, doguConfig.DoguName))
		}

		errs = append(errs, doguConfig.validate())
	}

	errs = append(errs, config.Global.validate())

	return errors.Join(errs...)
}

func (config CombinedDoguConfig) validate() error {
	var errs []error
	err := config.Config.validate(config.DoguName)
	if err != nil {
		errs = append(errs, fmt.Errorf("config for dogu %q invalid: %w", config.DoguName, err))
	}

	err = config.SensitiveConfig.validate(config.DoguName)
	if err != nil {
		errs = append(errs, fmt.Errorf("sensitive config for dogu %q invalid: %w", config.DoguName, err))
	}

	return nil
}

func (config DoguConfig) validate(doguName common.SimpleDoguName) error {
	var errs []error
	for configKey := range config.Present {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("present dogu config key invalid: %w", err))
		}

		// validate that all keys are of the same dogu
		if doguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("present %s does not match superordinate dogu name %q", configKey, doguName))
		}
	}

	for _, configKey := range config.Absent {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("absent dogu config key invalid: %w", err))
		}

		// absent keys cannot be present
		_, isPresent := config.Present[configKey]
		if isPresent {
			errs = append(errs, fmt.Errorf("%s cannot be present and absent at the same time", configKey))
		}

		// validate that all keys are of the same dogu
		if doguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("absent %s does not match superordinate dogu name %q", configKey, doguName))
		}
	}

	return errors.Join(errs...)
}

func (config SensitiveDoguConfig) validate(doguName common.SimpleDoguName) error {
	var errs []error
	for configKey := range config.Present {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("present dogu config key invalid: %w", err))
		}

		// validate that all keys are of the same dogu
		if doguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("present %s does not match superordinate dogu name %q", configKey, doguName))
		}
	}

	for _, configKey := range config.Absent {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("absent dogu config key invalid: %w", err))
		}

		// absent keys cannot be present
		_, isPresent := config.Present[configKey]
		if isPresent {
			errs = append(errs, fmt.Errorf("%s cannot be present and absent at the same time", configKey))
		}

		// validate that all keys are of the same dogu
		if doguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("absent %s does not match superordinate dogu name %q", configKey, doguName))
		}
	}

	return errors.Join(errs...)
}

func (config GlobalConfig) validate() error {
	var errs []error
	for configKey := range config.Present {

		// empty key is not allowed
		if string(configKey) == "" {
			errs = append(errs, fmt.Errorf("key for present global config should not be empty"))
		}
	}

	for _, configKey := range config.Absent {

		// empty key is not allowed
		if string(configKey) == "" {
			errs = append(errs, fmt.Errorf("key for absent global config should not be empty"))
		}

		// absent keys cannot be present
		_, isPresent := config.Present[configKey]
		if isPresent {
			errs = append(errs, fmt.Errorf("global config key %s cannot be present and absent at the same time", configKey))
		}
	}

	return errors.Join(errs...)
}
