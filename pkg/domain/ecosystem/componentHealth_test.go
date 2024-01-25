package ecosystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComponentHealthResult_String(t *testing.T) {
	tests := []struct {
		name         string
		healthStates map[HealthStatus][]ComponentName
		contains     []string
		notContains  []string
	}{
		{
			name:         "no components should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]ComponentName{},
			contains:     []string{"0 components are unhealthy: "},
		},
		{
			name: "only available components should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]ComponentName{
				AvailableHealthStatus: {"k8s-dogu-operator"},
			},
			contains:    []string{"0 components are unhealthy: "},
			notContains: []string{"k8s-dogu-operator"},
		},
		{
			name: "any components not available should be unhealthy",
			healthStates: map[HealthStatus][]ComponentName{
				AvailableHealthStatus:    {"k8s-blueprint-operator"},
				UnavailableHealthStatus:  {"k8s-etcd", "k8s-dogu-operator"},
				NotInstalledHealthStatus: {"k8s-service-discovery"},
				"other":                  {"k8s-component-operator"},
			},
			contains: []string{
				"4 components are unhealthy: ",
				"k8s-etcd",
				"k8s-dogu-operator",
				"k8s-service-discovery",
				"k8s-component-operator",
			},
			notContains: []string{"ks8-blueprint-operator"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComponentHealthResult{ComponentsByStatus: tt.healthStates}
			actual := result.String()
			for _, contains := range tt.contains {
				assert.Contains(t, actual, contains)
			}
			for _, notContains := range tt.notContains {
				assert.NotContains(t, actual, notContains)
			}
		})
	}
}

func TestComponentHealthResult_AllHealthy(t *testing.T) {
	tests := []struct {
		name         string
		healthStates map[HealthStatus][]ComponentName
		want         bool
	}{
		{
			name:         "should be healthy if empty",
			healthStates: map[HealthStatus][]ComponentName{},
			want:         true,
		},
		{
			name: "should be healthy if all are available",
			healthStates: map[HealthStatus][]ComponentName{
				AvailableHealthStatus: {"k8s-blueprint-operator", "k8s-etcd", "k8s-service-discovery"},
			},
			want: true,
		},
		{
			name: "should not be healthy if one is not available",
			healthStates: map[HealthStatus][]ComponentName{
				AvailableHealthStatus:   {"k8s-blueprint-operator", "k8s-etcd", "k8s-service-discovery"},
				UnavailableHealthStatus: {"k8s-dogu-operator"},
			},
			want: false,
		},
		{
			name: "should not be healthy if multiple are not available",
			healthStates: map[HealthStatus][]ComponentName{
				AvailableHealthStatus:    {"k8s-blueprint-operator", "k8s-etcd", "k8s-service-discovery"},
				UnavailableHealthStatus:  {"k8s-dogu-operator", "k8s-component-operator"},
				PendingHealthStatus:      {"k8s-longhorn"},
				NotInstalledHealthStatus: {"k8s-backup-operator"},
				"other":                  {"k8s-velero"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComponentHealthResult{ComponentsByStatus: tt.healthStates}
			assert.Equalf(t, tt.want, result.AllHealthy(), "AllHealthy()")
		})
	}
}

func TestCalculateComponentHealthResult(t *testing.T) {
	type args struct {
		installedComponents map[string]*ComponentInstallation
		requiredComponents  []RequiredComponent
	}
	tests := []struct {
		name string
		args args
		want ComponentHealthResult
	}{
		{
			name: "result should be empty for no required and no installed components",
			args: args{
				installedComponents: map[string]*ComponentInstallation{},
				requiredComponents:  []RequiredComponent{},
			},
			want: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]ComponentName{}},
		},
		{
			name: "result should contain components that are not installed but required",
			args: args{
				installedComponents: map[string]*ComponentInstallation{},
				requiredComponents:  []RequiredComponent{{Name: "k8s-etcd"}, {Name: "k8s-service-discovery"}},
			},
			want: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]ComponentName{
				NotInstalledHealthStatus: {"k8s-etcd", "k8s-service-discovery"},
			}},
		},
		{
			name: "result should contain any components with their health state",
			args: args{
				installedComponents: map[string]*ComponentInstallation{
					"k8s-blueprint-operator": {Name: "k8s-blueprint-operator", Health: AvailableHealthStatus},
					"k8s-dogu-operator":      {Name: "k8s-dogu-operator", Health: UnavailableHealthStatus},
					"k8s-component-operator": {Name: "k8s-component-operator", Health: UnavailableHealthStatus},
					"k8s-longhorn":           {Name: "k8s-longhorn", Health: PendingHealthStatus},
					"k8s-velero":             {Name: "k8s-velero", Health: "other"},
				},
				requiredComponents: []RequiredComponent{},
			},
			want: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]ComponentName{
				AvailableHealthStatus:   {"k8s-blueprint-operator"},
				UnavailableHealthStatus: {"k8s-dogu-operator", "k8s-component-operator"},
				PendingHealthStatus:     {"k8s-longhorn"},
				"other":                 {"k8s-velero"},
			}},
		},
		{
			name: "result should contain any components with their health state and components that are not installed but required",
			args: args{
				installedComponents: map[string]*ComponentInstallation{
					"k8s-blueprint-operator": {Name: "k8s-blueprint-operator", Health: AvailableHealthStatus},
					"k8s-dogu-operator":      {Name: "k8s-dogu-operator", Health: UnavailableHealthStatus},
					"k8s-component-operator": {Name: "k8s-component-operator", Health: UnavailableHealthStatus},
					"k8s-longhorn":           {Name: "k8s-longhorn", Health: PendingHealthStatus},
					"k8s-velero":             {Name: "k8s-velero", Health: "other"},
				},
				requiredComponents: []RequiredComponent{{Name: "k8s-etcd"}, {Name: "k8s-service-discovery"}},
			},
			want: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]ComponentName{
				AvailableHealthStatus:    {"k8s-blueprint-operator"},
				UnavailableHealthStatus:  {"k8s-dogu-operator", "k8s-component-operator"},
				PendingHealthStatus:      {"k8s-longhorn"},
				"other":                  {"k8s-velero"},
				NotInstalledHealthStatus: {"k8s-etcd", "k8s-service-discovery"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CalculateComponentHealthResult(tt.args.installedComponents, tt.args.requiredComponents), "CalculateComponentHealthResult(%v, %v)", tt.args.installedComponents, tt.args.requiredComponents)
		})
	}
}
