package domain

import (
	"fmt"
	"slices"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"golang.org/x/exp/maps"
)

// DoguDiffs contains the Diff for all expected Dogus to the current ecosystem.DoguInstallations.
type DoguDiffs []DoguDiff

func (diffs DoguDiffs) HasChanges() bool {
	for _, diff := range diffs {
		if diff.HasChanges() {
			return true
		}
	}
	return false
}

// DoguDiff represents the Diff for a single expected Dogu to the current ecosystem.DoguInstallation.
type DoguDiff struct {
	DoguName      cescommons.SimpleName
	Actual        DoguDiffState
	Expected      DoguDiffState
	NeededActions []Action
}

// DoguDiffState contains all fields to make a diff for dogus in respect to another DoguDiffState.
type DoguDiffState struct {
	Namespace        cescommons.Namespace
	Version          *core.Version
	InstalledVersion *core.Version
	Absent           bool
	MinVolumeSize    *ecosystem.VolumeSize
	StorageClassName *string
	AdditionalMounts []ecosystem.AdditionalMount
}

func (diff DoguDiff) HasChanges() bool {
	return len(diff.NeededActions) != 0
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
		"{Version: %q, Namespace: %q, Absent: %t}",
		diff.getSafeVersionString(),
		diff.Namespace,
		diff.Absent,
	)
}

func (diff *DoguDiffState) getSafeVersionString() string {
	if diff.Version != nil {
		return diff.Version.String()
	} else {
		return ""
	}
}

// determineDoguDiffs creates DoguDiffs for all dogus in the blueprint and all installed dogus as well.
// see determineDoguDiff for more information.
func determineDoguDiffs(blueprintDogus []Dogu, installedDogus map[cescommons.SimpleName]*ecosystem.DoguInstallation) []DoguDiff {
	var doguDiffs = map[cescommons.SimpleName]DoguDiff{}
	for _, blueprintDogu := range blueprintDogus {
		installedDogu := installedDogus[blueprintDogu.Name.SimpleName]
		determinedDoguDiff := determineDoguDiff(&blueprintDogu, installedDogu)
		// only add changes to diff
		if determinedDoguDiff != nil {
			doguDiffs[blueprintDogu.Name.SimpleName] = *determinedDoguDiff
		}
	}
	for _, installedDogu := range installedDogus {
		_, found := FindDoguByName(blueprintDogus, installedDogu.Name.SimpleName)
		// Only create DoguDiff if the installed dogu is not found in the blueprint.
		// If the installed dogu is in blueprint the DoguDiff was already determined above.
		if !found {
			determinedDoguDiff := determineDoguDiff(nil, installedDogu)
			// only add changes to diff
			if determinedDoguDiff != nil {
				doguDiffs[installedDogu.Name.SimpleName] = *determinedDoguDiff
			}
		}
	}
	return maps.Values(doguDiffs)
}

// determineDoguDiff creates a DoguDiff out of a Dogu from the blueprint and the ecosystem.DoguInstallation in the ecosystem.
// if the Dogu is nil (was not in the blueprint), the actual state is also the expected state.
// if the installedDogu is nil, it is considered to be not installed currently.
// returns a DoguDiff
func determineDoguDiff(blueprintDogu *Dogu, installedDogu *ecosystem.DoguInstallation) *DoguDiff {
	var expectedState, actualState DoguDiffState
	var doguName cescommons.SimpleName = "" // either blueprintDogu or installedDogu could be nil

	if installedDogu == nil {
		actualState = DoguDiffState{
			Absent: true,
		}
	} else {
		doguName = installedDogu.Name.SimpleName
		actualState = DoguDiffState{
			Namespace:        installedDogu.Name.Namespace,
			Version:          &installedDogu.Version,
			InstalledVersion: &installedDogu.InstalledVersion,
			MinVolumeSize:    installedDogu.MinVolumeSize,
			StorageClassName: installedDogu.StorageClassName,
			AdditionalMounts: installedDogu.AdditionalMounts,
		}
	}

	if blueprintDogu == nil {
		expectedState = actualState
	} else {
		doguName = blueprintDogu.Name.SimpleName
		expectedState = DoguDiffState{
			Namespace:        blueprintDogu.Name.Namespace,
			Version:          blueprintDogu.Version,
			Absent:           blueprintDogu.Absent,
			MinVolumeSize:    blueprintDogu.MinVolumeSize,
			StorageClassName: blueprintDogu.StorageClassName,
			AdditionalMounts: blueprintDogu.AdditionalMounts,
		}
	}

	actions := getNeededDoguActions(expectedState, actualState)
	if len(actions) == 0 {
		return nil
	}

	return &DoguDiff{
		DoguName:      doguName,
		Expected:      expectedState,
		Actual:        actualState,
		NeededActions: actions,
	}
}

func getNeededDoguActions(expected DoguDiffState, actual DoguDiffState) []Action {
	if expected.Absent == actual.Absent {
		if expected.Absent {
			return []Action{}
		} else {
			return getActionsForPresentDoguDiffs(expected, actual)
		}
	} else {
		// actual state is always the opposite
		if expected.Absent {
			return []Action{ActionUninstall}
		} else {
			return []Action{ActionInstall}
		}
	}
}

func getActionsForPresentDoguDiffs(expected DoguDiffState, actual DoguDiffState) []Action {
	var neededActions []Action
	if expected.Namespace != actual.Namespace {
		neededActions = append(neededActions, ActionSwitchDoguNamespace)
	}

	neededActions = appendActionForMinVolumeSize(neededActions, expected.MinVolumeSize, actual.MinVolumeSize)
	neededActions = appendActionForAdditionalMounts(neededActions, expected.AdditionalMounts, actual.AdditionalMounts)

	if expected.Version != nil && actual.Version != nil && expected.Version.IsNewerThan(*actual.Version) {
		neededActions = append(neededActions, ActionUpgrade)
	} else if expected.Version != nil && actual.Version.IsNewerThan(*expected.Version) {
		// if downgrades are allowed is not important here.
		// Downgrades can be rejected later, so forcing downgrades via a flag can be implemented without changing this code here.
		neededActions = append(neededActions, ActionDowngrade)
	}

	return neededActions
}

func appendActionForMinVolumeSize(actions []Action, expectedSize *ecosystem.VolumeSize, actualSize *ecosystem.VolumeSize) []Action {
	// if expected > actual = update needed
	if expectedSize == nil {
		return actions
	} else if actualSize == nil || expectedSize.Cmp(*actualSize) > 0 {
		return append(actions, ActionUpdateDoguResourceMinVolumeSize)
	}
	return actions
}

func appendActionForAdditionalMounts(actions []Action, expectedMounts []ecosystem.AdditionalMount, actualMounts []ecosystem.AdditionalMount) []Action {
	if !areAdditionalMountsEqual(expectedMounts, actualMounts) {
		return append(actions, ActionUpdateAdditionalMounts)
	}
	return actions
}

// areAdditionalMountsEqual compare the additional mounts without order
func areAdditionalMountsEqual(first []ecosystem.AdditionalMount, second []ecosystem.AdditionalMount) bool {
	if len(first) != len(second) {
		return false
	}
	for _, mount := range first {
		if !slices.Contains(second, mount) {
			return false
		}
	}
	return true
}
