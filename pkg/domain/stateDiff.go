package domain

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

type StateDiff struct {
	DoguDiffs []DoguDiff
}

type DoguDiff struct {
	DoguName     string
	Actual       DoguDiffState
	Expected     DoguDiffState
	NeededAction Action
}

type DoguDiffState struct {
	Namespace         string
	Version           core.Version
	InstallationState TargetState
}

type Action string

const (
	ActionNone            = "none"
	ActionInstall         = "install"
	ActionUninstall       = "uninstall"
	ActionUpgrade         = "upgrade"
	ActionDowngrade       = "downgrade"
	ActionSwitchNamespace = "namespace switch"
)

func determineDoguDiffs(spec *BlueprintSpec, installedDogus map[string]*ecosystem.DoguInstallation) []DoguDiff {
	var doguDiffs []DoguDiff
	for _, specDogu := range spec.EffectiveBlueprint.Dogus {
		installedDogu := installedDogus[specDogu.Name]
		doguDiffs = append(doguDiffs, determineDoguDiff(&specDogu, installedDogu))
	}
	return doguDiffs
}

// determineDoguDiff creates a DoguDiff out of a Dogu from the blueprint and the ecosystem.DoguInstallation in the ecosystem.
// if the Dogu is nil (was not in the blueprint), the actual state is also the expected state.
// if the installedDogu is nil, it is considered to be not installed currently.
// returns a DoguDiff
func determineDoguDiff(blueprintDogu *Dogu, installedDogu *ecosystem.DoguInstallation) DoguDiff {
	var expectedState, actualState DoguDiffState
	doguName := "" // either blueprintDogu or installedDogu could be nil

	if installedDogu == nil {
		actualState = DoguDiffState{
			InstallationState: TargetStateAbsent,
		}
	} else {
		doguName = installedDogu.Name
		actualState = DoguDiffState{
			Namespace:         installedDogu.Namespace,
			Version:           installedDogu.Version,
			InstallationState: TargetStatePresent,
		}
	}

	if blueprintDogu == nil {
		expectedState = actualState
	} else {
		doguName = blueprintDogu.Name
		expectedState = DoguDiffState{
			Namespace:         blueprintDogu.Namespace,
			Version:           blueprintDogu.Version,
			InstallationState: blueprintDogu.TargetState,
		}
	}

	return DoguDiff{
		DoguName:     doguName,
		Expected:     expectedState,
		Actual:       actualState,
		NeededAction: getNeededAction(expectedState, actualState),
	}
}

func getNeededAction(expected DoguDiffState, actual DoguDiffState) Action {
	if expected.InstallationState == actual.InstallationState {
		switch expected.InstallationState {
		case TargetStatePresent:
			// dogu should state installed, but maybe an upgrade, downgrade or a namespace switch?
			if expected.Namespace != actual.Namespace {
				return ActionSwitchNamespace
			}
			if expected.Version.IsNewerThan(actual.Version) {
				return ActionUpgrade
			} else if expected.Version.IsEqualTo(actual.Version) {
				return ActionNone
			} else { // os older
				return ActionDowngrade
			}
		case TargetStateAbsent:
			return ActionNone
		}
	} else {
		// actual state is always the opposite
		switch expected.InstallationState {
		case TargetStatePresent:
			return ActionInstall
		case TargetStateAbsent:
			return ActionUninstall
		}
	}
	// all cases should be handled above, but if new fields are added, this is a safe fallback for any bugs.
	return ActionNone
}
