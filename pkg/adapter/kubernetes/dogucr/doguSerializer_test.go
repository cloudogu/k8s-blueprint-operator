package dogucr

import (
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

var postgresDoguName = cescommons.QualifiedName{
	Namespace:  cescommons.Namespace("official"),
	SimpleName: cescommons.SimpleName("postgresql"),
}
var ldapDoguName = cescommons.QualifiedName{
	Namespace:  "official",
	SimpleName: "ldap",
}
var volSize25G = resource.MustParse("25G")
var defaultVolSize = resource.MustParse(v2.DefaultVolumeSize)

func Test_parseDoguCR(t *testing.T) {
	type args struct {
		cr *v2.Dogu
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
			args: args{cr: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: crResourceVersion,
				},
				Spec: v2.DoguSpec{
					Name:      "official/postgresql",
					Version:   version3214.Raw,
					Resources: v2.DoguResources{},
					UpgradeConfig: v2.UpgradeConfig{
						AllowNamespaceSwitch: true,
					},
				},
				Status: v2.DoguStatus{
					Status: v2.DoguStatusInstalled,
					Health: v2.AvailableHealthStatus,
				},
			}},
			want: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3214,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				MinVolumeSize:      defaultVolSize,
				PersistenceContext: persistenceContext,
			},
			wantErr: false,
		},
		{
			name: "cannot parse version",
			args: args{cr: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v2.DoguSpec{
					Name:      "official/postgresql",
					Version:   "vxyz",
					Resources: v2.DoguResources{},
					UpgradeConfig: v2.UpgradeConfig{
						AllowNamespaceSwitch: false,
					},
				},
				Status: v2.DoguStatus{
					Status: v2.DoguStatusInstalled,
					Health: v2.AvailableHealthStatus,
				},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse additional mounts",
			args: args{cr: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v2.DoguSpec{
					Name:      "official/postgresql",
					Version:   version3214.Raw,
					Resources: v2.DoguResources{},
					AdditionalMounts: []v2.DataMount{
						{
							SourceType: v2.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  "subfolder",
						},
						{
							SourceType: v2.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  "secsubfolder",
						},
					},
				},
			}},
			want: &ecosystem.DoguInstallation{
				Name:               postgresDoguName,
				Version:            version3214,
				PersistenceContext: persistenceContext,
				MinVolumeSize:      defaultVolSize,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configmap",
						Volume:     "volume",
						Subfolder:  "subfolder",
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "secret",
						Volume:     "secvolume",
						Subfolder:  "secsubfolder",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse deprecated DataVolumeSize",
			args: args{cr: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: version3214.Raw,
					Resources: v2.DoguResources{
						DataVolumeSize: "25G",
					},
				},
			}},
			want: &ecosystem.DoguInstallation{
				Name:               postgresDoguName,
				Version:            version3214,
				PersistenceContext: persistenceContext,
				MinVolumeSize:      volSize25G,
			},
			wantErr: false,
		},
		{
			name: "prefer MinDataVolumeSize over DataVolumeSize",
			args: args{cr: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: version3214.Raw,
					Resources: v2.DoguResources{
						DataVolumeSize:    "20G",
						MinDataVolumeSize: volSize25G,
					},
				},
			}},
			want: &ecosystem.DoguInstallation{
				Name:               postgresDoguName,
				Version:            version3214,
				PersistenceContext: persistenceContext,
				MinVolumeSize:      volSize25G,
			},
			wantErr: false,
		},
		{
			name: "parse MinDataVolumeSize without DataVolumeSize",
			args: args{cr: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: version3214.Raw,
					Resources: v2.DoguResources{
						MinDataVolumeSize: volSize25G,
					},
				},
			}},
			want: &ecosystem.DoguInstallation{
				Name:               postgresDoguName,
				Version:            version3214,
				PersistenceContext: persistenceContext,
				MinVolumeSize:      volSize25G,
			},
			wantErr: false,
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
				t.Errorf("parseDoguCR() \ngot = %+v \nwant= %+v", got, tt.want)
			}
		})
	}
}

