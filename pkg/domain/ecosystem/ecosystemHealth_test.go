package ecosystem

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/stretchr/testify/assert"
)

func TestHealthResult_String(t *testing.T) {
	type fields struct {
		DoguHealth      DoguHealthResult
		ComponentHealth ComponentHealthResult
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should print dogu and component health results with no unhealthy",
			fields: fields{
				DoguHealth:      DoguHealthResult{},
				ComponentHealth: ComponentHealthResult{},
			},
			want: "0 dogu(s) are unhealthy: \n  0 component(s) are unhealthy: ",
		},
		{
			name: "should print dogu and component health results with unhealthy",
			fields: fields{
				DoguHealth:      DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{UnavailableHealthStatus: {"nginx-ingress"}}},
				ComponentHealth: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]common.SimpleComponentName{UnavailableHealthStatus: {"k8s-etcd"}}},
			},
			want: "1 dogu(s) are unhealthy: nginx-ingress\n  1 component(s) are unhealthy: k8s-etcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HealthResult{
				DoguHealth:      tt.fields.DoguHealth,
				ComponentHealth: tt.fields.ComponentHealth,
			}
			assert.Equalf(t, tt.want, result.String(), "String()")
		})
	}
}

func TestHealthResult_AllHealthy(t *testing.T) {
	type fields struct {
		DoguHealth      DoguHealthResult
		ComponentHealth ComponentHealthResult
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should be healthy if no dogus or components are unavailable",
			fields: fields{
				DoguHealth:      DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{AvailableHealthStatus: {"nginx-ingress"}}},
				ComponentHealth: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]common.SimpleComponentName{AvailableHealthStatus: {"k8s-etcd"}}},
			},
			want: true,
		},
		{
			name: "should be unhealthy if dogus are unavailable",
			fields: fields{
				DoguHealth:      DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{UnavailableHealthStatus: {"nginx-ingress"}}},
				ComponentHealth: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]common.SimpleComponentName{AvailableHealthStatus: {"k8s-etcd"}}},
			},
			want: false,
		},
		{
			name: "should be unhealthy if components are unavailable",
			fields: fields{
				DoguHealth:      DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{AvailableHealthStatus: {"nginx-ingress"}}},
				ComponentHealth: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]common.SimpleComponentName{UnavailableHealthStatus: {"k8s-etcd"}}},
			},
			want: false,
		},
		{
			name: "should be unhealthy if dogus and components are unavailable",
			fields: fields{
				DoguHealth:      DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{UnavailableHealthStatus: {"nginx-ingress"}}},
				ComponentHealth: ComponentHealthResult{ComponentsByStatus: map[HealthStatus][]common.SimpleComponentName{UnavailableHealthStatus: {"k8s-etcd"}}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HealthResult{
				DoguHealth:      tt.fields.DoguHealth,
				ComponentHealth: tt.fields.ComponentHealth,
			}
			assert.Equalf(t, tt.want, result.AllHealthy(), "AllHealthy()")
		})
	}
}
