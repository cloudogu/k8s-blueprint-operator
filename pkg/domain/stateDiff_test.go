package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_determineDoguDiff(t *testing.T) {
	type args struct {
		blueprintDogu *Dogu
		installedDogu *ecosystem.DoguInstallation
	}
	tests := []struct {
		name string
		args args
		want DoguDiff
	}{
		{
			name: "equal, no action",
			args: args{
				blueprintDogu: &Dogu{
					Namespace:   "official",
					Name:        "postgresql",
					Version:     version3_2_1_1,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Namespace: "official",
					Name:      "postgresql",
					Version:   version3_2_1_1,
				},
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				NeededAction: ActionNone,
			},
		},
		{
			name: "install",
			args: args{
				blueprintDogu: &Dogu{
					Namespace:   "official",
					Name:        "postgresql",
					Version:     version3_2_1_1,
					TargetState: TargetStatePresent,
				},
				installedDogu: nil,
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					InstallationState: TargetStateAbsent,
				},
				Expected: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				NeededAction: ActionInstall,
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintDogu: &Dogu{
					Name:        "postgresql",
					TargetState: TargetStateAbsent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Namespace: "official",
					Name:      "postgresql",
					Version:   version3_2_1_1,
				},
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					InstallationState: TargetStateAbsent,
				},
				NeededAction: ActionUninstall,
			},
		},
		{
			name: "namespace switch",
			args: args{
				blueprintDogu: &Dogu{
					Namespace:   "premium",
					Name:        "postgresql",
					Version:     version3_2_1_1,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Namespace: "official",
					Name:      "postgresql",
					Version:   version3_2_1_1,
				},
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "premium",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				NeededAction: ActionSwitchNamespace,
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintDogu: &Dogu{
					Namespace:   "official",
					Name:        "postgresql",
					Version:     version3_2_1_2,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Namespace: "official",
					Name:      "postgresql",
					Version:   version3_2_1_1,
				},
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_2,
					InstallationState: TargetStatePresent,
				},
				NeededAction: ActionUpgrade,
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintDogu: &Dogu{
					Namespace:   "official",
					Name:        "postgresql",
					Version:     version3_2_1_1,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Namespace: "official",
					Name:      "postgresql",
					Version:   version3_2_1_2,
				},
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_2,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				NeededAction: ActionDowngrade,
			},
		},
		{
			name: "ignore present dogu, no action",
			args: args{
				blueprintDogu: nil,
				installedDogu: &ecosystem.DoguInstallation{
					Namespace: "official",
					Name:      "postgresql",
					Version:   version3_2_1_1,
				},
			},
			want: DoguDiff{
				DoguName: "postgresql",
				Actual: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "official",
					Version:           version3_2_1_1,
					InstallationState: TargetStatePresent,
				},
				NeededAction: ActionNone,
			},
		},
		{
			name: "should stay absent, no action",
			args: args{
				blueprintDogu: nil,
				installedDogu: nil,
			},
			want: DoguDiff{
				DoguName: "",
				Actual: DoguDiffState{
					InstallationState: TargetStateAbsent,
				},
				Expected: DoguDiffState{
					InstallationState: TargetStateAbsent,
				},
				NeededAction: ActionNone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determineDoguDiff(tt.args.blueprintDogu, tt.args.installedDogu), "determineDoguDiff(%v, %v)", tt.args.blueprintDogu, tt.args.installedDogu)
		})
	}
}

