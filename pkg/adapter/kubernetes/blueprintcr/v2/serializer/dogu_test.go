package serializer

import (
	"fmt"
	bpv2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

func TestConvertDogus(t *testing.T) {
	type args struct {
		dogus []bpv2.Dogu
	}
	proxyBodySize := resource.MustParse("1G")
	volumeSize := resource.MustParse("1Gi")
	tests := []struct {
		name    string
		args    args
		want    []domain.Dogu
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{dogus: nil},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "empty list",
			args:    args{dogus: []bpv2.Dogu{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal dogu",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent}},
			wantErr: assert.NoError,
		},
		{
			name:    "absent dogu",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Absent: true}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, TargetState: domain.TargetStateAbsent}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with max proxy body size",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false, PlatformConfig: bpv2.PlatformConfig{ReverseProxyConfig: bpv2.ReverseProxyConfig{MaxBodySize: "1G"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &proxyBodySize}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with proxy rewrite target",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false, PlatformConfig: bpv2.PlatformConfig{ReverseProxyConfig: bpv2.ReverseProxyConfig{RewriteTarget: "/"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{RewriteTarget: "/"}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with proxy additional config",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false, PlatformConfig: bpv2.PlatformConfig{ReverseProxyConfig: bpv2.ReverseProxyConfig{AdditionalConfig: "additional"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{AdditionalConfig: "additional"}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with invalid proxy body size",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false, PlatformConfig: bpv2.PlatformConfig{ReverseProxyConfig: bpv2.ReverseProxyConfig{MaxBodySize: "1GE"}}}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "dogu with min volume size",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, PlatformConfig: bpv2.PlatformConfig{ResourceConfig: bpv2.ResourceConfig{MinVolumeSize: "1Gi"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, MinVolumeSize: volumeSize}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with invalid volume size",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, PlatformConfig: bpv2.PlatformConfig{ResourceConfig: bpv2.ResourceConfig{MinVolumeSize: "1GIE"}}}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "no namespace",
			args:    args{dogus: []bpv2.Dogu{{Name: "postgres", Version: version3211.Raw}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unparsable version",
			args:    args{dogus: []bpv2.Dogu{{Name: "official/postgres", Version: "1."}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "should convert additionalMounts",
			args: args{dogus: []bpv2.Dogu{{
				Name:    "official/postgres",
				Version: version3211.Raw,
				PlatformConfig: bpv2.PlatformConfig{
					AdditionalMountsConfig: []bpv2.AdditionalMount{
						{
							SourceType: bpv2.DataSourceConfigMap,
							Name:       "configMap",
							Volume:     "volume",
							Subfolder:  "subfolder",
						},
						{
							SourceType: bpv2.DataSourceSecret,
							Name:       "sec",
							Volume:     "secvolume",
							Subfolder:  "secsubfolder",
						},
					},
				},
			}}},
			want: []domain.Dogu{{
				Name:        cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:     version3211,
				TargetState: domain.TargetStatePresent,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configMap",
						Volume:     "volume",
						Subfolder:  "subfolder",
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "sec",
						Volume:     "secvolume",
						Subfolder:  "secsubfolder",
					},
				}}},
			wantErr: assert.NoError,
		},
		{
			name: "should return nil slice if dogu contains an nil slice",
			args: args{dogus: []bpv2.Dogu{{
				Name:    "official/postgres",
				Version: version3211.Raw,
				PlatformConfig: bpv2.PlatformConfig{
					AdditionalMountsConfig: nil,
				},
			}}},
			want: []domain.Dogu{{
				Name:             cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:          version3211,
				TargetState:      domain.TargetStatePresent,
				AdditionalMounts: nil,
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDogus(tt.args.dogus)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertDogus(%v)", tt.args.dogus)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertDogus(%v)", tt.args.dogus)
		})
	}
}

func TestConvertMaskDogus(t *testing.T) {
	type args struct {
		dogus []bpv2.MaskDogu
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.MaskDogu
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{dogus: nil},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "empty list",
			args:    args{dogus: []bpv2.MaskDogu{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal dogu",
			args:    args{dogus: []bpv2.MaskDogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false}}},
			want:    []domain.MaskDogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent}},
			wantErr: assert.NoError,
		},
		{
			name:    "absent dogu",
			args:    args{dogus: []bpv2.MaskDogu{{Name: "official/postgres", Absent: true}}},
			want:    []domain.MaskDogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, TargetState: domain.TargetStateAbsent}},
			wantErr: assert.NoError,
		},
		{
			name:    "no namespace",
			args:    args{dogus: []bpv2.MaskDogu{{Name: "postgres", Version: version3211.Raw}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unparsable version",
			args:    args{dogus: []bpv2.MaskDogu{{Name: "official/postgres", Version: "1."}}},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertMaskDogus(tt.args.dogus)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertMaskDogus(%v)", tt.args.dogus)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertMaskDogus(%v)", tt.args.dogus)
		})
	}
}

