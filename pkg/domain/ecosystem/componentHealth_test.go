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
				AvailableHealthStatus: {ComponentName("k8s-dogu-operator")},
			},
			contains:    []string{"0 components are unhealthy: "},
			notContains: []string{"k8s-dogu-operator"},
		},
		{
			name: "any components not available should be unhealthy",
			healthStates: map[HealthStatus][]ComponentName{
				AvailableHealthStatus:    {ComponentName("k8s-blueprint-operator")},
				UnavailableHealthStatus:  {ComponentName("k8s-etcd"), ComponentName("k8s-dogu-operator")},
				NotInstalledHealthStatus: {ComponentName("k8s-service-discovery")},
				"other":                  {ComponentName("k8s-component-operator")},
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
