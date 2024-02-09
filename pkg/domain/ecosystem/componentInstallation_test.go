package ecosystem

import (
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testNamespace             = "k8s"
	testComponentName         = "k8s-dogu-operator"
	testLonghornComponentName = "k8s-longhorn"
)

var (
	testVersion1, _ = semver.NewVersion("1.0.0")
	testVersion2, _ = semver.NewVersion("2.0.0")
)

func TestInstallComponent(t *testing.T) {
	type args struct {
		namespace     string
		componentName string
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
				namespace:     testNamespace,
				componentName: testComponentName,
				version:       testVersion1,
			},
			want: &ComponentInstallation{
				Name:                  testComponentName,
				DistributionNamespace: testNamespace,
				Version:               testVersion1,
			},
		},
		{
			name: "longhorn should always be deployed in longhorn-system",
			args: args{
				namespace:     testNamespace,
				componentName: testLonghornComponentName,
				version:       testVersion1,
			},
			want: &ComponentInstallation{
				Name:                  testLonghornComponentName,
				DistributionNamespace: testNamespace,
				Version:               testVersion1,
				DeployNamespace:       "longhorn-system",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, InstallComponent(tt.args.namespace, tt.args.componentName, tt.args.version), "InstallComponent(%v, %v, %v)", tt.args.namespace, tt.args.componentName, tt.args.version)
		})
	}
}

func TestComponentInstallation_Upgrade(t *testing.T) {
	type fields struct {
		Name                  string
		DistributionNamespace string
		DeployNamespace       string
		Version               *semver.Version
		Status                string
		ValuesYamlOverwrite   string
		MappedValues          map[string]string
		PersistenceContext    map[string]interface{}
		Health                HealthStatus
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
				Name:                  tt.fields.Name,
				DistributionNamespace: tt.fields.DistributionNamespace,
				DeployNamespace:       tt.fields.DeployNamespace,
				Version:               tt.fields.Version,
				Status:                tt.fields.Status,
				ValuesYamlOverwrite:   tt.fields.ValuesYamlOverwrite,
				MappedValues:          tt.fields.MappedValues,
				PersistenceContext:    tt.fields.PersistenceContext,
				Health:                tt.fields.Health,
			}
			ci.Upgrade(tt.args.version)
			assert.Equal(t, tt.args.version, ci.Version)
		})
	}
}
