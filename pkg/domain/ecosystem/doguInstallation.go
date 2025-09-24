package ecosystem

import (
	"fmt"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DoguInstallation represents an installed or to be installed dogu in the ecosystem.
type DoguInstallation struct {
	// Name identifies the dogu by simple dogu name and namespace.
	Name cescommons.QualifiedName
	// Version is the desired version of the dogu
	Version core.Version
	// Status is the installation status of the dogu in the ecosystem
	Status string
	// Health is the current health status of the dogu in the ecosystem
	Health HealthStatus
	// InstalledVersion is the current version of the dogu
	InstalledVersion core.Version
	// StartedAt contains the time of the last restart of the dogu.
	StartedAt metav1.Time
	// UpgradeConfig contains configuration for dogu upgrades
	UpgradeConfig UpgradeConfig
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}
	// MinVolumeSize is the minimum storage of the dogu. This field is optional and can be nil to indicate that no
	// storage is needed.
	MinVolumeSize *VolumeSize
	// ReverseProxyConfig defines configuration for the ecosystem reverse proxy. This field is optional.
	ReverseProxyConfig *ReverseProxyConfig
	// AdditionalMounts provides the possibility to mount additional data into the dogu.
	AdditionalMounts []AdditionalMount
}

// TODO: Unused constants needed?
const (
	DoguStatusNotInstalled = ""
	DoguStatusInstalling   = "installing"
	DoguStatusUpgrading    = "upgrading"
	DoguStatusDeleting     = "deleting"
	DoguStatusInstalled    = "installed"
	DoguStatusPVCResizing  = "resizing PVC"
)

// Specific Nginx annotations. In future those annotations will be replaced be generalized fields in the dogu cr.
// The dogu-operator or service-discovery will interpret them.
const (
	NginxIngressAnnotationBodySize         = "nginx.ingress.kubernetes.io/proxy-body-size"
	NginxIngressAnnotationRewriteTarget    = "nginx.ingress.kubernetes.io/rewrite-target"
	NginxIngressAnnotationAdditionalConfig = "nginx.ingress.kubernetes.io/configuration-snippet"
)

type ReverseProxyConfig struct {
	MaxBodySize      *BodySize
	RewriteTarget    RewriteTarget
	AdditionalConfig AdditionalConfig
}

func (r *ReverseProxyConfig) IsEmpty() bool {
	return r == nil || (r.MaxBodySize == nil && r.RewriteTarget == nil && r.AdditionalConfig == nil)
}

// UpgradeConfig contains configuration hints regarding aspects during the upgrade of dogus.
type UpgradeConfig struct {
	// AllowNamespaceSwitch lets a dogu switch its dogu namespace during an upgrade. The dogu must be technically the
	// same dogu which did reside in a different namespace. The remote dogu's version must be equal to or greater than
	// the version of the local dogu.
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch,omitempty"`
}

type DataSourceType string

const (
	// DataSourceConfigMap mounts a config map as a data source.
	DataSourceConfigMap DataSourceType = "ConfigMap"
	// DataSourceSecret mounts a secret as a data source.
	DataSourceSecret DataSourceType = "Secret"
)

// AdditionalMount is a description of what data should be mounted to a specific Dogu volume (already defined in dogu.json).
type AdditionalMount struct {
	// SourceType defines where the data is coming from.
	// Valid options are:
	//   ConfigMap - data stored in a kubernetes ConfigMap.
	//   Secret - data stored in a kubernetes Secret.
	SourceType DataSourceType
	// Name is the name of the data source.
	Name string
	// Volume is the name of the volume to which the data should be mounted. It is defined in the respective dogu.json.
	Volume string
	// Subfolder defines a subfolder in which the data should be put within the volume.
	// +optional
	Subfolder *string
}

// InstallDogu is a factory for new DoguInstallation's.
func InstallDogu(
	name cescommons.QualifiedName,
	version *core.Version,
	minVolumeSize *VolumeSize,
	reverseProxyConfig *ReverseProxyConfig,
	additionalMounts []AdditionalMount) *DoguInstallation {

	doguVersion := core.Version{}
	if version != nil {
		doguVersion = *version
	}

	return &DoguInstallation{
		Name:               name,
		Version:            doguVersion,
		UpgradeConfig:      UpgradeConfig{AllowNamespaceSwitch: false},
		MinVolumeSize:      minVolumeSize,
		ReverseProxyConfig: reverseProxyConfig,
		AdditionalMounts:   additionalMounts,
	}
}

func (dogu *DoguInstallation) IsHealthy() bool {
	return dogu.Health == AvailableHealthStatus
}

func (dogu *DoguInstallation) IsVersionUpToDate() bool {
	return dogu.Version.IsEqualTo(dogu.InstalledVersion)
}

func (dogu *DoguInstallation) IsConfigUpToDate(globalConfigUpdateTime *metav1.Time, doguConfigUpdateTime *metav1.Time) bool {
	return !dogu.StartedAt.Before(globalConfigUpdateTime) && !dogu.StartedAt.Before(doguConfigUpdateTime)
}

func (dogu *DoguInstallation) Upgrade(newVersion *core.Version) {
	dogu.Version = core.Version{}
	if newVersion != nil {
		dogu.Version = *newVersion
	}

	dogu.UpgradeConfig.AllowNamespaceSwitch = false
}

func (dogu *DoguInstallation) SwitchNamespace(newNamespace cescommons.Namespace, isNamespaceSwitchAllowed bool) error {
	if !isNamespaceSwitchAllowed {
		return fmt.Errorf("not allowed to switch dogu namespace from %q to %q", dogu.Name.Namespace, newNamespace)
	}
	dogu.Name.Namespace = newNamespace
	dogu.UpgradeConfig.AllowNamespaceSwitch = true
	return nil
}

func (dogu *DoguInstallation) UpdateProxyBodySize(value *BodySize) {
	if dogu.ReverseProxyConfig == nil {
		dogu.ReverseProxyConfig = &ReverseProxyConfig{}
	}
	dogu.ReverseProxyConfig.MaxBodySize = value
}

func (dogu *DoguInstallation) UpdateMinVolumeSize(size *VolumeSize) {
	dogu.MinVolumeSize = size
}

func (dogu *DoguInstallation) UpdateProxyRewriteTarget(value RewriteTarget) {
	if dogu.ReverseProxyConfig == nil {
		dogu.ReverseProxyConfig = &ReverseProxyConfig{}
	}
	dogu.ReverseProxyConfig.RewriteTarget = value
}

func (dogu *DoguInstallation) UpdateProxyAdditionalConfig(value AdditionalConfig) {
	if dogu.ReverseProxyConfig == nil {
		dogu.ReverseProxyConfig = &ReverseProxyConfig{}
	}
	dogu.ReverseProxyConfig.AdditionalConfig = value
}

func (dogu *DoguInstallation) UpdateAdditionalMounts(mounts []AdditionalMount) {
	dogu.AdditionalMounts = mounts
}
