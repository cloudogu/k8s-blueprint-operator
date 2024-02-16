package domain

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
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

// StateDiffDeterminedEvent is an abstract type that aggregates shared statistic functionality over different stateDiffDetermined events.
type StateDiffDeterminedEvent struct {
	// DiffCount contains the total number of detected actions.
	DiffCount int
	// EventSubject contains the subject of the dected actions, f. i. "dogu diffs".
	EventSubject string
	// ToInstall contains the number of actions that introduce a new subject, f. i. the number of dogus to be installed.
	ToInstall int
	// ToUpgrade contains the number of actions that change an already existing subject, f. i. the number of dogus to be upgraded.
	ToUpgrade int
	// ToUninstall contains the number of actions that remove an already existing subject, f. i. the number of dogus to be deleted.
	ToUninstall int
	// OtherActions contains the number of actions that do not fall into the other three category, f. i. the number of dogus being ignored
	OtherActions int
}

func (s StateDiffDeterminedEvent) buildMessage() string {
	return fmt.Sprintf("state diff determined: %d %s (%d to install, %d to upgrade, %d to delete, %d others)",
		s.DiffCount, s.EventSubject, s.ToInstall, s.ToUpgrade, s.ToUninstall, s.OtherActions)
}

// StateDiffComponentDeterminedEvent provides event information over detected changes regarding components.
type StateDiffComponentDeterminedEvent struct {
	StateDiffDeterminedEvent
}

func newStateDiffComponentEvent(componentDiffs ComponentDiffs) StateDiffComponentDeterminedEvent {
	install, upgrade, uninstall, other := componentDiffs.Statistics()

	return StateDiffComponentDeterminedEvent{StateDiffDeterminedEvent: StateDiffDeterminedEvent{
		DiffCount:    len(componentDiffs),
		EventSubject: "component diffs",
		ToInstall:    install,
		ToUpgrade:    upgrade,
		ToUninstall:  uninstall,
		OtherActions: other,
	}}
}

// Name contains the StateDiffComponentDeterminedEvent display name.
func (s StateDiffComponentDeterminedEvent) Name() string {
	return "StateDiffComponentDetermined"
}

// Message contains the StateDiffComponentDeterminedEvent's statistics message.
func (s StateDiffComponentDeterminedEvent) Message() string {
	return s.buildMessage()
}

// StateDiffDoguDeterminedEvent provides event information over detected changes regarding dogus.
type StateDiffDoguDeterminedEvent struct {
	StateDiffDeterminedEvent
}

func newStateDiffDoguEvent(doguDiffs DoguDiffs) StateDiffDoguDeterminedEvent {
	install, upgrade, uninstall, other := doguDiffs.Statistics()

	return StateDiffDoguDeterminedEvent{StateDiffDeterminedEvent: StateDiffDeterminedEvent{
		DiffCount:    len(doguDiffs),
		EventSubject: "dogu diffs",
		ToInstall:    install,
		ToUpgrade:    upgrade,
		ToUninstall:  uninstall,
		OtherActions: other,
	}}
}

// Name contains the StateDiffDoguDeterminedEvent display name.
func (s StateDiffDoguDeterminedEvent) Name() string {
	return "StateDiffDoguDetermined"
}

// Message contains the StateDiffDoguDeterminedEvent's statistics message.
func (s StateDiffDoguDeterminedEvent) Message() string {
	return s.buildMessage()
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
