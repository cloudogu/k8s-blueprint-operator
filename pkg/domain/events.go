package domain

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"slices"
	"strings"
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

type BlueprintSpecStaticallyValidatedEvent struct{}

func (b BlueprintSpecStaticallyValidatedEvent) Name() string {
	return "BlueprintSpecStaticallyValidated"
}

func (b BlueprintSpecStaticallyValidatedEvent) Message() string {
	return ""
}

type BlueprintSpecValidatedEvent struct{}

func (b BlueprintSpecValidatedEvent) Name() string {
	return "BlueprintSpecValidated"
}

func (b BlueprintSpecValidatedEvent) Message() string {
	return ""
}

type EffectiveBlueprintCalculatedEvent struct {
	Result EffectiveBlueprint
}

func (e EffectiveBlueprintCalculatedEvent) Name() string {
	return "EffectiveBlueprintCalculated"
}

func (e EffectiveBlueprintCalculatedEvent) Message() string {
	return ""
}

type GlobalConfigDiffDeterminedEvent struct {
	GlobalConfigDiffs GlobalConfigDiffs
}

func (e GlobalConfigDiffDeterminedEvent) Name() string {
	return "GlobalConfigDiffDetermined"
}

func (e GlobalConfigDiffDeterminedEvent) Message() string {
	var stringPerAction []string
	var actionsCounter int
	for action, amount := range e.GlobalConfigDiffs.countByAction() {
		stringPerAction = append(stringPerAction, fmt.Sprintf("%q: %d", action, amount))
		actionsCounter += amount
	}
	slices.Sort(stringPerAction)
	return fmt.Sprintf("global config diff determined: %d actions (%s)", actionsCounter, strings.Join(stringPerAction, ", "))
}

type DoguConfigDiffDeterminedEvent struct {
	CombinedDogusConfigDiffs map[common.SimpleDoguName]CombinedDoguConfigDiffs
}

func (e DoguConfigDiffDeterminedEvent) Name() string {
	return "DoguConfigDiffDetermined"
}

func (e DoguConfigDiffDeterminedEvent) Message() string {
	var stringPerAction []string
	var actionsCounter int
	for action, amount := range countByAction(e.CombinedDogusConfigDiffs) {
		stringPerAction = append(stringPerAction, fmt.Sprintf("%q: %d", action, amount))
		actionsCounter += amount
	}
	slices.Sort(stringPerAction)
	return fmt.Sprintf("dogu config diff determined: %d actions (%s)", actionsCounter, strings.Join(stringPerAction, ", "))
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

// StateDiffDoguDeterminedEvent provides event information over detected changes regarding dogus.
type StateDiffDoguDeterminedEvent struct {
	doguDiffs DoguDiffs
}

func newStateDiffDoguEvent(doguDiffs DoguDiffs) StateDiffDoguDeterminedEvent {
	return StateDiffDoguDeterminedEvent{
		doguDiffs: doguDiffs,
	}
}

// Name contains the StateDiffDoguDeterminedEvent display name.
func (s StateDiffDoguDeterminedEvent) Name() string {
	return "StateDiffDoguDetermined"
}

const groupedDoguProxyAction = "update reverse proxy"

// Message contains the StateDiffDoguDeterminedEvent's statistics message.
func (s StateDiffDoguDeterminedEvent) Message() string {
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

	message, amount := getActionAmountMessage(amountActions)

	return fmt.Sprintf("dogu state diff determined: %d actions (%s)", amount, message)
}

type EcosystemHealthyUpfrontEvent struct {
	doguHealthIgnored      bool
	componentHealthIgnored bool
}

func (d EcosystemHealthyUpfrontEvent) Name() string {
	return "EcosystemHealthyUpfront"
}

func (d EcosystemHealthyUpfrontEvent) Message() string {
	return fmt.Sprintf("dogu health ignored: %t; component health ignored: %t", d.doguHealthIgnored, d.componentHealthIgnored)
}

type EcosystemUnhealthyUpfrontEvent struct {
	HealthResult ecosystem.HealthResult
}

func (d EcosystemUnhealthyUpfrontEvent) Name() string {
	return "EcosystemUnhealthyUpfront"
}

func (d EcosystemUnhealthyUpfrontEvent) Message() string {
	return d.HealthResult.String()
}

type BlueprintDryRunEvent struct{}

func (b BlueprintDryRunEvent) Name() string {
	return "BlueprintDryRun"
}

func (b BlueprintDryRunEvent) Message() string {
	return "Executed blueprint in dry run mode. Remove flag to continue"
}

type BlueprintApplicationPreProcessedEvent struct {
}

func (e BlueprintApplicationPreProcessedEvent) Name() string {
	return "BlueprintApplicationPreProcessed"
}

func (e BlueprintApplicationPreProcessedEvent) Message() string {
	return "maintenance mode activated"
}

type InProgressEvent struct{}

func (e InProgressEvent) Name() string {
	return "InProgress"
}

func (e InProgressEvent) Message() string {
	return ""
}

type BlueprintAppliedEvent struct{}

func (e BlueprintAppliedEvent) Name() string {
	return "BlueprintApplied"
}

func (e BlueprintAppliedEvent) Message() string {
	return "waiting for ecosystem health"
}

type EcosystemHealthyAfterwardsEvent struct{}

func (e EcosystemHealthyAfterwardsEvent) Name() string {
	return "EcosystemHealthyAfterwards"
}

func (e EcosystemHealthyAfterwardsEvent) Message() string {
	return ""
}

type EcosystemUnhealthyAfterwardsEvent struct {
	HealthResult ecosystem.HealthResult
}

func (e EcosystemUnhealthyAfterwardsEvent) Name() string {
	return "EcosystemUnhealthyAfterwards"
}

func (e EcosystemUnhealthyAfterwardsEvent) Message() string {
	return e.HealthResult.String()
}

type SensitiveConfigDataCensoredEvent struct{}

func (e SensitiveConfigDataCensoredEvent) Name() string {
	return "sensitiveConfigDataCensored"
}

func (e SensitiveConfigDataCensoredEvent) Message() string {
	return ""
}

type ExecutionFailedEvent struct {
	err error
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
	return "maintenance mode deactivated"
}

type ApplyRegistryConfigEvent struct{}

func (e ApplyRegistryConfigEvent) Name() string {
	return "ApplyRegistryConfig"
}

func (e ApplyRegistryConfigEvent) Message() string {
	return "apply registry config"
}

type ApplyRegistryConfigFailedEvent struct {
	err error
}

func (e ApplyRegistryConfigFailedEvent) Name() string {
	return "ApplyDoguConfigFailed"
}

func (e ApplyRegistryConfigFailedEvent) Message() string {
	return e.err.Error()
}

type RegistryConfigAppliedEvent struct{}

func (e RegistryConfigAppliedEvent) Name() string {
	return "RegistryConfigApplied"
}

func (e RegistryConfigAppliedEvent) Message() string {
	return "registry config applied"
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
