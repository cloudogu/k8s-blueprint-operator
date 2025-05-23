package dogucr

import (
	"encoding/json"
	"errors"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
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
	doguName, nameErr := cescommons.QualifiedNameFromString(cr.Spec.Name)

	minVolumeSize, volumeSizeErr := parseMinDataVolumeSize(cr)

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
		MinVolumeSize:      minVolumeSize,
		ReverseProxyConfig: reverseProxyConfigEntries,
		PersistenceContext: persistenceContext,
		AdditionalMounts:   parseAdditionalMounts(cr.Spec.AdditionalMounts),
	}, nil
}

func parseMinDataVolumeSize(cr *v2.Dogu) (*ecosystem.VolumeSize, error) {
	emptySize := resource.Quantity{}
	//only use the deprecated value, if the new value is not set
	if cr.Spec.Resources.MinDataVolumeSize == emptySize {
		return ecosystem.GetQuantityReference(cr.Spec.Resources.DataVolumeSize)
	} else {
		return &cr.Spec.Resources.MinDataVolumeSize, nil
	}
}

func parseAdditionalMounts(mounts []v2.DataMount) []ecosystem.AdditionalMount {
	var result []ecosystem.AdditionalMount
	for _, m := range mounts {
		result = append(result, ecosystem.AdditionalMount{
			SourceType: ecosystem.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		})
	}
	return result
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
	minVolumeSize := resource.Quantity{}
	if dogu.MinVolumeSize != nil {
		minVolumeSize = *dogu.MinVolumeSize
	}
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
				// always set MinDataVolumeSize instead of the deprecated DataVolumeSize
				MinDataVolumeSize: minVolumeSize,
			},
			SupportMode: false,
			UpgradeConfig: v2.UpgradeConfig{
				AllowNamespaceSwitch: dogu.UpgradeConfig.AllowNamespaceSwitch,
				ForceUpgrade:         false,
			},
			AdditionalIngressAnnotations: getNginxIngressAnnotations(dogu.ReverseProxyConfig),
			AdditionalMounts:             toDoguCRAdditionalMounts(dogu.AdditionalMounts),
		},
		Status: v2.DoguStatus{},
	}
}

func toDoguCRAdditionalMounts(mounts []ecosystem.AdditionalMount) []v2.DataMount {
	var result []v2.DataMount
	for _, m := range mounts {
		result = append(result, v2.DataMount{
			SourceType: v2.DataSourceType(m.SourceType),
			Name:       m.Name,
			Volume:     m.Volume,
			Subfolder:  m.Subfolder,
		})
	}
	return result
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
	// do not use omitempty, because we cannot delete things then
	Name                         string             `json:"name"`
	Version                      string             `json:"version"`
	Resources                    doguResourcesPatch `json:"resources"`
	SupportMode                  bool               `json:"supportMode"`
	UpgradeConfig                upgradeConfigPatch `json:"upgradeConfig"`
	AdditionalIngressAnnotations map[string]string  `json:"additionalIngressAnnotations"`
	AdditionalMounts             []v2.DataMount     `json:"additionalMounts"`
}

type upgradeConfigPatch struct {
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch"`
	ForceUpgrade         bool `json:"forceUpgrade"`
}

// DoguResources defines the physical resources used by the dogu.
type doguResourcesPatch struct {
	// DataVolumeSize
	// Deprecated: use MinDataVolumeSize instead. Only set it to correct possibly wrong dogu CRs
	DataVolumeSize    string            `json:"dataVolumeSize"`
	MinDataVolumeSize resource.Quantity `json:"minDataVolumeSize"`
}

func toDoguCRPatch(dogu *ecosystem.DoguInstallation) *doguCRPatch {
	minVolumeSize := resource.Quantity{}
	if dogu.MinVolumeSize != nil {
		minVolumeSize = *dogu.MinVolumeSize
	}
	return &doguCRPatch{
		Spec: doguSpecPatch{
			Name:    dogu.Name.String(),
			Version: dogu.Version.Raw,
			Resources: doguResourcesPatch{
				// remove the deprecated value from the dogu CR and replace it with the new one
				DataVolumeSize:    "",
				MinDataVolumeSize: minVolumeSize,
			},
			AdditionalIngressAnnotations: getNginxIngressAnnotations(dogu.ReverseProxyConfig),
			// always set this to false as a dogu cannot start in support mode
			SupportMode: false,
			UpgradeConfig: upgradeConfigPatch{
				AllowNamespaceSwitch: dogu.UpgradeConfig.AllowNamespaceSwitch,
				// this is a useful default as long as blueprints itself have no forceUpgrade flag implemented
				ForceUpgrade: false,
			},
			AdditionalMounts: toDoguCRAdditionalMounts(dogu.AdditionalMounts),
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
