package ecosystem

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/stretchr/testify/assert"
)

func TestHealthResult_String(t *testing.T) {
	type fields struct {
		DoguHealth DoguHealthResult
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should print dogu health results with no unhealthy",
			fields: fields{
				DoguHealth: DoguHealthResult{},
			},
			want: "0 dogu(s) are unhealthy: ",
		},
		{
			name: "should print dogu health results with unhealthy",
			fields: fields{
				DoguHealth: DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{UnavailableHealthStatus: {"nginx-ingress"}}},
			},
			want: "1 dogu(s) are unhealthy: nginx-ingress",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HealthResult{
				DoguHealth: tt.fields.DoguHealth,
			}
			assert.Equalf(t, tt.want, result.String(), "String()")
		})
	}
}

func TestHealthResult_AllHealthy(t *testing.T) {
	type fields struct {
		DoguHealth DoguHealthResult
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should be healthy if no dogus are unavailable",
			fields: fields{
				DoguHealth: DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{AvailableHealthStatus: {"nginx-ingress"}}},
			},
			want: true,
		},
		{
			name: "should be unhealthy if dogus are unavailable",
			fields: fields{
				DoguHealth: DoguHealthResult{DogusByStatus: map[HealthStatus][]cescommons.SimpleName{UnavailableHealthStatus: {"nginx-ingress"}}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HealthResult{
				DoguHealth: tt.fields.DoguHealth,
			}
			assert.Equalf(t, tt.want, result.AllHealthy(), "AllHealthy()")
		})
	}
}
