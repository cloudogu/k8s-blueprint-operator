package componentcr

import (
	_ "embed"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	compV1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const (
	testDeployNamespace = "ecosystem"
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
	valuesOverwrite := map[string]interface{}{"key": "value", "key1": map[string]string{"key": "value"}}
	expectedValuesOverwrite := map[string]interface{}{"key": "value", "key1": map[string]interface{}{"key": "value"}}
	valuesOverwriteYAMLBytes, err := yaml.Marshal(valuesOverwrite)
	require.NoError(t, err)

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
						Namespace:       testDeployNamespace,
						ResourceVersion: testResourceVersion,
					},
					Spec: compV1.ComponentSpec{
						Namespace:           testNamespace,
						Name:                testComponentNameRaw,
						Version:             testVersion1.String(),
						DeployNamespace:     "longhorn-system",
						ValuesYamlOverwrite: string(valuesOverwriteYAMLBytes),
					},
					Status: compV1.ComponentStatus{
						Status:           testStatus,
						Health:           testHealthStatus,
						InstalledVersion: testVersion1.String(),
					},
				},
			},
			want: &ecosystem.ComponentInstallation{
				Name:            testComponentName,
				ExpectedVersion: testVersion1,
				ActualVersion:   testVersion1,
				Status:          testStatus,
				Health:          ecosystem.HealthStatus(testHealthStatus),
				DeployConfig: map[string]interface{}{
					"deployNamespace": "longhorn-system",
					"overwriteConfig": expectedValuesOverwrite,
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
			name: "should return expected version parse error",
			args: args{
				cr: &compV1.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name: testComponentNameRaw,
					},
					Spec: compV1.ComponentSpec{
						Namespace: testNamespace,
						Name:      testComponentNameRaw,
						Version:   "fsdfsd",
					},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "should return actual version parse error",
			args: args{
				cr: &compV1.Component{
					ObjectMeta: metav1.ObjectMeta{
						Name:            testComponentNameRaw,
						ResourceVersion: testResourceVersion,
					},
					Spec: compV1.ComponentSpec{
						Namespace: testNamespace,
						Name:      testComponentNameRaw,
						Version:   testVersion1.String(),
					},
					Status: compV1.ComponentStatus{
						InstalledVersion: "fsdfsd",
					},
				},
			},
			want: &ecosystem.ComponentInstallation{
				Name:            testComponentName,
				ExpectedVersion: testVersion1,
				ActualVersion:   nil,
				DeployConfig:    map[string]interface{}{},
				PersistenceContext: map[string]interface{}{
					componentInstallationRepoContextKey: componentInstallationRepoContext{testResourceVersion},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return error on missing resource name",
			args: args{
				cr: &compV1.Component{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: compV1.ComponentSpec{
						Namespace: testNamespace,
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
						Namespace:           testNamespace,
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
					Name:            testComponentName,
					ExpectedVersion: testVersion1,
					Status:          testStatus,
					Health:          ecosystem.HealthStatus(testHealthStatus),
					PersistenceContext: map[string]interface{}{
						componentInstallationRepoContextKey: componentInstallationRepoContext{testResourceVersion},
					},
					DeployConfig: map[string]interface{}{
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
					Namespace:           testNamespace,
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
	testDeployNamespace := "longhorn-system"
	testDeployConfig := "key: value\n"
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
					Name:            testComponentName,
					ExpectedVersion: testVersion1,
					DeployConfig: map[string]interface{}{
						"deployNamespace": "longhorn-system",
						"overwriteConfig": map[string]interface{}{"key": "value"},
					},
				},
			},
			want: &componentCRPatch{
				Spec: componentSpecPatch{
					Namespace:           testNamespace,
					Name:                testComponentNameRaw,
					Version:             testVersion1.String(),
					DeployNamespace:     &testDeployNamespace,
					ValuesYamlOverwrite: &testDeployConfig,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should return error on wrong package config type",
			args: args{
				component: &ecosystem.ComponentInstallation{
					Name:            testComponentName,
					ExpectedVersion: testVersion1,
					DeployConfig: map[string]interface{}{
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
					Name:            testComponentName,
					ExpectedVersion: testVersion1,
					DeployConfig: map[string]interface{}{
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
					Name:            testComponentName,
					ExpectedVersion: testVersion1,
					DeployConfig: map[string]interface{}{
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
