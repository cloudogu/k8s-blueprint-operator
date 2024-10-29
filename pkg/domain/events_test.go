package domain

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
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
			name:            "blueprint dry run",
			event:           BlueprintDryRunEvent{},
			expectedName:    "BlueprintDryRun",
			expectedMessage: "Executed blueprint in dry run mode. Remove flag to continue",
		},
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
			name:            "ecosystem healthy",
			event:           EcosystemHealthyUpfrontEvent{},
			expectedName:    "EcosystemHealthyUpfront",
			expectedMessage: "dogu health ignored: false; component health ignored: false",
		},
		{
			name:            "ignore dogu health",
			event:           EcosystemHealthyUpfrontEvent{doguHealthIgnored: true},
			expectedName:    "EcosystemHealthyUpfront",
			expectedMessage: "dogu health ignored: true; component health ignored: false",
		},
		{
			name:            "ignore component health",
			event:           EcosystemHealthyUpfrontEvent{componentHealthIgnored: true},
			expectedName:    "EcosystemHealthyUpfront",
			expectedMessage: "dogu health ignored: false; component health ignored: true",
		},
		{
			name: "ecosystem unhealthy upfront",
			event: EcosystemUnhealthyUpfrontEvent{
				HealthResult: ecosystem.HealthResult{
					DoguHealth: ecosystem.DoguHealthResult{
						DogusByStatus: map[ecosystem.HealthStatus][]common.SimpleDoguName{
							ecosystem.AvailableHealthStatus:   {"postgresql"},
							ecosystem.UnavailableHealthStatus: {"ldap"},
							ecosystem.PendingHealthStatus:     {"admin"},
						},
					},
				},
			},
			expectedName:    "EcosystemUnhealthyUpfront",
			expectedMessage: "ecosystem health:\n  2 dogu(s) are unhealthy: admin, ldap\n  0 component(s) are unhealthy: ",
		},
		{
			name:            "effective blueprint calculated",
			event:           EffectiveBlueprintCalculatedEvent{},
			expectedName:    "EffectiveBlueprintCalculated",
			expectedMessage: "",
		},
		{
			name: "dogu state diff determined",
			event: newStateDiffDoguEvent(
				DoguDiffs{
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUpgrade, ActionUpdateDoguResourceMinVolumeSize, ActionUpdateDoguProxyBodySize, ActionUpdateDoguProxyRewriteTarget, ActionUpdateDoguProxyAdditionalConfig}},
					{NeededActions: []Action{ActionDowngrade}},
				}),
			expectedName:    "StateDiffDoguDetermined",
			expectedMessage: "dogu state diff determined: 11 actions (\"downgrade\": 1, \"install\": 2, \"uninstall\": 3, \"update resource minimum volume size\": 1, \"update reverse proxy\": 3, \"upgrade\": 1)",
		},
		{
			name: "component state diff determined",
			event: newStateDiffComponentEvent(
				ComponentDiffs{
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUpgrade, ActionUpdateComponentDeployConfig, ActionSwitchComponentNamespace}},
					{NeededActions: []Action{ActionDowngrade}},
				}),
			expectedName:    "StateDiffComponentDetermined",
			expectedMessage: "component state diff determined: 9 actions (\"component namespace switch\": 1, \"downgrade\": 1, \"install\": 2, \"uninstall\": 3, \"update component package config\": 1, \"upgrade\": 1)",
		},
		{
			name: "global config diff determined",
			event: GlobalConfigDiffDeterminedEvent{GlobalConfigDiffs: GlobalConfigDiffs{
				{NeededAction: ConfigActionNone},
				{NeededAction: ConfigActionNone},
				{NeededAction: ConfigActionSet},
				{NeededAction: ConfigActionRemove},
			}},
			expectedName:    "GlobalConfigDiffDetermined",
			expectedMessage: "global config diff determined: 4 actions (\"none\": 2, \"remove\": 1, \"set\": 1)",
		},
		{
			name: "dogu config diff determined",
			event: DoguConfigDiffDeterminedEvent{CombinedDogusConfigDiffs: map[common.SimpleDoguName]CombinedDoguConfigDiffs{
				"dogu1": {
					DoguConfigDiff: []DoguConfigEntryDiff{
						{NeededAction: ConfigActionNone},
						{NeededAction: ConfigActionSet},
						{NeededAction: ConfigActionRemove},
					},
					SensitiveDoguConfigDiff: []SensitiveDoguConfigEntryDiff{
						{NeededAction: ConfigActionNone},
						{NeededAction: ConfigActionSet},
						{NeededAction: ConfigActionRemove},
					},
				},
			}},
			expectedName:    "DoguConfigDiffDetermined",
			expectedMessage: "dogu config diff determined: 6 actions (\"none\": 2, \"remove\": 2, \"set\": 2)",
		},
		{
			name:            "blueprint application pre-processed",
			event:           BlueprintApplicationPreProcessedEvent{},
			expectedName:    "BlueprintApplicationPreProcessed",
			expectedMessage: "maintenance mode activated",
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
			expectedMessage: "maintenance mode deactivated",
		},
		{
			name:            "execution failed",
			event:           ExecutionFailedEvent{err: fmt.Errorf("test-error")},
			expectedName:    "ExecutionFailed",
			expectedMessage: "test-error",
		},
		{
			name:            "apply registry config",
			event:           ApplyRegistryConfigEvent{},
			expectedName:    "ApplyRegistryConfig",
			expectedMessage: "apply registry config",
		},
		{
			name:            "registry config applied",
			event:           RegistryConfigAppliedEvent{},
			expectedName:    "RegistryConfigApplied",
			expectedMessage: "registry config applied",
		},
		{
			name:            "registry config apply failed",
			event:           ApplyRegistryConfigFailedEvent{fmt.Errorf("test-error")},
			expectedName:    "ApplyDoguConfigFailed",
			expectedMessage: "test-error",
		},
		{
			name:            "await self upgrade",
			event:           AwaitSelfUpgradeEvent{},
			expectedName:    "AwaitSelfUpgrade",
			expectedMessage: "the operator awaits an upgrade for itself before other changes will be applied",
		},
		{
			name:            "self upgrade completed",
			event:           SelfUpgradeCompletedEvent{},
			expectedName:    "SelfUpgradeCompleted",
			expectedMessage: "if a self upgrade was necessary, it was successful",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedName, tt.event.Name())
			assert.Equal(t, tt.expectedMessage, tt.event.Message())
		})
	}
}
