package domain

import (
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/Masterminds/semver/v3"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

var (
	testComponentName = common.QualifiedComponentName{
		Namespace:  "k8s",
		SimpleName: "my-component",
	}
	blueprintOperatorSimpleName = common.SimpleComponentName("k8s-blueprint-operator")
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
				blueprintComponent: mockTargetComponent(compVersion3211, false, nil),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				NeededActions: nil,
			},
		},
		{
			name: "install",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, false, nil),
				installedComponent: nil,
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState("", nil, true, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				NeededActions: []Action{ActionInstall},
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintComponent: mockTargetComponent(nil, true, nil),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, nil, true, nil),
				NeededActions: []Action{ActionUninstall},
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3212, false, nil),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3212, false, nil),
				NeededActions: []Action{ActionUpgrade},
			},
		},
		{
			name: "update package config",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, false, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, false, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				NeededActions: []Action{ActionUpdateComponentDeployConfig},
			},
		},
		{
			name: "update package config and upgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3212, false, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				installedComponent: mockComponentInstallation(compVersion3211),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3212, false, map[string]interface{}{"deployNamespace": "k8s-longhorn"}),
				NeededActions: []Action{ActionUpdateComponentDeployConfig, ActionUpgrade},
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintComponent: mockTargetComponent(compVersion3211, false, nil),
				installedComponent: mockComponentInstallation(compVersion3212),
			},
			want: ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3212, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
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
				Actual:        mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				Expected:      mockComponentDiffState(testDistributionNamespace, compVersion3211, false, nil),
				NeededActions: nil,
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
				Actual:        ComponentDiffState{Absent: true},
				Expected:      ComponentDiffState{Absent: true},
				NeededActions: nil,
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

func TestComponentDiff_String(t *testing.T) {
	actual := ComponentDiffState{
		Version: compVersion3211,
		Absent:  false,
	}
	expected := ComponentDiffState{
		Version: compVersion3212,
		Absent:  false,
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
		Namespace: "k8s",
		Version:   compVersion3211,
		Absent:    false,
	}

	assert.Equal(t, `{Namespace: "k8s", Version: "3.2.1-1", InstallationState: "present"}`, diff.String())
}

func mockTargetComponent(version *semver.Version, absent bool, deployConfig ecosystem.DeployConfig) *Component {
	return &Component{
		Name:         testComponentName,
		Version:      version,
		Absent:       absent,
		DeployConfig: deployConfig,
	}
}

func mockComponentInstallation(version *semver.Version) *ecosystem.ComponentInstallation {
	return &ecosystem.ComponentInstallation{
		Name:            testComponentName,
		ExpectedVersion: version,
	}
}

func mockComponentDiffState(namespace common.ComponentNamespace, version *semver.Version, absent bool, deployConfig ecosystem.DeployConfig) ComponentDiffState {
	return ComponentDiffState{
		Namespace:    namespace,
		Version:      version,
		Absent:       absent,
		DeployConfig: deployConfig,
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
						Name:    testComponentName,
						Version: compVersion3211,
						Absent:  true,
					},
				},
				installedComponents: nil,
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Absent: true,
					},
					Expected: ComponentDiffState{
						Namespace: testComponentName.Namespace,
						Version:   compVersion3211,
						Absent:    false,
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
						Name:            testComponentName,
						ExpectedVersion: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Namespace: testComponentName.Namespace,
						Version:   compVersion3211,
						Absent:    false,
					},
					Expected: ComponentDiffState{
						Namespace: testComponentName.Namespace,
						Version:   compVersion3211,
						Absent:    false,
					},
					NeededActions: nil,
				},
			},
		},
		{
			name: "determine distribution namespace switch",
			args: args{
				blueprintComponents: []Component{
					{
						Name:    common.QualifiedComponentName{Namespace: "k8s-testing", SimpleName: "my-component"},
						Version: compVersion3211,
						Absent:  true,
					},
				},
				installedComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
					testComponentName.SimpleName: {
						Name:            common.QualifiedComponentName{Namespace: "k8s", SimpleName: "my-component"},
						ExpectedVersion: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Version:   compVersion3211,
						Absent:    false,
						Namespace: testDistributionNamespace,
					},
					Expected: ComponentDiffState{
						Version:   compVersion3211,
						Absent:    false,
						Namespace: testChangeDistributionNamespace,
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
						Name:    testComponentName,
						Version: compVersion3212,
						Absent:  true,
					},
				},
				installedComponents: map[common.SimpleComponentName]*ecosystem.ComponentInstallation{
					testComponentName.SimpleName: {
						Name:            testComponentName,
						ExpectedVersion: compVersion3211,
					},
				},
			},
			want: []ComponentDiff{
				{
					Name: testComponentName.SimpleName,
					Actual: ComponentDiffState{
						Version:   compVersion3211,
						Namespace: testComponentName.Namespace,
						Absent:    false,
					},
					Expected: ComponentDiffState{
						Version:   compVersion3212,
						Namespace: testComponentName.Namespace,
						Absent:    false,
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
		Namespace common.ComponentNamespace
		Version   *semver.Version
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
				Namespace: tt.fields.Namespace,
				Version:   tt.fields.Version,
			}
			assert.Equalf(t, tt.want, diff.getSafeVersionString(), "getSafeVersionString()")
		})
	}
}

