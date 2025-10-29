package domain

import (
	"errors"
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	libconfig "github.com/cloudogu/k8s-registry-lib/config"
)

type Config struct {
	Dogus  DoguConfig
	Global GlobalConfigEntries
}

type ConfigEntry struct {
	// Key is the configuration key name
	Key libconfig.Key
	// Absent indicates whether this key should be deleted (true) or set (false)
	Absent bool
	// Value is used for regular (non-sensitive) configuration entries
	// Mutually exclusive with SecretRef and ConfigRef
	Value *libconfig.Value
	// Sensitive indicates whether this config is sensitive and should be stored securely (true) or not (false)
	Sensitive bool
	// SecretRef is used for sensitive configuration entries
	// Mutually exclusive with Value and ConfigRef
	SecretRef *SensitiveValueRef
	// ConfigRef is used for configuration entries
	// Mutually exclusive with Value and SecretRef
	ConfigRef *ConfigValueRef
}

type DoguConfig map[cescommons.SimpleName]DoguConfigEntries
type DoguConfigEntries ConfigEntries
type GlobalConfigEntries ConfigEntries
type ConfigEntries []ConfigEntry

type SensitiveValueRef struct {
	// SecretName is the name of the secret, from which the config key should be loaded.
	// The secret must be in the same namespace.
	SecretName string `json:"secretName"`
	// SecretKey is the name of the key within the secret given by SecretName.
	// The value is used as the value for the sensitive config key.
	SecretKey string `json:"secretKey"`
}

type ConfigValueRef struct {
	// ConfigMapName is the name of the config map, from which the config key should be loaded.
	// The config map must be in the same namespace.
	ConfigMapName string `json:"configMapName"`
	// ConfigMapKey is the name of the key within the config map given by ConfigMapName.
	// The value is used as the value for the config key.
	ConfigMapKey string `json:"configMapKey"`
}

func (config GlobalConfigEntries) GetGlobalConfigKeys() []common.GlobalConfigKey {
	var keys []common.GlobalConfigKey
	for _, entry := range config {
		keys = append(keys, entry.Key)
	}
	return keys
}

func (config Config) GetDoguConfigKeys() []common.DoguConfigKey {
	return config.getDoguKeysBySensitivity(false)
}

func (config Config) GetSensitiveConfigReferences() map[common.DoguConfigKey]SensitiveValueRef {
	refs := map[common.DoguConfigKey]SensitiveValueRef{}
	for doguName, doguConfig := range config.Dogus {
		for _, entry := range doguConfig {
			if entry.SecretRef != nil {
				key := common.DoguConfigKey{
					DoguName: doguName,
					Key:      entry.Key,
				}
				refs[key] = *entry.SecretRef
			}
		}
	}
	return refs
}

func (config Config) GetConfigReferences() map[common.DoguConfigKey]ConfigValueRef {
	refs := map[common.DoguConfigKey]ConfigValueRef{}
	for doguName, doguConfig := range config.Dogus {
		for _, entry := range doguConfig {
			if entry.ConfigRef != nil {
				key := common.DoguConfigKey{
					DoguName: doguName,
					Key:      entry.Key,
				}
				refs[key] = *entry.ConfigRef
			}
		}
	}
	return refs
}

func (config Config) GetSensitiveGlobalConfigReferences() (map[common.GlobalConfigKey]SensitiveValueRef, error) {
	refs := map[common.GlobalConfigKey]SensitiveValueRef{}
	for _, entry := range config.Global {
		if entry.SecretRef != nil {
			refs[entry.Key] = *entry.SecretRef
		}
		if entry.Sensitive {
			return map[common.GlobalConfigKey]SensitiveValueRef{}, fmt.Errorf("global config values cannot be sensitive")
		}
	}
	return refs, nil
}

func (config Config) GetGlobalConfigReferences() map[common.GlobalConfigKey]ConfigValueRef {
	refs := map[common.GlobalConfigKey]ConfigValueRef{}
	for _, entry := range config.Global {
		if entry.ConfigRef != nil {
			refs[entry.Key] = *entry.ConfigRef
		}
	}
	return refs
}

func (config Config) GetSensitiveDoguConfigKeys() []common.DoguConfigKey {
	return config.getDoguKeysBySensitivity(true)
}

