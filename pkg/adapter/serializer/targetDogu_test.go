package serializer

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	version3211, _ = core.ParseVersion("3.2.1-1")
)

func TestConvertDogus(t *testing.T) {
	type args struct {
		dogus []TargetDogu
	}
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
			want:    []domain.Dogu{{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent}},
			wantErr: assert.NoError,
		},
		// {
		// 	name:    "dogu with max proxy body size",
		// 	args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{MaxBodySize: "1.2G"}}}}},
		// 	want:    []domain.Dogu{{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: resource.MustParse("1200M")}}},
		// 	wantErr: assert.NoError,
		// },
		// {
		// 	name:    "dogu with max proxy body size should converted to decimal prefix",
		// 	args:    args{dogus: []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{MaxBodySize: "1Gi"}}}}},
		// 	want:    []domain.Dogu{{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: resource.MustParse("1073.74M")}}},
		// 	wantErr: assert.NoError,
		// },
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
		// {
		// 	name:    "ok",
		// 	args:    args{dogus: []domain.Dogu{{Name: common.QualifiedDoguName{Namespace: "official", SimpleName: "postgres"}, Version: version3211, TargetState: domain.TargetStatePresent, ReverseProxyConfig: ecosystem.ReverseProxyConfig{MaxBodySize: resource.MustParse("100M")}}}},
		// 	want:    []TargetDogu{{Name: "official/postgres", Version: version3211.Raw, TargetState: "present", PlatformConfig: PlatformConfig{ReverseProxyConfig: ReverseProxyConfig{MaxBodySize: "100M"}}}},
		// 	wantErr: assert.NoError,
		// },
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
