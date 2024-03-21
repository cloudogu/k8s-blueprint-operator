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
		deployConfig  DeployConfig
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
				deployConfig:  map[string]interface{}{"deployNamespace": "longhorn-system"},
			},
			want: &ComponentInstallation{
				Name:            testComponentName,
				ExpectedVersion: testVersion1,
				DeployConfig:    map[string]interface{}{"deployNamespace": "longhorn-system"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, InstallComponent(tt.args.componentName, tt.args.version, tt.args.deployConfig), "InstallComponent(%v, %v, %v)", tt.args.componentName, tt.args.version)
		})
	}
}

func TestComponentInstallation_Upgrade(t *testing.T) {
	type fields struct {
		Name               common.QualifiedComponentName
		Version            *semver.Version
		Status             string
		PersistenceContext map[string]interface{}
		Health             HealthStatus
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
				Name:               tt.fields.Name,
				ExpectedVersion:    tt.fields.Version,
				Status:             tt.fields.Status,
				PersistenceContext: tt.fields.PersistenceContext,
				Health:             tt.fields.Health,
			}
			ci.Upgrade(tt.args.version)
			assert.Equal(t, tt.args.version, ci.ExpectedVersion)
		})
	}
}

func TestComponentInstallation_UpdateDeployConfig(t *testing.T) {
	t.Run("should set config", func(t *testing.T) {
		// given
		sut := ComponentInstallation{}
		config := map[string]interface{}{"key": "value"}

		// when
		sut.UpdateDeployConfig(config)

		// then
		assert.Equal(t, DeployConfig(config), sut.DeployConfig)
	})
}
