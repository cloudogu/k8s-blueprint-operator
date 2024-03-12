package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/Masterminds/semver/v3"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
)

var (
	testComponentName = common.QualifiedComponentName{
		Namespace:  "k8s",
		SimpleName: "my-component",
	}
)

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
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent, nil),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				NeededActions: []Action{ActionNone},
			},
		},
		{
			name: "install",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent, nil),
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState("", nil, TargetStateAbsent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				NeededActions: []Action{ActionInstall},
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintComponent: mockTargetComponent(nil, TargetStateAbsent, nil),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, nil, TargetStateAbsent, nil),
				NeededActions: []Action{ActionUninstall},
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3212, TargetStatePresent, nil),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3212, TargetStatePresent, nil),
				NeededActions: []Action{ActionUpgrade},
			},
		},
		{
			name: "update package config",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				NeededActions: []Action{ActionUpdateComponentPackageConfig},
			},
		},
		{
			name: "update package config and upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3212, TargetStatePresent, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3212, TargetStatePresent, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				NeededActions: []Action{ActionUpdateComponentPackageConfig, ActionUpgrade},
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, TargetStatePresent, nil),
				installedComponent: mockComponentInstallation(compVersion3212),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3212, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				NeededActions: []Action{ActionDowngrade},
			},
		},
		{
			name: "ignore present component, no action",
			args: args{
				blueprintComponent: nil,
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, TargetStatePresent, nil),
				NeededActions: []Action{ActionNone},
			},
		},
		{
			name: "should stay absent, no action", // this is empty set comparison is weird and should basically not occur
			args: args{
				blueprintComponent: nil,
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:          "",
				Actual:        ComponentDiffState{InstallationState: TargetStateAbsent},
				Expected:      ComponentDiffState{InstallationState: TargetStateAbsent},
				NeededActions: []Action{ActionNone},
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
		name                      string
		dd                        ComponentDiffs
		wantToInstall             int
		wantToUpgrade             int
		wantToUninstall           int
		wantToUpdateNamespace     int
		wantToUpdatePackageConfig int
		wantOther                 int
	}{
		{
			name:                      "0 overall",
			dd:                        ComponentDiffs{},
			wantToInstall:             0,
			wantToUpgrade:             0,
			wantToUninstall:           0,
			wantToUpdateNamespace:     0,
			wantToUpdatePackageConfig: 0,
			wantOther:                 0,
		},
		{
			name: "4 to install, 3 to upgrade, 2 to uninstall, 2 to update namespace, 3 to update package config, 3 other",
			dd: ComponentDiffs{
				{NeededActions: []Action{ActionNone}},
				{NeededActions: []Action{ActionInstall}},
				{NeededActions: []Action{ActionUninstall}},
				{NeededActions: []Action{ActionInstall}},
				{NeededActions: []Action{ActionUpgrade, ActionSwitchComponentNamespace, ActionUpdateComponentPackageConfig}},
				{NeededActions: []Action{ActionInstall}},
				{NeededActions: []Action{ActionDowngrade}},
				{NeededActions: []Action{ActionUninstall}},
				{NeededActions: []Action{ActionInstall}},
				{NeededActions: []Action{ActionUpgrade, ActionSwitchComponentNamespace, ActionUpdateComponentPackageConfig}},
				{NeededActions: []Action{ActionUpgrade, ActionUpdateComponentPackageConfig}},
			},
			wantToInstall:             4,
			wantToUpgrade:             3,
			wantToUninstall:           2,
			wantToUpdateNamespace:     2,
			wantToUpdatePackageConfig: 3,
			wantOther:                 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToInstall, gotToUpgrade, gotToUninstall, gotToUpdateNamespace, gotToUpdatePackageConfig, gotOther := tt.dd.Statistics()
			assert.Equalf(t, tt.wantToInstall, gotToInstall, "Statistics()")
			assert.Equalf(t, tt.wantToUpgrade, gotToUpgrade, "Statistics()")
			assert.Equalf(t, tt.wantToUninstall, gotToUninstall, "Statistics()")
			assert.Equalf(t, tt.wantToUpdateNamespace, gotToUpdateNamespace, "Statistics()")
			assert.Equalf(t, tt.wantToUpdatePackageConfig, gotToUpdatePackageConfig, "Statistics()")
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
		Name:          testComponentName.SimpleName,
		Actual:        actual,
		Expected:      expected,
		NeededActions: []Action{ActionInstall},
	}

	assert.Equal(t, "{"+
		"Name: \"my-component\", "+
		"Actual: {Namespace: \"\", Version: \"3.2.1-1\", InstallationState: \"present\"}, "+
		"Expected: {Namespace: \"\", Version: \"3.2.1-2\", InstallationState: \"present\"}, "+
		"NeededActions: [\"install\"]"+
		"}", diff.String())
}

func TestComponentDiffState_String(t *testing.T) {
	diff := &ComponentDiffState{
		Namespace:         "k8s",
		Version:           compVersion3211,
		InstallationState: TargetStatePresent,
	}

	assert.Equal(t, `{Namespace: "k8s", Version: "3.2.1-1", InstallationState: "present"}`, diff.String())
}

func mockTargetComponent(version *semver.Version, state TargetState, packageConfig ecosystem.PackageConfig) *Component {
	return &Component{
		Name:          testComponentName,
		Version:       version,
		TargetState:   state,
		PackageConfig: packageConfig,
	}
}

func mockComponentInstallation(version *semver.Version) *ecosystem.ComponentInstallation {
	return &ecosystem.ComponentInstallation{
		Name:    testComponentName,
		Version: version,
	}
}

func mockComponentDiffState(namespace common.ComponentNamespace, version *semver.Version, state TargetState, packageConfig ecosystem.PackageConfig) ComponentDiffState {
	return ComponentDiffState{
		Namespace:         namespace,
		Version:           version,
		InstallationState: state,
		PackageConfig:     packageConfig,
	}
}

func Test_determineComponentDiffs(t *testing.T) {
	type args struct {
		logger              logr.Logger
		blueprintComponents []Component
		installedComponents map[common.SimpleComponentName]*ecosystem.ComponentInstallation
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
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						InstallationState: TargetStateAbsent,
					},
					Expected: ComponentDiffState{
						Namespace:         testComponentName.Namespace,
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
					},
					NeededActions: []Action{ActionInstall},
				},
			},
		},
		{
			name: "an installed component which is not in the blueprint",
			args: args{
				blueprintComponents: nil,
				installedComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
					testComponentName.SimpleName: {
						Name:    testComponentName,
						Version: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Namespace:         testComponentName.Namespace,
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
					},
					Expected: ComponentDiffState{
						Namespace:         testComponentName.Namespace,
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
					},
					NeededActions: []Action{ActionNone},
				},
			},
		},
		{
			name: "determine distribution namespace switch",
			args: args{
				blueprintComponents: []Component{
					{
						Name:        common.QualifiedComponentName{Namespace: "k8s-testing", SimpleName: "my-component"},
						Version:     compVersion3211,
						TargetState: TargetStatePresent,
					},
				},
				installedComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
					testComponentName.SimpleName: {
						Name:    common.QualifiedComponentName{Namespace: "k8s", SimpleName: "my-component"},
						Version: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
						Namespace:         testDistributionNamespace,
					},
					Expected: ComponentDiffState{
						Version:           compVersion3211,
						InstallationState: TargetStatePresent,
						Namespace:         testChangeDistributionNamespace,
					},
					NeededActions: []Action{ActionSwitchComponentNamespace},
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
				installedComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
					testComponentName.SimpleName: {
						Name:    testComponentName,
						Version: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Version:           compVersion3211,
						Namespace:         testComponentName.Namespace,
						InstallationState: TargetStatePresent,
					},
					Expected: ComponentDiffState{
						Version:           compVersion3212,
						Namespace:         testComponentName.Namespace,
						InstallationState: TargetStatePresent,
					},
					NeededActions: []Action{ActionUpgrade},
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

func TestComponentDiffState_getSafeVersionString(t *testing.T) {
	version1, _ := semver.NewVersion("1.0.0")

	type fields struct {
		Namespace         common.ComponentNamespace
		Version           *semver.Version
		InstallationState TargetState
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success",
			fields: fields{Version: version1},
			want:   "1.0.0",
		},
		{
			name:   "should return empty string and no panic on nil version",
			fields: fields{Version: nil},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := &ComponentDiffState{
				Namespace:         tt.fields.Namespace,
				Version:           tt.fields.Version,
				InstallationState: tt.fields.InstallationState,
			}
			assert.Equalf(t, tt.want, diff.getSafeVersionString(), "getSafeVersionString()")
		})
	}
}
