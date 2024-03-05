package domain

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"golang.org/x/exp/maps"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

// DoguDiffs contains the Diff for all expected Dogus to the current ecosystem.DoguInstallations.
type DoguDiffs []DoguDiff

// Statistics aggregates various figures about the required actions of the DoguDiffs.
func (dd DoguDiffs) Statistics() (toInstall, toUpgrade, toUninstall, toUpdateReverseProxyConfig, toUpdateResourceConfig, other int) {
	for _, doguDiff := range dd {
		for _, action := range doguDiff.NeededActions {
			switch action {
			case ActionInstall:
				toInstall += 1
			case ActionUpgrade:
				toUpgrade += 1
			case ActionUninstall:
				toUninstall += 1
			case ActionUpdateDoguProxyBodySize, ActionUpdateDoguProxyRewriteTarget, ActionUpdateDoguProxyAdditionalConfig:
				toUpdateReverseProxyConfig += 1
			case ActionUpdateDoguResourceMinVolumeSize:
				toUpdateResourceConfig += 1
			default:
				other += 1
			}
		}
	}
	return
}

// DoguDiff represents the Diff for a single expected Dogu to the current ecosystem.DoguInstallation.
type DoguDiff struct {
	DoguName      common.SimpleDoguName
	Actual        DoguDiffState
	Expected      DoguDiffState
	NeededActions []Action
}

// DoguDiffState contains all fields to make a diff for dogus in respect to another DoguDiffState.
type DoguDiffState struct {
	Namespace          common.DoguNamespace
	Version            core.Version
	InstallationState  TargetState
	MinVolumeSize      ecosystem.VolumeSize
	ReverseProxyConfig ecosystem.ReverseProxyConfig
}

// String returns a string representation of the DoguDiff.
func (diff *DoguDiff) String() string {
	return fmt.Sprintf(
		"{DoguName: %q, Actual: %s, Expected: %s, NeededActions: %q}",
		diff.DoguName,
		diff.Actual.String(),
		diff.Expected.String(),
		diff.NeededActions,
	)
}

// String returns a string representation of the DoguDiffState.
func (diff *DoguDiffState) String() string {
	return fmt.Sprintf(
		"{Version: %q, Namespace: %q, InstallationState: %q}",
		diff.Version.Raw,
		diff.Namespace,
		diff.InstallationState,
	)
}

// determineDoguDiffs creates DoguDiffs for all dogus in the blueprint and all installed dogus as well.
// see determineDoguDiff for more information.
func determineDoguDiffs(blueprintDogus []Dogu, installedDogus map[common.SimpleDoguName]*ecosystem.DoguInstallation) []DoguDiff {
	var doguDiffs = map[common.SimpleDoguName]DoguDiff{}
	for _, blueprintDogu := range blueprintDogus {
		installedDogu := installedDogus[blueprintDogu.Name.SimpleName]
		doguDiffs[blueprintDogu.Name.SimpleName] = determineDoguDiff(&blueprintDogu, installedDogu)
	}
	for _, installedDogu := range installedDogus {
		blueprintDogu, notFound := FindDoguByName(blueprintDogus, installedDogu.Name.SimpleName)

		if notFound == nil {
			doguDiffs[installedDogu.Name.SimpleName] = determineDoguDiff(&blueprintDogu, installedDogu)
		} else {
			// if no dogu with this name in blueprint, use nil to indicate that
			doguDiffs[installedDogu.Name.SimpleName] = determineDoguDiff(nil, installedDogu)
		}
	}
	return maps.Values(doguDiffs)
}

