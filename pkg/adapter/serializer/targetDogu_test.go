package serializer

import (
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

var (
	version3211, _ = core.ParseVersion("3.2.1-1")
)

func TestConvertDogus(t *testing.T) {
	type args struct {
		dogus []TargetDogu
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
			args:    args{dogus: []TargetDogu{}},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "normal dogu",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present"}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with max proxy body size",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{MaxBodySize: "1G"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &proxyBodySize}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with proxy rewrite target",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{RewriteTarget: "/"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{RewriteTarget: "/"}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with proxy additional config",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{AdditionalConfig: "additional"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{AdditionalConfig: "additional"}}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with invalid proxy body size",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{MaxBodySize: "1GE"}}}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "dogu with min volume size",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ResourceConfig: ResourceConfig{MinVolumeSize: "1Gi"}}}}},
			want:    []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, MinVolumeSize: &volumeSize}},
			wantErr: assert.NoError,
		},
		{
			name:    "dogu with invalid volume size",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ResourceConfig: ResourceConfig{MinVolumeSize: "1GIE"}}}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "no namespace",
			args:    args{dogus: []TargetDogu{{Name: "postgres", Version: version3211.Raw, TargetState: "present"}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unparsable version",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: "1.", TargetState: "present"}}},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "unknown target state",
			args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "unknown"}}},
			want:    nil,
			wantErr: assert.Error,
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

func TestConvertToDoguDTOs(t *testing.T) {
	type args struct {
		dogus []domain.Dogu
	}
	bodySize := resource.MustParse("100M")
	volumeSize := resource.MustParse("1G")
	tests := []struct {
		name    string
		args    args
		want    []TargetDogu
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{},
			want:    []TargetDogu{},
			wantErr: assert.NoError,
		},
		{
			name:    "empty list",
			args:    args{dogus: []domain.Dogu{}},
			want:    []TargetDogu{},
			wantErr: assert.NoError,
		},
		{
			name:    "ok",
			args:    args{dogus: []domain.Dogu{{Name: cescommons.QualifiedName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, MinVolumeSize: &volumeSize, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: &bodySize, RewriteTarget: "/", AdditionalConfig: "additional"}}}},
			want:    []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ResourceConfig: ResourceConfig{MinVolumeSize: "1G"}, ReverseProxyConfig: ReverseProxyConfig{MaxBodySize: "100M", RewriteTarget: "/", AdditionalConfig: "additional"}}}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToDoguDTOs(tt.args.dogus)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertToDoguDTOs(%v)", tt.args.dogus)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertToDoguDTOs(%v)", tt.args.dogus)
		})
	}
}
