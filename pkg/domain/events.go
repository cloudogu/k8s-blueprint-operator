package domain

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/util"
)

type Event interface {
	Name() string
	Message() string
}

type BlueprintSpecInvalidEvent struct {
	ValidationError error
}

func (b BlueprintSpecInvalidEvent) Name() string {
	return "BlueprintSpecInvalid"
}

func (b BlueprintSpecInvalidEvent) Message() string {
	return b.ValidationError.Error()
}

type MissingConfigReferencesEvent struct {
	err error
}

func NewMissingConfigReferencesEvent(err error) MissingConfigReferencesEvent {
	return MissingConfigReferencesEvent{err: err}
}

func (e MissingConfigReferencesEvent) Name() string {
	return "MissingConfigReferences"
}

func (e MissingConfigReferencesEvent) Message() string {
	return e.err.Error()
}

// StateDiffComponentDeterminedEvent provides event information over detected changes regarding components.
type StateDiffComponentDeterminedEvent struct {
	componentDiffs []ComponentDiff
}

func newStateDiffComponentEvent(componentDiffs ComponentDiffs) StateDiffComponentDeterminedEvent {
	return StateDiffComponentDeterminedEvent{
		componentDiffs: componentDiffs,
	}
}

// Name contains the StateDiffComponentDeterminedEvent display name.
func (s StateDiffComponentDeterminedEvent) Name() string {
	return "StateDiffComponentDetermined"
}

// Message contains the StateDiffComponentDeterminedEvent's statistics message.
func (s StateDiffComponentDeterminedEvent) Message() string {
	var amountActions = map[Action]int{}
	for _, diff := range s.componentDiffs {
		for _, action := range diff.NeededActions {
			amountActions[action]++
		}
	}

	message, amount := getActionAmountMessage(amountActions)

	return fmt.Sprintf("component state diff determined: %d actions (%s)", amount, message)
}

func getActionAmountMessage(amountActions map[Action]int) (message string, totalAmount int) {
	var messages []string
	for action, amount := range amountActions {
		messages = append(messages, fmt.Sprintf("%q: %d", action, amount))
		totalAmount += amount
	}
	slices.Sort(messages)
	message = strings.Join(messages, ", ")
	return
}

// StateDiffDeterminedEvent provides event information over detected changes regarding dogus.
type StateDiffDeterminedEvent struct {
	doguDiffs            DoguDiffs
	GlobalConfigDiffs    GlobalConfigDiffs
	DoguConfigDiffs      map[cescommons.SimpleName]DoguConfigDiffs
	SensitiveConfigDiffs map[cescommons.SimpleName]SensitiveDoguConfigDiffs
}

func newStateDiffEvent(stateDiff StateDiff) StateDiffDeterminedEvent {
	return StateDiffDeterminedEvent{
		doguDiffs:            stateDiff.DoguDiffs,
		DoguConfigDiffs:      stateDiff.DoguConfigDiffs,
		GlobalConfigDiffs:    stateDiff.GlobalConfigDiffs,
		SensitiveConfigDiffs: stateDiff.SensitiveDoguConfigDiffs,
	}
}

// Name contains the StateDiffDoguDeterminedEvent display name.
func (s StateDiffDeterminedEvent) Name() string {
	return "StateDiffDetermined"
}

const groupedDoguProxyAction = "update reverse proxy"

// Message contains the StateDiffDoguDeterminedEvent's statistics message.
func (s StateDiffDeterminedEvent) Message() string {
	var amountActions = map[Action]int{}
	for _, diff := range s.doguDiffs {
		for _, action := range diff.NeededActions {
			if action.IsDoguProxyAction() {
				amountActions[groupedDoguProxyAction]++
			} else {
				amountActions[action]++
			}
		}
	}

	doguMessage, doguAmount := getActionAmountMessage(amountActions)

	return fmt.Sprintf("state diff determined:\n  %s\n  %d dogu actions (%s)", s.generateConfigChangeCounter(), doguAmount, doguMessage)
}

func (s StateDiffDeterminedEvent) generateConfigChangeCounter() string {
	configActions := util.Map(s.GlobalConfigDiffs, func(entryDiff GlobalConfigEntryDiff) ConfigAction {
		return entryDiff.NeededAction
	})
	for _, doguDiff := range s.DoguConfigDiffs {
		configActions = append(configActions, util.Map(doguDiff, func(entryDiff DoguConfigEntryDiff) ConfigAction {
			return entryDiff.NeededAction
		})...)
	}
	for _, doguDiff := range s.SensitiveConfigDiffs {
		configActions = append(configActions, util.Map(doguDiff, func(entryDiff SensitiveDoguConfigEntryDiff) ConfigAction {
			return entryDiff.NeededAction
		})...)
	}

	var stringPerAction []string
	var actionsCounter int
	for action, amount := range countByAction(configActions) {
		stringPerAction = append(stringPerAction, fmt.Sprintf("%q: %d", action, amount))
		if action != ConfigActionNone {
			actionsCounter += amount
		}
	}
	slices.Sort(stringPerAction)
	return fmt.Sprintf("%d config changes (%s)", actionsCounter, strings.Join(stringPerAction, ", "))
}

