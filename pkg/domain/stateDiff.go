package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
)

// StateDiff represents the diff between the defined state in the effective blueprint and the actual state in the ecosystem.
// If there is a state in the ecosystem, which is not represented in the effective blueprint, then the expected state is the actual state.
type StateDiff struct {
	DoguDiffs         DoguDiffs
	ComponentDiffs    ComponentDiffs
	DoguConfigDiffs   map[common.SimpleDoguName]CombinedDoguConfigDiffs
	GlobalConfigDiffs GlobalConfigDiffs
}

func (diff StateDiff) GetDoguConfigDiffsByAction() map[ConfigAction]DoguConfigDiffs {
	configDiffsByAction := map[ConfigAction]DoguConfigDiffs{}

	for _, combinedConfig := range diff.DoguConfigDiffs {
		for key, value := range combinedConfig.DoguConfigDiff.GetDoguConfigDiffByAction() {
			configDiffsByAction[key] = append(configDiffsByAction[key], value...)
		}
	}

	return configDiffsByAction
}

func (diff StateDiff) GetSensitiveDoguConfigDiffsByAction() map[ConfigAction]SensitiveDoguConfigDiffs {
	configDiffsByAction := map[ConfigAction]SensitiveDoguConfigDiffs{}

	for _, combinedConfig := range diff.DoguConfigDiffs {
		for key, value := range combinedConfig.SensitiveDoguConfigDiff.GetSensitiveDoguConfigDiffByAction() {
			configDiffsByAction[key] = append(configDiffsByAction[key], value...)
		}
	}

	return configDiffsByAction
}

// Action represents a needed Action for a dogu to reach the expected state.
type Action string

const (
	ActionInstall                         = "install"
	ActionUninstall                       = "uninstall"
	ActionUpgrade                         = "upgrade"
	ActionDowngrade                       = "downgrade"
	ActionSwitchDoguNamespace             = "dogu namespace switch"
	ActionUpdateDoguProxyBodySize         = "update proxy body size"
	ActionUpdateDoguProxyRewriteTarget    = "update proxy rewrite target"
	ActionUpdateDoguProxyAdditionalConfig = "update proxy additional config"
	ActionUpdateDoguResourceMinVolumeSize = "update resource minimum volume size"
	ActionSwitchComponentNamespace        = "component namespace switch"
	ActionUpdateComponentDeployConfig     = "update component package config"
)

func (a Action) IsDoguProxyAction() bool {
	return a == ActionUpdateDoguProxyBodySize || a == ActionUpdateDoguProxyAdditionalConfig || a == ActionUpdateDoguProxyRewriteTarget
}