func Test_toDoguCR(t *testing.T) {
	tests := []struct {
		name string
		dogu *ecosystem.DoguInstallation
		want *v2.Dogu
	}{
		{
			name: "ok",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3214,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
			},
			want: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
					Labels: map[string]string{
						"app":       "ces",
						"dogu.name": "postgresql",
					},
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: version3214.Raw,
					Resources: v2.DoguResources{
						DataVolumeSize: "",
					},
					UpgradeConfig: v2.UpgradeConfig{
						AllowNamespaceSwitch: true,
						ForceUpgrade:         false,
					},
					AdditionalIngressAnnotations: nil,
				},
				Status: v2.DoguStatus{},
			},
		},
		{
			name: "convert additional mounts",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3214,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{
						SourceType: ecosystem.DataSourceConfigMap,
						Name:       "configmap",
						Volume:     "volume",
						Subfolder:  "subfolder",
					},
					{
						SourceType: ecosystem.DataSourceSecret,
						Name:       "secret",
						Volume:     "secvolume",
						Subfolder:  "secsubfolder",
					},
				},
			},
			want: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
					Labels: map[string]string{
						"app":       "ces",
						"dogu.name": "postgresql",
					},
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: version3214.Raw,
					AdditionalMounts: []v2.DataMount{
						{
							SourceType: v2.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  "subfolder",
						},
						{
							SourceType: v2.DataSourceSecret,
							Name:       "secret",
							Volume:     "secvolume",
							Subfolder:  "secsubfolder",
						},
					},
				},
				Status: v2.DoguStatus{},
			},
		},
		{
			name: "set MinDataVolumeSize instead of deprecated DataVolumeSize",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3214,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				MinVolumeSize: volSize25G,
			},
			want: &v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
					Labels: map[string]string{
						"app":       "ces",
						"dogu.name": "postgresql",
					},
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: version3214.Raw,
					Resources: v2.DoguResources{
						MinDataVolumeSize: volSize25G,
					},
					UpgradeConfig: v2.UpgradeConfig{
						AllowNamespaceSwitch: true,
						ForceUpgrade:         false,
					},
					AdditionalIngressAnnotations: nil,
				},
				Status: v2.DoguStatus{},
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
				Version: version3214,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
			},
			want: &doguCRPatch{
				Spec: doguSpecPatch{
					Name:    "official/postgresql",
					Version: version3214.Raw,
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
	quantity2 := resource.MustParse("2Gi")
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
				Version: version3214,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				MinVolumeSize: quantity2,
				AdditionalMounts: []ecosystem.AdditionalMount{
					{SourceType: ecosystem.DataSourceConfigMap, Name: "test", Volume: "volume", Subfolder: "subfolder"},
				},
			},
			want:    "{\"spec\":{\"name\":\"official/postgresql\",\"version\":\"3.2.1-4\",\"resources\":{\"dataVolumeSize\":\"\",\"minDataVolumeSize\":\"2Gi\"},\"supportMode\":false,\"upgradeConfig\":{\"allowNamespaceSwitch\":true,\"forceUpgrade\":false},\"additionalIngressAnnotations\":null,\"additionalMounts\":[{\"sourceType\":\"ConfigMap\",\"name\":\"test\",\"volume\":\"volume\",\"subfolder\":\"subfolder\"}]}}",
			wantErr: assert.NoError,
		},
		{
			name: "test default minVolumeSize",
			dogu: &ecosystem.DoguInstallation{
				Name:    postgresDoguName,
				Version: version3214,
				Status:  ecosystem.DoguStatusInstalled,
				Health:  ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				AdditionalMounts: []ecosystem.AdditionalMount{
					{SourceType: ecosystem.DataSourceConfigMap, Name: "test", Volume: "volume", Subfolder: "subfolder"},
				},
			},
			want:    "{\"spec\":{\"name\":\"official/postgresql\",\"version\":\"3.2.1-4\",\"resources\":{\"dataVolumeSize\":\"\",\"minDataVolumeSize\":\"0\"},\"supportMode\":false,\"upgradeConfig\":{\"allowNamespaceSwitch\":true,\"forceUpgrade\":false},\"additionalIngressAnnotations\":null,\"additionalMounts\":[{\"sourceType\":\"ConfigMap\",\"name\":\"test\",\"volume\":\"volume\",\"subfolder\":\"subfolder\"}]}}",
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
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNginxIngressAnnotations(tt.args.config), "getNginxIngressAnnotations(%v)", tt.args.config)
		})
	}
}

func Test_parseDoguAdditionalIngressAnnotationsCR(t *testing.T) {
	quantity1 := resource.MustParse("1G")
	type args struct {
		annotations v2.IngressAnnotations
	}
	tests := []struct {
		name    string
		args    args
		want    ecosystem.ReverseProxyConfig
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should parse annotations",
			args: args{
				annotations: v2.IngressAnnotations{
					"nginx.ingress.kubernetes.io/proxy-body-size":       "1G",
					"nginx.ingress.kubernetes.io/rewrite-target":        "/",
					"nginx.ingress.kubernetes.io/configuration-snippet": "snippet",
				},
			},
			want: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &quantity1,
				RewriteTarget:    "/",
				AdditionalConfig: "snippet",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "should return internal error on invalid quantity",
			args: args{
				annotations: v2.IngressAnnotations{
					"nginx.ingress.kubernetes.io/proxy-body-size": "1GG",
				},
			},
			want: ecosystem.ReverseProxyConfig{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "failed to parse quantity \"1GG\"")
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDoguAdditionalIngressAnnotationsCR(tt.args.annotations)
			if !tt.wantErr(t, err, fmt.Sprintf("parseDoguAdditionalIngressAnnotationsCR(%v)", tt.args.annotations)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseDoguAdditionalIngressAnnotationsCR(%v)", tt.args.annotations)
		})
	}
}

func Test_getNginxIngressAnnotations1(t *testing.T) {
	quantity := resource.MustParse("1M")
	type args struct {
		config ecosystem.ReverseProxyConfig
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "should parse config",
			args: args{config: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &quantity,
				RewriteTarget:    "/",
				AdditionalConfig: "additional",
			}},
			want: map[string]string{
				"nginx.ingress.kubernetes.io/proxy-body-size":       "1M",
				"nginx.ingress.kubernetes.io/rewrite-target":        "/",
				"nginx.ingress.kubernetes.io/configuration-snippet": "additional",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNginxIngressAnnotations(tt.args.config), "getNginxIngressAnnotations(%v)", tt.args.config)
		})
	}
}
