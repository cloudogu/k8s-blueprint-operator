package domain

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/util"
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

type DogusHealthyEvent struct{}

func (d DogusHealthyEvent) Name() string {
	return "DogusHealthy"
}

func (d DogusHealthyEvent) Message() string {
	return ""
}

type IgnoreDoguHealthEvent struct{}

func (i IgnoreDoguHealthEvent) Name() string {
	return "IgnoreDoguHealth"
}

func (i IgnoreDoguHealthEvent) Message() string {
	return "ignore dogu health flag is set; ignoring dogu health"
}

type DogusUnhealthyEvent struct {
	HealthResult ecosystem.DoguHealthResult
}

func (d DogusUnhealthyEvent) Name() string {
	return "DogusUnhealthy"
}

func (d DogusUnhealthyEvent) Message() string {
	unhealthyDogus := util.Map(d.HealthResult.UnhealthyDogus, ecosystem.UnhealthyDogu.String)
	slices.Sort(unhealthyDogus)
	return fmt.Sprintf("%d dogus are unhealthy: %s", len(unhealthyDogus),
		strings.Join(unhealthyDogus, ", "))
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

type ExecutionFailedEvent struct {
	err error
}

func (e ExecutionFailedEvent) Name() string {
	return "ExecutionFailed"
}

func (e ExecutionFailedEvent) Message() string {
	return e.err.Error()
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

type CompletedEvent struct{}

func (e CompletedEvent) Name() string {
	return "completed"
}

func (e CompletedEvent) Message() string {
	return ""
}
