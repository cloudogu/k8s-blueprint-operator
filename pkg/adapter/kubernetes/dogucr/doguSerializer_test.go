package dogucr

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

var postgresDoguName = common.QualifiedDoguName{
	Namespace:  common.DoguNamespace("official"),
	SimpleName: common.SimpleDoguName("postgresql"),
}

func Test_parseDoguCR(t *testing.T) {
	type args struct {
		cr *v1.Dogu
	}
	tests := []struct {
		name    string
		args    args
		want    *ecosystem.DoguInstallation
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{cr: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			args: args{cr: &v1.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: crResourceVersion,
				},
				Spec: v1.DoguSpec{
					Name:      "official/postgresql",
					Version:   version3_2_1_4.Raw,
					Resources: v1.DoguResources{},
					UpgradeConfig: v1.UpgradeConfig{
						AllowNamespaceSwitch: true,
					},
				},
				Status: v1.DoguStatus{
					Status: v1.DoguStatusInstalled,
					Health: v1.AvailableHealthStatus,
				},
			}},
			want: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3_2_1_4,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				PersistenceContext: persistenceContext,
			},
			wantErr: false,
		},
		{
			name: "cannot parse version",
			args: args{cr: &v1.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v1.DoguSpec{
					Name:      "official/postgresql",
					Version:   "vxyz",
					Resources: v1.DoguResources{},
					UpgradeConfig: v1.UpgradeConfig{
						AllowNamespaceSwitch: false,
					},
				},
				Status: v1.DoguStatus{
					Status: v1.DoguStatusInstalled,
					Health: v1.AvailableHealthStatus,
				},
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDoguCR(tt.args.cr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDoguCR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDoguCR() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDoguCR(t *testing.T) {
	tests := []struct {
		name string
		dogu *ecosystem.DoguInstallation
		want *v1.Dogu
	}{
		{
			name: "ok",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3_2_1_4,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
			},
			want: &v1.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
					Labels: map[string]string{
						"app":       "ces",
						"dogu.name": "postgresql",
					},
				},
				Spec: v1.DoguSpec{
					Name:    "official/postgresql",
					Version: version3_2_1_4.Raw,
					Resources: v1.DoguResources{
						DataVolumeSize: "",
					},
					UpgradeConfig: v1.UpgradeConfig{
						AllowNamespaceSwitch: true,
						ForceUpgrade:         false,
					},
					AdditionalIngressAnnotations: nil,
				},
				Status: v1.DoguStatus{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toDoguCR(tt.dogu)
			assert.Equalf(t, tt.want, got, "toDoguCR(%v)", tt.dogu)
		})
	}
}

func Test_toDoguCRPatch(t *testing.T) {
	tests := []struct {
		name string
		dogu *ecosystem.DoguInstallation
		want *doguCRPatch
	}{
		{
			name: "ok",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3_2_1_4,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
			},
			want: &doguCRPatch{
				Spec: doguSpecPatch{
					Name:    "official/postgresql",
					Version: version3_2_1_4.Raw,
					Resources: doguResourcesPatch{
						DataVolumeSize: "0",
					},
					UpgradeConfig: upgradeConfigPatch{
						AllowNamespaceSwitch: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toDoguCRPatch(tt.dogu)
			assert.Equalf(t, tt.want, got, "toDoguCR(%v)", tt.dogu)
		})
	}
}

func Test_toDoguCRPatchBytes(t *testing.T) {
	tests := []struct {
		name    string
		dogu    *ecosystem.DoguInstallation
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			// TODO check ReverseProxy
			name: "ok",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3_2_1_4,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				MinVolumeSize: resource.MustParse("2Gi"),
			},
			want:    "{\"spec\":{\"name\":\"official/postgresql\",\"version\":\"3.2.1-4\",\"resources\":{\"dataVolumeSize\":\"2Gi\"},\"supportMode\":false,\"upgradeConfig\":{\"allowNamespaceSwitch\":true,\"forceUpgrade\":false},\"additionalIngressAnnotations\":null}}",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toDoguCRPatchBytes(tt.dogu)
			if !tt.wantErr(t, err, fmt.Sprintf("toDoguCRPatchBytes(%v)", tt.dogu)) {
				return
			}
			assert.Equalf(t, tt.want, string(got), "toDoguCRPatchBytes(%v)", tt.dogu)
		})
	}
}

func Test_getNginxIngressAnnotations(t *testing.T) {
	type args struct {
		config ecosystem.ReverseProxyConfig
	}
	zeroQuantity := resource.MustParse("0")
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "should set proxy body size on zero quantity",
			args: args{config: ecosystem.ReverseProxyConfig{MaxBodySize: &zeroQuantity}},
			want: map[string]string{ecosystem.NginxIngressAnnotationBodySize: "0"},
		},
		{
			name: "should not set proxy body size on nil",
			args: args{config: ecosystem.ReverseProxyConfig{}},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNginxIngressAnnotations(tt.args.config), "getNginxIngressAnnotations(%v)", tt.args.config)
		})
	}
}
