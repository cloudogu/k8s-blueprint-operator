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

type StateDiffDeterminedEvent struct {
	StateDiff StateDiff
}

func (s StateDiffDeterminedEvent) Name() string {
	return "StateDiffDetermined"
}

func (s StateDiffDeterminedEvent) Message() string {
	doguInstall, doguUpgrade, doguUninstall, others := s.StateDiff.DoguDiffs.Statistics()
	// TODO how should the component state diff be incorporated into this event output? An own event for each of them, or one for all ⚔️?
	compInstall, compoUpgrade, compUninstall, compOthers := s.StateDiff.ComponentDiffs.Statistics()

	return fmt.Sprintf("state diff determined: %d dogu diffs (%d to install, %d to upgrade, %d to delete, %d others); %d component diffs (%d to install, %d to upgrade, %d to delete, %d others)",
		len(s.StateDiff.DoguDiffs), doguInstall, doguUpgrade, doguUninstall, others,
		len(s.StateDiff.ComponentDiffs), compInstall, compoUpgrade, compUninstall, compOthers,
	)
}

type EcosystemHealthyUpfrontEvent struct {
	doguHealthIgnored bool
}

func (d EcosystemHealthyUpfrontEvent) Name() string {
	return "EcosystemHealthyUpfront"
}

func (d EcosystemHealthyUpfrontEvent) Message() string {
	return fmt.Sprintf("dogu health ignored: %t", d.doguHealthIgnored)
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

type InProgressEvent struct{}

func (e InProgressEvent) Name() string {
	return "InProgress"
}

func (e InProgressEvent) Message() string {
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

type CompletedEvent struct{}

func (e CompletedEvent) Name() string {
	return "completed"
}

func (e CompletedEvent) Message() string {
	return ""
}
