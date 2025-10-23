package domain

import (
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
)

// StateDiff represents the diff between the defined state in the effective blueprint and the actual state in the ecosystem.
// If there is a state in the ecosystem, which is not represented in the effective blueprint, then the expected state is the actual state.
type StateDiff struct {
	DoguDiffs                DoguDiffs
	DoguConfigDiffs          map[cescommons.SimpleName]DoguConfigDiffs
	SensitiveDoguConfigDiffs map[cescommons.SimpleName]SensitiveDoguConfigDiffs
	GlobalConfigDiffs        GlobalConfigDiffs
}

// Action represents a needed Action for a dogu to reach the expected state.
type Action = bpv2.DoguAction

const (
	// ActionInstall means the dogu is to be installed
	ActionInstall = bpv2.DoguActionInstall
	// ActionUninstall means the dogu is to be uninstalled
	ActionUninstall = bpv2.DoguActionUninstall
	// ActionUpgrade means an upgrade needs to be performed for the dogu
	ActionUpgrade = bpv2.DoguActionUpgrade
	// ActionDowngrade means a downgrade needs to be performed for the dogu
	ActionDowngrade = bpv2.DoguActionDowngrade
	// ActionSwitchDoguNamespace means the dogu should be pulled from a different dogu registry namespace
	ActionSwitchDoguNamespace = bpv2.DoguActionSwitchNamespace
	// ActionUpdateDoguReverseProxyConfig means the reverse proxy config of the dogu needs to be updated
	ActionUpdateDoguReverseProxyConfig = bpv2.DoguActionUpdateReverseProxyConfig
	// ActionUpdateDoguResourceMinVolumeSize means the minimum volume size of the dogu needs to be changed
	ActionUpdateDoguResourceMinVolumeSize = bpv2.DoguActionUpdateResourceMinVolumeSize
	// ActionUpdateAdditionalMounts means the additional mounts should be updated for the dogu
	ActionUpdateAdditionalMounts = bpv2.DoguActionUpdateAdditionalMounts
)

func (diff StateDiff) HasChanges() bool {
	return diff.DoguDiffs.HasChanges() ||
		diff.HasConfigChanges()
}

func (diff StateDiff) HasConfigChanges() bool {
	return diff.GlobalConfigDiffs.HasChanges() ||
		diff.HasDoguConfigChanges() ||
		diff.HasSensitiveDoguConfigChanges()
}

func (diff StateDiff) HasDoguConfigChanges() bool {
	for _, configDiff := range diff.DoguConfigDiffs {
		if configDiff.HasChanges() {
			return true
		}
	}
	return false
}

func (diff StateDiff) HasSensitiveDoguConfigChanges() bool {
	for _, configDiff := range diff.SensitiveDoguConfigDiffs {
		if configDiff.HasChanges() {
			return true
		}
	}
	return false
}
