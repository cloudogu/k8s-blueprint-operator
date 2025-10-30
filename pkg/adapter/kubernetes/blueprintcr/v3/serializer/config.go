package serializer

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	libconfig "github.com/cloudogu/k8s-registry-lib/config"
	"k8s.io/utils/ptr"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func ConvertToConfigDTO(config domain.Config) *bpv3.Config {
	if config.IsEmpty() {
		return nil
	}

	var dogus map[string][]bpv3.ConfigEntry
	// we check for empty values to make good use of default values
	// this makes testing easier
	if len(config.Dogus) != 0 {
		dogus = make(map[string][]bpv3.ConfigEntry, len(config.Dogus))
		for doguName, doguConfig := range config.Dogus {
			dogus[string(doguName)] = convertToDoguConfigDTO(doguConfig)
		}
	}

	return &bpv3.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDTO(config.Global),
	}
}

func ConvertToConfigDomain(config *bpv3.Config) domain.Config {
	if config == nil {
		return domain.Config{}
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

	return domain.Config{
		Dogus:  dogus,
		Global: convertToGlobalConfigDomain(config.Global),
	}
}

func convertToDoguConfigDTO(config domain.DoguConfigEntries) []bpv3.ConfigEntry {
	return convertToConfigEntriesDTO(domain.ConfigEntries(config))
}

func convertToDoguConfigEntriesDomain(config []bpv3.ConfigEntry) domain.DoguConfigEntries {
	return domain.DoguConfigEntries(convertToConfigEntriesDomain(config))
}

func convertToGlobalConfigDTO(config domain.GlobalConfigEntries) []bpv3.ConfigEntry {
	return convertToConfigEntriesDTO(domain.ConfigEntries(config))
}

func convertToConfigEntriesDTO(config domain.ConfigEntries) []bpv3.ConfigEntry {
	if len(config) == 0 {
		return nil
	}

	result := make([]bpv3.ConfigEntry, len(config))

	for i, domainEntry := range config {
		var absent *bool
		if domainEntry.Absent {
			absent = &domainEntry.Absent
		}

		var sensitive *bool
		var value *string
		if domainEntry.Sensitive {
			sensitive = &domainEntry.Sensitive
		} else {
			// only set value if not sensitive
			value = (*string)(domainEntry.Value)
		}

		var secretRef *bpv3.Reference
		if domainEntry.SecretRef != nil {
			secretRef = &bpv3.Reference{
				Name: domainEntry.SecretRef.SecretName,
				Key:  domainEntry.SecretRef.SecretKey,
			}
		}
		var configRef *bpv3.Reference
		if domainEntry.ConfigRef != nil {
			configRef = &bpv3.Reference{
				Name: domainEntry.ConfigRef.ConfigMapName,
				Key:  domainEntry.ConfigRef.ConfigMapKey,
			}
		}

		result[i] = bpv3.ConfigEntry{
			Key:       domainEntry.Key.String(),
			Absent:    absent,
			Value:     value,
			Sensitive: sensitive,
			SecretRef: secretRef,
			ConfigRef: configRef,
		}
	}

	return result
}

func convertToGlobalConfigDomain(config []bpv3.ConfigEntry) domain.GlobalConfigEntries {
	return domain.GlobalConfigEntries(convertToConfigEntriesDomain(config))
}

func convertToConfigEntriesDomain(config []bpv3.ConfigEntry) domain.ConfigEntries {
	if len(config) == 0 {
		return nil
	}

	result := make([]domain.ConfigEntry, len(config))

	for i, v2Entry := range config {
		var secretRef *domain.SensitiveValueRef
		if v2Entry.SecretRef != nil {
			secretRef = &domain.SensitiveValueRef{
				SecretName: v2Entry.SecretRef.Name,
				SecretKey:  v2Entry.SecretRef.Key,
			}
		}

		var configRef *domain.ConfigValueRef
		if v2Entry.ConfigRef != nil {
			configRef = &domain.ConfigValueRef{
				ConfigMapName: v2Entry.ConfigRef.Name,
				ConfigMapKey:  v2Entry.ConfigRef.Key,
			}
		}

		result[i] = domain.ConfigEntry{
			Key:       libconfig.Key(v2Entry.Key),
			Absent:    ptr.Deref(v2Entry.Absent, false),
			Value:     (*libconfig.Value)(v2Entry.Value),
			Sensitive: ptr.Deref(v2Entry.Sensitive, false),
			SecretRef: secretRef,
			ConfigRef: configRef,
		}
	}

	return result
}
