package configref

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapRefReader struct {
	configMapClient configMapClient
}

func NewConfigMapRefReader(configMapClient configMapClient) *ConfigMapRefReader {
	return &ConfigMapRefReader{
		configMapClient: configMapClient,
	}
}

func (reader *ConfigMapRefReader) GetValues(ctx context.Context, refs map[common.DoguConfigKey]domain.ConfigValueRef) (map[common.DoguConfigKey]common.DoguConfigValue, error) {
	configMapsByName, configMapErrors := reader.loadNeededConfigMaps(ctx, maps.Values(refs))
	config, keyErrors := reader.loadKeysFromConfigMaps(refs, configMapsByName)

	// combine errors so that the user gets info about not found configMaps and missing keys in existing configMaps
	err := errors.Join(configMapErrors, keyErrors)
	if err != nil {
		err = fmt.Errorf("could not load config via references: %w", err)
		return nil, err
	}

	return config, nil
}

func (reader *ConfigMapRefReader) loadKeysFromConfigMaps(
	refs map[common.DoguConfigKey]domain.ConfigValueRef,
	configMapsByName map[string]*v1.ConfigMap,
) (map[common.DoguConfigKey]common.DoguConfigValue, error) {
	var errs []error
	loadedConfig := map[common.DoguConfigKey]common.DoguConfigValue{}

	for configKey, ref := range refs {
		configMap, found := configMapsByName[ref.ConfigMapName]
		if !found {
			// no error here, because we already have an error for missing configMaps in the loadNeededConfigMaps function
			// we want error messages for missing keys too, even if a configMap does not exist
			continue
		}
		configValue, err := reader.loadKeyFromConfigMap(configMap, ref.ConfigMapKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		loadedConfig[configKey] = configValue
	}
	return loadedConfig, errors.Join(errs...)
}

func (reader *ConfigMapRefReader) loadKeyFromConfigMap(configMap *v1.ConfigMap, key string) (common.DoguConfigValue, error) {
	// do not use the StringData field, it is a write-only field in K8s
	value, exists := configMap.Data[key]
	if !exists {
		return "", domainservice.NewNotFoundError(
			nil,
			"referenced key %q in configMap %q does not exist", key, configMap.Name,
		)
	}
	return common.DoguConfigValue(value), nil
}

func (reader *ConfigMapRefReader) loadNeededConfigMaps(
	ctx context.Context,
	refs iter.Seq[domain.ConfigValueRef],
) (map[string]*v1.ConfigMap, error) {
	configMapsByName := map[string]*v1.ConfigMap{}
	var errs []error

	for ref := range refs {
		_, alreadyLoaded := configMapsByName[ref.ConfigMapName]
		if alreadyLoaded {
			continue
		}
		configMap, err := reader.loadConfigMap(ctx, ref.ConfigMapName)
		if err != nil {
			errs = append(errs, err)
		}
		// also save nil entries, so that we do not try to load this configMap again
		configMapsByName[ref.ConfigMapName] = configMap
	}
	// delete nil entries
	maps.DeleteFunc(configMapsByName, func(s string, configMap *v1.ConfigMap) bool {
		return configMap == nil
	})
	return configMapsByName, errors.Join(errs...)
}

func (reader *ConfigMapRefReader) loadConfigMap(ctx context.Context, name string) (*v1.ConfigMap, error) {
	configMap, err := reader.configMapClient.Get(ctx, name, metav1.GetOptions{})
	if configMap == nil || err != nil {
		return nil, domainservice.NewNotFoundError(
			err, "referenced configMap %q does not exist", name,
		)
	}
	return configMap, nil
}
