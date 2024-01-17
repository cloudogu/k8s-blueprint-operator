package domain

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvents(t *testing.T) {
	tests := []struct {
		name            string
		event           Event
		expectedName    string
		expectedMessage string
	}{
		{
			name:            "blueprint spec invalid",
			event:           BlueprintSpecInvalidEvent{ValidationError: assert.AnError},
			expectedName:    "BlueprintSpecInvalid",
			expectedMessage: assert.AnError.Error(),
		},
		{
			name:            "blueprint spec statically validated",
			event:           BlueprintSpecStaticallyValidatedEvent{},
			expectedName:    "BlueprintSpecStaticallyValidated",
			expectedMessage: "",
		},
		{
			name:            "blueprint spec validated",
			event:           BlueprintSpecValidatedEvent{},
			expectedName:    "BlueprintSpecValidated",
			expectedMessage: "",
		},
		{
			name:            "dogus healthy",
			event:           DogusHealthyEvent{},
			expectedName:    "DogusHealthy",
			expectedMessage: "",
		},
		{
			name: "dogus unhealthy",
			event: DogusUnhealthyEvent{HealthResult: ecosystem.DoguHealthResult{UnhealthyDogus: []ecosystem.UnhealthyDogu{
				{
					Namespace: "official",
					Name:      "ldap",
					Version:   mustParseVersion("1.2.3-1"),
					Health:    "unavailable",
				},
				{
					Namespace: "premium",
					Name:      "admin",
					Version:   mustParseVersion("3.4.5-3"),
					Health:    "broken",
				},
			}}},
			expectedName:    "DogusUnhealthy",
			expectedMessage: "2 dogus are unhealthy: official/ldap:1.2.3-1 is unavailable, premium/admin:3.4.5-3 is broken",
		},
		{
			name:            "ignore dogu health",
			event:           IgnoreDoguHealthEvent{},
			expectedName:    "IgnoreDoguHealth",
			expectedMessage: "ignore dogu health flag is set; ignoring dogu health",
		},
		{
			name:            "effective blueprint calculated",
			event:           EffectiveBlueprintCalculatedEvent{},
			expectedName:    "EffectiveBlueprintCalculated",
			expectedMessage: "",
		},
		{
			name: "state diff determined",
			event: StateDiffDeterminedEvent{StateDiff: StateDiff{DoguDiffs: DoguDiffs{
				{NeededAction: ActionInstall},
				{NeededAction: ActionUninstall},
				{NeededAction: ActionNone},
				{NeededAction: ActionInstall},
				{NeededAction: ActionUninstall},
				{NeededAction: ActionUninstall},
				{NeededAction: ActionUpgrade},
				{NeededAction: ActionDowngrade},
			}}},
			expectedName:    "StateDiffDetermined",
			expectedMessage: "state diff determined: 8 dogu diffs (2 to install, 1 to upgrade, 3 to delete, 2 others)",
		},
		{
			name:            "In progress",
			event:           InProgressEvent{},
			expectedName:    "InProgress",
			expectedMessage: "",
		},
		{
			name:            "blueprint applied",
			event:           BlueprintAppliedEvent{},
			expectedName:    "BlueprintApplied",
			expectedMessage: "waiting for ecosystem health",
		},
		{
			name:            "completed",
			event:           CompletedEvent{},
			expectedName:    "completed",
			expectedMessage: "",
		},
		{
			name:            "execution failed",
			event:           ExecutionFailedEvent{err: fmt.Errorf("test-error")},
			expectedName:    "ExecutionFailed",
			expectedMessage: "test-error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedName, tt.event.Name())
			assert.Equal(t, tt.expectedMessage, tt.event.Message())
		})
	}
}
