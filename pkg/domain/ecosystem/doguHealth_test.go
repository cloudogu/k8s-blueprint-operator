package ecosystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateDoguHealthResult(t *testing.T) {
	tests := []struct {
		name  string
		dogus []*DoguInstallation
		want  DoguHealthResult
	}{
		{
			name: "",
			dogus: []*DoguInstallation{
				{
					Name:   "postgresql",
					Health: AvailableHealthStatus,
				},
				{
					Name:   "postfix",
					Health: UnavailableHealthStatus,
				},
				{
					Name:   "ldap",
					Health: PendingHealthStatus,
				},
			},
			want: DoguHealthResult{
				DogusByStatus: map[HealthStatus][]DoguName{
					AvailableHealthStatus:   {"postgresql"},
					UnavailableHealthStatus: {"postfix"},
					PendingHealthStatus:     {"ldap"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CalculateDoguHealthResult(tt.dogus), "CalculateDoguHealthResult(%v)", tt.dogus)
		})
	}
}

func TestDoguHealthResult_String(t *testing.T) {
	tests := []struct {
		name         string
		healthStates map[HealthStatus][]DoguName
		contains     []string
		notContains  []string
	}{
		{
			name:         "no dogus should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]DoguName{},
			contains:     []string{"0 dogu(s) are unhealthy: "},
		},
		{
			name: "only available dogus should result in 0 components unhealthy",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus: {"nginx-ingress"},
			},
			contains:    []string{"0 dogu(s) are unhealthy: "},
			notContains: []string{"nginx-ingress"},
		},
		{
			name: "any dogus not available should be unhealthy",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus:   {"nginx-static"},
				UnavailableHealthStatus: {"postgresql", "redmine"},
				"other":                 {"scm"},
			},
			contains: []string{
				"3 dogu(s) are unhealthy: ",
				"postgresql",
				"redmine",
				"scm",
			},
			notContains: []string{"nginx-static"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DoguHealthResult{DogusByStatus: tt.healthStates}
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

func TestDoguHealthResult_AllHealthy(t *testing.T) {
	tests := []struct {
		name         string
		healthStates map[HealthStatus][]DoguName
		want         bool
	}{
		{
			name:         "should be healthy if empty",
			healthStates: map[HealthStatus][]DoguName{},
			want:         true,
		},
		{
			name: "should be healthy if all are available",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus: {"nginx-ingress", "nginx-static", "postfix"},
			},
			want: true,
		},
		{
			name: "should not be healthy if one is not available",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus:   {"nginx-ingress", "nginx-static", "postfix"},
				UnavailableHealthStatus: {"ldap"},
			},
			want: false,
		},
		{
			name: "should not be healthy if multiple are not available",
			healthStates: map[HealthStatus][]DoguName{
				AvailableHealthStatus:   {"nginx-ingress", "nginx-static", "postfix"},
				UnavailableHealthStatus: {"ldap", "redmine"},
				PendingHealthStatus:     {"postgresql"},
				"other":                 {"plantuml"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DoguHealthResult{DogusByStatus: tt.healthStates}
			assert.Equalf(t, tt.want, result.AllHealthy(), "AllHealthy()")
		})
	}
}
