package domain

import "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

// StateDiff represents the diff between the defined state in the effective blueprint and the actual state in the ecosystem.
// If there is a state in the ecosystem, which is not represented in the effective blueprint, then the expected state is the actual state.
type StateDiff struct {
	DoguDiffs         DoguDiffs
	ComponentDiffs    ComponentDiffs
	DoguConfigDiffs   map[common.SimpleDoguName]CombinedDoguConfigDiffs
	GlobalConfigDiffs GlobalConfigDiffs
}

// Action represents a needed Action for a dogu to reach the expected state.
type Action string

const (
	ActionNone                            = "none"
	ActionInstall                         = "install"
	ActionUninstall                       = "uninstall"
	ActionUpgrade                         = "upgrade"
	ActionDowngrade                       = "downgrade"
	ActionSwitchDoguNamespace             = "dogu namespace switch"
	ActionSwitchComponentNamespace        = "component namespace switch"
	ActionUpdateDoguProxyBodySize         = "update proxy body size"
	ActionUpdateDoguProxyRewriteTarget    = "update proxy rewrite target"
	ActionUpdateDoguProxyAdditionalConfig = "update proxy additional config"
	ActionUpdateDoguResourceMinVolumeSize = "update resource minimum volume size"
)