func Test_determineDoguDiffs(t *testing.T) {
	type args struct {
		blueprintDogus []Dogu
		installedDogus map[string]*ecosystem.DoguInstallation
	}
	tests := []struct {
		name string
		args args
		want []DoguDiff
	}{
		{
			name: "no dogus",
			args: args{
				blueprintDogus: nil,
				installedDogus: nil,
			},
			want: []DoguDiff{},
		},
		{
			name: "a not installed dogu in the blueprint",
			args: args{
				blueprintDogus: []Dogu{
					{
						Namespace:   "official",
						Name:        "postgresql",
						Version:     version3_2_1_1,
						TargetState: TargetStatePresent,
					},
				},
				installedDogus: nil,
			},
			want: []DoguDiff{
				{
					DoguName: "postgresql",
					Actual: DoguDiffState{
						InstallationState: TargetStateAbsent,
					},
					Expected: DoguDiffState{
						Namespace:         "official",
						Version:           version3_2_1_1,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionInstall,
				},
			},
		},
		{
			name: "an installed dogu which is not in the blueprint",
			args: args{
				blueprintDogus: nil,
				installedDogus: map[string]*ecosystem.DoguInstallation{
					"postgresql": {
						Namespace: "official",
						Name:      "postgresql",
						Version:   version3_2_1_1,
					},
				},
			},
			want: []DoguDiff{
				{
					DoguName: "postgresql",
					Actual: DoguDiffState{
						Namespace:         "official",
						Version:           version3_2_1_1,
						InstallationState: TargetStatePresent,
					},
					Expected: DoguDiffState{
						Namespace:         "official",
						Version:           version3_2_1_1,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionNone,
				},
			},
		},
		{
			name: "an installed dogu which is also in the blueprint",
			args: args{
				blueprintDogus: []Dogu{
					{
						Namespace:   "official",
						Name:        "postgresql",
						Version:     version3_2_1_2,
						TargetState: TargetStatePresent,
					},
				},
				installedDogus: map[string]*ecosystem.DoguInstallation{
					"postgresql": {
						Namespace: "official",
						Name:      "postgresql",
						Version:   version3_2_1_1,
					},
				},
			},
			want: []DoguDiff{
				{
					DoguName: "postgresql",
					Actual: DoguDiffState{
						Namespace:         "official",
						Version:           version3_2_1_1,
						InstallationState: TargetStatePresent,
					},
					Expected: DoguDiffState{
						Namespace:         "official",
						Version:           version3_2_1_2,
						InstallationState: TargetStatePresent,
					},
					NeededAction: ActionUpgrade,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, determineDoguDiffs(tt.args.blueprintDogus, tt.args.installedDogus), "determineDoguDiffs(%v, %v)", tt.args.blueprintDogus, tt.args.installedDogus)
		})
	}
}

func TestDoguDiffs_Statistics(t *testing.T) {
	tests := []struct {
		name            string
		dd              DoguDiffs
		wantToInstall   int
		wantToUpgrade   int
		wantToUninstall int
		wantOther       int
	}{
		{
			name:            "0 overall",
			dd:              DoguDiffs{},
			wantToInstall:   0,
			wantToUpgrade:   0,
			wantToUninstall: 0,
			wantOther:       0,
		},
		{
			name: "4 to install, 3 to upgrade, 2 to uninstall, 3 other",
			dd: DoguDiffs{
				{NeededAction: ActionNone},
				{NeededAction: ActionInstall},
				{NeededAction: ActionUninstall},
				{NeededAction: ActionInstall},
				{NeededAction: ActionUpgrade},
				{NeededAction: ActionSwitchNamespace},
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
			wantOther:       3,
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

func TestDoguDiff_String(t *testing.T) {
	actual := DoguDiffState{
		Namespace:         "official",
		Version:           version3_2_1_1,
		InstallationState: TargetStatePresent,
	}
	expected := DoguDiffState{
		Namespace:         "premium",
		Version:           version3_2_1_2,
		InstallationState: TargetStatePresent,
	}
	diff := &DoguDiff{
		DoguName:     "postgresql",
		Actual:       actual,
		Expected:     expected,
		NeededAction: ActionInstall,
	}

	assert.Equal(t, "{"+
		"DoguName: \"postgresql\", "+
		"Actual: {Version: \"3.2.1-1\", Namespace: \"official\", InstallationState: \"present\"}, "+
		"Expected: {Version: \"3.2.1-2\", Namespace: \"premium\", InstallationState: \"present\"}, "+
		"NeededAction: \"install\""+
		"}", diff.String())
}
func TestDoguDiffState_String(t *testing.T) {
	diff := &DoguDiffState{
		Namespace:         "official",
		Version:           version3_2_1_1,
		InstallationState: TargetStatePresent,
	}

	assert.Equal(t, "{Version: \"3.2.1-1\", Namespace: \"official\", InstallationState: \"present\"}", diff.String())
}