func TestComponentDiffs_GetComponentDiffByName(t *testing.T) {
	t.Run("find diff", func(t *testing.T) {
		blueprintOpDiff := ComponentDiff{
			Name:          blueprintOperatorSimpleName,
			Actual:        ComponentDiffState{},
			Expected:      ComponentDiffState{},
			NeededActions: []Action{ActionUninstall},
		}
		diffs := ComponentDiffs{
			blueprintOpDiff,
			ComponentDiff{
				Name:          testComponentName.SimpleName,
				Actual:        ComponentDiffState{},
				Expected:      ComponentDiffState{},
				NeededActions: []Action{ActionUninstall},
			},
		}

		foundDiff := diffs.GetComponentDiffByName(blueprintOperatorSimpleName)

		assert.Equal(t, blueprintOpDiff, foundDiff)
	})

	t.Run("don't find diff", func(t *testing.T) {
		diffs := ComponentDiffs{}

		foundDiff := diffs.GetComponentDiffByName(blueprintOperatorSimpleName)

		assert.Equal(t, ComponentDiff{}, foundDiff)
	})
}

func TestComponentDiff_HasChanges(t *testing.T) {
	t.Run("no change", func(t *testing.T) {
		diff := ComponentDiff{
			NeededActions: []Action{},
		}
		assert.False(t, diff.HasChanges())
	})
	t.Run("change for any", func(t *testing.T) {
		diff := ComponentDiff{
			NeededActions: []Action{ActionInstall},
		}
		assert.True(t, diff.HasChanges())
	})
}

func TestComponentDiff_IsExpectedVersion(t *testing.T) {
	tests := []struct {
		name     string
		expected *semver.Version
		actual   *semver.Version
		want     bool
	}{
		{
			name:     "equal",
			expected: semver.MustParse("1.0"),
			actual:   semver.MustParse("1.0"),
			want:     true,
		},
		{
			name:     "equal dev versions",
			expected: semver.MustParse("0.2.0-dev"),
			actual:   semver.MustParse("0.2.0-dev"),
			want:     true,
		},
		{
			name:     "higher expected",
			expected: semver.MustParse("1.1"),
			actual:   semver.MustParse("1.0"),
			want:     false,
		},
		{
			name:     "nothing expected",
			expected: nil,
			actual:   semver.MustParse("1.0"),
			want:     true,
		},
		{
			name:     "nothing installed",
			expected: semver.MustParse("1.0"),
			actual:   nil,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := ComponentDiff{Expected: ComponentDiffState{Version: tt.expected}}
			assert.Equalf(t, tt.want, diff.IsExpectedVersion(tt.actual), "{version:%v}.IsExpectedVersion(%v)", tt.expected, tt.actual)
		})
	}
}
