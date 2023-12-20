package dogucr

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

func Test_parseDoguCR(t *testing.T) {
	type args struct {
		cr *v1.Dogu
	}
	tests := []struct {
		name    string
		args    args
		want    ecosystem.DoguInstallation
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{cr: nil},
			want:    ecosystem.DoguInstallation{},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{cr: &v1.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v1.DoguSpec{
					Name:      "official/postgresql",
					Version:   version3_2_1_4.Raw,
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
			want: ecosystem.DoguInstallation{
				Namespace:     "official",
				Name:          "postgresql",
				Version:       version3_2_1_4,
				Status:        ecosystem.DoguStatusInstalled,
				Health:        ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{},
				PersistenceContext: map[string]interface{}{
					doguInstallationRepoContextKey: doguInstallationRepoContext{resourceVersion: "abc"},
				},
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
			want:    ecosystem.DoguInstallation{},
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
