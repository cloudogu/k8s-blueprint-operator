package debugmodecr

import (
	"reflect"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	v1 "github.com/cloudogu/k8s-debug-mode-cr-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_parseDoguCR(t *testing.T) {
	type args struct {
		cr *v1.DebugMode
	}
	tests := []struct {
		name    string
		args    args
		want    *ecosystem.DebugMode
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
			args: args{cr: &v1.DebugMode{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            debugModeSingletonCRName,
					ResourceVersion: crResourceVersion,
				},
				Spec: v1.DebugModeSpec{
					DeactivateTimestamp: metav1.Now(),
					TargetLogLevel:      "DEBUG",
				},
				Status: v1.DebugModeStatus{
					Phase: v1.DebugModeStatusSet,
				},
			}},
			want: &ecosystem.DebugMode{
				Phase: "SetDebugMode",
			},
			wantErr: false,
		},
		{
			name:    "nil cr is internal Error",
			args:    args{cr: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "accepts empty phase",
			args: args{cr: &v1.DebugMode{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            debugModeSingletonCRName,
					ResourceVersion: crResourceVersion,
				},
				Spec: v1.DebugModeSpec{
					DeactivateTimestamp: metav1.Now(),
					TargetLogLevel:      "DEBUG",
				},
				Status: v1.DebugModeStatus{
					Phase: "",
				},
			}},
			want: &ecosystem.DebugMode{
				Phase: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDebugModeCR(tt.args.cr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDebugModeCR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDebugModeCR() \ngot = %+v \nwant= %+v", got, tt.want)
			}
		})
	}
}
