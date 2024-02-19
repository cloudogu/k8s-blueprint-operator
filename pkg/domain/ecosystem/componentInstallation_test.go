package ecosystem

import (
	"github.com/Masterminds/semver/v3"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testNamespace             = "k8s"
	testComponentName         = k8sK8sDoguOperator
	testLonghornComponentName = "k8s-longhorn"
)

var (
	testVersion1, _ = semver.NewVersion("1.0.0")
	testVersion2, _ = semver.NewVersion("2.0.0")
)

func TestInstallComponent(t *testing.T) {
	type args struct {
		componentName common.QualifiedComponentName
		version       *semver.Version
	}
	tests := []struct {
		name string
		args args
		want *ComponentInstallation
	}{
		{
			name: "success",
			args: args{
				componentName: testComponentName,
				version:       testVersion1,
			},
			want: &ComponentInstallation{
				Name:    testComponentName,
				Version: testVersion1,
			},
		},
		{
			name: "longhorn should always be deployed in longhorn-system",
			args: args{
				componentName: k8sK8sLonghorn,
				version:       testVersion1,
			},
			want: &ComponentInstallation{
				Name:            k8sK8sLonghorn,
				Version:         testVersion1,
				DeployNamespace: "longhorn-system",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, InstallComponent(tt.args.componentName, tt.args.version), "InstallComponent(%v, %v, %v)", tt.args.componentName, tt.args.version)
		})
	}
}

func TestComponentInstallation_Upgrade(t *testing.T) {
	type fields struct {
		Name                common.QualifiedComponentName
		DeployNamespace     string
		Version             *semver.Version
		Status              string
		ValuesYamlOverwrite string
		MappedValues        map[string]string
		PersistenceContext  map[string]interface{}
		Health              HealthStatus
	}
	type args struct {
		version *semver.Version
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "should set the version parameter in struct",
			fields: fields{
				Version: testVersion1,
			},
			args: args{
				version: testVersion2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &ComponentInstallation{
				Name:                tt.fields.Name,
				DeployNamespace:     tt.fields.DeployNamespace,
				Version:             tt.fields.Version,
				Status:              tt.fields.Status,
				ValuesYamlOverwrite: tt.fields.ValuesYamlOverwrite,
				MappedValues:        tt.fields.MappedValues,
				PersistenceContext:  tt.fields.PersistenceContext,
				Health:              tt.fields.Health,
			}
			ci.Upgrade(tt.args.version)
			assert.Equal(t, tt.args.version, ci.Version)
		})
	}
}
