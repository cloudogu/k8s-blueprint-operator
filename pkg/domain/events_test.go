package domain

import (
	"fmt"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
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
			name:            "ecosystem healthy",
			event:           EcosystemHealthyEvent{},
			expectedName:    "EcosystemHealthy",
			expectedMessage: "dogu health ignored: false",
		},
		{
			name:            "ignore dogu health",
			event:           EcosystemHealthyEvent{doguHealthIgnored: true},
			expectedName:    "EcosystemHealthy",
			expectedMessage: "dogu health ignored: true",
		},
		{
			name: "ecosystem unhealthy upfront",
			event: EcosystemUnhealthyEvent{
				HealthResult: ecosystem.HealthResult{
					DoguHealth: ecosystem.DoguHealthResult{
						DogusByStatus: map[ecosystem.HealthStatus][]cescommons.SimpleName{
							ecosystem.AvailableHealthStatus:   {"postgresql"},
							ecosystem.UnavailableHealthStatus: {"ldap"},
							ecosystem.PendingHealthStatus:     {"admin"},
						},
					},
				},
			},
			expectedName:    "EcosystemUnhealthy",
			expectedMessage: "Ecosystem became unhealthy. Reason:\n  2 dogu(s) are unhealthy: admin, ldap",
		},
		{
			name: "dogu state diff determined",
			event: newStateDiffEvent(
				StateDiff{DoguDiffs: DoguDiffs{
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUninstall}},
					{NeededActions: []Action{ActionUpgrade, ActionUpdateDoguResourceMinVolumeSize, ActionUpdateDoguProxyBodySize, ActionUpdateDoguProxyRewriteTarget, ActionUpdateDoguProxyAdditionalConfig}},
					{NeededActions: []Action{ActionDowngrade}},
				}}),
			expectedName:    "StateDiffDetermined",
			expectedMessage: "state diff determined:\n  0 config changes ()\n  11 dogu actions (\"downgrade\": 1, \"install\": 2, \"uninstall\": 3, \"update resource minimum volume size\": 1, \"update reverse proxy\": 3, \"upgrade\": 1)",
		},
		{
			name: "config diff determined",
			event: newStateDiffEvent(StateDiff{
				GlobalConfigDiffs: GlobalConfigDiffs{
					{NeededAction: ConfigActionNone},
					{NeededAction: ConfigActionNone},
					{NeededAction: ConfigActionSet},
					{NeededAction: ConfigActionRemove},
				},
				DoguConfigDiffs: map[cescommons.SimpleName]DoguConfigDiffs{
					"dogu1": []DoguConfigEntryDiff{
						{NeededAction: ConfigActionNone},
						{NeededAction: ConfigActionSet},
						{NeededAction: ConfigActionRemove},
					},
				},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
					"dogu1": []SensitiveDoguConfigEntryDiff{
						{NeededAction: ConfigActionNone},
						{NeededAction: ConfigActionSet},
						{NeededAction: ConfigActionRemove},
					},
				},
			}),
			expectedName:    "StateDiffDetermined",
			expectedMessage: "state diff determined:\n  6 config changes (\"none\": 4, \"remove\": 3, \"set\": 3)\n  0 dogu actions ()",
		},
		{
			name: "config and dogu diff determined",
			event: newStateDiffEvent(StateDiff{
				DoguDiffs: DoguDiffs{
					{NeededActions: []Action{ActionInstall}},
					{NeededActions: []Action{ActionUninstall}},
				},
				GlobalConfigDiffs: GlobalConfigDiffs{
					{NeededAction: ConfigActionSet},
					{NeededAction: ConfigActionRemove},
				},
			}),
			expectedName:    "StateDiffDetermined",
			expectedMessage: "state diff determined:\n  2 config changes (\"remove\": 1, \"set\": 1)\n  2 dogu actions (\"install\": 1, \"uninstall\": 1)",
		},
		{
			name: "config references missing",
			event: NewMissingConfigReferencesEvent(
				assert.AnError,
			),
			expectedName:    "MissingConfigReferences",
			expectedMessage: assert.AnError.Error(),
		},
		{
			name: "dogus applied",
			event: DogusAppliedEvent{
				Diffs: DoguDiffs{
					{
						DoguName: "jenkins",
						NeededActions: []Action{
							ActionUpgrade, ActionSwitchDoguNamespace,
						},
					},
				},
			},
			expectedName:    "DogusApplied",
			expectedMessage: "dogus applied: \"jenkins\": [upgrade, dogu namespace switch]",
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
		{
			name:            "apply ecosystem config",
			event:           ApplyEcosystemConfigEvent{},
			expectedName:    "ApplyEcosystemConfig",
			expectedMessage: "apply ecosystem config",
		},
		{
			name:            "ecosystem config applied",
			event:           EcosystemConfigAppliedEvent{},
			expectedName:    "EcosystemConfigApplied",
			expectedMessage: "ecosystem config applied",
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
