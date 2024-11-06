package dogucr

import (
	"encoding/json"
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v2 "github.com/cloudogu/k8s-dogu-operator/v3/api/v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func parseDoguCR(cr *v2.Dogu) (*ecosystem.DoguInstallation, error) {
	if cr == nil {
		return nil, &domainservice.InternalError{
			WrappedError: nil,
			Message:      "cannot parse dogu CR as it is nil",
		}
	}
	// parse dogu fields
	version, versionErr := core.ParseVersion(cr.Spec.Version)
	doguName, nameErr := cescommons.QualifiedDoguNameFromString(cr.Spec.Name)

	volumeSize, volumeSizeErr := ecosystem.GetQuantityReference(cr.Spec.Resources.DataVolumeSize)

	reverseProxyConfigEntries, proxyErr := parseDoguAdditionalIngressAnnotationsCR(cr.Spec.AdditionalIngressAnnotations)

	err := errors.Join(versionErr, nameErr, volumeSizeErr, proxyErr)
	if err != nil {
		return nil, &domainservice.InternalError{
			WrappedError: err,
			Message:      "cannot load dogu CR as it cannot be parsed correctly",
		}
	}

	// parse persistence context
	persistenceContext := make(map[string]interface{}, 1)
	persistenceContext[doguInstallationRepoContextKey] = doguInstallationRepoContext{
		resourceVersion: cr.GetResourceVersion(),
	}
	return &ecosystem.DoguInstallation{
		Name:               doguName,
		Version:            version,
		Status:             cr.Status.Status,
		Health:             ecosystem.HealthStatus(cr.Status.Health),
		UpgradeConfig:      ecosystem.UpgradeConfig{AllowNamespaceSwitch: cr.Spec.UpgradeConfig.AllowNamespaceSwitch},
		MinVolumeSize:      volumeSize,
		ReverseProxyConfig: reverseProxyConfigEntries,
		PersistenceContext: persistenceContext,
	}, nil
}

func parseDoguAdditionalIngressAnnotationsCR(annotations v2.IngressAnnotations) (ecosystem.ReverseProxyConfig, error) {
	reverseProxyConfig := ecosystem.ReverseProxyConfig{}

	reverseProxyBodySize, ok := annotations[ecosystem.NginxIngressAnnotationBodySize]
	if ok {
		// Sizes for Nginx can be specified in bytes, kilobytes (suffixes k and K) or megabytes (suffixes m and M), for example, “1024”, “8k”, “1m” in Decimal SI.
		// Since the actual dogu-operator and service-discovery just use this format we can expect that the values for the volume size in are safe to set in the doguinstallation.
		// Formats “1024”, “8k”, “1m” can be parsed by resource.Quantity
		// See: [Documentation](https://nginx.org/en/docs/syntax.html)
		quantity, err := resource.ParseQuantity(reverseProxyBodySize)
		if err != nil {
			return ecosystem.ReverseProxyConfig{}, domainservice.NewInternalError(err, "failed to parse quantity %q", reverseProxyBodySize)
		}
		reverseProxyConfig.MaxBodySize = &quantity
	}

	reverseProxyConfig.RewriteTarget = ecosystem.RewriteTarget(annotations[ecosystem.NginxIngressAnnotationRewriteTarget])
	reverseProxyConfig.AdditionalConfig = ecosystem.AdditionalConfig(annotations[ecosystem.NginxIngressAnnotationAdditionalConfig])

	return reverseProxyConfig, nil
}

func toDoguCR(dogu *ecosystem.DoguInstallation) *v2.Dogu {
	return &v2.Dogu{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: string(dogu.Name.SimpleName),
			Labels: map[string]string{
				"app":       "ces",
				"dogu.name": string(dogu.Name.SimpleName),
			},
		},
		Spec: v2.DoguSpec{
			Name:    dogu.Name.String(),
			Version: dogu.Version.Raw,
			Resources: v2.DoguResources{
				DataVolumeSize: ecosystem.GetQuantityString(dogu.MinVolumeSize),
			},
			SupportMode: false,
			UpgradeConfig: v2.UpgradeConfig{
				AllowNamespaceSwitch: dogu.UpgradeConfig.AllowNamespaceSwitch,
				ForceUpgrade:         false,
			},
			AdditionalIngressAnnotations: getNginxIngressAnnotations(dogu.ReverseProxyConfig),
		},
		Status: v2.DoguStatus{},
	}
}

func getNginxIngressAnnotations(config ecosystem.ReverseProxyConfig) map[string]string {
	annotations := v2.IngressAnnotations{}
	maxBodySize := config.MaxBodySize
	if maxBodySize != nil {
		annotations[ecosystem.NginxIngressAnnotationBodySize] = maxBodySize.String()
	}

	rewriteTarget := config.RewriteTarget
	if rewriteTarget != "" {
		annotations[ecosystem.NginxIngressAnnotationRewriteTarget] = string(rewriteTarget)
	}

	additionalConfig := config.AdditionalConfig
	if additionalConfig != "" {
		annotations[ecosystem.NginxIngressAnnotationAdditionalConfig] = string(additionalConfig)
	}

	// Use nil here to delete existing annotation from the cr.
	if len(annotations) == 0 {
		return nil
	}

	return annotations
}

type doguCRPatch struct {
	Spec doguSpecPatch `json:"spec"`
}

type doguSpecPatch struct {
	Name                         string             `json:"name"`
	Version                      string             `json:"version"`
	Resources                    doguResourcesPatch `json:"resources"`
	SupportMode                  bool               `json:"supportMode"`
	UpgradeConfig                upgradeConfigPatch `json:"upgradeConfig"`
	AdditionalIngressAnnotations map[string]string  `json:"additionalIngressAnnotations"`
}

type upgradeConfigPatch struct {
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch"`
	ForceUpgrade         bool `json:"forceUpgrade"`
}

// DoguResources defines the physical resources used by the dogu.
type doguResourcesPatch struct {
	DataVolumeSize string `json:"dataVolumeSize"`
}

func toDoguCRPatch(dogu *ecosystem.DoguInstallation) *doguCRPatch {
	return &doguCRPatch{
		Spec: doguSpecPatch{
			Name:    dogu.Name.String(),
			Version: dogu.Version.Raw,
			Resources: doguResourcesPatch{
				DataVolumeSize: ecosystem.GetQuantityString(dogu.MinVolumeSize),
			},
			AdditionalIngressAnnotations: getNginxIngressAnnotations(dogu.ReverseProxyConfig),
			// always set this to false as a dogu cannot start in support mode
			SupportMode: false,
			UpgradeConfig: upgradeConfigPatch{
				AllowNamespaceSwitch: dogu.UpgradeConfig.AllowNamespaceSwitch,
				// this is a useful default as long as blueprints itself have no forceUpgrade flag implemented
				ForceUpgrade: false,
			},
		},
	}
}

func toDoguCRPatchBytes(dogu *ecosystem.DoguInstallation) ([]byte, error) {
	crPatch := toDoguCRPatch(dogu)
	patch, err := json.Marshal(crPatch)
	if err != nil {
		return []byte{}, &domainservice.InternalError{
			WrappedError: err,
			Message:      fmt.Sprintf("cannot patch dogu CR for dogu %q", dogu.Name),
		}
	}
	return patch, nil
}
