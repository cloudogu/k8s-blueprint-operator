package domain

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func Test_determineDoguDiff(t *testing.T) {
	proxyBodySize := resource.MustParse("1M")
	volumeSize1 := resource.MustParse("1Gi")
	volumeSize2 := resource.MustParse("2Gi")

	type args struct {
		blueprintDogu *Dogu
		installedDogu *ecosystem.DoguInstallation
	}
	quantity100M := resource.MustParse("100M")
	quantity10M := resource.MustParse("10M")
	quantity100MPtr := &quantity100M
	quantity10MPtr := &quantity10M
	tests := []struct {
		name string
		args args
		want DoguDiff
	}{
		{
			name: "equal, no action",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3211,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				NeededActions: nil,
			},
		},
		{
			name: "install",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3211,
					TargetState: TargetStatePresent,
				},
				installedDogu: nil,
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					InstallationState: TargetStateAbsent,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				NeededActions: []Action{ActionInstall},
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					TargetState: TargetStateAbsent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					InstallationState: TargetStateAbsent,
				},
				NeededActions: []Action{ActionUninstall},
			},
		},
		{
			name: "namespace switch",
			args: args{
				blueprintDogu: &Dogu{
					Name:        premiumNexus,
					Version:     version3211,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "premium",
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				NeededActions: []Action{ActionSwitchDoguNamespace},
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
				},
				NeededActions: []Action{ActionUpgrade},
			},
		},
		{
			name: "multiple update actions",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize:      &proxyBodySize,
						AdditionalConfig: "additional",
						RewriteTarget:    "/",
					},
					MinVolumeSize: &volumeSize2,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:          officialNexus,
					Version:       version3211,
					MinVolumeSize: &volumeSize1,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
					MinVolumeSize:     &volumeSize1,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize:      &proxyBodySize,
						AdditionalConfig: "additional",
						RewriteTarget:    "/",
					},
					MinVolumeSize: &volumeSize2,
				},
				NeededActions: []Action{ActionUpdateDoguResourceMinVolumeSize, ActionUpdateDoguProxyBodySize, ActionUpdateDoguProxyRewriteTarget, ActionUpdateDoguProxyAdditionalConfig, ActionUpgrade},
			},
		},
		{
			name: "downgrade",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3211,
					TargetState: TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				NeededActions: []Action{ActionDowngrade},
			},
		},
		{
			name: "ignore present dogu, no action",
			args: args{
				blueprintDogu: nil,
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3211,
					InstallationState: TargetStatePresent,
				},
				NeededActions: nil,
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
				NeededActions: []Action{},
			},
		},
		{
			name: "update proxy body size",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity100MPtr,
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity100MPtr,
					},
				},
				NeededActions: []Action{ActionUpdateDoguProxyBodySize},
			},
		},
		{
			name: "update if proxy body size changed",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity100MPtr,
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity10MPtr,
					},
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity10MPtr,
					},
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity100MPtr,
					},
				},
				NeededActions: []Action{ActionUpdateDoguProxyBodySize},
			},
		},
		{
			name: "no update if body sizes are nil",
			args: args{
				blueprintDogu: &Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				Expected: DoguDiffState{
					Namespace:         officialNamespace,
					Version:           version3212,
					InstallationState: TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				NeededActions: nil,
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
		installedDogus map[common.SimpleDoguName]*ecosystem.DoguInstallation
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
						Name:        officialNexus,
						Version:     version3211,
						TargetState: TargetStatePresent,
					},
				},
				installedDogus: nil,
			},
			want: []DoguDiff{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						InstallationState: TargetStateAbsent,
					},
					Expected: DoguDiffState{
						Namespace:         officialNamespace,
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					NeededActions: []Action{ActionInstall},
				},
			},
		},
		{
			name: "an installed dogu which is not in the blueprint",
			args: args{
				blueprintDogus: nil,
				installedDogus: map[common.SimpleDoguName]*ecosystem.DoguInstallation{
					"postgresql": {
						Name:    officialNexus,
						Version: version3211,
					},
				},
			},
			want: []DoguDiff{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						Namespace:         officialNamespace,
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					Expected: DoguDiffState{
						Namespace:         officialNamespace,
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					NeededActions: nil,
				},
			},
		},
		{
			name: "an installed dogu which is also in the blueprint",
			args: args{
				blueprintDogus: []Dogu{
					{
						Name:        officialNexus,
						Version:     version3212,
						TargetState: TargetStatePresent,
					},
				},
				installedDogus: map[common.SimpleDoguName]*ecosystem.DoguInstallation{
					"nexus": {
						Name:    officialNexus,
						Version: version3211,
					},
				},
			},
			want: []DoguDiff{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						Namespace:         officialNamespace,
						Version:           version3211,
						InstallationState: TargetStatePresent,
					},
					Expected: DoguDiffState{
						Namespace:         officialNamespace,
						Version:           version3212,
						InstallationState: TargetStatePresent,
					},
					NeededActions: []Action{ActionUpgrade},
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

func TestDoguDiff_String(t *testing.T) {
	actual := DoguDiffState{
		Namespace:         "official",
		Version:           version3211,
		InstallationState: TargetStatePresent,
	}
	expected := DoguDiffState{
		Namespace:         "premium",
		Version:           version3212,
		InstallationState: TargetStatePresent,
	}
	diff := &DoguDiff{
		DoguName:      "postgresql",
		Actual:        actual,
		Expected:      expected,
		NeededActions: []Action{ActionUpgrade, ActionUpdateDoguResourceMinVolumeSize},
	}

	assert.Equal(t, "{"+
		"DoguName: \"postgresql\", "+
		"Actual: {Version: \"3.2.1-1\", Namespace: \"official\", InstallationState: \"present\"}, "+
		"Expected: {Version: \"3.2.1-2\", Namespace: \"premium\", InstallationState: \"present\"}, "+
		"NeededActions: [\"upgrade\" \"update resource minimum volume size\"]"+
		"}", diff.String())
}
func TestDoguDiffState_String(t *testing.T) {
	diff := &DoguDiffState{
		Namespace:         "official",
		Version:           version3211,
		InstallationState: TargetStatePresent,
	}

	assert.Equal(t, "{Version: \"3.2.1-1\", Namespace: \"official\", InstallationState: \"present\"}", diff.String())
}
