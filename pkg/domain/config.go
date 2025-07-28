package domain

import (
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
	"golang.org/x/exp/maps"
	"slices"
)

type Config struct {
	Dogus  map[cescommons.SimpleName]CombinedDoguConfig
	Global GlobalConfig
}

type CombinedDoguConfig struct {
	DoguName        cescommons.SimpleName
	Config          DoguConfig
	SensitiveConfig SensitiveDoguConfig
}

type DoguConfig struct {
	Present map[common.DoguConfigKey]common.DoguConfigValue
	Absent  []common.DoguConfigKey
}

type SensitiveDoguConfig struct {
	Present map[common.DoguConfigKey]SensitiveValueRef
	Absent  []common.DoguConfigKey
}

type SensitiveValueRef struct {
	// SecretName is the name of the secret, from which the config key should be loaded.
	// The secret must be in the same namespace.
	SecretName string `json:"secretName"`
	// SecretKey is the name of the key within the secret given by SecretName.
	// The value is used as the value for the sensitive config key.
	SecretKey string `json:"secretKey"`
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

func (config Config) GetSensitiveConfigReferences() map[common.SensitiveDoguConfigKey]SensitiveValueRef {
	refs := map[common.SensitiveDoguConfigKey]SensitiveValueRef{}
	for _, doguConfig := range config.Dogus {
		for key, ref := range doguConfig.SensitiveConfig.Present {
			refs[key] = ref
		}
	}
	return refs
}

func (config Config) GetSensitiveDoguConfigKeys() []common.SensitiveDoguConfigKey {
	var keys []common.SensitiveDoguConfigKey
	for _, doguConfig := range config.Dogus {
		keys = append(keys, maps.Keys(doguConfig.SensitiveConfig.Present)...)
		keys = append(keys, doguConfig.SensitiveConfig.Absent...)
	}
	return keys
}

// GetDogusWithChangedConfig returns a list of dogus for which possible config changes are needed.
func (config Config) GetDogusWithChangedConfig() []cescommons.SimpleName {
	var dogus []cescommons.SimpleName
	for dogu, doguConfig := range config.Dogus {
		if len(doguConfig.Config.Present) != 0 || len(doguConfig.Config.Absent) != 0 {
			dogus = append(dogus, dogu)
		}
	}
	return dogus
}

// GetDogusWithChangedSensitiveConfig returns a list of dogus for which possible sensitive config changes are needed.
func (config Config) GetDogusWithChangedSensitiveConfig() []cescommons.SimpleName {
	var dogus []cescommons.SimpleName
	for dogu, doguConfig := range config.Dogus {
		if len(doguConfig.SensitiveConfig.Present) != 0 || len(doguConfig.SensitiveConfig.Absent) != 0 {
			dogus = append(dogus, dogu)
		}
	}
	return dogus
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
		//TODO: check if this mapping is still necessary
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

func validateDoguConfigKeys(keys []common.DoguConfigKey, referencedDoguName cescommons.SimpleName) error {
	var errs []error
	for _, configKey := range keys {
		err := configKey.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("dogu config key invalid: %w", err))
		}

		// validate that all keys are of the same dogu
		if referencedDoguName != configKey.DoguName {
			errs = append(errs, fmt.Errorf("key %s does not match superordinate dogu name %q", configKey, referencedDoguName))
		}
	}
	return errors.Join(errs...)
}

func validateNoDuplicates(presentKeys []common.DoguConfigKey, absentKeys []common.DoguConfigKey) error {
	var errs []error
	// no present keys in absent
	for _, presentKey := range presentKeys {
		if slices.Contains(absentKeys, presentKey) {
			errs = append(errs, fmt.Errorf("key %s cannot be present and absent at the same time", presentKey))
		}
	}
	// no absent keys in present
	for _, absentKey := range absentKeys {
		if slices.Contains(presentKeys, absentKey) {
			errs = append(errs, fmt.Errorf("key %s cannot be present and absent at the same time", absentKey))
		}
	}

	// no absent duplicates
	absentDuplicates := util.GetDuplicates(absentKeys)
	if len(absentDuplicates) > 0 {
		errs = append(errs, fmt.Errorf("absent dogu config should not contain duplicate keys: %v", absentDuplicates))
	}

	// present keys cannot have duplicates because of their data structure

	return errors.Join(errs...)
}

func (config DoguConfig) validate(referencedDoguName cescommons.SimpleName) error {
	var errs []error

	presentKeyErr := validateDoguConfigKeys(maps.Keys(config.Present), referencedDoguName)
	if presentKeyErr != nil {
		errs = append(errs, fmt.Errorf("dogu config is invalid: %w", presentKeyErr))
	}

	absentKeyErr := validateDoguConfigKeys(config.Absent, referencedDoguName)
	if absentKeyErr != nil {
		errs = append(errs, fmt.Errorf("absent dogu config is invalid: %w", presentKeyErr))
	}

	duplicatesErr := validateNoDuplicates(maps.Keys(config.Present), config.Absent)
	if duplicatesErr != nil {
		errs = append(errs, fmt.Errorf("absent dogu config is invalid: %w", presentKeyErr))
	}

	return errors.Join(errs...)
}

func (config SensitiveDoguConfig) validate(referencedDoguName cescommons.SimpleName) error {
	var errs []error

	presentKeyErr := validateDoguConfigKeys(maps.Keys(config.Present), referencedDoguName)
	if presentKeyErr != nil {
		errs = append(errs, fmt.Errorf("sensitive dogu config is invalid: %w", presentKeyErr))
	}

	absentKeyErr := validateDoguConfigKeys(config.Absent, referencedDoguName)
	if absentKeyErr != nil {
		errs = append(errs, fmt.Errorf("sensitive absent dogu config is invalid: %w", presentKeyErr))
	}

	duplicatesErr := validateNoDuplicates(maps.Keys(config.Present), config.Absent)
	if duplicatesErr != nil {
		errs = append(errs, fmt.Errorf("dogu config is invalid: %w", presentKeyErr))
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
