package config

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

const healthConfigMapName = "k8s-blueprint-operator-health-config"
const componentHealthConfigKey = "components"

type HealthConfigRepository struct {
	cmClient configMapInterface
}

func NewHealthConfigRepository(cmClient corev1.ConfigMapInterface) *HealthConfigRepository {
	return &HealthConfigRepository{cmClient: cmClient}
}

func (h *HealthConfigRepository) Get(ctx context.Context) (domain.HealthConfig, error) {
	configMap, err := h.cmClient.Get(ctx, healthConfigMapName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return domain.HealthConfig{}, &domainservice.NotFoundError{
			WrappedError: err,
			Message:      fmt.Sprintf("could not find health config map %q", healthConfigMapName),
		}
	} else if err != nil {
		return domain.HealthConfig{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("failed to get config map %q", healthConfigMapName),
		}
	}

	componentHealthConfigStr, exists := configMap.Data[componentHealthConfigKey]
	if !exists {
		return domain.HealthConfig{}, &domainservice.NotFoundError{
			WrappedError: nil,
			Message:      fmt.Sprintf("could not find component health config in config map %q", healthConfigMapName),
		}
	}

	var componentHealthConfig domain.ComponentHealthConfig
	err = yaml.Unmarshal([]byte(componentHealthConfigStr), &componentHealthConfig)
	if err != nil {
		return domain.HealthConfig{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to parse component health config",
		}
	}

	return domain.HealthConfig{
		ComponentHealthConfig: componentHealthConfig,
	}, nil
}
