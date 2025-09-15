package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	libconfig "github.com/cloudogu/k8s-registry-lib/config"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToConfigDTO(config *domain.Config) *v2.Config {
	if config == nil {
		return nil
	}

	var dogus map[string][]v2.ConfigEntry
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[string][]v2.ConfigEntry, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[string(doguName)] = convertToDoguConfigDTO(doguConfig)
		}
	}

	return &v2.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config *v2.Config) *domain.Config {
	if config == nil {
		return nil
	}
	var dogus map[cescommons.SimpleName]domain.DoguConfigEntries
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[cescommons.SimpleName]domain.DoguConfigEntries, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[cescommons.SimpleName(doguName)] = convertToDoguConfigEntriesDomain(doguConfig)
		}
	}

	return &domain.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDomain(config.Global),
	}
}

func convertToDoguConfigDTO(config domain.DoguConfigEntries) []v2.ConfigEntry {
	return convertToConfigEntriesDTO(domain.ConfigEntries(config))
}

func convertToDoguConfigEntriesDomain(config []v2.ConfigEntry) domain.DoguConfigEntries {
	return domain.DoguConfigEntries(convertToConfigEntriesDomain(config))
}

func convertToGlobalConfigDTO(config domain.GlobalConfigEntries) []v2.ConfigEntry {
	return convertToConfigEntriesDTO(domain.ConfigEntries(config))
}

func convertToConfigEntriesDTO(config domain.ConfigEntries) []v2.ConfigEntry {
	if config == nil || len(config) == 0 {
		return nil
	}

	result := make([]v2.ConfigEntry, len(config))

	for i, domainEntry := range config {
		var absent *bool
		if domainEntry.Absent {
			absent = &domainEntry.Absent
		}

		var sensitive *bool
		if domainEntry.Sensitive {
			sensitive = &domainEntry.Sensitive
		}

		var secretRef *v2.SecretReference
		if domainEntry.SecretRef != nil {
			secretRef = &v2.SecretReference{
				Name: domainEntry.SecretRef.SecretName,
				Key:  domainEntry.SecretRef.SecretKey,
			}
		}

		result[i] = v2.ConfigEntry{
			Key:       domainEntry.Key.String(),
			Absent:    absent,
			Value:     (*string)(domainEntry.Value),
			Sensitive: sensitive,
			SecretRef: secretRef,
		}
	}

	return result
}

func convertToGlobalConfigDomain(config []v2.ConfigEntry) domain.GlobalConfigEntries {
	return domain.GlobalConfigEntries(convertToConfigEntriesDomain(config))
}

func convertToConfigEntriesDomain(config []v2.ConfigEntry) domain.ConfigEntries {
	if config == nil || len(config) == 0 {
		return nil
	}

	result := make([]domain.ConfigEntry, len(config))

	for i, v2Entry := range config {
		absent := false
		if v2Entry.Absent != nil {
			absent = *v2Entry.Absent
		}

		sensitive := false
		if v2Entry.Sensitive != nil {
			sensitive = *v2Entry.Sensitive
		}

		var secretRef *domain.SensitiveValueRef
		if v2Entry.SecretRef != nil {
			secretRef = &domain.SensitiveValueRef{
				SecretName: v2Entry.SecretRef.Name,
				SecretKey:  v2Entry.SecretRef.Key,
			}
		}

		result[i] = domain.ConfigEntry{
			Key:       libconfig.Key(v2Entry.Key),
			Absent:    absent,
			Value:     (*libconfig.Value)(v2Entry.Value),
			Sensitive: sensitive,
			SecretRef: secretRef,
		}
	}

	return result
}
