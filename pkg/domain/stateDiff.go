package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
)

// StateDiff represents the diff between the defined state in the effective blueprint and the actual state in the ecosystem.
// If there is a state in the ecosystem, which is not represented in the effective blueprint, then the expected state is the actual state.
type StateDiff struct {
	DoguDiffs                DoguDiffs
	ComponentDiffs           ComponentDiffs
	DoguConfigDiffs          map[cescommons.SimpleName]DoguConfigDiffs
	SensitiveDoguConfigDiffs map[cescommons.SimpleName]SensitiveDoguConfigDiffs
	GlobalConfigDiffs        GlobalConfigDiffs
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
	ActionUpdateAdditionalMounts          = "update additional mounts"
)

func (a Action) IsDoguProxyAction() bool {
	return a == ActionUpdateDoguProxyBodySize || a == ActionUpdateDoguProxyAdditionalConfig || a == ActionUpdateDoguProxyRewriteTarget
}

func (diff StateDiff) HasChanges() bool {
	return diff.DoguDiffs.HasChanges() ||
		diff.ComponentDiffs.HasChanges() ||
		diff.GlobalConfigDiffs.HasChanges() ||
		diff.HasDoguConfigChanges()
}

func (diff StateDiff) HasDoguConfigChanges() bool {
	for _, configDiff := range diff.DoguConfigDiffs {
		if configDiff.HasChanges() {
			return true
		}
	}
	for _, configDiff := range diff.SensitiveDoguConfigDiffs {
		if configDiff.HasChanges() {
			return true
		}
	}
	return false
}
