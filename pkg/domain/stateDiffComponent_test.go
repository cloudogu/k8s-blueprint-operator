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
				blueprintComponent: mockTargetComponent(version3211, TargetStatePresent),
				installedComponent: mockComponentInstallation(version3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3211, TargetStatePresent),
				Expected:     mockComponentDiffState(version3211, TargetStatePresent),
				NeededAction: ActionNone,
			},
		},
		{
			name: "install",
			args: args{
				blueprintComponent: mockTargetComponent(version3211, TargetStatePresent),
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(core.Version{}, TargetStateAbsent),
				Expected:     mockComponentDiffState(version3211, TargetStatePresent),
				NeededAction: ActionInstall,
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintComponent: mockTargetComponent(core.Version{}, TargetStateAbsent),
				installedComponent: mockComponentInstallation(version3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3211, TargetStatePresent),
				Expected:     mockComponentDiffState(core.Version{}, TargetStateAbsent),
				NeededAction: ActionUninstall,
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(version3212, TargetStatePresent),
				installedComponent: mockComponentInstallation(version3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3211, TargetStatePresent),
				Expected:     mockComponentDiffState(version3212, TargetStatePresent),
				NeededAction: ActionUpgrade,
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintComponent: mockTargetComponent(version3211, TargetStatePresent),
				installedComponent: mockComponentInstallation(version3212),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3212, TargetStatePresent),
				Expected:     mockComponentDiffState(version3211, TargetStatePresent),
				NeededAction: ActionDowngrade,
			},
		},
		{
			name: "ignore present component, no action",
			args: args{
				blueprintComponent: nil,
				installedComponent: mockComponentInstallation(version3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(version3211, TargetStatePresent),
				Expected:     mockComponentDiffState(version3211, TargetStatePresent),
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

func TestComponentDiffs_Statistics(t *testing.T) {
	tests := []struct {
		name            string
		dd              ComponentDiffs
		wantToInstall   int
		wantToUpgrade   int
		wantToUninstall int
		wantOther       int
	}{
		{
			name:            "0 overall",
			dd:              ComponentDiffs{},
			wantToInstall:   0,
			wantToUpgrade:   0,
			wantToUninstall: 0,
			wantOther:       0,
		},
		{
			name: "4 to install, 3 to upgrade, 2 to uninstall, 3 other",
			dd: ComponentDiffs{
				{NeededAction: ActionNone},
				{NeededAction: ActionInstall},
				{NeededAction: ActionUninstall},
				{NeededAction: ActionInstall},
				{NeededAction: ActionUpgrade},
				{NeededAction: ActionInstall},
				{NeededAction: ActionDowngrade},
				{NeededAction: ActionUninstall},
				{NeededAction: ActionInstall},
				{NeededAction: ActionUpgrade},
				{NeededAction: ActionUpgrade},
			},
			wantToInstall:   4,
			wantToUpgrade:   3,
			wantToUninstall: 2,
			wantOther:       2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToInstall, gotToUpgrade, gotToUninstall, gotOther := tt.dd.Statistics()
			assert.Equalf(t, tt.wantToInstall, gotToInstall, "Statistics()")
			assert.Equalf(t, tt.wantToUpgrade, gotToUpgrade, "Statistics()")
			assert.Equalf(t, tt.wantToUninstall, gotToUninstall, "Statistics()")
			assert.Equalf(t, tt.wantOther, gotOther, "Statistics()")
		})
	}
}

func TestComponentDiff_String(t *testing.T) {
	actual := ComponentDiffState{
		Version:           version3211,
		InstallationState: TargetStatePresent,
	}
	expected := ComponentDiffState{
		Version:           version3212,
		InstallationState: TargetStatePresent,
	}
	diff := &ComponentDiff{
		Name:         testComponentName,
		Actual:       actual,
		Expected:     expected,
		NeededAction: ActionInstall,
	}

	assert.Equal(t, "{"+
		"Name: \"my-component\", "+
		"Actual: {Version: \"3.2.1-1\", InstallationState: \"present\"}, "+
		"Expected: {Version: \"3.2.1-2\", InstallationState: \"present\"}, "+
		"NeededAction: \"install\""+
		"}", diff.String())
}

func TestComponentDiffState_String(t *testing.T) {
	diff := &ComponentDiffState{
		Version:           version3211,
		InstallationState: TargetStatePresent,
	}

	assert.Equal(t, "{Version: \"3.2.1-1\", InstallationState: \"present\"}", diff.String())
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

func Test_determineComponentDiffs(t *testing.T) {
	type args struct {
		logger              logr.Logger
		blueprintComponents []Component
		installedComponents map[string]*ecosystem.ComponentInstallation
	}
	tests := []struct {
		name string
		args args
		want []ComponentDiff
	}{
		{
			name: "no components",
			args: args{
				blueprintComponents: nil,
				installedComponents: nil,
			},
			want: []ComponentDiff{},
		},
		{
			name: "a not installed component in the blueprint",
			args: args{
				blueprintComponents: []Component{
					{
						Name:        testComponentName,
						Version:     version3211,
						TargetState: TargetStatePresent,
					},
				},
				installedComponents: nil,
			},
			want: []ComponentDiff{
				{
					Name: testComponentName,
					Actual: ComponentDiffState{
						InstallationState: TargetStateAbsent,
					},
					Expected: ComponentDiffState{
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionInstall,
				},
			},
		},
		{
			name: "an installed component which is not in the blueprint",
			args: args{
				blueprintComponents: nil,
				installedComponents: map[string]*ecosystem.ComponentInstallation{
					testComponentName: {
						Name:    testComponentName,
						Version: version3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName,
					Actual: ComponentDiffState{
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					Expected: ComponentDiffState{
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionNone,
				},
			},
		},
		{
			name: "an installed component which is also in the blueprint",
			args: args{
				blueprintComponents: []Component{
					{
						Name:        testComponentName,
						Version:     version3212,
						TargetState: TargetStatePresent,
					},
				},
				installedComponents: map[string]*ecosystem.ComponentInstallation{
					testComponentName: {
						Name:    testComponentName,
						Version: version3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName,
					Actual: ComponentDiffState{
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					Expected: ComponentDiffState{
						Version:           version3212,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionUpgrade,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determineComponentDiffs(tt.args.logger, tt.args.blueprintComponents, tt.args.installedComponents), "determineComponentDiffs(%v, %v, %v)", tt.args.logger, tt.args.blueprintComponents, tt.args.installedComponents)
		})
	}
}
