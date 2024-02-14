package componentcr

import (
	_ "embed"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

const (
	testNamespace       = "ecosystem"
	testStatus          = ecosystem.ComponentStatusNotInstalled
	testHealthStatus    = compV1.AvailableHealthStatus
	testResourceVersion = "1"
)

var (
	//go:embed testdata/testPatch
	testPatchBytes  []byte
	testVersion1, _ = semver.NewVersion("1.0.0-1")
)

func Test_parseComponentCR(t *testing.T) {
	type args struct {
		cr *compV1.Component
	}
	tests := []struct {
		name    string
		args    args
		want    *ecosystem.ComponentInstallation
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				cr: &compV1.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name:            testComponentNameRaw,
						Namespace:       testNamespace,
						ResourceVersion: testResourceVersion,
					},
					Spec: compV1.ComponentSpec{
						Namespace: testDistributionNamespace,
						Name:      testComponentNameRaw,
						Version:   testVersion1.String(),
					},
					Status: compV1.ComponentStatus{
						Status: testStatus,
						Health: testHealthStatus,
					},
				},
			},
			want: &ecosystem.ComponentInstallation{
				Name:    testComponentName,
				Version: testVersion1,
				Status:  testStatus,
				Health:  ecosystem.HealthStatus(testHealthStatus),
				PersistenceContext: map[string]interface{}{
					componentInstallationRepoContextKey: componentInstallationRepoContext{testResourceVersion},
				},
			},
			wantErr: false,
		},
		{
			name: "should return error on nil component",
			args: args{
				cr: nil,
			},
			wantErr: true,
		},
		{
			name: "should return error version parse error",
			args: args{
				cr: &compV1.Component{
					Spec: compV1.ComponentSpec{
						Version: "fsdfsd",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseComponentCR(tt.args.cr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseComponentCR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseComponentCR() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toComponentCR(t *testing.T) {
	type args struct {
		componentInstallation *ecosystem.ComponentInstallation
	}
	tests := []struct {
		name string
		args args
		want *compV1.Component
	}{
		{
			name: "success",
			args: args{
				componentInstallation: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
					Status:  testStatus,
					Health:  ecosystem.HealthStatus(testHealthStatus),
					PersistenceContext: map[string]interface{}{
						componentInstallationRepoContextKey: componentInstallationRepoContext{testResourceVersion},
					},
				},
			},
			want: &compV1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name: testComponentNameRaw,
					Labels: map[string]string{
						ComponentNameLabelKey:    testComponentNameRaw,
						ComponentVersionLabelKey: testVersion1.String(),
					},
				},
				Spec: compV1.ComponentSpec{
					Namespace: testDistributionNamespace,
					Name:      testComponentNameRaw,
					Version:   testVersion1.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toComponentCR(tt.args.componentInstallation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toComponentCR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toComponentCRPatch(t *testing.T) {
	type args struct {
		component *ecosystem.ComponentInstallation
	}
	tests := []struct {
		name string
		args args
		want *componentCRPatch
	}{
		{
			name: "success",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
				},
			},
			want: &componentCRPatch{
				Spec: componentSpecPatch{
					Namespace: testDistributionNamespace,
					Name:      testComponentNameRaw,
					Version:   testVersion1.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toComponentCRPatch(tt.args.component); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toComponentCRPatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toComponentCRPatchBytes(t *testing.T) {
	type args struct {
		component *ecosystem.ComponentInstallation
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
				},
			},
			want:    testPatchBytes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toComponentCRPatchBytes(tt.args.component)
			if (err != nil) != tt.wantErr {
				t.Errorf("toComponentCRPatchBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toComponentCRPatchBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