func TestConvertToDoguDTOs(t *testing.T) {
	type args struct {
		dogus []domain.Dogu
	}
	bodySize := resource.MustParse("100M")
	volumeSize := resource.MustParse("1G")
	tests := []struct {
		name string
		args args
		want []bpv2.Dogu
	}{
		{
			name: "nil",
			args: args{},
			want: []bpv2.Dogu{},
		},
		{
			name: "empty list",
			args: args{dogus: []domain.Dogu{}},
			want: []bpv2.Dogu{},
		},
		{
			name: "ok",
			args: args{dogus: []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, MinVolumeSize: volumeSize, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &bodySize, RewriteTarget: "/", AdditionalConfig: "additional"}}}},
			want: []bpv2.Dogu{{Name: "official/postgres", Version: version3211.Raw, Absent: false, PlatformConfig: bpv2.PlatformConfig{ResourceConfig: bpv2.ResourceConfig{MinVolumeSize: "1G"}, ReverseProxyConfig: bpv2.ReverseProxyConfig{MaxBodySize: "100M", RewriteTarget: "/", AdditionalConfig: "additional"}}}},
		},
		{
			name: "additionalMountsConfig",
			args: args{dogus: []domain.Dogu{{
				Name:               cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:            version3211,
				TargetState:        domain.TargetStatePresent,
				MinVolumeSize:      volumeSize,
				ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &bodySize, RewriteTarget: "/", AdditionalConfig: "additional"},
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configMap",
						Volume:     "volume",
						Subfolder:  "subfolder",
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "sec",
						Volume:     "secvolume",
						Subfolder:  "secsubfolder",
					},
				},
			}}},
			want: []bpv2.Dogu{{
				Name:    "official/postgres",
				Version: version3211.Raw,
				Absent:  false,
				PlatformConfig: bpv2.PlatformConfig{
					ResourceConfig:     bpv2.ResourceConfig{MinVolumeSize: "1G"},
					ReverseProxyConfig: bpv2.ReverseProxyConfig{MaxBodySize: "100M", RewriteTarget: "/", AdditionalConfig: "additional"},
					AdditionalMountsConfig: []bpv2.AdditionalMount{
						{
							SourceType: bpv2.DataSourceConfigMap,
							Name:       "configMap",
							Volume:     "volume",
							Subfolder:  "subfolder",
						},
						{
							SourceType: bpv2.DataSourceSecret,
							Name:       "sec",
							Volume:     "secvolume",
							Subfolder:  "secsubfolder",
						},
					},
				}}},
		},
		{
			name: "should return nil slice if dogu contains an nil slice",
			args: args{dogus: []domain.Dogu{{
				Name:               cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:            version3211,
				TargetState:        domain.TargetStatePresent,
				MinVolumeSize:      volumeSize,
				ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &bodySize, RewriteTarget: "/", AdditionalConfig: "additional"},
				AdditionalMounts:   nil,
			}}},
			want: []bpv2.Dogu{{
				Name:    "official/postgres",
				Version: version3211.Raw,
				Absent:  false,
				PlatformConfig: bpv2.PlatformConfig{
					ResourceConfig:         bpv2.ResourceConfig{MinVolumeSize: "1G"},
					ReverseProxyConfig:     bpv2.ReverseProxyConfig{MaxBodySize: "100M", RewriteTarget: "/", AdditionalConfig: "additional"},
					AdditionalMountsConfig: nil,
				}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToDoguDTOs(tt.args.dogus)
			assert.Equalf(t, tt.want, got, "ConvertToDoguDTOs(%v)", tt.args.dogus)
		})
	}
}