type EcosystemHealthyEvent struct {
	doguHealthIgnored      bool
	componentHealthIgnored bool
}

func (d EcosystemHealthyEvent) Name() string {
	return "EcosystemHealthy"
}

func (d EcosystemHealthyEvent) Message() string {
	return fmt.Sprintf("dogu health ignored: %t; component health ignored: %t", d.doguHealthIgnored, d.componentHealthIgnored)
}

type EcosystemUnhealthyEvent struct {
	HealthResult ecosystem.HealthResult
}

func (d EcosystemUnhealthyEvent) Name() string {
	return "EcosystemUnhealthy"
}

func (d EcosystemUnhealthyEvent) Message() string {
	return "Ecosystem became unhealthy (up-to-date list is in the EcosystemHealthy condition):\n  " + d.HealthResult.String()
}

type ComponentsAppliedEvent struct {
	Diffs ComponentDiffs
}

func (e ComponentsAppliedEvent) Name() string {
	return "ComponentsApplied"
}

func (e ComponentsAppliedEvent) Message() string {
	var buffer bytes.Buffer
	buffer.WriteString("components applied: ")
	var details []string
	for _, diff := range e.Diffs {
		actionsAsStrings := util.Map(diff.NeededActions, func(action Action) string {
			return string(action)
		})
		actions := strings.Join(actionsAsStrings, ", ")
		details = append(details, fmt.Sprintf("%q: [%v]", diff.Name, actions))
	}
	buffer.WriteString(strings.Join(details, ", "))
	return buffer.String()
}

type DogusAppliedEvent struct {
	Diffs DoguDiffs
}

func (e DogusAppliedEvent) Name() string {
	return "DogusApplied"
}

func (e DogusAppliedEvent) Message() string {
	var buffer bytes.Buffer
	buffer.WriteString("dogus applied: ")
	var details []string
	for _, diff := range e.Diffs {
		actionsAsStrings := util.Map(diff.NeededActions, func(action Action) string {
			return string(action)
		})
		actions := strings.Join(actionsAsStrings, ", ")
		details = append(details, fmt.Sprintf("%q: [%v]", diff.DoguName, actions))
	}
	buffer.WriteString(strings.Join(details, ", "))
	return buffer.String()
}

type DogusNotUpToDateEvent struct {
	DogusNotUpToDate []cescommons.SimpleName
}

func (e DogusNotUpToDateEvent) Name() string {
	return "DogusNotUpToDate"
}

func (e DogusNotUpToDateEvent) Message() string {
	dogusNotUpToDate := util.Map(e.DogusNotUpToDate, func(dogu cescommons.SimpleName) string { return string(dogu) })
	slices.Sort(dogusNotUpToDate)
	return fmt.Sprintf("%d dogu(s) not up to date yet: %s", len(dogusNotUpToDate), strings.Join(dogusNotUpToDate, ", "))
}

type BlueprintAppliedEvent struct{}

func (e BlueprintAppliedEvent) Name() string {
	return "BlueprintApplied"
}

func (e BlueprintAppliedEvent) Message() string {
	return "waiting for ecosystem health"
}

type BlueprintStoppedEvent struct{}

func (e BlueprintStoppedEvent) Name() string {
	return "BlueprintStopped"
}

func (e BlueprintStoppedEvent) Message() string {
	return "Blueprint is set as stopped and will not be applied. Remove flag to continue"
}

type ExecutionFailedEvent struct {
	err error
}

func NewExecutionFailedEvent(err error) ExecutionFailedEvent {
	return ExecutionFailedEvent{err: err}
}

func (e ExecutionFailedEvent) Name() string {
	return "ExecutionFailed"
}

func (e ExecutionFailedEvent) Message() string {
	return e.err.Error()
}

type CompletedEvent struct{}

func (e CompletedEvent) Name() string {
	return "completed"
}

func (e CompletedEvent) Message() string {
	return ""
}

type ApplyEcosystemConfigEvent struct{}

func (e ApplyEcosystemConfigEvent) Name() string {
	return "ApplyEcosystemConfig"
}

func (e ApplyEcosystemConfigEvent) Message() string {
	return "apply ecosystem config"
}

type EcosystemConfigAppliedEvent struct{}

func (e EcosystemConfigAppliedEvent) Name() string {
	return "EcosystemConfigApplied"
}

func (e EcosystemConfigAppliedEvent) Message() string {
	return "ecosystem config applied"
}

type AwaitSelfUpgradeEvent struct{}

func (e AwaitSelfUpgradeEvent) Name() string {
	return "AwaitSelfUpgrade"
}

func (e AwaitSelfUpgradeEvent) Message() string {
	return "the operator awaits an upgrade for itself before other changes will be applied"
}

type SelfUpgradeCompletedEvent struct{}

func (e SelfUpgradeCompletedEvent) Name() string {
	return "SelfUpgradeCompleted"
}

func (e SelfUpgradeCompletedEvent) Message() string {
	return "if a self upgrade was necessary, it was successful"
}