// determineDoguDiff creates a DoguDiff out of a Dogu from the blueprint and the ecosystem.DoguInstallation in the ecosystem.
// if the Dogu is nil (was not in the blueprint), the actual state is also the expected state.
// if the installedDogu is nil, it is considered to be not installed currently.
// returns a DoguDiff
func determineDoguDiff(blueprintDogu *Dogu, installedDogu *ecosystem.DoguInstallation) DoguDiff {
	var expectedState, actualState DoguDiffState
	var doguName common.SimpleDoguName = "" // either blueprintDogu or installedDogu could be nil

	if installedDogu == nil {
		actualState = DoguDiffState{
			InstallationState: TargetStateAbsent,
		}
	} else {
		doguName = installedDogu.Name.SimpleName
		actualState = DoguDiffState{
			Namespace:          installedDogu.Name.Namespace,
			Version:            installedDogu.Version,
			InstallationState:  TargetStatePresent,
			MinVolumeSize:      installedDogu.MinVolumeSize,
			ReverseProxyConfig: installedDogu.ReverseProxyConfig,
		}
	}

	if blueprintDogu == nil {
		expectedState = actualState
	} else {
		doguName = blueprintDogu.Name.SimpleName
		expectedState = DoguDiffState{
			Namespace:          blueprintDogu.Name.Namespace,
			Version:            blueprintDogu.Version,
			InstallationState:  blueprintDogu.TargetState,
			MinVolumeSize:      blueprintDogu.MinVolumeSize,
			ReverseProxyConfig: blueprintDogu.ReverseProxyConfig,
		}
	}

	return DoguDiff{
		DoguName:      doguName,
		Expected:      expectedState,
		Actual:        actualState,
		NeededActions: getNeededDoguActions(expectedState, actualState),
	}
}

func getNeededDoguActions(expected DoguDiffState, actual DoguDiffState) []Action {
	var neededActions []Action
	if expected.InstallationState == actual.InstallationState {
		switch expected.InstallationState {
		case TargetStatePresent:
			// dogu should stay installed, but maybe it needs an upgrade, downgrade or a namespace switch?
			if expected.Namespace != actual.Namespace {
				neededActions = append(neededActions, ActionSwitchDoguNamespace)
			}
			if expected.MinVolumeSize != actual.MinVolumeSize {
				if expected.MinVolumeSize.Cmp(actual.MinVolumeSize) == 1 {
					neededActions = append(neededActions, ActionUpdateDoguResourceMinVolumeSize)
				}
				// TODO check for downgrade
			}

			proxyBodySizeAction := getActionForProxyBodySizes(expected.ReverseProxyConfig.MaxBodySize, actual.ReverseProxyConfig.MaxBodySize)
			if proxyBodySizeAction != "" {
				neededActions = append(neededActions, proxyBodySizeAction)
			}

			if expected.ReverseProxyConfig.RewriteTarget != actual.ReverseProxyConfig.RewriteTarget {
				neededActions = append(neededActions, ActionUpdateDoguProxyRewriteTarget)
			}
			if expected.ReverseProxyConfig.AdditionalConfig != actual.ReverseProxyConfig.AdditionalConfig {
				neededActions = append(neededActions, ActionUpdateDoguProxyAdditionalConfig)
			}
			if expected.Version.IsNewerThan(actual.Version) {
				neededActions = append(neededActions, ActionUpgrade)
			} else if expected.Version.IsEqualTo(actual.Version) {
				if len(neededActions) == 0 {
					neededActions = append(neededActions, ActionNone)
				}
			} else { // is older
				// if downgrades are allowed is not important here.
				// Downgrades can be rejected later, so forcing downgrades via a flag can be implemented without changing this code here.
				neededActions = append(neededActions, ActionDowngrade)
			}

			return neededActions
		case TargetStateAbsent:
			return append(neededActions, ActionNone)
		}
	} else {
		// actual state is always the opposite
		switch expected.InstallationState {
		case TargetStatePresent:
			return append(neededActions, ActionInstall)
		case TargetStateAbsent:
			return append(neededActions, ActionUninstall)
		}
	}
	// all cases should be handled above, but if new fields are added, this is a safe fallback for any bugs.
	return append(neededActions, ActionNone)
}

func getActionForProxyBodySizes(expectedProxyBodySize *ecosystem.BodySize, actualProxyBodySize *ecosystem.BodySize) Action {
	if expectedProxyBodySize == nil && actualProxyBodySize == nil {
		return ""
	} else if proxyBodySizeIdentityChanged(expectedProxyBodySize, actualProxyBodySize) {
		return ActionUpdateDoguProxyBodySize
	} else {
		if expectedProxyBodySize.Cmp(*actualProxyBodySize) != 0 {
			return ActionUpdateDoguProxyBodySize
		}
	}
	return ""
}

func proxyBodySizeIdentityChanged(expectedProxyBodySize *ecosystem.BodySize, actualProxyBodySize *ecosystem.BodySize) bool {
	return (expectedProxyBodySize != nil && actualProxyBodySize == nil) || (expectedProxyBodySize == nil && actualProxyBodySize != nil)
}
