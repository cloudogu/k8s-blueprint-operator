package healthconfig

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/blueprint-lib/v2"
	"time"

	"k8s.io/api/core/v1"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

const (
	healthConfigMapName      = "k8s-blueprint-operator-health-config"
	componentHealthConfigKey = "components"
	waitHealthConfigKey      = "wait"
)

var defaultHealthConfig = healthConfig{
	Components: componentHealthConfig{
		Required: nil,
	},
	Wait: waitHealthConfig{
		Timeout:  duration{10 * time.Minute},
		Interval: duration{10 * time.Second},
	},
}

type HealthConfigProvider struct {
	cmClient configMapInterface
}

func NewHealthConfigProvider(cmClient corev1.ConfigMapInterface) *HealthConfigProvider {
	return &HealthConfigProvider{cmClient: cmClient}
}

func (h *HealthConfigProvider) GetWaitConfig(ctx context.Context) (ecosystem.WaitConfig, error) {
	config, err := h.getAll(ctx)
	if err != nil {
		return ecosystem.WaitConfig{}, err
	}

	return convertToWaitConfigDomain(config.Wait), nil
}

func convertToWaitConfigDomain(config waitHealthConfig) ecosystem.WaitConfig {
	return ecosystem.WaitConfig{
		Timeout:  config.Timeout.Duration,
		Interval: config.Interval.Duration,
	}
}

func (h *HealthConfigProvider) GetRequiredComponents(ctx context.Context) ([]ecosystem.RequiredComponent, error) {
	config, err := h.getAll(ctx)
	if err != nil {
		return nil, err
	}

	return util.Map(config.Components.Required, convertToRequiredComponentDomain), nil
}

func convertToRequiredComponentDomain(component requiredComponent) ecosystem.RequiredComponent {
	return ecosystem.RequiredComponent{Name: v2.SimpleComponentName(component.Name)}
}

func (h *HealthConfigProvider) getAll(ctx context.Context) (healthConfig, error) {
	configMap, err := h.cmClient.Get(ctx, healthConfigMapName, metav1.GetOptions{})
	if k8sErr.IsNotFound(err) {
		return defaultHealthConfig, nil
	} else if err != nil {
		return healthConfig{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("failed to get config map %q", healthConfigMapName),
		}
	}

	components, componentsErr := parseComponentConfig(configMap)
	wait, waitErr := parseWaitConfig(configMap)

	return healthConfig{
		Components: components,
		Wait:       wait,
	}, errors.Join(componentsErr, waitErr)
}

func parseWaitConfig(configMap *v1.ConfigMap) (waitHealthConfig, error) {
	waitConfigStr, exists := configMap.Data[waitHealthConfigKey]
	if !exists {
		return defaultHealthConfig.Wait, nil
	}

	var wait waitHealthConfig
	err := yaml.Unmarshal([]byte(waitConfigStr), &wait)
	if err != nil {
		return waitHealthConfig{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to parse wait health config",
		}
	}

	if wait.Interval.Duration == 0 {
		wait.Interval = defaultHealthConfig.Wait.Interval
	}

	if wait.Timeout.Duration == 0 {
		wait.Timeout = defaultHealthConfig.Wait.Timeout
	}

	return wait, nil
}

func parseComponentConfig(configMap *v1.ConfigMap) (componentHealthConfig, error) {
	componentHealthConfigStr, exists := configMap.Data[componentHealthConfigKey]
	if !exists {
		return defaultHealthConfig.Components, nil
	}

	var components componentHealthConfig
	err := yaml.Unmarshal([]byte(componentHealthConfigStr), &components)
	if err != nil {
		return componentHealthConfig{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to parse component health config",
		}
	}
	return components, nil
}
