package dogucr

import (
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

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
				Namespace: "official",
				Name:      "postgresql",
				Version:   version3_2_1_4,
				Status:    ecosystem.DoguStatusInstalled,
				Health:    ecosystem.AvailableHealthStatus,
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
		name    string
		dogu    *ecosystem.DoguInstallation
		want    *v1.Dogu
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			dogu: &ecosystem.DoguInstallation{
				Namespace: "official",
				Name:      "postgresql",
				Version:   version3_2_1_4,
				Status:    ecosystem.DoguStatusInstalled,
				Health:    ecosystem.AvailableHealthStatus,
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
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toDoguCR(tt.dogu)
			if !tt.wantErr(t, err, fmt.Sprintf("toDoguCR(%v)", tt.dogu)) {
				return
			}
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
				Namespace: "official",
				Name:      "postgresql",
				Version:   version3_2_1_4,
				Status:    ecosystem.DoguStatusInstalled,
				Health:    ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
			},
			want: &doguCRPatch{
				Spec: doguSpecPatch{
					Name:    "official/postgresql",
					Version: version3_2_1_4.Raw,
					//Resources: doguResourcesPatch{
					//	DataVolumeSize: "",
					//},
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
			name: "ok",
			dogu: &ecosystem.DoguInstallation{
				Namespace: "official",
				Name:      "postgresql",
				Version:   version3_2_1_4,
				Status:    ecosystem.DoguStatusInstalled,
				Health:    ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
			},
			want:    "{\"spec\":{\"name\":\"official/postgresql\",\"version\":\"3.2.1-4\",\"supportMode\":false,\"upgradeConfig\":{\"allowNamespaceSwitch\":true,\"forceUpgrade\":false}}}",
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
