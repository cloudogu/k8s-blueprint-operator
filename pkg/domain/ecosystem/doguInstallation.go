package ecosystem

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
)

// DoguInstallation represents an installed or to be installed dogu in the ecosystem.
type DoguInstallation struct {
	// Name identifies the dogu by simple dogu name and namespace.
	Name common.QualifiedDoguName
	// Version is the version of the dogu
	Version core.Version
	// Status is the installation status of the dogu in the ecosystem
	Status string
	// Health is the current health status of the dogu in the ecosystem
	Health HealthStatus
	// UpgradeConfig contains configuration for dogu upgrades
	UpgradeConfig UpgradeConfig
	// PersistenceContext can hold generic values needed for persistence with repositories, e.g. version counters or transaction contexts.
	// This field has a generic map type as the values within it highly depend on the used type of repository.
	// This field should be ignored in the whole domain.
	PersistenceContext map[string]interface{}

	MinVolumeSize      VolumeSize
	ReverseProxyConfig ReverseProxyConfig
}

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

// UpgradeConfig contains configuration hints regarding aspects during the upgrade of dogus.
type UpgradeConfig struct {
	// AllowNamespaceSwitch lets a dogu switch its dogu namespace during an upgrade. The dogu must be technically the
	// same dogu which did reside in a different namespace. The remote dogu's version must be equal to or greater than
	// the version of the local dogu.
	AllowNamespaceSwitch bool `json:"allowNamespaceSwitch,omitempty"`
}

// InstallDogu is a factory for new DoguInstallation's.
func InstallDogu(name common.QualifiedDoguName, version core.Version, minVolumeSize VolumeSize, reverseProxyConfig ReverseProxyConfig) *DoguInstallation {
	return &DoguInstallation{
		Name:               name,
		Version:            version,
		UpgradeConfig:      UpgradeConfig{AllowNamespaceSwitch: false},
		MinVolumeSize:      minVolumeSize,
		ReverseProxyConfig: reverseProxyConfig,
	}
}

func (dogu *DoguInstallation) IsHealthy() bool {
	return dogu.Health == AvailableHealthStatus
}

func (dogu *DoguInstallation) Upgrade(newVersion core.Version) {
	dogu.Version = newVersion
	dogu.UpgradeConfig.AllowNamespaceSwitch = false
}

func (dogu *DoguInstallation) SwitchNamespace(newNamespace common.DoguNamespace, isNamespaceSwitchAllowed bool) error {
	if !isNamespaceSwitchAllowed {
		return fmt.Errorf("not allowed to switch dogu namespace from %q to %q", dogu.Name.Namespace, newNamespace)
	}
	dogu.Name.Namespace = newNamespace
	dogu.UpgradeConfig.AllowNamespaceSwitch = true
	return nil
}

func (dogu *DoguInstallation) UpdateProxyBodySize(value *BodySize) {
	dogu.ReverseProxyConfig.MaxBodySize = value
}

func (dogu *DoguInstallation) UpdateMinVolumeSize(size VolumeSize) {
	dogu.MinVolumeSize = size
}

func (dogu *DoguInstallation) UpdateProxyRewriteTarget(value RewriteTarget) {
	dogu.ReverseProxyConfig.RewriteTarget = value
}

func (dogu *DoguInstallation) UpdateProxyAdditionalConfig(value AdditionalConfig) {
	dogu.ReverseProxyConfig.AdditionalConfig = value
}
