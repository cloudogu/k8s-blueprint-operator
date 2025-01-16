package domain

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"slices"

	bpv2 "github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

type configValidator struct {
	config bpv2.Config
}

func newConfigValidator(config bpv2.Config) *configValidator {
	return &configValidator{config: config}
}

func (configValidator *configValidator) validate() error {
	var errs []error
	for doguName, doguConfig := range configValidator.config.Dogus {
		if doguName != doguConfig.DoguName {
			errs = append(errs, fmt.Errorf("dogu name %q in map and dogu name %q in value are not equal", doguName, doguConfig.DoguName))
		}
		newConfigValidator()
		errs = append(errs, doguConfig.validate())
	}
	errs = append(errs, configValidator.config.Global.validate())

	return errors.Join(errs...)
}

type combinedDoguConfigValidator struct {
}

func newCombinedDoguConfigValidator() *combinedDoguConfigValidator {
	return &combinedDoguConfigValidator{}
}

func (configValidator *configValidator) validateCombined(config bpv2.Config) error {
	err := errors.Join(
		configValidator.config.Config.validate(config.DoguName),
		configValidator.config.SensitiveConfig.validate(config.DoguName),
		configValidator.config.validateConflictingConfigKeys(),
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

func (config bpv2.DoguConfig) validate(referencedDoguName cescommons.SimpleName) error {
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
