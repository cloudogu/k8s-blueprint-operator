package componentcr

import (
	_ "embed"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		wantErr assert.ErrorAssertionFunc
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
						Namespace:           testDistributionNamespace,
						Name:                testComponentNameRaw,
						Version:             testVersion1.String(),
						DeployNamespace:     "longhorn-system",
						ValuesYamlOverwrite: "key: value",
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
				PackageConfig: map[string]interface{}{
					"deployNamespace": "longhorn-system",
					// TODO check cluster error
					"overwriteConfig": map[string]interface{}{"key": "value"},
				},
				PersistenceContext: map[string]interface{}{
					componentInstallationRepoContextKey: componentInstallationRepoContext{testResourceVersion},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return error on nil component",
			args: args{
				cr: nil,
			},
			wantErr: assert.Error,
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
			wantErr: assert.Error,
		},
		{
			name: "should return error on missing resource name",
			args: args{
				cr: &compV1.Component{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: compV1.ComponentSpec{
						Namespace: testDistributionNamespace,
						Name:      testComponentNameRaw,
						Version:   testVersion1.String(),
					},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "should return error on wrong package config value",
			args: args{
				cr: &compV1.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name: "name",
					},
					Spec: compV1.ComponentSpec{
						Namespace:           testDistributionNamespace,
						Name:                testComponentNameRaw,
						Version:             testVersion1.String(),
						DeployNamespace:     "longhorn-system",
						ValuesYamlOverwrite: "no yaml object",
					},
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseComponentCR(tt.args.cr)
			if !tt.wantErr(t, err, fmt.Sprintf("parseComponentCR(%v)", tt.args.cr)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_toComponentCR(t *testing.T) {
	type args struct {
		componentInstallation *ecosystem.ComponentInstallation
	}
	tests := []struct {
		name    string
		args    args
		want    *compV1.Component
		wantErr assert.ErrorAssertionFunc
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
					PackageConfig: map[string]interface{}{
						"deployNamespace": "longhorn-system",
						"overwriteConfig": map[string]interface{}{"key": "value"},
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
					Namespace:           testDistributionNamespace,
					Name:                testComponentNameRaw,
					Version:             testVersion1.String(),
					DeployNamespace:     "longhorn-system",
					ValuesYamlOverwrite: "key: value\n",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toComponentCR(tt.args.componentInstallation)
			if !tt.wantErr(t, err, fmt.Sprintf("toComponentCR(%v)", tt.args.componentInstallation)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_toComponentCRPatch(t *testing.T) {
	type args struct {
		component *ecosystem.ComponentInstallation
	}
	tests := []struct {
		name    string
		args    args
		want    *componentCRPatch
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
					PackageConfig: map[string]interface{}{
						"deployNamespace": "longhorn-system",
						"overwriteConfig": map[string]interface{}{"key": "value"},
					},
				},
			},
			want: &componentCRPatch{
				Spec: componentSpecPatch{
					Namespace:           testDistributionNamespace,
					Name:                testComponentNameRaw,
					Version:             testVersion1.String(),
					DeployNamespace:     "longhorn-system",
					ValuesYamlOverwrite: "key: value\n",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return error on wrong package config type",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
					PackageConfig: map[string]interface{}{
						"deployNamespace": map[string]interface{}{"no": "string"},
					},
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toComponentCRPatch(tt.args.component)
			if !tt.wantErr(t, err, fmt.Sprintf("toComponentCRPatch(%v)", tt.args.component)) {
				return
			}
			assert.Equal(t, tt.want, got)
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
					PackageConfig: map[string]interface{}{
						"deployNamespace": "longhorn-system",
						"overwriteConfig": map[string]interface{}{"key": "value"},
					},
				},
			},
			want:    testPatchBytes,
			wantErr: assert.NoError,
		},
		{
			name: "should return error on error creating patch",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:    testComponentName,
					Version: testVersion1,
					PackageConfig: map[string]interface{}{
						"deployNamespace": map[string]interface{}{"no": "string"},
					},
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toComponentCRPatchBytes(tt.args.component)
			if !tt.wantErr(t, err, fmt.Sprintf("toComponentCRPatchBytes(%v)", tt.args.component)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
