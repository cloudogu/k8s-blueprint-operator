package domain

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	additionalConfig = "additional"
	rewriteConfig    = "/"
	subfolder2       = "secsubfolder"
	subfolder3       = "different_subfolder"
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
		want *DoguDiff
	}{
		{
			name: "equal, no action",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3211,
					Absent:  false,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3211,
				},
			},
			want: nil,
		},
		{
			name: "install",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3211,
					Absent:  false,
				},
				installedDogu: nil,
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Absent: true,
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3211,
					Absent:    false,
				},
				NeededActions: []Action{ActionInstall},
			},
		},
		{
			name: "uninstall",
			args: args{
				blueprintDogu: &Dogu{
					Name:   officialNexus,
					Absent: true,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3211,
					InstalledVersion: version3211,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3211,
					InstalledVersion: &version3211,
					Absent:           false,
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Absent:    true,
				},
				NeededActions: []Action{ActionUninstall},
			},
		},
		{
			name: "namespace switch",
			args: args{
				blueprintDogu: &Dogu{
					Name:    premiumNexus,
					Version: &version3211,
					Absent:  false,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3211,
					InstalledVersion: version3211,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3211,
					InstalledVersion: &version3211,
					Absent:           false,
				},
				Expected: DoguDiffState{
					Namespace: "premium",
					Version:   &version3211,
					Absent:    false,
				},
				NeededActions: []Action{ActionSwitchDoguNamespace},
			},
		},
		{
			name: "upgrade",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					Absent:  false,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3211,
					InstalledVersion: version3211,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3211,
					InstalledVersion: &version3211,
					Absent:           false,
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3212,
					Absent:    false,
				},
				NeededActions: []Action{ActionUpgrade},
			},
		},
		{
			name: "reset version",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3211,
					Absent:  false,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3211,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3211,
					Absent:           false,
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3211,
					Absent:    false,
				},
				NeededActions: []Action{ActionResetVersion},
			},
		},
		{
			name: "update minVolSize if actual < expected",
			args: args{
				blueprintDogu: &Dogu{
					Name:          officialNexus,
					Version:       &version3212,
					MinVolumeSize: &quantity100M,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3212,
					MinVolumeSize:    &quantity10M,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3212,
					MinVolumeSize:    &quantity10M,
				},
				Expected: DoguDiffState{
					Namespace:     officialNamespace,
					Version:       &version3212,
					MinVolumeSize: &quantity100M,
				},
				NeededActions: []Action{ActionUpdateDoguResourceMinVolumeSize},
			},
		},
		{
			name: "don't update minVolSize if actual == expected",
			args: args{
				blueprintDogu: &Dogu{
					Name:          officialNexus,
					Version:       &version3212,
					MinVolumeSize: &quantity100M,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:          officialNexus,
					Version:       version3212,
					MinVolumeSize: &quantity100M,
				},
			},
			want: nil,
		},
		{
			name: "don't update minVolSize if actual > expected",
			args: args{
				blueprintDogu: &Dogu{
					Name:          officialNexus,
					Version:       &version3212,
					MinVolumeSize: &quantity10M,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:          officialNexus,
					Version:       version3212,
					MinVolumeSize: &quantity100M,
				},
			},
			want: nil,
		},
		{
			name: "multiple update actions",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize:      &proxyBodySize,
						AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig),
						RewriteTarget:    ecosystem.RewriteTarget(rewriteConfig),
					},
					MinVolumeSize: &volumeSize2,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3211,
					InstalledVersion: version3211,
					MinVolumeSize:    &volumeSize1,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3211,
					InstalledVersion: &version3211,
					MinVolumeSize:    &volumeSize1,
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize:      &proxyBodySize,
						AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig),
						RewriteTarget:    ecosystem.RewriteTarget(rewriteConfig),
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
					Name:    officialNexus,
					Version: &version3211,
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3212,
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3212,
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3211,
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
			want: nil,
		},
		{
			name: "should stay absent, no action",
			args: args{
				blueprintDogu: nil,
				installedDogu: nil,
			},
			want: nil,
		},
		{
			name: "update proxy body size",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity100MPtr,
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: nil,
					},
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3212,
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
					Name:    officialNexus,
					Version: &version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity100MPtr,
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity10MPtr,
					},
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3212,
					ReverseProxyConfig: ecosystem.ReverseProxyConfig{
						MaxBodySize: quantity10MPtr,
					},
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3212,
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
					Name:    officialNexus,
					Version: &version3212,
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
			want: nil,
		},
		{
			name: "no action if additional mounts are equal",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "no action if additional mounts are equal but order is different",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:    officialNexus,
					Version: version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "needs update action for additional mounts if the size is different",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
					},
				},
				NeededActions: []Action{ActionUpdateAdditionalMounts},
			},
		},
		{
			name: "needs update action for additional mounts if an element is different",
			args: args{
				blueprintDogu: &Dogu{
					Name:    officialNexus,
					Version: &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder3,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
				installedDogu: &ecosystem.DoguInstallation{
					Name:             officialNexus,
					Version:          version3212,
					InstalledVersion: version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
			},
			want: &DoguDiff{
				DoguName: "nexus",
				Actual: DoguDiffState{
					Namespace:        officialNamespace,
					Version:          &version3212,
					InstalledVersion: &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
				Expected: DoguDiffState{
					Namespace: officialNamespace,
					Version:   &version3212,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  subfolder3,
						},
						{
							SourceType: ecosystem.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  subfolder2,
						},
					},
				},
				NeededActions: []Action{ActionUpdateAdditionalMounts},
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
				blueprintDogus: []Dogu{
					{
						Name:    officialNexus,
						Version: &version3211,
						Absent:  false,
					},
				},
				installedDogus: nil,
			},
			want: []DoguDiff{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						Absent: true,
					},
					Expected: DoguDiffState{
						Namespace: officialNamespace,
						Version:   &version3211,
						Absent:    false,
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
			want: []DoguDiff{},
		},
		{
			name: "an installed dogu which is also in the blueprint",
			args: args{
				blueprintDogus: []Dogu{
					{
						Name:    officialNexus,
						Version: &version3212,
						Absent:  false,
					},
				},
				installedDogus: map[cescommons.SimpleName]*ecosystem.DoguInstallation{
					"nexus": {
						Name:             officialNexus,
						Version:          version3211,
						InstalledVersion: version3211,
					},
				},
			},
			want: []DoguDiff{
				{
					DoguName: "nexus",
					Actual: DoguDiffState{
						Namespace:        officialNamespace,
						Version:          &version3211,
						InstalledVersion: &version3211,
						Absent:           false,
					},
					Expected: DoguDiffState{
						Namespace: officialNamespace,
						Version:   &version3212,
						Absent:    false,
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
		Namespace: "official",
		Version:   &version3211,
		Absent:    false,
	}
	expected := DoguDiffState{
		Namespace: "premium",
		Version:   &version3212,
		Absent:    false,
	}
	diff := &DoguDiff{
		DoguName:      "postgresql",
		Actual:        actual,
		Expected:      expected,
		NeededActions: []Action{ActionUpgrade, ActionUpdateDoguResourceMinVolumeSize},
	}

	assert.Equal(t, "{"+
		"DoguName: \"postgresql\", "+
		"Actual: {Version: \"3.2.1-1\", Namespace: \"official\", Absent: false}, "+
		"Expected: {Version: \"3.2.1-2\", Namespace: \"premium\", Absent: false}, "+
		"NeededActions: [\"upgrade\" \"update resource minimum volume size\"]"+
		"}", diff.String())
}
func TestDoguDiffState_String(t *testing.T) {
	diff := &DoguDiffState{
		Namespace: "official",
		Version:   &version3211,
		Absent:    false,
	}

	assert.Equal(t, "{Version: \"3.2.1-1\", Namespace: \"official\", Absent: false}", diff.String())
}
