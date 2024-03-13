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
					{NeededActions: []Action{ActionNone}},
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUpgrade, ActionUpdateDoguResourceMinVolumeSize, ActionUpdateDoguProxyBodySize, ActionUpdateDoguProxyRewriteTarget, ActionUpdateDoguProxyAdditionalConfig}},
					{NeededActions: []Action{ActionDowngrade}},
				}),
			expectedName:    "StateDiffDoguDetermined",
			expectedMessage: "state diff determined: 8 dogu diffs (2 to install, 1 to upgrade, 3 to delete, 2 others)\ndogu config diffs: (1 to update resource config, 3 to update reverse proxy config)",
		},
		{
			name: "component state diff determined",
			event: newStateDiffComponentEvent(
				ComponentDiffs{
					{NeededAction: ActionInstall},
					{NeededAction: ActionUninstall},
					{NeededAction: ActionNone},
					{NeededAction: ActionInstall},
					{NeededAction: ActionUninstall},
					{NeededAction: ActionUninstall},
					{NeededAction: ActionUpgrade},
					{NeededAction: ActionDowngrade},
				}),
			expectedName:    "StateDiffComponentDetermined",
			expectedMessage: "state diff determined: 8 component diffs (2 to install, 1 to upgrade, 3 to delete, 2 others)",
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
						{NeededAction: ConfigActionSetToEncrypt},
						{NeededAction: ConfigActionRemove},
					},
				},
			}},
			expectedName:    "DoguConfigDiffDetermined",
			expectedMessage: "dogu config diff determined: 6 actions (\"none\": 2, \"remove\": 2, \"set\": 1, \"setToEncrypt\": 1)",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedName, tt.event.Name())
			assert.Equal(t, tt.expectedMessage, tt.event.Message())
		})
	}
}
