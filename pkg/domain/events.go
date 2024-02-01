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
	toInstall, toUpgrade, toUninstall, others := s.StateDiff.DoguDiffs.Statistics()
	return fmt.Sprintf("state diff determined: %d dogu diffs (%d to install, %d to upgrade, %d to delete, %d others)",
		len(s.StateDiff.DoguDiffs), toInstall, toUpgrade, toUninstall, others)
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
