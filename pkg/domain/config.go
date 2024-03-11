package domain

import (
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
	"golang.org/x/exp/maps"
	"slices"
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

func (config GlobalConfig) GetGlobalConfigKeys() []common.GlobalConfigKey {
	var keys []common.GlobalConfigKey
	keys = append(keys, maps.Keys(config.Present)...)
	keys = append(keys, config.Absent...)
	return keys
}

func (config Config) GetDoguConfigKeys() []common.DoguConfigKey {
	var keys []common.DoguConfigKey
	for _, doguConfig := range config.Dogus {
		keys = append(keys, maps.Keys(doguConfig.Config.Present)...)
		keys = append(keys, doguConfig.Config.Absent...)
	}
	return keys
}

func (config Config) GetSensitiveDoguConfigKeys() []common.SensitiveDoguConfigKey {
	var keys []common.SensitiveDoguConfigKey
	for _, doguConfig := range config.Dogus {
		keys = append(keys, maps.Keys(doguConfig.SensitiveConfig.Present)...)
		keys = append(keys, doguConfig.SensitiveConfig.Absent...)
	}
	return keys
}

// censorValues censors all sensitive configuration data to make them unrecognisable.
func (config Config) censorValues() Config {
	for _, doguConfig := range config.Dogus {
		for k := range doguConfig.SensitiveConfig.Present {
			doguConfig.SensitiveConfig.Present[k] = censorValue
		}
	}
	return config
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
	err := errors.Join(
		config.Config.validate(config.DoguName),
		config.SensitiveConfig.validate(config.DoguName),
		config.validateConflictingConfigKeys(),
	)

	if err != nil {
		return fmt.Errorf("config for dogu %q is invalid: %w", config.DoguName, err)
	}
	return nil
}

// validateConflictingConfigKeys checks that there are no conflicting keys in normal config and sensitive config.
// This is a problem as both config types are loaded via the same API in dogus at the moment.
func (config CombinedDoguConfig) validateConflictingConfigKeys() error {
	var normalKeys []common.DoguConfigKey
	normalKeys = append(normalKeys, maps.Keys(config.Config.Present)...)
	normalKeys = append(normalKeys, config.Config.Absent...)
	var sensitiveKeys []common.SensitiveDoguConfigKey
	sensitiveKeys = append(sensitiveKeys, maps.Keys(config.SensitiveConfig.Present)...)
	sensitiveKeys = append(sensitiveKeys, config.SensitiveConfig.Absent...)

	var errorList []error

	for _, sensitiveKey := range sensitiveKeys {
		keyToSearch := common.DoguConfigKey{
			DoguName: sensitiveKey.DoguName,
			Key:      sensitiveKey.Key,
		}
		if slices.Contains(normalKeys, keyToSearch) {
			errorList = append(errorList, fmt.Errorf("dogu config key %s cannot be in normal and sensitive configuration at the same time", keyToSearch))
		}
	}
	return errors.Join(errorList...)
}

func (config DoguConfig) validate(referencedDoguName common.SimpleDoguName) error {
	var errs []error

	for configKey := range config.Present {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("present dogu config key invalid: %w", err))
		}

		// validate that all keys are of the same dogu
		if referencedDoguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("present %s does not match superordinate dogu name %q", configKey, referencedDoguName))
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
		if referencedDoguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("absent %s does not match superordinate dogu name %q", configKey, referencedDoguName))
		}
	}

	absentDuplicates := util.GetDuplicates(config.Absent)
	if len(absentDuplicates) > 0 {
		errs = append(errs, fmt.Errorf("absent dogu config should not contain duplicate keys: %v", absentDuplicates))
	}

	return errors.Join(errs...)
}

func (config SensitiveDoguConfig) validate(referencedDoguName common.SimpleDoguName) error {
	var errs []error
	for configKey := range config.Present {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("present dogu config key invalid: %w", err))
		}

		// validate that all keys are of the same dogu
		if referencedDoguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("present %s does not match superordinate dogu name %q", configKey, referencedDoguName))
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
		if referencedDoguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("absent %s does not match superordinate dogu name %q", configKey, referencedDoguName))
		}
	}

	absentDuplicates := util.GetDuplicates(config.Absent)
	if len(absentDuplicates) > 0 {
		errs = append(errs, fmt.Errorf("absent dogu config should not contain duplicate keys: %v", absentDuplicates))
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
			errs = append(errs, fmt.Errorf("global config key %q cannot be present and absent at the same time", configKey))
		}
	}

	absentDuplicates := util.GetDuplicates(config.Absent)
	if len(absentDuplicates) > 0 {
		errs = append(errs, fmt.Errorf("absent global config should not contain duplicate keys: %v", absentDuplicates))
	}

	return errors.Join(errs...)
}