func (config Config) getDoguKeysBySensitivity(isSensitive bool) []common.DoguConfigKey {
	var keys []common.DoguConfigKey
	for doguName, doguConfig := range config.Dogus {
		for _, entry := range doguConfig {
			if entry.Sensitive == isSensitive {
				keys = append(keys, common.DoguConfigKey{
					DoguName: doguName,
					Key:      entry.Key,
				})
			}
		}
	}
	return keys
}

// GetDogusWithChangedConfig returns a list of dogus for which possible config changes are needed.
func (config Config) GetDogusWithChangedConfig() []cescommons.SimpleName {
	return config.getDogusWithChangedConfigBySenisitivity(false)
}

// GetDogusWithChangedSensitiveConfig returns a list of dogus for which possible sensitive config changes are needed.
func (config Config) GetDogusWithChangedSensitiveConfig() []cescommons.SimpleName {
	return config.getDogusWithChangedConfigBySenisitivity(true)
}

func (config Config) getDogusWithChangedConfigBySenisitivity(isSensitive bool) []cescommons.SimpleName {
	var dogus []cescommons.SimpleName
	for dogu, doguConfig := range config.Dogus {
		for _, entry := range doguConfig {
			if entry.Sensitive == isSensitive {
				dogus = append(dogus, dogu)
				break // we only need to know if this dogu has at least one desired config value, so we can break to next outer loop
			}
		}
	}
	return dogus
}

func (config Config) IsEmpty() bool {
	return len(config.Dogus) == 0 && len(config.Global) == 0
}

func (config Config) validate() error {
	var errs []error
	for doguName, doguConfig := range config.Dogus {
		if doguName == "" {
			errs = append(errs, fmt.Errorf("dogu name for dogu config should not be empty"))

		}
		errs = append(errs, doguConfig.validate(doguName))
	}
	errs = append(errs, config.Global.validate())

	return errors.Join(errs...)
}

func (config DoguConfigEntries) validate(doguName cescommons.SimpleName) error {
	var allErrs error
	for _, entry := range config {
		allErrs = errors.Join(allErrs, entry.validate())
	}
	allErrs = errors.Join(allErrs, ConfigEntries(config).validateConflictingConfigKeys())

	if allErrs != nil {
		return fmt.Errorf("config for dogu %q is invalid: %w", doguName, allErrs)
	}
	return nil
}

// validateConflictingConfigKeys checks that there are no conflicting keys in normal config and sensitive config.
// This is a problem as both config types are loaded via the same API in dogus at the moment.
func (config ConfigEntries) validateConflictingConfigKeys() error {
	var errorList []error
	seen := make(map[libconfig.Key]struct{})
	for _, entry := range config {
		if _, exists := seen[entry.Key]; exists {
			errorList = append(errorList, fmt.Errorf("duplicate dogu config Key found: %s", entry.Key))
		}
		seen[entry.Key] = struct{}{}
	}
	return errors.Join(errorList...)
}

// Validate ensures ConfigEntry has valid state
func (config ConfigEntry) validate() error {
	var errs []error

	if config.Key == "" {
		errs = append(errs, fmt.Errorf("key for config should not be empty"))
	}

	if config.Absent {
		if config.Value != nil || config.SecretRef != nil {
			errs = append(errs, fmt.Errorf("absent entries cannot have value or secretRef"))
		}
		return errors.Join(errs...)
	}

	// For present entries, exactly one of Value or SecretRef must be set
	hasValue := config.Value != nil
	hasSecretRef := config.SecretRef != nil

	if hasValue && hasSecretRef {
		errs = append(errs, fmt.Errorf("config entries can have either a value or a secretRef"))
	}

	if hasValue && config.Sensitive {
		errs = append(errs, fmt.Errorf("sensitive config entries are not allowed to have normal values"))
	}

	return errors.Join(errs...)
}

func (config GlobalConfigEntries) validate() error {
	var allErrs error
	for _, entry := range config {
		allErrs = errors.Join(allErrs, entry.validateGlobal())
	}
	allErrs = errors.Join(allErrs, ConfigEntries(config).validateConflictingConfigKeys())

	if allErrs != nil {
		return fmt.Errorf("global config is invalid: %w", allErrs)
	}
	return nil
}

func (config ConfigEntry) validateGlobal() error {
	var errs []error

	if config.Key == "" {
		errs = append(errs, fmt.Errorf("key for global config should not be empty"))
	}

	if config.Absent {
		if config.Value != nil {
			errs = append(errs, fmt.Errorf("absent entries cannot have value"))
		}
		return errors.Join(errs...)
	}

	return errors.Join(errs...)
}
