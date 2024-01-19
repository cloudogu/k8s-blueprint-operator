package domain

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

const testComponentName = "my-component"

func Test_determineComponentDiff(t *testing.T) {
	type args struct {
		logger             logr.Logger
		blueprintComponent *Component
		installedComponent *ecosystem.ComponentInstallation
	}
	tests := []struct {
		name string
		args args
		want ComponentDiff
	}{
		{
			name: "equal, no action",
			args: args{
				blueprintComponent: mockTargetComponent(version3_2_1_1, TargetStatePresent),
				installedComponent: mockComponentInstallation(version3_2_1_1),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				Expected:     mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				NeededAction: ActionNone,
			},
		},
		{
			name: "install",
			args: args{
				blueprintComponent: mockTargetComponent(version3_2_1_1, TargetStatePresent),
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(core.Version{}, TargetStateAbsent),
				Expected:     mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				NeededAction: ActionInstall,
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintComponent: mockTargetComponent(core.Version{}, TargetStateAbsent),
				installedComponent: mockComponentInstallation(version3_2_1_1),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				Expected:     mockComponentDiffState(core.Version{}, TargetStateAbsent),
				NeededAction: ActionUninstall,
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(version3_2_1_2, TargetStatePresent),
				installedComponent: mockComponentInstallation(version3_2_1_1),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				Expected:     mockComponentDiffState(version3_2_1_2, TargetStatePresent),
				NeededAction: ActionUpgrade,
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintComponent: mockTargetComponent(version3_2_1_1, TargetStatePresent),
				installedComponent: mockComponentInstallation(version3_2_1_2),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3_2_1_2, TargetStatePresent),
				Expected:     mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				NeededAction: ActionDowngrade,
			},
		},
		{
			name: "ignore present dogu, no action",
			args: args{
				blueprintComponent: nil,
				installedComponent: mockComponentInstallation(version3_2_1_1),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				Expected:     mockComponentDiffState(version3_2_1_1, TargetStatePresent),
				NeededAction: ActionNone,
			},
		},
		{
			name: "should stay absent, no action", // this is empty set comparison is weird and should basically not occur
			args: args{
				blueprintComponent: nil,
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:         "",
				Actual:       ComponentDiffState{InstallationState: TargetStateAbsent},
				Expected:     ComponentDiffState{InstallationState: TargetStateAbsent},
				NeededAction: ActionNone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determineComponentDiff(tt.args.logger, tt.args.blueprintComponent, tt.args.installedComponent), "determineComponentDiff(%v, %v, %v)", tt.args.logger, tt.args.blueprintComponent, tt.args.installedComponent)
		})
	}
}

func mockTargetComponent(version core.Version, state TargetState) *Component {
	return &Component{
		Name:        testComponentName,
		Version:     version,
		TargetState: state,
	}
}

func mockComponentInstallation(version core.Version) *ecosystem.ComponentInstallation {
	return &ecosystem.ComponentInstallation{
		Name:    testComponentName,
		Version: version,
	}
}

func mockComponentDiffState(version core.Version, state TargetState) ComponentDiffState {
	return ComponentDiffState{
		Version:           version,
		InstallationState: state,
	}
}
