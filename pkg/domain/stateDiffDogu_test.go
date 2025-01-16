package domain

import (
	"github.com/cloudogu/blueprint-lib/v2"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func Test_determineDoguDiff(t *testing.T) {
	proxyBodySize := resource.MustParse("1M")
	volumeSize1 := resource.MustParse("1Gi")
	volumeSize2 := resource.MustParse("2Gi")

	type args struct {
		blueprintDogu *v2.Dogu
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
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3211,
					TargetState: v2.TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				NeededActions: nil,
			},
		},
		{
			name: "install",
			args: args{
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3211,
					TargetState: v2.TargetStatePresent,
				},
				installedDogu: nil,
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					InstallationState: v2.TargetStateAbsent,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				NeededActions: []Action{ActionInstall},
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					TargetState: v2.TargetStateAbsent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					InstallationState: v2.TargetStateAbsent,
				},
				NeededActions: []Action{ActionUninstall},
			},
		},
		{
			name: "namespace switch",
			args: args{
				blueprintDogu: &v2.Dogu{
					Name:        premiumNexus,
					Version:     version3211,
					TargetState: v2.TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         "premium",
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				NeededActions: []Action{ActionSwitchDoguNamespace},
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: v2.TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
				},
				NeededActions: []Action{ActionUpgrade},
			},
		},
		{
			name: "multiple update actions",
			args: args{
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: v2.TargetStatePresent,
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
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
					MinVolumeSize:     &volumeSize1,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
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
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3211,
					TargetState: v2.TargetStatePresent,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
				},
			},
			want: DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
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
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3211,
					InstallationState: v2.TargetStatePresent,
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
					InstallationState: v2.TargetStateAbsent,
				},
				Expected: DoguDiffState{
					InstallationState: v2.TargetStateAbsent,
				},
				NeededActions: []Action{},
			},
		},
		{
			name: "update proxy body size",
			args: args{
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: v2.TargetStatePresent,
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
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
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
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: v2.TargetStatePresent,
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
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity10MPtr,
					},
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
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
				blueprintDogu: &v2.Dogu{
					Name:        officialNexus,
					Version:     version3212,
					TargetState: v2.TargetStatePresent,
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
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				Expected: DoguDiffState{
					Namespace:         v2.officialNamespace,
					Version:           version3212,
					InstallationState: v2.TargetStatePresent,
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
		blueprintDogus []v2.Dogu
		installedDogus map[cescommons.SimpleName]*ecosystem.DoguInstallation
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
				blueprintDogus: []v2.Dogu{
					{
						Name:        officialNexus,
						Version:     version3211,
						TargetState: v2.TargetStatePresent,
					},
				},
				installedDogus: nil,
			},
			want: []DoguDiff{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						InstallationState: v2.TargetStateAbsent,
					},
					Expected: DoguDiffState{
						Namespace:         v2.officialNamespace,
						Version:           version3211,
						InstallationState: v2.TargetStatePresent,
					},
					NeededActions: []Action{ActionInstall},
				},
			},
		},
		{
			name: "an installed dogu which is not in the blueprint",
			args: args{
				blueprintDogus: nil,
				installedDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{
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
						Namespace:         v2.officialNamespace,
						Version:           version3211,
						InstallationState: v2.TargetStatePresent,
					},
					Expected: DoguDiffState{
						Namespace:         v2.officialNamespace,
						Version:           version3211,
						InstallationState: v2.TargetStatePresent,
					},
					NeededActions: nil,
				},
			},
		},
		{
			name: "an installed dogu which is also in the blueprint",
			args: args{
				blueprintDogus: []v2.Dogu{
					{
						Name:        officialNexus,
						Version:     version3212,
						TargetState: v2.TargetStatePresent,
					},
				},
				installedDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{
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
						Namespace:         v2.officialNamespace,
						Version:           version3211,
						InstallationState: v2.TargetStatePresent,
					},
					Expected: DoguDiffState{
						Namespace:         v2.officialNamespace,
						Version:           version3212,
						InstallationState: v2.TargetStatePresent,
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
		InstallationState: v2.TargetStatePresent,
	}
	expected := DoguDiffState{
		Namespace:         "premium",
		Version:           version3212,
		InstallationState: v2.TargetStatePresent,
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
		InstallationState: v2.TargetStatePresent,
	}

	assert.Equal(t, "{Version: \"3.2.1-1\", Namespace: \"official\", InstallationState: \"present\"}", diff.String())
}
