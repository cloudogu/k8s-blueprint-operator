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
			expectedMessage: "dogu health ignored: false; component health ignored: false",
		},
		{
			name:            "ignore dogu health",
			event:           EcosystemHealthyEvent{doguHealthIgnored: true},
			expectedName:    "EcosystemHealthy",
			expectedMessage: "dogu health ignored: true; component health ignored: false",
		},
		{
			name:            "ignore component health",
			event:           EcosystemHealthyEvent{componentHealthIgnored: true},
			expectedName:    "EcosystemHealthy",
			expectedMessage: "dogu health ignored: false; component health ignored: true",
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
			expectedMessage: "ecosystem health:\n  2 dogu(s) are unhealthy: admin, ldap\n  0 component(s) are unhealthy: ",
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
			expectedMessage: "global config diff determined: 2 changes (\"none\": 2, \"remove\": 1, \"set\": 1)",
		},
		{
			name: "dogu config diff determined",
			event: DoguConfigDiffDeterminedEvent{
				DoguConfigDiffs: map[cescommons.SimpleName]DoguConfigDiffs{
					"dogu1": []DoguConfigEntryDiff{
						{NeededAction: ConfigActionNone},
						{NeededAction: ConfigActionSet},
						{NeededAction: ConfigActionRemove},
					},
				},
			},
			expectedName:    "DoguConfigDiffDetermined",
			expectedMessage: "dogu config diff determined: 2 changes (\"none\": 1, \"remove\": 1, \"set\": 1)",
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
			name: "sensitive dogu config diff determined",
			event: SensitiveDoguConfigDiffDeterminedEvent{
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]SensitiveDoguConfigDiffs{
					"dogu1": []SensitiveDoguConfigEntryDiff{
						{NeededAction: ConfigActionNone},
						{NeededAction: ConfigActionSet},
						{NeededAction: ConfigActionRemove},
					},
				},
			},
			expectedName:    "SensitiveDoguConfigDiffDetermined",
			expectedMessage: "sensitive dogu config diff determined: 2 changes (\"none\": 1, \"remove\": 1, \"set\": 1)",
		},
		{
			name: "components applied",
			event: ComponentsAppliedEvent{
				Diffs: ComponentDiffs{
					{
						Name: "dogu-operator",
						NeededActions: []Action{
							ActionUpgrade, ActionSwitchComponentNamespace,
						},
					},
				},
			},
			expectedName:    "ComponentsApplied",
			expectedMessage: "components applied: \"dogu-operator\": [upgrade, component namespace switch]",
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
			name:            "applying ecosystem config failed",
			event:           ApplyEcosystemConfigFailedEvent{fmt.Errorf("test-error")},
			expectedName:    "ApplyEcosystemConfigFailed",
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
