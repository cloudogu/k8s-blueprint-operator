package domain

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/Masterminds/semver/v3"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

const testComponentName = "my-component"

var (
	compVersion3211 = semver.MustParse("3.2.1-1")
)

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
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(compVersion3211, TargetStatePresent),
				Expected:     mockComponentDiffState(compVersion3211, TargetStatePresent),
				NeededAction: ActionNone,
			},
		},
		{
			name: "install",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent),
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(nil, TargetStateAbsent),
				Expected:     mockComponentDiffState(compVersion3211, TargetStatePresent),
				NeededAction: ActionInstall,
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintComponent: mockTargetComponent(nil, TargetStateAbsent),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(compVersion3211, TargetStatePresent),
				Expected:     mockComponentDiffState(nil, TargetStateAbsent),
				NeededAction: ActionUninstall,
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3212, TargetStatePresent),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(compVersion3211, TargetStatePresent),
				Expected:     mockComponentDiffState(compVersion3212, TargetStatePresent),
				NeededAction: ActionUpgrade,
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent),
				installedComponent: mockComponentInstallation(compVersion3212),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(compVersion3212, TargetStatePresent),
				Expected:     mockComponentDiffState(compVersion3211, TargetStatePresent),
				NeededAction: ActionDowngrade,
			},
		},
		{
			name: "ignore present component, no action",
			args: args{
				blueprintComponent: nil,
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:         testComponentName,
				Actual:       mockComponentDiffState(compVersion3211, TargetStatePresent),
				Expected:     mockComponentDiffState(compVersion3211, TargetStatePresent),
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
			compDiff, err := determineComponentDiff(tt.args.blueprintComponent, tt.args.installedComponent)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, compDiff, "determineComponentDiff(%v, %v, %v)", tt.args.logger, tt.args.blueprintComponent, tt.args.installedComponent)
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
		Version:           compVersion3211,
		InstallationState: TargetStatePresent,
	}
	expected := ComponentDiffState{
		Version:           compVersion3212,
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
		"Actual: {DistributionNamespace: \"\", Version: \"3.2.1-1\", InstallationState: \"present\"}, "+
		"Expected: {DistributionNamespace: \"\", Version: \"3.2.1-2\", InstallationState: \"present\"}, "+
		"NeededAction: \"install\""+
		"}", diff.String())
}

func TestComponentDiffState_String(t *testing.T) {
	diff := &ComponentDiffState{
		DistributionNamespace: "k8s",
		Version:               compVersion3211,
		InstallationState:     TargetStatePresent,
	}

	assert.Equal(t, `{DistributionNamespace: "k8s", Version: "3.2.1-1", InstallationState: "present"}`, diff.String())
}

func mockTargetComponent(version *semver.Version, state TargetState) *Component {
	return &Component{
		Name:        testComponentName,
		Version:     version,
		TargetState: state,
	}
}

func mockComponentInstallation(version *semver.Version) *ecosystem.ComponentInstallation {
	return &ecosystem.ComponentInstallation{
		Name:    testComponentName,
		Version: version,
	}
}

func mockComponentDiffState(version *semver.Version, state TargetState) ComponentDiffState {
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
						Version:     compVersion3211,
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
						Version:           compVersion3211,
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
						Version: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName,
					Actual: ComponentDiffState{
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
					},
					Expected: ComponentDiffState{
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionNone,
				},
			},
		},
		{
			name: "determine upgrade for an installed component which is also in the blueprint",
			args: args{
				blueprintComponents: []Component{
					{
						Name:        testComponentName,
						Version:     compVersion3212,
						TargetState: TargetStatePresent,
					},
				},
				installedComponents: map[string]*ecosystem.ComponentInstallation{
					testComponentName: {
						Name:    testComponentName,
						Version: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName,
					Actual: ComponentDiffState{
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
					},
					Expected: ComponentDiffState{
						Version:           compVersion3212,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionUpgrade,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compDiffs, err := determineComponentDiffs(tt.args.blueprintComponents, tt.args.installedComponents)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, compDiffs, "determineComponentDiffs(%v, %v, %v)", tt.args.logger, tt.args.blueprintComponents, tt.args.installedComponents)
		})
	}
}
