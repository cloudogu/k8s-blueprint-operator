package serializer

import (
	"fmt"
	"testing"

	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
)

var (
	wrongVersion        = "1."
	rewriteTarget       = "/"
	additionalConfig    = "additional"
	volumeSize          = resource.MustParse("1Gi")
	volumeSizeString    = volumeSize.String()
	proxyBodySize       = resource.MustParse("1G")
	proxyBodySizeString = proxyBodySize.String()
	subfolder           = "subfolder"
	subfolder2          = "secsubfolder"
)

func TestConvertDogus(t *testing.T) {
	type args struct {
		dogus []bpv3.Dogu
	}

	wrongBodySize := "1GE"
	wrongVolumeSize := "1GIE"

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
			args:    args{dogus: []bpv3.Dogu{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal dogu",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: &version3211, Absent: false}},
			wantErr: assert.NoError,
		},
		{
			name:    "absent dogu",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Absent: &trueVar}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Absent: true}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with max proxy body size",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar, PlatformConfig: &bpv3.PlatformConfig{ReverseProxyConfig: &bpv3.ReverseProxyConfig{MaxBodySize: &proxyBodySizeString}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: &version3211, Absent: false, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &proxyBodySize}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with proxy rewrite target",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar, PlatformConfig: &bpv3.PlatformConfig{ReverseProxyConfig: &bpv3.ReverseProxyConfig{RewriteTarget: &rewriteTarget}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: &version3211, Absent: false, ReverseProxyConfig: ecosystem.ReverseProxyConfig{RewriteTarget: ecosystem.RewriteTarget(rewriteTarget)}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with proxy additional config",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar, PlatformConfig: &bpv3.PlatformConfig{ReverseProxyConfig: &bpv3.ReverseProxyConfig{AdditionalConfig: &additionalConfig}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: &version3211, Absent: false, ReverseProxyConfig: ecosystem.ReverseProxyConfig{AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig)}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with invalid proxy body size",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar, PlatformConfig: &bpv3.PlatformConfig{ReverseProxyConfig: &bpv3.ReverseProxyConfig{MaxBodySize: &wrongBodySize}}}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "dogu with min volume size",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, PlatformConfig: &bpv3.PlatformConfig{ResourceConfig: &bpv3.ResourceConfig{MinVolumeSize: &volumeSizeString}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: &version3211, Absent: false, MinVolumeSize: &volumeSize}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with invalid volume size",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, PlatformConfig: &bpv3.PlatformConfig{ResourceConfig: &bpv3.ResourceConfig{MinVolumeSize: &wrongVolumeSize}}}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "no namespace",
			args:    args{dogus: []bpv3.Dogu{{Name: "postgres", Version: &version3211.Raw}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unparsable version",
			args:    args{dogus: []bpv3.Dogu{{Name: "official/postgres", Version: &wrongVersion}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "should convert additionalMounts",
			args: args{dogus: []bpv3.Dogu{{
				Name:    "official/postgres",
				Version: &version3211.Raw,
				PlatformConfig: &bpv3.PlatformConfig{
					AdditionalMountsConfig: []bpv3.AdditionalMount{
						{
							SourceType: bpv3.DataSourceConfigMap,
							Name:       "configMap",
							Volume:     "volume",
							Subfolder:  &subfolder,
						},
						{
							SourceType: bpv3.DataSourceSecret,
							Name:       "sec",
							Volume:     "secvolume",
							Subfolder:  &subfolder2,
						},
					},
				},
			}}},
			want: []domain.Dogu{{
				Name:    cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version: &version3211,
				Absent:  false,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configMap",
						Volume:     "volume",
						Subfolder:  subfolder,
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "sec",
						Volume:     "secvolume",
						Subfolder:  subfolder2,
					},
				}}},
			wantErr: assert.NoError,
		},
		{
			name: "should return nil slice if dogu contains an nil slice",
			args: args{dogus: []bpv3.Dogu{{
				Name:    "official/postgres",
				Version: &version3211.Raw,
				PlatformConfig: &bpv3.PlatformConfig{
					AdditionalMountsConfig: nil,
				},
			}}},
			want: []domain.Dogu{{
				Name:             cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:          &version3211,
				Absent:           false,
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
		dogus []bpv3.MaskDogu
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
			args:    args{dogus: []bpv3.MaskDogu{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal dogu",
			args:    args{dogus: []bpv3.MaskDogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar}}},
			want:    []domain.MaskDogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, Absent: false}},
			wantErr: assert.NoError,
		},
		{
			name:    "absent dogu",
			args:    args{dogus: []bpv3.MaskDogu{{Name: "official/postgres", Absent: &trueVar}}},
			want:    []domain.MaskDogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Absent: true}},
			wantErr: assert.NoError,
		},
		{
			name:    "no namespace",
			args:    args{dogus: []bpv3.MaskDogu{{Name: "postgres", Version: &version3211.Raw}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unparsable version",
			args:    args{dogus: []bpv3.MaskDogu{{Name: "official/postgres", Version: &wrongVersion}}},
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
	tests := []struct {
		name string
		args args
		want []bpv3.Dogu
	}{
		{
			name: "nil",
			args: args{},
			want: []bpv3.Dogu{},
		},
		{
			name: "empty list",
			args: args{dogus: []domain.Dogu{}},
			want: []bpv3.Dogu{},
		},
		{
			name: "ok",
			args: args{dogus: []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: &version3211, Absent: false, MinVolumeSize: &volumeSize, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &proxyBodySize, RewriteTarget: ecosystem.RewriteTarget(rewriteTarget), AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig)}}}},
			want: []bpv3.Dogu{{Name: "official/postgres", Version: &version3211.Raw, Absent: &falseVar, PlatformConfig: &bpv3.PlatformConfig{ResourceConfig: &bpv3.ResourceConfig{MinVolumeSize: &volumeSizeString}, ReverseProxyConfig: &bpv3.ReverseProxyConfig{MaxBodySize: &proxyBodySizeString, RewriteTarget: &rewriteTarget, AdditionalConfig: &additionalConfig}}}},
		},
		{
			name: "additionalMountsConfig",
			args: args{dogus: []domain.Dogu{{
				Name:               cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:            &version3211,
				Absent:             false,
				MinVolumeSize:      &volumeSize,
				ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &proxyBodySize, RewriteTarget: ecosystem.RewriteTarget(rewriteTarget), AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig)},
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configMap",
						Volume:     "volume",
						Subfolder:  subfolder,
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "sec",
						Volume:     "secvolume",
						Subfolder:  subfolder2,
					},
				},
			}}},
			want: []bpv3.Dogu{{
				Name:    "official/postgres",
				Version: &version3211.Raw,
				Absent:  &falseVar,
				PlatformConfig: &bpv3.PlatformConfig{
					ResourceConfig:     &bpv3.ResourceConfig{MinVolumeSize: &volumeSizeString},
					ReverseProxyConfig: &bpv3.ReverseProxyConfig{MaxBodySize: &proxyBodySizeString, RewriteTarget: &rewriteTarget, AdditionalConfig: &additionalConfig},
					AdditionalMountsConfig: []bpv3.AdditionalMount{
						{
							SourceType: bpv3.DataSourceConfigMap,
							Name:       "configMap",
							Volume:     "volume",
							Subfolder:  &subfolder,
						},
						{
							SourceType: bpv3.DataSourceSecret,
							Name:       "sec",
							Volume:     "secvolume",
							Subfolder:  &subfolder2,
						},
					},
				}}},
		},
		{
			name: "should return nil slice if dogu contains an nil slice",
			args: args{dogus: []domain.Dogu{{
				Name:               cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"},
				Version:            &version3211,
				Absent:             false,
				MinVolumeSize:      &volumeSize,
				ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &proxyBodySize, RewriteTarget: ecosystem.RewriteTarget(rewriteTarget), AdditionalConfig: ecosystem.AdditionalConfig(additionalConfig)},
				AdditionalMounts:   nil,
			}}},
			want: []bpv3.Dogu{{
				Name:    "official/postgres",
				Version: &version3211.Raw,
				Absent:  &falseVar,
				PlatformConfig: &bpv3.PlatformConfig{
					ResourceConfig:         &bpv3.ResourceConfig{MinVolumeSize: &volumeSizeString},
					ReverseProxyConfig:     &bpv3.ReverseProxyConfig{MaxBodySize: &proxyBodySizeString, RewriteTarget: &rewriteTarget, AdditionalConfig: &additionalConfig},
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
